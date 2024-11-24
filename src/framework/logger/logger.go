package logger

import (
	"mask_api_gin/src/framework/config"

	"log"
)

var logWriter *Logger

// InitLogger 初始程序日志
func InitLogger() {
	fileDir := config.Get("logger.fileDir").(string)
	fileName := config.Get("logger.fileName").(string)
	level := config.Get("logger.level").(int)
	maxDay := config.Get("logger.maxDay").(int)
	maxSize := config.Get("logger.maxSize").(int)

	newLog, err := NewLogger(config.Env(), fileDir, fileName, level, maxDay, maxSize)
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}

	logWriter = newLog
}

// Close 关闭程序日志写入
func Close() {
	logWriter.Close()
}

// Infof 信息级日志
func Infof(format string, v ...any) {
	logWriter.Infof(format, v...)
}

// Warnf 告警级日志
func Warnf(format string, v ...any) {
	logWriter.Warnf(format, v...)
}

// Errorf 错误级日志
func Errorf(format string, v ...any) {
	logWriter.Errorf(format, v...)
}

// Fatalf 抛出错误并退出程序
func Fatalf(format string, v ...any) {
	log.Fatalf(format, v...)
}
