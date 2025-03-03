package logger

import "fmt"

const (
	RESET_LOG_COLOR = "\033[0m"
	ERROR_LOG_COLOR = "\033[31m"

	//for possible future extension of the interface and more logging types
	SUCCESS_LOG_COLOR = "\033[32m"
	WARNING_LOG_COLOR = "\033[33m"
	INFO_LOG_COLOR    = "\033[34m"
)

type Logger interface {
	Message(message string)
	Error(err error)
	Close()
}

var Log Logger

func UseLogger(logger Logger) error {
	if Log != nil {
		return fmt.Errorf("already using one logger")
	}
	Log = logger
	return nil
}

func Message(message string) {
	fmt.Println(message)
	Log.Message(message)
}

func Error(err error) {
	fmt.Println(ERROR_LOG_COLOR + err.Error() + RESET_LOG_COLOR)
	Log.Error(err)
}

func Close() {
	Log.Close()
}
