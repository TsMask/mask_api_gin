package logger

import (
	"mask_api_gin/src/internal/config"

	"log"
)

var logWriter *Logger

// InitLogger 初始程序日志
func InitLogger() {
	fileDir := config.GetStringOr("logger.fileDir", "")
	fileName := config.GetStringOr("logger.fileName", "")
	level := config.GetIntOr("logger.level", 0)
	maxDay := config.GetIntOr("logger.maxDay", 0)
	maxSize := config.GetIntOr("logger.maxSize", 0)

	newLog, err := NewLogger(config.Env(), fileDir, fileName, level, maxDay, maxSize)
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}

	logWriter = newLog
}

// Close 关闭程序日志写入
func Close() {
	if logWriter != nil {
		logWriter.Close()
	}
}

// Infof 信息级日志
func Infof(format string, v ...any) {
	if logWriter == nil {
		log.Printf("[INFO] "+format, v...)
		return
	}
	logWriter.Infof(format, v...)
}

// Warnf 告警级日志
func Warnf(format string, v ...any) {
	if logWriter == nil {
		log.Printf("[WARN] "+format, v...)
		return
	}
	logWriter.Warnf(format, v...)
}

// Errorf 错误级日志
func Errorf(format string, v ...any) {
	if logWriter == nil {
		log.Printf("[ERROR] "+format, v...)
		return
	}
	logWriter.Errorf(format, v...)
}

// Fatalf 抛出错误并退出程序
func Fatalf(format string, v ...any) {
	log.Fatalf(format, v...)
}
