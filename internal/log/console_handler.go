package log

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"sync"
	"time"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorGray   = "\033[37m"
	colorBold   = "\033[1m"
)

type ConsoleHandler struct {
	mu       sync.Mutex
	w        io.Writer
	level    slog.Level
	useColor bool
}

func NewConsoleHandler(w io.Writer, level slog.Level, useColor bool) *ConsoleHandler {
	return &ConsoleHandler{
		w:        w,
		level:    level,
		useColor: useColor,
	}
}

func (h *ConsoleHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level
}

func (h *ConsoleHandler) Handle(ctx context.Context, record slog.Record) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	timeStr := fmt.Sprintf("[%s]", record.Time.Format("06-01-02 15:04:05"))

	levelStr := formatLevel(record.Level)
	if h.useColor {
		levelStr = colorizeLevel(record.Level, levelStr)
	}

	builder := strings.Builder{}
	builder.WriteString(timeStr)
	builder.WriteString(" ")
	builder.WriteString(levelStr)
	builder.WriteString(" ")

	component := ""
	record.Attrs(func(attr slog.Attr) bool {
		if attr.Key == "component" {
			component = attr.Value.String()
			return false
		}
		return true
	})

	if component != "" {
		componentStr := fmt.Sprintf("[%-4s]", component)
		if h.useColor {
			componentStr = colorPurple + componentStr + colorReset
		}
		builder.WriteString(componentStr)
		builder.WriteString(" ")
	}

	if h.useColor {
		builder.WriteString(colorBold)
	}
	builder.WriteString(record.Message)
	if h.useColor {
		builder.WriteString(colorReset)
	}

	hasAttrs := false
	var operation, state string

	record.Attrs(func(attr slog.Attr) bool {
		if attr.Key == "operation" {
			operation = attr.Value.String()
		} else if attr.Key == "state" {
			state = attr.Value.String()
		}
		return true
	})

	if operation != "" && state != "" {
		builder.WriteString(" | ")
		if h.useColor {
			builder.WriteString(colorCyan)
		}
		builder.WriteString(operation)
		if h.useColor {
			builder.WriteString(colorReset)
		}
		builder.WriteString(" -> ")
		builder.WriteString(state)
		hasAttrs = true
	}

	record.Attrs(func(attr slog.Attr) bool {
		if attr.Key != "component" && attr.Key != "operation" && attr.Key != "state" &&
			attr.Key != slog.TimeKey && attr.Key != slog.LevelKey && attr.Key != slog.MessageKey &&
			attr.Key != "duration_seconds" {
			if !hasAttrs {
				builder.WriteString(" |")
				hasAttrs = true
			} else {
				builder.WriteString(",")
			}
			builder.WriteString(fmt.Sprintf(" %s=%s", attr.Key, formatAttrValue(attr.Value)))
		}
		return true
	})

	builder.WriteString("\n")

	_, err := io.WriteString(h.w, builder.String())
	return err
}

func (h *ConsoleHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *ConsoleHandler) WithGroup(name string) slog.Handler {
	return h
}

func formatLevel(level slog.Level) string {
	var levelStr string
	switch {
	case level >= slog.LevelError:
		levelStr = "[ERROR]"
	case level >= slog.LevelWarn:
		levelStr = "[WARN ]"
	case level >= slog.LevelInfo:
		levelStr = "[INFO ]"
	default:
		levelStr = "[DEBUG]"
	}
	return levelStr
}

func colorizeLevel(level slog.Level, levelStr string) string {
	var color string
	switch {
	case level >= slog.LevelError:
		color = colorRed
	case level >= slog.LevelWarn:
		color = colorYellow
	case level >= slog.LevelInfo:
		color = colorGreen
	default:
		color = colorBlue
	}
	return color + levelStr + colorReset
}

func formatAttrValue(v slog.Value) string {
	switch v.Kind() {
	case slog.KindString:
		return fmt.Sprintf("%q", v.String())
	case slog.KindTime:
		return v.Time().Format(time.RFC3339)
	case slog.KindAny:
		if v.Any() == nil {
			return "<nil>"
		}
		return fmt.Sprintf("%v", v.Any())
	default:
		return v.String()
	}
}
