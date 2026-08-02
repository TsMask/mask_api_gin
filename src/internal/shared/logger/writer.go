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

// parseLevelMap 日志等级映射表，数值越大优先级越高：silent < info < warn < error
var parseLevelMap = map[string]int{
	"SILENT": 0,
	"INFO":   1,
	"WARN":   2,
	"ERROR":  3,
}

// Logger 日志器对象
type Logger struct {
	Env        string        `json:"env"`      // 运行环境
	Level      string        `json:"level"`    // 日志等级阈值
	MaxDay     int           `json:"maxDay"`   // 保留最长天数
	MaxSize    int64         `json:"maxSize"`  // 单文件最大字节数
	FilePath   string        `json:"filePath"` // 文件路径
	FileName   string        `json:"fileName"` // 文件名
	fileHandle *os.File      `json:"-"`        // 文件句柄
	LogDate    string        `json:"logDate"`  // 当前日志所属日期
	logNew     *log.Logger   `json:"-"`        // 日志实例
	bufWriter  *bufio.Writer `json:"-"`        // 缓存写入器
	mu         sync.Mutex    `json:"-"`        // 保护写入与分割互斥，避免并发占用句柄
	closed     atomic.Bool   `json:"-"`        // 是否已关闭，阻止关闭后再分割/写入
}

// createFile 创建日志文件与缓存写入器
func createFile(fileDir string, fileName string) (*os.File, *bufio.Writer, error) {
	logFilePath := filepath.Join(fileDir, fileName)
	if err := os.MkdirAll(filepath.Dir(logFilePath), 0755); err != nil {
		return nil, nil, fmt.Errorf("failed to mkdir logger dir: %v", err)
	}
	fileHandle, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open log file: %v", err)
	}

	bufWriter := bufio.NewWriterSize(fileHandle, 4*1024)
	return fileHandle, bufWriter, nil
}

// checkFile 定时检查并执行文件分割，自行调度下一次检查。
func (l *Logger) checkFile() {
	if l.closed.Load() {
		return
	}

	// 整段在锁内：读取文件信息 + 判断 + 分割，避免与写入/关闭并发；
	// LogDate 在锁内更新，保证状态一致。
	var removeBefore time.Time
	needRemove := false

	l.mu.Lock()
	fileInfo, err := l.fileHandle.Stat()
	if err != nil {
		fmt.Fprintf(os.Stderr, "logger: failed to get log file info: %v\n", err)
	} else {
		currTime := time.Now()
		// 任何轮转后都按 MaxDay 清理过期归档，阈值统一为 now - MaxDay
		removeBefore = currTime.AddDate(0, 0, -l.MaxDay)
		switch {
		case l.LogDate != currTime.Format(time.DateOnly):
			// 跨天：把当前文件归档为"昨天"对应日期
			l.LogDate = currTime.Format(time.DateOnly)
			l.rotateLocked(currTime.AddDate(0, 0, -1).Format(time.DateOnly))
			needRemove = true
		case fileInfo.Size() >= l.MaxSize:
			// 超过单文件大小上限，按时间戳归档
			l.rotateLocked(currTime.Format("2006-01-02_150405"))
			needRemove = true
		}
	}
	l.mu.Unlock()

	// 过期清理只依赖初始化后不可变的字段，锁外执行，避免长时间持锁阻塞写入
	if needRemove {
		l.removeOldFile(removeBefore)
	}

	// 无论是否触发分割都继续定时，避免一次 Stat 失败就永久停止分割
	if !l.closed.Load() {
		time.AfterFunc(1*time.Minute, l.checkFile)
	}
}

// rotateLocked 把当前日志文件归档并重新打开基础文件。调用方必须持有 l.mu。
func (l *Logger) rotateLocked(timeFormat string) {
	if l.closed.Load() {
		return
	}

	// 必须先刷新缓冲区：否则缓冲数据丢失，且 Windows 上缓冲器仍引用已关闭的句柄，
	// 导致句柄占用、文件无法重命名分割
	if l.bufWriter != nil {
		_ = l.bufWriter.Flush()
	}

	// 关闭当前文件句柄（句柄释放后方可重命名）
	if err := l.fileHandle.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "logger: failed to close log file: %v\n", err)
	}

	archivePath := filepath.Join(l.FilePath, fmt.Sprintf("%s.%s", l.FileName, timeFormat))
	basePath := filepath.Join(l.FilePath, l.FileName)

	// 归档重命名
	if err := os.Rename(basePath, archivePath); err != nil {
		fmt.Fprintf(os.Stderr, "logger: failed to rename log file: %v\n", err)
	}

	// 重新打开基础文件（必须是基础文件名，而非归档文件名，否则后续写入全部进入归档文件）
	fileHandle, bufWriter, err := createFile(l.FilePath, l.FileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "logger: failed to open log file: %v\n", err)
		return
	}
	writer := io.Writer(bufWriter)
	if l.Env == "local" {
		writer = io.MultiWriter(bufWriter, os.Stderr)
	}

	// 重新绑定句柄、缓冲器与日志输出器
	l.fileHandle = fileHandle
	l.bufWriter = bufWriter
	l.logNew.SetOutput(writer)
}

// removeOldFile 删除早于 oldFileDate 的归档文件。
// 仅依赖初始化后不可变的字段，无需持有 l.mu。
func (l *Logger) removeOldFile(oldFileDate time.Time) {
	if l.MaxDay < 1 {
		return
	}

	files, err := os.ReadDir(l.FilePath)
	if err != nil {
		l.Errorf("logger removeOldFile ReadDir err: %v", err.Error())
		return
	}

	prefix := l.FileName + "."
	for _, file := range files {
		// 归档文件名格式：fileName.date 或 fileName.date_time，截取 fileName. 之后的后缀
		suffix := strings.TrimPrefix(file.Name(), prefix)
		if suffix == file.Name() {
			continue
		}
		// 日期部分固定 10 位（2006-01-02），后缀不足则不是归档文件，跳过避免越界
		if len(suffix) < 10 {
			continue
		}
		fileDate, err := time.Parse("2006-01-02", suffix[:10])
		if err != nil {
			continue
		}

		// 文件日期早于给定阈值则删除
		if fileDate.Before(oldFileDate) {
			if err := os.Remove(filepath.Join(l.FilePath, file.Name())); err != nil {
				l.Errorf("logger removeOldFile Remove err: %v", err.Error())
			}
		}
	}
}

// writeLog 写入日志
func (l *Logger) writeLog(level, format string, args ...any) {
	// 等级过滤：低于阈值不输出
	levelVal, okLevel := parseLevelMap[level]
	if !okLevel {
		return
	}
	threshold, okThreshold := parseLevelMap[l.Level]
	if !okThreshold || levelVal < threshold {
		return
	}

	logMsg := fmt.Sprintf("[%s] %s\n", level, fmt.Sprintf(format, args...))

	// 关闭检查置于锁内，避免 Close 在检查与写入之间关闭句柄导致写入已关闭的 fd
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.closed.Load() {
		return
	}
	_ = l.logNew.Output(4, logMsg)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.writeLog("INFO", format, args...)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.writeLog("WARN", format, args...)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.writeLog("ERROR", format, args...)
}

// Fatal 记录致命错误并立即落盘（不受 level 阈值限制），随后由调用方退出进程。
// 始终同时写入文件与 stderr：进程即将退出，stderr 保证 systemd/journald、容器
// 编排等通过 stderr 采集致命错误，不依赖日志文件是否可读。
// 用独立输出器实现，避免与常规日志的 env 路由耦合，且 local 下不与 logNew 重复。
func (l *Logger) Fatal(format string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.closed.Load() {
		return
	}
	logMsg := fmt.Sprintf("[FATAL] %s\n", fmt.Sprintf(format, args...))
	out := io.MultiWriter(l.bufWriter, os.Stderr)
	fatalLogger := log.New(out, "", log.LstdFlags|log.Lshortfile)
	_ = fatalLogger.Output(3, logMsg)
	if l.bufWriter != nil {
		_ = l.bufWriter.Flush()
	}
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
