package log

import (
	"fmt"
	"os"
	"time"
)

type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
	LogLevelNone
)

var (
	currentLevel = LogLevelInfo
)

func SetLevel(level LogLevel) {
	currentLevel = level
}

func GetLevel() LogLevel {
	return currentLevel
}

func Debug(format string, args ...any) {
	if currentLevel <= LogLevelDebug {
		logMessage("DEBUG", format, args...)
	}
}

func Info(format string, args ...any) {
	if currentLevel <= LogLevelInfo {
		logMessage("INFO ", format, args...)
	}
}

func Warn(format string, args ...any) {
	if currentLevel <= LogLevelWarn {
		logMessage("WARN ", format, args...)
	}
}

func Error(format string, args ...any) {
	if currentLevel <= LogLevelError {
		logMessage("ERROR", format, args...)
	}
}

func Fatal(format string, args ...any) {
	logMessage("FATAL", format, args...)
	os.Exit(1)
}

func logMessage(level string, format string, args ...any) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	message := fmt.Sprintf(format, args...)
	fmt.Fprintf(os.Stderr, "[%s] [%s]: %s\n", timestamp, level, message)
}
