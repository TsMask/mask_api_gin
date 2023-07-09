package logger

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"
)

const (
	info int = iota
	debug
	error
	warn
)

var logMapping = map[int]string{
	0: "INFO",
	1: "DEBUG",
	2: "ERROR",
	3: "WARN",
}

func logWithLevel(level int, format string, v ...interface{}) {
	if level < 0 {
		return
	}

	// 文件行号
	_, file, line, _ := runtime.Caller(2)
	file = filepath.Base(file)
	prefix := fmt.Sprintf("%s [%s:%d] ", logMapping[level], file, line)

	log.SetPrefix(prefix)       // 设置日志前缀
	log.SetFlags(log.LstdFlags) // 设置日期和时间格式

	log.Printf(format+"\n", v...)
}

func Infof(format string, v ...interface{}) {
	logWithLevel(info, format, v...)
}

func Debugf(format string, v ...interface{}) {
	logWithLevel(debug, format, v...)
}

func Errorf(format string, v ...interface{}) {
	logWithLevel(error, format, v...)
}

func Warnf(format string, v ...interface{}) {
	logWithLevel(warn, format, v...)
}

func Panicf(format string, v ...interface{}) {
	log.Fatalf(format, v...)
}
