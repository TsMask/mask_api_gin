package logger

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// Logger 日志器对象
type Logger struct {
	env         string         // 运行环境
	filePath    string         // 文件路径
	fileName    string         // 文件名
	level       int            // 日志等级标识
	maxDay      int            // 保留最长天数
	maxSize     int64          // 文件最大空间
	fileHandle  *os.File       // 文件实例
	bufWriter   *bufio.Writer  // 缓冲写入器
	logger      *log.Logger    // 日志实例
	logLevelMap map[int]string // 日志等级标识名
	logDay      int            // 日志当前日
	mu          sync.Mutex     // 保护写入与分割互斥，避免并发占用句柄
	closed      atomic.Bool    // 是否已关闭，阻止关闭后再分割/写入
}

const (
	LogLevelSilent = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
)

// NewLogger 实例日志器对象
func NewLogger(env, fileDir, fileName string, level, maxDay, maxSize int) (*Logger, error) {
	logFilePath := filepath.Join(fileDir, fileName)
	if err := os.MkdirAll(filepath.Dir(logFilePath), 0755); err != nil {
		return nil, fmt.Errorf("failed to mkdir logger dir: %v", err)
	}
	fileHandle, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %v", err)
	}

	bufWriter := bufio.NewWriterSize(fileHandle, 64*1024)

	writer := io.Writer(bufWriter)
	if env == "local" {
		writer = io.MultiWriter(bufWriter, os.Stderr)
	}

	logger := log.New(writer, "", log.LstdFlags|log.Lshortfile)

	logLevelMap := map[int]string{
		LogLevelSilent: "SILENT",
		LogLevelInfo:   "INFO",
		LogLevelWarn:   "WARN",
		LogLevelError:  "ERROR",
	}

	stdLogger := &Logger{
		env:         env,
		filePath:    fileDir,
		fileName:    fileName,
		level:       level,
		maxDay:      maxDay,
		maxSize:     int64(maxSize * 1024 * 1024),
		fileHandle:  fileHandle,
		bufWriter:   bufWriter,
		logger:      logger,
		logLevelMap: logLevelMap,
		logDay:      time.Now().Day(),
	}

	go stdLogger.checkFile()

	return stdLogger, nil
}

// checkFile 检查文件分割，自定时调用
func (l *Logger) checkFile() {
	if l.closed.Load() {
		return
	}

	// 加锁读取文件信息，避免与 rotateFile/Close 并发访问句柄
	l.mu.Lock()
	fileInfo, err := l.fileHandle.Stat()
	l.mu.Unlock()
	if err != nil {
		fmt.Fprintf(os.Stderr, "logger: failed to get log file info: %v\n", err)
	} else {
		currTime := time.Now()
		if l.logDay != currTime.Day() {
			l.logDay = currTime.Day()
			l.rotateFile(currTime.AddDate(0, 0, -1).Format("2006-01-02"))
			// 移除超过保存最长天数的文件
			l.removeOldFile(currTime.AddDate(0, 0, -l.maxDay))
		} else if fileInfo.Size() >= l.maxSize {
			l.rotateFile(currTime.Format("2006-01-02_150405"))
		} else if time.Since(fileInfo.ModTime()).Hours() > 24 {
			l.rotateFile(fileInfo.ModTime().Format("2006-01-02"))
		}
	}

	// 无论是否触发分割都继续定时，避免一次 Stat 失败就永久停止分割
	if !l.closed.Load() {
		time.AfterFunc(1*time.Minute, l.checkFile)
	}
}

// rotateFile 检查文件大小进行分割
func (l *Logger) rotateFile(timeFormat string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.closed.Load() {
		return
	}

	// 必须先刷新缓冲区：否则缓冲数据丢失，且缓冲器仍引用已关闭的句柄，
	// Windows 上会导致句柄占用、文件无法重命名分割
	if l.bufWriter != nil {
		_ = l.bufWriter.Flush()
	}

	// 关闭当前文件句柄（句柄释放后方可重命名）
	if err := l.fileHandle.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "logger: failed to close log file: %v\n", err)
	}

	newFileName := fmt.Sprintf("%s.%s", l.fileName, timeFormat)
	newFilePath := filepath.Join(l.filePath, newFileName)
	oldFilePath := filepath.Join(l.filePath, l.fileName)

	// 重命名
	if err := os.Rename(oldFilePath, newFilePath); err != nil {
		fmt.Fprintf(os.Stderr, "logger: failed to rename log file: %v\n", err)
	}

	// 新文件句柄
	fileHandle, err := os.OpenFile(oldFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Fprintf(os.Stderr, "logger: failed to open log file: %v\n", err)
		return
	}

	l.fileHandle = fileHandle

	// 重新创建缓冲写入器
	l.bufWriter = bufio.NewWriterSize(fileHandle, 64*1024)

	// 重新设置 logger 的 writer
	writer := io.Writer(l.bufWriter)
	if l.env == "local" {
		writer = io.MultiWriter(l.bufWriter, os.Stderr)
	}
	l.logger.SetOutput(writer)
}

// RemoveOldFile 删除旧文件
func (l *Logger) removeOldFile(oldFileDate time.Time) {
	// 遍历目标文件夹中的文件
	files, err := os.ReadDir(l.filePath)
	if err != nil {
		l.Errorf("logger RemoveOldFile ReadDir err: %v", err.Error())
		return
	}

	prefix := l.fileName + "."
	for _, file := range files {
		// 分割文件名格式：fileName.date 或 fileName.date_time，截取 fileName. 之后的后缀
		suffix := strings.TrimPrefix(file.Name(), prefix)
		if suffix == file.Name() {
			continue
		}
		// 日期部分固定 10 位（2006-01-02），后缀不足则不是分割文件，跳过避免越界
		if len(suffix) < 10 {
			continue
		}
		fileDate, err := time.Parse("2006-01-02", suffix[:10])
		if err != nil {
			continue
		}

		// 判断文件日期是否在给定日期之前
		if fileDate.Before(oldFileDate) {
			if err := os.Remove(filepath.Join(l.filePath, file.Name())); err != nil {
				l.Errorf("logger RemoveOldFile Remove err: %v", err.Error())
				continue
			}
		}
	}
}

// writeLog 写入chan
func (l *Logger) writeLog(level int, format string, args ...interface{}) {
	// 日志等级小于指定等级不输出文件
	if level < l.level {
		return
	}
	if l.closed.Load() {
		return
	}

	logMsg := fmt.Sprintf("[%s] %s\n", l.logLevelMap[level], fmt.Sprintf(format, args...))
	l.mu.Lock()
	_ = l.logger.Output(4, logMsg)
	l.mu.Unlock()
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.writeLog(LogLevelInfo, format, args...)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.writeLog(LogLevelWarn, format, args...)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.writeLog(LogLevelError, format, args...)
}

// Close 日志关闭
func (l *Logger) Close() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.closed.Store(true)

	if l.bufWriter != nil {
		_ = l.bufWriter.Flush()
	}

	if err := l.fileHandle.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "logger: failed to close log file: %v\n", err)
	}
}
