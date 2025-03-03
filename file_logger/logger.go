package file_logger

import (
	"fmt"
	"os"
)

type FileLogger struct {
	logFile *os.File
}

func NewFileLogger(logFileName string) (fl *FileLogger) {
	var err error
	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Printf("Failed to open/create the log file: %v\n", err)
		os.Exit(1)
	}

	return &FileLogger{logFile: logFile}
}

func (fl FileLogger) Message(message string) {
	_, _ = fl.logFile.WriteString(fmt.Sprintf("%s\n", message))

}

func (fl FileLogger) Error(err error) {
	_, _ = fl.logFile.WriteString(fmt.Sprintf("ERROR: %s\n", err))
}

func (fl FileLogger) Close() {
	fl.logFile.Close()
}
