package errors

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/titembaatar/sway.flem/internal/log"
)

// ErrorHandler manages error reporting and recovery
type ErrorHandler struct {
	// Settings
	ExitOnFatal    bool
	VerboseLogging bool
	DebugMode      bool

	// Statistics
	ErrorCount   int
	WarningCount int
	FatalErrors  []error
}

// NewErrorHandler creates a new error handler
func NewErrorHandler(exitOnFatal, verbose, debug bool) *ErrorHandler {
	return &ErrorHandler{
		ExitOnFatal:    exitOnFatal,
		VerboseLogging: verbose,
		DebugMode:      debug,
	}
}

// Handle processes an error and decides how to respond
func (h *ErrorHandler) Handle(err error) bool {
	if err == nil {
		return false
	}

	var appErr *AppError
	if errors.As(err, &appErr) {
		return h.handleAppError(appErr)
	}

	// Treat untyped errors as regular errors (non-fatal)
	h.ErrorCount++
	log.Error("Error: %v", err)
	return false
}

// HandleWithMessage processes an error with additional context
func (h *ErrorHandler) HandleWithMessage(err error, message string) bool {
	if err == nil {
		return false
	}

	wrappedErr := Wrap(err, message)
	return h.Handle(wrappedErr)
}

// handleAppError processes a typed AppError
func (h *ErrorHandler) handleAppError(err *AppError) bool {
	switch err.Severity {
	case SeverityFatal:
		h.FatalErrors = append(h.FatalErrors, err)
		h.logFatalError(err)
		return h.handleFatal(err)

	case SeverityError:
		h.ErrorCount++
		h.logError(err)
		return false

	case SeverityWarning:
		h.WarningCount++
		h.logWarning(err)
		return false

	default:
		// Unknown severity, treat as error
		h.ErrorCount++
		log.Error("Error: %v", err)
		return false
	}
}

// handleFatal processes a fatal error
func (h *ErrorHandler) handleFatal(err *AppError) bool {
	if h.ExitOnFatal {
		fmt.Fprintf(os.Stderr, "Fatal error: %v\n", err)
		if err.Suggestion != "" {
			fmt.Fprintf(os.Stderr, "Suggestion: %s\n", err.Suggestion)
		}
		os.Exit(1)
	}

	return true
}

// logFatalError logs a fatal error
func (h *ErrorHandler) logFatalError(err *AppError) {
	log.Fatal("%v", err)
	if err.Suggestion != "" && h.VerboseLogging {
		log.Info("Suggestion: %s", err.Suggestion)
	}
}

// logError logs a regular error
func (h *ErrorHandler) logError(err *AppError) {
	log.Error("%v", err)
	if err.Suggestion != "" && h.VerboseLogging {
		log.Info("Suggestion: %s", err.Suggestion)
	}
}

// logWarning logs a warning
func (h *ErrorHandler) logWarning(err *AppError) {
	log.Warn("%v", err)
	if err.Suggestion != "" && h.VerboseLogging {
		log.Info("Suggestion: %s", err.Suggestion)
	}
}

// SummarizeErrors prints a summary of encountered errors
func (h *ErrorHandler) SummarizeErrors() {
	if h.ErrorCount == 0 && h.WarningCount == 0 {
		return
	}

	if h.ErrorCount > 0 {
		log.Info("Encountered %d error(s)", h.ErrorCount)
	}

	if h.WarningCount > 0 {
		log.Info("Encountered %d warning(s)", h.WarningCount)
	}

	if h.DebugMode && len(h.FatalErrors) > 0 {
		log.Debug("Fatal errors encountered:")
		for i, err := range h.FatalErrors {
			log.Debug("  %d: %v", i+1, err)
		}
	}
}

// HasErrors returns true if errors were encountered
func (h *ErrorHandler) HasErrors() bool {
	return h.ErrorCount > 0
}

// HasWarnings returns true if warnings were encountered
func (h *ErrorHandler) HasWarnings() bool {
	return h.WarningCount > 0
}

// ResetCounts resets the error and warning counters
func (h *ErrorHandler) ResetCounts() {
	h.ErrorCount = 0
	h.WarningCount = 0
	h.FatalErrors = nil
}

// FormatErrorsForUser formats errors in a user-friendly way
func FormatErrorsForUser(errs []error) string {
	if len(errs) == 0 {
		return ""
	}

	var messages []string
	for _, err := range errs {
		var appErr *AppError
		if errors.As(err, &appErr) {
			msg := appErr.Message
			if appErr.Suggestion != "" {
				msg += fmt.Sprintf(" (%s)", appErr.Suggestion)
			}
			messages = append(messages, msg)
		} else {
			messages = append(messages, err.Error())
		}
	}

	return strings.Join(messages, "\n")
}
