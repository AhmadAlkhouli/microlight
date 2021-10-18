package logger

import (
	"log"
	"os"
)

type Logger struct {
	Info  *log.Logger
	Error *log.Logger
	Debug *log.Logger
}

func Create() *Logger {
	var (
		info  = log.New(os.Stdout, "Info:\t", log.Ldate|log.Ltime|log.Lshortfile)
		error = log.New(os.Stderr, "Error:\t", log.Ldate|log.Ltime|log.Lshortfile)
		debug = log.New(os.Stderr, "Debug:\t", log.Ldate|log.Ltime|log.Lshortfile)
	)
	return &Logger{Info: info, Error: error, Debug: debug}
}
