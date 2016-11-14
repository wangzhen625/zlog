package zlog

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

const (
	DEBUG_LOG = iota
	NOTICE_LOG
	INFO_LOG
	ERROR_LOG
)

var severityName = []string{
	DEBUG_LOG:  " Debug",
	NOTICE_LOG: "Notice",
	INFO_LOG:   " Info ",
	ERROR_LOG:  " Error",
}

const default_call_depth int = 2

type Logger struct {
	root_path      string
	log_file       *os.File
	log_level      int
	depth          int
	next_file_time time.Time
}

var programName string
var logger Logger

func InitLogger(root_path string, level ...int) {

	//get the program name as log file prefix
	file, _ := exec.LookPath(os.Args[0])
	programName = filepath.Base(file)

	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	logger = Logger{}
	logger.depth = default_call_depth
	logger.root_path = root_path

	var levelEnum = 0
	if len(level) > 0 {
		levelEnum = level[0]
	}
	if levelEnum != DEBUG_LOG && levelEnum != NOTICE_LOG && levelEnum != INFO_LOG && levelEnum != ERROR_LOG {
		panic("Logger is not supported")
	}
	logger.log_level = levelEnum

	err := logger.getLogFile()
	if err != nil {
		panic(err)
	}
}

func (logger *Logger) SetCallDepth(depth int) {
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

func Notice(format string, args ...interface{}) {
	if NOTICE_LOG < logger.log_level {
		return
	}

	logger.logFormat(NOTICE_LOG, fmt.Sprintf(format, args...))
}

func (logger *Logger) getLogFile() error {
	root_path := logger.root_path
	flag, err := IsExist(root_path)

	if err != nil {

		panic(err)
	}

	if flag == false {
		os.MkdirAll(root_path, os.ModeDir)
	}

	date := time.Unix(time.Now().Unix(), 0).Format("2006-01-02")
	next_time := time.Unix(time.Now().Unix()+(24*3600), 0)
	next_time = time.Date(next_time.Year(), next_time.Month(), next_time.Day(), 0, 0, 0, 0, next_time.Location())
	logger.next_file_time = next_time

	log_path := fmt.Sprintf("%s/%s.%s.log", root_path, programName, date)
	file, err := os.OpenFile(log_path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)

	if file == nil {
		return errors.New("open log file failed")
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
	if now.Unix() > logger.next_file_time.Unix() {
		if err := logger.getLogFile(); err != nil {
			panic(err)
		}
	}

	time := time.Unix(now.Unix(), 0).Format("2006-01-02 15:04:05")

	_, file, line, ok := runtime.Caller(logger.depth)
	if ok == false {
		panic(errors.New("get the line failed"))
	}
	_, err := Write(logger.log_file, fmt.Sprintf("%s [%s]: %s (%s:%d) \n", time, severityName[level], log, file, line))
	if err != nil {
		panic(err)
	}
}

func Close() error {
	return logger.log_file.Close()
}
