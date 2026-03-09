package logger

import (
	"fmt"
	"log"
	"os"
)

var appLogger *log.Logger

func init() {
	// Create logs directory if it doesn't exist
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		os.Mkdir("logs", 0755)
	}
	logFile, err := os.OpenFile("logs/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Could not create log file: %v\n", err)
		return
	}
	// Concise format: date time message
	appLogger = log.New(logFile, "", log.Ldate|log.Ltime)
}

func Println(v ...interface{}) {
	fmt.Println(v...) // Print to CLI
	if appLogger != nil {
		appLogger.Println(v...)
	}
}

func Printf(format string, v ...interface{}) {
	fmt.Printf(format, v...) // Print to CLI
	if appLogger != nil {
		appLogger.Printf(format, v...)
	}
}
