package errors

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/titembaatar/sway.flem/internal/log"
)

type ErrorHandler struct {
	ExitOnFatal    bool
	VerboseLogging bool
	DebugMode      bool
	ErrorCount     int
	WarningCount   int
	FatalErrors    []error
}

func NewErrorHandler(exitOnFatal, verbose, debug bool) *ErrorHandler {
	return &ErrorHandler{
		ExitOnFatal:    exitOnFatal,
		VerboseLogging: verbose,
		DebugMode:      debug,
	}
}

func (h *ErrorHandler) Handle(err error) bool {
	if err == nil {
		return false
	}

	var appErr *AppError
	if errors.As(err, &appErr) {
		return h.handleAppError(appErr)
	}

	h.ErrorCount++
	log.Error("Error: %v", err)
	return false
}

func (h *ErrorHandler) HandleWithMessage(err error, message string) bool {
	if err == nil {
		return false
	}

	wrappedErr := Wrap(err, message)
	return h.Handle(wrappedErr)
}

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
		h.ErrorCount++
		log.Error("Error: %v", err)
		return false
	}
}

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

func (h *ErrorHandler) logFatalError(err *AppError) {
	log.Fatal("%v", err)
	if err.Suggestion != "" && h.VerboseLogging {
		log.Info("Suggestion: %s", err.Suggestion)
	}
}

func (h *ErrorHandler) logError(err *AppError) {
	log.Error("%v", err)
	if err.Suggestion != "" && h.VerboseLogging {
		log.Info("Suggestion: %s", err.Suggestion)
	}
}

func (h *ErrorHandler) logWarning(err *AppError) {
	log.Warn("%v", err)
	if err.Suggestion != "" && h.VerboseLogging {
		log.Info("Suggestion: %s", err.Suggestion)
	}
}

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

func (h *ErrorHandler) HasErrors() bool {
	return h.ErrorCount > 0
}

func (h *ErrorHandler) HasWarnings() bool {
	return h.WarningCount > 0
}

func (h *ErrorHandler) ResetCounts() {
	h.ErrorCount = 0
	h.WarningCount = 0
	h.FatalErrors = nil
}

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
