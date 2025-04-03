package log

import (
	"fmt"
	"os"
	"strings"
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
	Component    = ""
)

const (
	ComponentConfig = "CONFIG"
	ComponentSway   = "SWAY"
	ComponentApp    = "APP"
	ComponentCore   = "CORE"
)

func SetLevel(level LogLevel) {
	currentLevel = level
}

func GetLevel() LogLevel {
	return currentLevel
}

func SetComponent(component string) {
	Component = strings.ToUpper(component)
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

func Operation(name string) *OperationLogger {
	return &OperationLogger{
		Name:      name,
		StartTime: time.Now(),
	}
}

type OperationLogger struct {
	Name      string
	StartTime time.Time
}

func (o *OperationLogger) Begin() {
	if currentLevel <= LogLevelInfo {
		logMessage("INFO ", "Starting operation: %s", o.Name)
	}
}

func (o *OperationLogger) End() {
	if currentLevel <= LogLevelInfo {
		elapsed := time.Since(o.StartTime)
		logMessage("INFO ", "Completed operation: %s (took %.2fs)", o.Name, elapsed.Seconds())
	}
}

func (o *OperationLogger) EndWithError(err error) {
	if currentLevel <= LogLevelError {
		elapsed := time.Since(o.StartTime)
		logMessage("ERROR", "Failed operation: %s (took %.2fs): %v", o.Name, elapsed.Seconds(), err)
	}
}

func logMessage(level string, format string, args ...any) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	message := fmt.Sprintf(format, args...)

	if Component != "" {
		fmt.Fprintf(os.Stderr, "[%s] [%s] [%s]: %s\n", timestamp, level, Component, message)
	} else {
		fmt.Fprintf(os.Stderr, "[%s] [%s]: %s\n", timestamp, level, message)
	}
}
