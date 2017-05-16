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
	date        string
	time        string
	hostname    string
	pid         string
	rootPath    string
	logFile     *os.File
}

var logFileProperty LogFileProperty

func init() {
	logFileProperty.programName = filepath.Base(os.Args[0])
	logFileProperty.hostname, _ = os.Hostname()
	logFileProperty.pid = strconv.Itoa(os.Getpid())
	//logFileProperty.date = time.Now()
}

func Write(file *os.File, content string) (bool, error) {
	_, err := file.WriteString(content)

	if err != nil {
		return false, err
	}
	return true, nil
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
	//os.MkdirAll(logger.root_path, os.ModeDir)
}

var onceLogDir sync.Once

func (logFileProperty *LogFileProperty) getLogFile() error {

	onceLogDir.Do(createLogDir)

	nextCreateFileTime = time.Now().Unix()/(24*3600)*(24*3600) + 16*3600

	logName := logName(time.Now())
	logPath := fmt.Sprintf("%s/%s", logFileProperty.rootPath, logName)

	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
	if file == nil {
		return errors.New("open log file failed:" + err.Error())
	}

	logFileProperty.logFile = file

	return err
}

func Close() error {
	return logFileProperty.logFile.Close()
}
