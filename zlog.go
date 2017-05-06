package zlog

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
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

const default_call_depth int = 2

type Logger struct {
	root_path string
	log_file  *os.File
	log_level int
	depth     int
}

var logger Logger

var next_create_file_time int64

func InitLogger(root_path string, level int) {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	logger = Logger{}
	logger.depth = default_call_depth
	logger.root_path = root_path

	if level < DEBUG_LOG || level > FATAL_LOG {
		panic("Logger is not supported")
	}
	logger.log_level = level

	err := logger.getLogFile()
	if err != nil {
		panic(err)
	}
}

// call after InitLogger function
// generally, you needn't change it
func SetCallDepth(depth int) {
	if depth > 0 {
		logger.depth = depth
	}
}

func Debug(format string, args ...interface{}) {
	if DEBUG_LOG < logger.log_level {
		return
	}
	logger.logFormat(DEBUG_LOG, fmt.Sprintf(format, args...))
}

func Info(format string, args ...interface{}) {
	if INFO_LOG < logger.log_level {
		return
	}

	logger.logFormat(INFO_LOG, fmt.Sprintf(format, args...))
}

func Error(format string, args ...interface{}) {
	if ERROR_LOG < logger.log_level {
		return
	}

	logger.logFormat(ERROR_LOG, fmt.Sprintf(format, args...))
}

func Trace(format string, args ...interface{}) {
	if TRACE_LOG < logger.log_level {
		return
	}

	logger.logFormat(TRACE_LOG, fmt.Sprintf(format, args...))
}

func Fatal(format string, args ...interface{}) {
	if FATAL_LOG < logger.log_level {
		return
	}

	logger.logFormat(FATAL_LOG, fmt.Sprintf(format, args...))
	panic("")
}

var once_log_dir sync.Once

func (logger *Logger) getLogFile() error {

	once_log_dir.Do(createLogDir)

	next_create_file_time = time.Now().Unix()/(24*3600)*(24*3600) + 16*3600

	log_name := logName(time.Now())
	log_path := fmt.Sprintf("%s/%s", logger.root_path, log_name)

	file, err := os.OpenFile(log_path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
	if file == nil {
		return errors.New("open log file failed:" + err.Error())
	}

	logger.log_file = file

	return err
}

func (logger *Logger) logFormat(level int, log string) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	now := time.Now()
	if now.Unix() > next_create_file_time {
		if err := logger.getLogFile(); err != nil {
			panic(err)
		}
	}

	time := time.Unix(now.Unix(), 0).Format("2006-01-02 15:04:05")
	time = fmt.Sprintf("%s.%09d", time, now.Nanosecond())
	_, file, line, ok := runtime.Caller(logger.depth)
	if ok == false {
		panic(errors.New("get the line failed"))
	}

	tmp := strings.Split(file, "/")
	file = tmp[len(tmp)-1]

	_, err := Write(logger.log_file, fmt.Sprintf("%s [%s]: %s (%s:%d) \n", time, severityName[level], log, file, line))
	if err != nil {
		panic(err)
	}
}

func Close() error {
	return logger.log_file.Close()
}
