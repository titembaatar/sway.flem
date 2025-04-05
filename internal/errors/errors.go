package errors

import (
	"errors"
	"fmt"
)

// Standard error types
var (
	// Global errors
	ErrNotImplemented = errors.New("not implemented")

	// System errors
	ErrSystemFailure   = errors.New("system operation failed")
	ErrCommandNotFound = errors.New("command not found")

	// Configuration errors
	ErrConfigNotFound = errors.New("configuration file not found")
	ErrConfigInvalid  = errors.New("invalid configuration")
	ErrConfigParse    = errors.New("failed to parse configuration")

	// Sway errors
	ErrSwayNotRunning  = errors.New("sway is not running")
	ErrSwayCommandFail = errors.New("sway command failed")
)

// ErrorSeverity determines how errors should be handled
type ErrorSeverity int

const (
	// SeverityFatal indicates that the program cannot continue
	SeverityFatal ErrorSeverity = iota
	// SeverityError indicates a significant error but the program can continue
	SeverityError
	// SeverityWarning indicates a minor issue that shouldn't stop execution
	SeverityWarning
)

// AppError wraps errors with additional context
type AppError struct {
	// Original error
	Err error
	// Message provides user-friendly context
	Message string
	// Category helps classify the error source
	Category string
	// Severity indicates how the error should be handled
	Severity ErrorSeverity
	// Source file information for config errors
	File     string
	Line     int
	Position int
	// Suggestion for how to fix the error (if available)
	Suggestion string
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.File != "" && e.Line > 0 {
		return fmt.Sprintf("[%s] %s at %s:%d: %v", e.Category, e.Message, e.File, e.Line, e.Err)
	}

	if e.Message != "" {
		return fmt.Sprintf("[%s] %s: %v", e.Category, e.Message, e.Err)
	}

	return fmt.Sprintf("[%s] %v", e.Category, e.Err)
}

// Unwrap returns the underlying error
func (e *AppError) Unwrap() error {
	return e.Err
}

// IsFatal returns true if the error is fatal
func (e *AppError) IsFatal() bool {
	return e.Severity == SeverityFatal
}

// WithSuggestion adds a suggestion to an existing error
func (e *AppError) WithSuggestion(suggestion string) *AppError {
	e.Suggestion = suggestion
	return e
}

// New creates a new AppError with the given error and message
func New(err error, message string) *AppError {
	return &AppError{
		Err:      err,
		Message:  message,
		Category: "General",
		Severity: SeverityError,
	}
}

// NewFatal creates a new fatal AppError
func NewFatal(err error, message string) *AppError {
	return &AppError{
		Err:      err,
		Message:  message,
		Category: "General",
		Severity: SeverityFatal,
	}
}

// NewWarning creates a new warning AppError
func NewWarning(err error, message string) *AppError {
	return &AppError{
		Err:      err,
		Message:  message,
		Category: "General",
		Severity: SeverityWarning,
	}
}

// WithCategory sets the category for an AppError
func (e *AppError) WithCategory(category string) *AppError {
	e.Category = category
	return e
}

// WithFile adds file information to an AppError
func (e *AppError) WithFile(file string, line, position int) *AppError {
	e.File = file
	e.Line = line
	e.Position = position
	return e
}

// Is checks if target is the same as the wrapped error
func (e *AppError) Is(target error) bool {
	return errors.Is(e.Err, target)
}

// HandleError determines how to handle an error based on severity
// Returns true if the program should exit
func HandleError(err error) bool {
	if err == nil {
		return false
	}

	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Severity == SeverityFatal
	}

	// Default to treating unknown errors as non-fatal but important
	return false
}

// Wrap adds context to an existing error
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}

	var appErr *AppError
	if errors.As(err, &appErr) {
		// Don't wrap it twice, just update the message
		if message != "" {
			appErr.Message = fmt.Sprintf("%s: %s", message, appErr.Message)
		}
		return appErr
	}

	return New(err, message)
}

// WrapIfNotNil wraps an error if it's not nil, otherwise returns nil
func WrapIfNotNil(err error, message string) error {
	if err == nil {
		return nil
	}
	return Wrap(err, message)
}
