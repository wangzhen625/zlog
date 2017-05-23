package zlog

import (
	"fmt"
	"io"
	"os"
	"path"
	"runtime"
	"sync"
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
	mu       sync.Mutex
}

var logger Logger

const defaultCallDepth int = 3

var message = make(chan string)

//50M roll back the file
var rollFileSize int64 = 1024 * 1024 * 50

func init() {
	logger.depth = defaultCallDepth
	logger.logLevel = TraceLevel
}
func InitLogger(rootPath string, level int) {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	if level < DebugLevel || level > FatalLevel {
		panic("Logger level is not supported")
	}

	logFileProperty.rootPath = rootPath
	err := logFileProperty.getLogFile(time.Now())
	if err != nil {
		panic(err)
	}

	go WriteMsg()
}

func SetOutput(out io.Writer) {

}

// call after InitLogger function
// generally, you needn't change it
func SetCallDepth(depth int) {
	if depth > 0 {
		logger.depth = depth
	}
}

func (logger *Logger) logFormat(level int, log string) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	fileTime, filename, line := makeLogHead()
	logger.mu.Lock()
	message <- fmt.Sprintf("%s [%s]: %s (%s:%d) \n", fileTime, severityName[level], log, filename, line)
	logger.mu.Unlock()
}

func makeLogHead() (headTime, fileName string, line int) {
	now := time.Now()
	fileTime := now.Format("20060102 15:04:05")
	fileTime = fmt.Sprintf("%s.%09d", fileTime, now.Nanosecond())
	_, filePath, line, ok := runtime.Caller(logger.depth)
	if ok == false {
		fileName = "xxx"
		line = 0
		//panic(errors.New("get the line failed"))
	}
	//tmp := strings.Split(file, "/")
	//file = tmp[len(tmp)-1]
	_, fileName = path.Split(filePath)

	return fileTime, fileName, line
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
