package log

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"time"
)

// LogLevel defines the severity of log messages
type LogLevel int

// Compatible log levels with original logger
const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
	LogLevelNone
)

// Component type identifiers
const (
	ComponentCore   = "CORE"
	ComponentSway   = "SWAY"
	ComponentConfig = "CONF"
	ComponentEnv    = "ENV "
	ComponentWs     = "WKSP"
	ComponentCon    = "CON "
	ComponentApp    = "APP "
)

// Global state
var (
	currentLevel     = LogLevelInfo
	currentComponent = ""
	logger           *slog.Logger
)

// init initializes the logger with default settings
func init() {
	// Check if we should use JSON logging
	useJSON := os.Getenv("FLEM_JSON_LOGS") == "true"

	if useJSON {
		// Initialize JSON logger
		handler := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
		logger = slog.New(handler)
	} else {
		// Initialize pretty text logger
		handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			Level: slog.LevelInfo,
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				// Customize time format
				if a.Key == slog.TimeKey {
					if t, ok := a.Value.Any().(time.Time); ok {
						return slog.String(slog.TimeKey, t.Format("2006-01-02 15:04:05"))
					}
				}
				return a
			},
		})
		logger = slog.New(handler)
	}
}

// SetOutput sets the output destination for the logger
func SetOutput(w io.Writer) {
	handler := slog.NewTextHandler(w, &slog.HandlerOptions{
		Level: convertLevel(currentLevel),
	})
	logger = slog.New(handler)
}

// SetLevel sets the logging level
func SetLevel(level LogLevel) {
	currentLevel = level

	// Update the logger's level
	handler := logger.Handler()
	if h, ok := handler.(*slog.TextHandler); ok {
		h.WithAttrs([]slog.Attr{
			slog.Int("level", int(convertLevel(level))),
		})
	}
}

// GetLevel returns the current logging level
func GetLevel() LogLevel {
	return currentLevel
}

// SetComponent sets the current component name
func SetComponent(component string) {
	currentComponent = component
}

// Debug logs a debug message
func Debug(format string, args ...any) {
	if currentLevel <= LogLevelDebug {
		logger.DebugContext(
			context.Background(),
			fmt.Sprintf(format, args...),
			slog.String("component", currentComponent),
		)
	}
}

// Info logs an informational message
func Info(format string, args ...any) {
	if currentLevel <= LogLevelInfo {
		logger.InfoContext(
			context.Background(),
			fmt.Sprintf(format, args...),
			slog.String("component", currentComponent),
		)
	}
}

// Warn logs a warning message
func Warn(format string, args ...any) {
	if currentLevel <= LogLevelWarn {
		logger.WarnContext(
			context.Background(),
			fmt.Sprintf(format, args...),
			slog.String("component", currentComponent),
		)
	}
}

// Error logs an error message
func Error(format string, args ...any) {
	if currentLevel <= LogLevelError {
		logger.ErrorContext(
			context.Background(),
			fmt.Sprintf(format, args...),
			slog.String("component", currentComponent),
		)
	}
}

// Fatal logs a fatal error message and exits
func Fatal(format string, args ...any) {
	logger.ErrorContext(
		context.Background(),
		fmt.Sprintf(format, args...),
		slog.String("component", currentComponent),
		slog.String("fatal", "true"),
	)
	os.Exit(1)
}

// Operation creates a new operation logger
func Operation(name string) *OperationLogger {
	return &OperationLogger{
		Name:      name,
		StartTime: time.Now(),
	}
}

// OperationLogger tracks an operation for logging
type OperationLogger struct {
	Name      string
	StartTime time.Time
}

// Begin logs the start of an operation
func (o *OperationLogger) Begin() {
	if currentLevel <= LogLevelInfo {
		logger.InfoContext(
			context.Background(),
			fmt.Sprintf("Starting operation: %s", o.Name),
			slog.String("component", currentComponent),
			slog.String("operation", o.Name),
			slog.String("state", "begin"),
		)
	}
}

// End logs the completion of an operation
func (o *OperationLogger) End() {
	if currentLevel <= LogLevelInfo {
		elapsed := time.Since(o.StartTime)
		logger.InfoContext(
			context.Background(),
			fmt.Sprintf("Completed operation: %s (took %.2fs)", o.Name, elapsed.Seconds()),
			slog.String("component", currentComponent),
			slog.String("operation", o.Name),
			slog.String("state", "end"),
			slog.Float64("duration_seconds", elapsed.Seconds()),
		)
	}
}

// EndWithError logs the failure of an operation
func (o *OperationLogger) EndWithError(err error) {
	if currentLevel <= LogLevelError {
		elapsed := time.Since(o.StartTime)
		logger.ErrorContext(
			context.Background(),
			fmt.Sprintf("Failed operation: %s (took %.2fs): %v", o.Name, elapsed.Seconds(), err),
			slog.String("component", currentComponent),
			slog.String("operation", o.Name),
			slog.String("state", "error"),
			slog.Float64("duration_seconds", elapsed.Seconds()),
			slog.String("error", err.Error()),
		)
	}
}

// Convert internal log level to slog level
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
		return slog.LevelError + 1 // Higher than Error
	default:
		return slog.LevelInfo
	}
}

// EnableJSONLogging switches the logger to output JSON format
func EnableJSONLogging() {
	handler := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level: convertLevel(currentLevel),
	})
	logger = slog.New(handler)
}
