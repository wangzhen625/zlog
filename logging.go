package zlog

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

const (
	DEBUG_LOG = iota
	TRACE_LOG
	INFO_LOG
	ERROR_LOG
	FATAL_LOG
)

var severityName = []string{
	DEBUG_LOG: "Debug",
	TRACE_LOG: "Trace",
	INFO_LOG:  " Info",
	ERROR_LOG: "Error",
	FATAL_LOG: "Fatal",
}

type Logger struct {
	logLevel int
	depth    int
}

var logger Logger

const defaultCallDepth int = 2

var nextCreateFileTime int64

var messages chan string

func InitLogger(rootPath string, level int) {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	logger = Logger{}
	logger.depth = defaultCallDepth
	logFileProperty.rootPath = rootPath

	if level < DEBUG_LOG || level > FATAL_LOG {
		panic("Logger is not supported")
	}
	logger.logLevel = level

	err := logFileProperty.getLogFile()
	if err != nil {
		panic(err)
	}

	messages = make(chan string)
}

// call after InitLogger function
// generally, you needn't change it
func SetCallDepth(depth int) {
	if depth > 0 {
		logger.depth = depth
	}
}

func Debug(format string, args ...interface{}) {
	if DEBUG_LOG < logger.logLevel {
		return
	}
	logger.logFormat(DEBUG_LOG, fmt.Sprintf(format, args...))
}

func Info(format string, args ...interface{}) {
	if INFO_LOG < logger.logLevel {
		return
	}

	logger.logFormat(INFO_LOG, fmt.Sprintf(format, args...))
}

func Error(format string, args ...interface{}) {
	if ERROR_LOG < logger.logLevel {
		return
	}

	logger.logFormat(ERROR_LOG, fmt.Sprintf(format, args...))
}

func Trace(format string, args ...interface{}) {
	if TRACE_LOG < logger.logLevel {
		return
	}

	logger.logFormat(TRACE_LOG, fmt.Sprintf(format, args...))
}

func Fatal(format string, args ...interface{}) {
	if FATAL_LOG < logger.logLevel {
		return
	}

	logger.logFormat(FATAL_LOG, fmt.Sprintf(format, args...))
	os.Exit(-1)
}

func (logger *Logger) logFormat(level int, log string) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	now := time.Now()
	if now.Unix() > nextCreateFileTime {
		if err := logFileProperty.getLogFile(); err != nil {
			panic(err)
		}
	}

	time := now.Format("20060102 15:04:05")
	time = fmt.Sprintf("%s.%09d", time, now.Nanosecond())
	_, file, line, ok := runtime.Caller(logger.depth)
	if ok == false {
		panic(errors.New("get the line failed"))
	}

	tmp := strings.Split(file, "/")
	file = tmp[len(tmp)-1]

	_, err := Write(logFileProperty.logFile, fmt.Sprintf("%s [%s]: %s (%s:%d) \n", time, severityName[level], log, file, line))
	if err != nil {
		panic(err)
	}
}
