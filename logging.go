package zlog

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"time"
)

const (
	DebugLevel = iota
	TraceLevel
	InfoLevel
	ErrorLevel
	FatalLevel
)

// return string base of log level
var severityName = []string{
	DebugLevel: "Debug",
	TraceLevel: "Trace",
	InfoLevel:  " Info",
	ErrorLevel: "Error",
	FatalLevel: "Fatal",
}

type Logger struct {
	logLevel int
	depth    int
	bufOne   bytes.Buffer
	bufTwo   bytes.Buffer
	curbuf   *bytes.Buffer
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

	if level < DebugLevel || level > FatalLevel {
		panic("Logger is not supported")
	}
	logger.logLevel = level

	err := logFileProperty.getLogFile()
	if err != nil {
		panic(err)
	}

	messages = make(chan string)
}

func SetOutput(out io.Writer) {

}

func timeoutFlush(timeout time.Duration) {

}

// call after InitLogger function
// generally, you needn't change it
func SetCallDepth(depth int) {
	if depth > 0 {
		logger.depth = depth
	}
}

func Debug(format string, args ...interface{}) {
	if DebugLevel < logger.logLevel {
		return
	}
	logger.logFormat(DebugLevel, fmt.Sprintf(format, args...))
}

func Info(format string, args ...interface{}) {
	if InfoLevel < logger.logLevel {
		return
	}

	logger.logFormat(InfoLevel, fmt.Sprintf(format, args...))
}

func Error(format string, args ...interface{}) {
	if ErrorLevel < logger.logLevel {
		return
	}

	logger.logFormat(ErrorLevel, fmt.Sprintf(format, args...))
}

func Trace(format string, args ...interface{}) {
	if TraceLevel < logger.logLevel {
		return
	}

	logger.logFormat(TraceLevel, fmt.Sprintf(format, args...))
}

func Fatal(format string, args ...interface{}) {
	if FatalLevel < logger.logLevel {
		return
	}

	logger.logFormat(FatalLevel, fmt.Sprintf(format, args...))
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
