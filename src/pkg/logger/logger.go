package logger

import (
	"fmt"
	"log"
	"os"
)

const (
	// Silent silent log level
	Silent int = iota
	// Info info log level
	Info
	// Warn warn log level
	Warn
	// Error error log level
	Error
)

var logMapping = map[int]string{
	0: "silent",
	1: "info",
	2: "warn",
	3: "error",
}

func logWithLevel(level int, format string, v ...interface{}) {
	if level <= Silent {
		return
	}

	stdLog := log.New(os.Stdout, "["+logMapping[level]+"] ", log.LstdFlags|log.Lshortfile)
	stdLog.Output(3, fmt.Sprintf(format, v...))
}

func Infof(format string, v ...interface{}) {
	logWithLevel(Info, format, v...)
}

func Warnf(format string, v ...interface{}) {
	logWithLevel(Warn, format, v...)
}

func Errorf(format string, v ...interface{}) {
	logWithLevel(Error, format, v...)
}

// Panicf 抛出错误并退出程序
func Panicf(format string, v ...interface{}) {
	log.Fatalf(format, v...)
}
