package logger

import "fmt"

const (
	DEBUG = iota
	INFO  = iota
	ERROR = iota
)

var Level = ERROR

func Error(msg ...interface{}) {
	if Level <= ERROR {
		fmt.Print("Error: ")
		fmt.Println(msg...)
	}
}

func Info(msg ...interface{}) {
	if Level <= INFO {
		fmt.Print("Info: ")
		fmt.Println(msg...)
	}
}

func Debug(msg ...interface{}) {
	if Level <= DEBUG {
		fmt.Print("Debug: ")
		fmt.Println(msg...)
	}
}
