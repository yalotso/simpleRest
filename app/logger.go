package app

import (
	"log"
)

type Logger interface {
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

type logger struct {
	infoLog  *log.Logger
	errorLog *log.Logger
}

func NewLogger(infoLog, errorLog *log.Logger) Logger {
	return &logger{
		infoLog:  infoLog,
		errorLog: errorLog,
	}
}

func (l *logger) Infof(format string, args ...interface{}) {
	l.infoLog.Printf(format, args...)
}

func (l *logger) Errorf(format string, args ...interface{}) {
	l.errorLog.Printf(format, args...)
}
