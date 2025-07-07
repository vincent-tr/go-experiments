package common

import (
	"fmt"
	"time"
)

var currentTime *time.Time

func SetCurrentTime(t time.Time) {
	currentTime = &t
}

type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarning
	LogLevelError
)

type Logger struct {
	name string
}

func NewLogger(name string) *Logger {
	return &Logger{name: name}
}

func (l *Logger) Log(level LogLevel, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)

	var now time.Time
	if currentTime != nil {
		now = *currentTime
	} else {
		now = time.Now()
	}

	var levelStr string
	switch level {
	case LogLevelDebug:
		levelStr = "\033[36mDEBUG\033[0m" // Cyan
	case LogLevelInfo:
		levelStr = "\033[32mINFO\033[0m" // Green
	case LogLevelWarning:
		levelStr = "\033[33mWARNING\033[0m" // Yellow
	case LogLevelError:
		levelStr = "\033[31mERROR\033[0m" // Red
	}

	fmt.Printf("%s [%s] %s: %s\n", now.Format("2006-01-02 15:04:05"), levelStr, l.name, msg)
}

func (l *Logger) Debug(format string, args ...interface{}) {
	l.Log(LogLevelDebug, format, args...)
}

func (l *Logger) Info(format string, args ...interface{}) {
	l.Log(LogLevelInfo, format, args...)
}

func (l *Logger) Warning(format string, args ...interface{}) {
	l.Log(LogLevelWarning, format, args...)
}

func (l *Logger) Error(format string, args ...interface{}) {
	l.Log(LogLevelError, format, args...)
}
