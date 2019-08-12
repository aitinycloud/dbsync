package logging

import (
	"fmt"
	"log"
	"os"
	"time"
)

var (
	//LogSavePath Log SavePath
	LogSavePath = "logs/"
	//LogSaveName Log SaveName
	LogSaveName = "log_"
	//LogFileExt Log FileExt
	LogFileExt = "log"
	//TimeFormat Time Format
	TimeFormat = "2006_01_02_15:04:05"
	//DateTimeFormat Date Time Format
	DateTimeFormat = "2006-01-02 15:04:05"
)

func getLogFilePath() string {
	return fmt.Sprintf("%s", LogSavePath)
}

func getLogFileFullPath() string {
	prefixPath := getLogFilePath()
	suffixPath := fmt.Sprintf("%s%s.%s", LogSaveName, time.Now().Format(TimeFormat), LogFileExt)

	return fmt.Sprintf("%s%s", prefixPath, suffixPath)
}

func openLogFile(filePath string) *os.File {
	_, err := os.Stat(filePath)
	switch {
	case os.IsNotExist(err):
		mkDir()
	case os.IsPermission(err):
		log.Fatalf("Permission :%v", err)
	}

	handle, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Fail to OpenFile :%v", err)
	}

	return handle
}

func mkDir() {
	//dir, _ := os.Getwd()
	//err := os.MkdirAll(dir+"/"+getLogFilePath(), os.ModePerm)
	err := os.MkdirAll(getLogFilePath(), os.ModePerm)
	if err != nil {
		panic(err)
	}
}
