package logger

import (
	"log"

	"github.com/spf13/viper"
)

var logWriter *Logger

// InitLogger 初始程序日志
func InitLogger() {
	env := viper.GetString("env")
	conf := viper.GetStringMap("logger")
	fileDir := conf["filedir"].(string)
	fileName := conf["filename"].(string)
	level := conf["level"].(int)
	maxDay := conf["maxday"].(int)
	maxSize := conf["maxsize"].(int)

	newLog, err := NewLogger(env, fileDir, fileName, level, maxDay, maxSize)
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
