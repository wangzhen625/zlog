package zlog

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

type LogFileProperty struct {
	programName string
	hostname    string
	pid         string
	rootPath    string
	file        *os.File
}

var logFileProperty LogFileProperty

func init() {
	logFileProperty.programName = filepath.Base(os.Args[0])
	logFileProperty.hostname, _ = os.Hostname()
	logFileProperty.pid = strconv.Itoa(os.Getpid())
}

const flushInterval = 2 * time.Second

func WriteFile() {
	t := time.NewTicker(flushInterval)
	for {
		select {
		case <-message:
			bufBytes := logger.readbuf.ptr.Bytes()
			logFileProperty.file.Write(bufBytes)
			logger.readbuf.ptr.Reset()
		case <-t.C:
			logger.switchBuf()
			bufBytes := logger.readbuf.ptr.Bytes()
			logFileProperty.file.Write(bufBytes)
			logger.readbuf.ptr.Reset()
		}
	}
}

func logName(t time.Time) (name string) {

	name = fmt.Sprintf("%s.%04d%02d%02d-%02d%02d%02d.%s.%s.log",
		logFileProperty.programName,
		t.Year(),
		t.Month(),
		t.Day(),
		t.Hour(),
		t.Minute(),
		t.Second(),
		logFileProperty.hostname,
		logFileProperty.pid,
	)
	return name
}

func createLogDir() {
	os.MkdirAll(logFileProperty.rootPath, 0777)
}

var onceLogDir sync.Once
var nextDayCreateFileTime int64

func (logFileProperty *LogFileProperty) getLogFile() error {

	onceLogDir.Do(createLogDir)

	nextDayCreateFileTime = time.Now().Unix()/(24*3600)*(24*3600) + 16*3600

	logName := logName(time.Now())
	logPath := fmt.Sprintf("%s/%s", logFileProperty.rootPath, logName)

	logFileProperty.file.Close()
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
	if file == nil {
		return errors.New("open log file failed:" + err.Error())
	}
	logFileProperty.file = file

	return err
}

func close() error {
	return logFileProperty.file.Close()
}

func (logFileProperty *LogFileProperty) getFileSize() int64 {
	fileInfo, _ := logFileProperty.file.Stat()
	return fileInfo.Size()
}

func fileInfoMonitor() {
	t := time.NewTicker(3 * time.Second)
	for {
		select {
		case <-t.C:
			if logFileProperty.getFileSize() > rollFileSize ||
				time.Now().Unix() > nextDayCreateFileTime {
				logFileProperty.getLogFile()
			}
		}
	}
}

func isExitDir(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}
