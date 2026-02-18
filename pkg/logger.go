package pkg

import (
	"fmt"
	"log"
	"os"
)

type Level string

const (
	DebugLevel Level = "DEBUG"
	InfoLevel  Level = "INFO"
	WarnLevel  Level = "WARN"
	ErrorLevel Level = "ERROR"
)

type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)
}

type Field struct {
	Key   string
	Value interface{}
}

func String(key, value string) Field {
	return Field{Key: key, Value: value}
}

func Int(key string, value int) Field {
	return Field{Key: key, Value: value}
}

func Err(err error) Field {
	return Field{Key: "error", Value: err}
}

type SimpleLogger struct {
	logger *log.Logger
	level  Level
}

func NewSimpleLogger(level Level) *SimpleLogger {
	return &SimpleLogger{
		logger: log.New(os.Stdout, "", log.LstdFlags),
		level:  level,
	}
}

func (l *SimpleLogger) Debug(msg string, fields ...Field) {
	if l.shouldLog(DebugLevel) {
		l.log(DebugLevel, msg, fields...)
	}
}

func (l *SimpleLogger) Info(msg string, fields ...Field) {
	if l.shouldLog(InfoLevel) {
		l.log(InfoLevel, msg, fields...)
	}
}

func (l *SimpleLogger) Warn(msg string, fields ...Field) {
	if l.shouldLog(WarnLevel) {
		l.log(WarnLevel, msg, fields...)
	}
}

func (l *SimpleLogger) Error(msg string, fields ...Field) {
	if l.shouldLog(ErrorLevel) {
		l.log(ErrorLevel, msg, fields...)
	}
}

func (l *SimpleLogger) Fatal(msg string, fields ...Field) {
	l.log(ErrorLevel, msg, fields...)
	os.Exit(1)
}

func (l *SimpleLogger) log(level Level, msg string, fields ...Field) {
	logMsg := fmt.Sprintf("[%s] %s", level, msg)

	if len(fields) > 0 {
		logMsg += " |"
		for _, f := range fields {
			logMsg += fmt.Sprintf(" %s=%v", f.Key, f.Value)
		}
	}

	l.logger.Println(logMsg)
}

func (l *SimpleLogger) shouldLog(level Level) bool {
	levels := map[Level]int{
		DebugLevel: 0,
		InfoLevel:  1,
		WarnLevel:  2,
		ErrorLevel: 3,
	}

	return levels[level] >= levels[l.level]
}

type NoOpLogger struct{}

func (l *NoOpLogger) Debug(msg string, fields ...Field) {}
func (l *NoOpLogger) Info(msg string, fields ...Field)  {}
func (l *NoOpLogger) Warn(msg string, fields ...Field)  {}
func (l *NoOpLogger) Error(msg string, fields ...Field) {}
func (l *NoOpLogger) Fatal(msg string, fields ...Field) {}

func NewNoOpLogger() Logger {
	return &NoOpLogger{}
}
