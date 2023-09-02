package logger

import (
	"fmt"
	"log"
	"os"
)

const (
	// Silent silent log level
	silent int = iota
	// Info info log level
	info
	// Warn warn log level
	warn
	// Error error log level
	err
)

var logLevelMap = map[int]string{
	silent: "silent",
	info:   "info",
	warn:   "warn",
	err:    "error",
}

func logWithLevel(level int, format string, v ...any) {
	if level <= silent {
		return
	}

	stdLog := log.New(os.Stdout, "["+logLevelMap[level]+"] ", log.LstdFlags|log.Lshortfile)
	stdLog.Output(3, fmt.Sprintf(format, v...))
}

func Infof(format string, v ...any) {
	logWithLevel(info, format, v...)
}

func Warnf(format string, v ...any) {
	logWithLevel(warn, format, v...)
}

func Errorf(format string, v ...any) {
	logWithLevel(err, format, v...)
}

// Fatalf 抛出错误并退出程序
func Fatalf(format string, v ...any) {
	log.Fatalf(format, v...)
}
