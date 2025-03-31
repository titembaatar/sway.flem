package log

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func captureOutput(f func()) string {
	// Save and restore original stderr
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	// Execute the function that produces output
	f()

	// Restore original stderr and read captured output
	w.Close()
	os.Stderr = oldStderr
	var buf bytes.Buffer
	io.Copy(&buf, r)

	return buf.String()
}

func TestLogLevels(t *testing.T) {
	tests := []struct {
		name      string
		logLevel  LogLevel
		logFunc   func(string, ...any)
		message   string
		expectLog bool
	}{
		{"Debug with DebugLevel", LogLevelDebug, Debug, "test debug", true},
		{"Info with DebugLevel", LogLevelDebug, Info, "test info", true},
		{"Warn with DebugLevel", LogLevelDebug, Warn, "test warn", true},
		{"Error with DebugLevel", LogLevelDebug, Error, "test error", true},

		{"Debug with InfoLevel", LogLevelInfo, Debug, "test debug", false},
		{"Info with InfoLevel", LogLevelInfo, Info, "test info", true},
		{"Warn with InfoLevel", LogLevelInfo, Warn, "test warn", true},
		{"Error with InfoLevel", LogLevelInfo, Error, "test error", true},

		{"Debug with WarnLevel", LogLevelWarn, Debug, "test debug", false},
		{"Info with WarnLevel", LogLevelWarn, Info, "test info", false},
		{"Warn with WarnLevel", LogLevelWarn, Warn, "test warn", true},
		{"Error with WarnLevel", LogLevelWarn, Error, "test error", true},

		{"Debug with ErrorLevel", LogLevelError, Debug, "test debug", false},
		{"Info with ErrorLevel", LogLevelError, Info, "test info", false},
		{"Warn with ErrorLevel", LogLevelError, Warn, "test warn", false},
		{"Error with ErrorLevel", LogLevelError, Error, "test error", true},

		{"Debug with NoneLevel", LogLevelNone, Debug, "test debug", false},
		{"Info with NoneLevel", LogLevelNone, Info, "test info", false},
		{"Warn with NoneLevel", LogLevelNone, Warn, "test warn", false},
		{"Error with NoneLevel", LogLevelNone, Error, "test error", false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Set the log level for this test
			SetLevel(tc.logLevel)

			// Capture output
			output := captureOutput(func() {
				tc.logFunc(tc.message)
			})

			// Check if output contains message
			contains := strings.Contains(output, tc.message)

			if tc.expectLog && !contains {
				t.Errorf("Expected log to contain '%s', but it didn't. Output: %s", tc.message, output)
			} else if !tc.expectLog && contains {
				t.Errorf("Expected log to not contain '%s', but it did. Output: %s", tc.message, output)
			}
		})
	}
}

func TestLogLevelGetterAndSetter(t *testing.T) {
	// Save original level
	originalLevel := GetLevel()
	defer SetLevel(originalLevel)

	// Test setting and getting level
	levels := []LogLevel{LogLevelDebug, LogLevelInfo, LogLevelWarn, LogLevelError, LogLevelNone}

	for _, level := range levels {
		SetLevel(level)
		if got := GetLevel(); got != level {
			t.Errorf("GetLevel() = %v, want %v", got, level)
		}
	}
}

func TestLogFormatting(t *testing.T) {
	// Save original level
	originalLevel := GetLevel()
	defer SetLevel(originalLevel)

	// Set to debug to capture all logs
	SetLevel(LogLevelDebug)

	// Test log message formatting
	output := captureOutput(func() {
		Debug("Test %s with %d", "debug", 123)
	})

	// Check if output contains formatted message
	expected := "Test debug with 123"
	if !strings.Contains(output, expected) {
		t.Errorf("Expected log to contain '%s', but it didn't. Output: %s", expected, output)
	}

	// Check if output contains log level
	if !strings.Contains(output, "[DEBUG]") {
		t.Errorf("Expected log to contain '[DEBUG]', but it didn't. Output: %s", output)
	}
}
