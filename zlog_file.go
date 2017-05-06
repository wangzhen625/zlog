package zlog

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func Write(file *os.File, content string) (bool, error) {
	_, err := file.WriteString(content)

	if err != nil {
		return false, err
	}
	return true, nil
}

//get the program name as log file prefix
var programName = filepath.Base(os.Args[0])

func logName(t time.Time) (name string) {
	name = fmt.Sprintf("%s.%04d-%02d-%02d.log",
		programName,
		t.Year(),
		t.Month(),
		t.Day(),
	)
	return name
}

func createLogDir() {
	os.MkdirAll(logger.root_path, 0777)
	//os.MkdirAll(logger.root_path, os.ModeDir)
}
