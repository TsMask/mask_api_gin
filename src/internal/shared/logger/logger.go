package logger

import (
	"io"
	"log"
	"os"
	"strings"
	"time"

	"mask_api_gin/src/internal/config"
)

var logger *Logger // 日志实例

// InitLogger 初始程序日志
func InitLogger() {
	// 重复初始化时先关闭旧实例，避免句柄与分割 goroutine 泄漏
	if logger != nil {
		logger.Close()
	}
	fileDir := config.GetStringOr("logger.fileDir", "/tmp")
	fileName := config.GetStringOr("logger.fileName", "app.log")
	level := config.GetStringOr("logger.level", "warn")
	maxDay := config.GetIntOr("logger.maxDay", 7)
	maxSize := config.GetIntOr("logger.maxSize", 10)
	env := config.Env()

	fileHandle, bufWriter, err := createFile(fileDir, fileName)
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}
	writer := io.Writer(bufWriter)
	if env == "local" {
		writer = io.MultiWriter(bufWriter, os.Stderr)
	}

	logNew := log.New(writer, "", log.LstdFlags|log.Lshortfile)
	logger = &Logger{
		Env:        env,
		FilePath:   fileDir,
		FileName:   fileName,
		Level:      strings.ToUpper(level),
		MaxDay:     maxDay,
		MaxSize:    int64(maxSize * 1024 * 1024),
		fileHandle: fileHandle,
		logNew:     logNew,
		bufWriter:  bufWriter,
		LogDate:    time.Now().Format(time.DateOnly),
	}

	go logger.checkFile()
}

// Close 关闭程序日志写入
func Close() {
	if logger != nil {
		logger.Close()
	}
}

// Infof 信息级日志
func Infof(format string, v ...any) {
	if logger == nil {
		log.Printf("[INFO] "+format, v...)
		return
	}
	logger.Infof(format, v...)
}

// Warnf 告警级日志
func Warnf(format string, v ...any) {
	if logger == nil {
		log.Printf("[WARN] "+format, v...)
		return
	}
	logger.Warnf(format, v...)
}

// Errorf 错误级日志
func Errorf(format string, v ...any) {
	if logger == nil {
		log.Printf("[ERROR] "+format, v...)
		return
	}
	logger.Errorf(format, v...)
}

// Fatalf 记录致命错误并落盘后退出程序。
func Fatalf(format string, v ...any) {
	if logger == nil {
		log.Fatalf(format, v...)
		return
	}
	logger.Fatal(format, v...)
	os.Exit(1)
}
