package log

import (
	"context"
	"fmt"
	"io"
	"log/slog"
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

const (
	ComponentCore   = "CORE"
	ComponentSway   = "SWAY"
	ComponentConfig = "CONF"
	ComponentEnv    = "ENV "
	ComponentWs     = "WKSP"
	ComponentCon    = "CON "
	ComponentApp    = "APP "
)

var (
	currentLevel     = LogLevelInfo
	currentComponent = ""
	logger           *slog.Logger
	useColorOutput   = true
)

func init() {
	if os.Getenv("NO_COLOR") != "" || os.Getenv("FLEM_NO_COLOR") != "" {
		useColorOutput = false
	}

	useJSON := os.Getenv("FLEM_JSON_LOGS") == "true"

	if useJSON {
		handler := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
		logger = slog.New(handler)
	} else {
		handler := NewConsoleHandler(os.Stderr, convertLevel(currentLevel), useColorOutput)
		logger = slog.New(handler)
	}
}

func SetOutput(w io.Writer) {
	if isUsingJSON() {
		handler := slog.NewJSONHandler(w, &slog.HandlerOptions{
			Level: convertLevel(currentLevel),
		})
		logger = slog.New(handler)
	} else {
		handler := NewConsoleHandler(w, convertLevel(currentLevel), useColorOutput)
		logger = slog.New(handler)
	}
}

func SetLevel(level LogLevel) {
	currentLevel = level

	if isUsingJSON() {
		handler := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
			Level: convertLevel(level),
		})
		logger = slog.New(handler)
	} else {
		handler := NewConsoleHandler(os.Stderr, convertLevel(level), useColorOutput)
		logger = slog.New(handler)
	}
}

func GetLevel() LogLevel {
	return currentLevel
}

func SetComponent(component string) {
	currentComponent = component
}

func SetColorOutput(enabled bool) {
	if useColorOutput == enabled {
		return
	}

	useColorOutput = enabled

	if !isUsingJSON() {
		handler := NewConsoleHandler(os.Stderr, convertLevel(currentLevel), useColorOutput)
		logger = slog.New(handler)
	}
}

func Debug(format string, args ...any) {
	if currentLevel <= LogLevelDebug {
		logger.DebugContext(
			context.Background(),
			fmt.Sprintf(format, args...),
			slog.String("component", currentComponent),
		)
	}
}

func Info(format string, args ...any) {
	if currentLevel <= LogLevelInfo {
		logger.InfoContext(
			context.Background(),
			fmt.Sprintf(format, args...),
			slog.String("component", currentComponent),
		)
	}
}

func Warn(format string, args ...any) {
	if currentLevel <= LogLevelWarn {
		logger.WarnContext(
			context.Background(),
			fmt.Sprintf(format, args...),
			slog.String("component", currentComponent),
		)
	}
}

func Error(format string, args ...any) {
	if currentLevel <= LogLevelError {
		logger.ErrorContext(
			context.Background(),
			fmt.Sprintf(format, args...),
			slog.String("component", currentComponent),
		)
	}
}

func Fatal(format string, args ...any) {
	logger.ErrorContext(
		context.Background(),
		fmt.Sprintf(format, args...),
		slog.String("component", currentComponent),
		slog.String("fatal", "true"),
	)
	os.Exit(1)
}

func Operation(name string) *OperationLogger {
	return &OperationLogger{
		Name:      name,
		StartTime: time.Now(),
		Level:     LogLevelInfo,
	}
}

type OperationLogger struct {
	Name      string
	StartTime time.Time
	Level     LogLevel
}

func (o *OperationLogger) WithLevel(level LogLevel) *OperationLogger {
	o.Level = level
	return o
}

func (o *OperationLogger) Begin() {
	if currentLevel <= o.Level {
		logger.InfoContext(
			context.Background(),
			fmt.Sprintf("Starting: %s", o.Name),
			slog.String("component", currentComponent),
			slog.String("operation", o.Name),
			slog.String("state", "begin"),
		)
	}
}

func (o *OperationLogger) End() {
	if currentLevel <= o.Level {
		elapsed := time.Since(o.StartTime)
		logger.InfoContext(
			context.Background(),
			fmt.Sprintf("Completed: %s (took %.2fs)", o.Name, elapsed.Seconds()),
			slog.String("component", currentComponent),
			slog.String("operation", o.Name),
			slog.String("state", "end"),
			slog.Float64("duration_seconds", elapsed.Seconds()),
		)
	}
}

func (o *OperationLogger) EndWithError(err error) {
	if currentLevel <= LogLevelError {
		elapsed := time.Since(o.StartTime)
		logger.ErrorContext(
			context.Background(),
			fmt.Sprintf("Failed: %s (took %.2fs): %v", o.Name, elapsed.Seconds(), err),
			slog.String("component", currentComponent),
			slog.String("operation", o.Name),
			slog.String("state", "error"),
			slog.Float64("duration_seconds", elapsed.Seconds()),
			slog.String("error", err.Error()),
		)
	}
}

func convertLevel(level LogLevel) slog.Level {
	switch level {
	case LogLevelDebug:
		return slog.LevelDebug
	case LogLevelInfo:
		return slog.LevelInfo
	case LogLevelWarn:
		return slog.LevelWarn
	case LogLevelError:
		return slog.LevelError
	case LogLevelNone:
		return slog.LevelError + 1
	default:
		return slog.LevelInfo
	}
}

func isUsingJSON() bool {
	_, isJSON := logger.Handler().(*slog.JSONHandler)
	return isJSON
}

func EnableJSONLogging() {
	handler := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level: convertLevel(currentLevel),
	})
	logger = slog.New(handler)
}

func DisableJSONLogging() {
	handler := NewConsoleHandler(os.Stderr, convertLevel(currentLevel), useColorOutput)
	logger = slog.New(handler)
}
