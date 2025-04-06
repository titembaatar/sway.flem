package errors

import (
	"errors"
	"fmt"
)

var (
	ErrNotImplemented  = errors.New("not implemented")
	ErrSystemFailure   = errors.New("system operation failed")
	ErrCommandNotFound = errors.New("command not found")
	ErrConfigNotFound  = errors.New("configuration file not found")
	ErrConfigInvalid   = errors.New("invalid configuration")
	ErrConfigParse     = errors.New("failed to parse configuration")
	ErrSwayNotRunning  = errors.New("sway is not running")
	ErrSwayCommandFail = errors.New("sway command failed")
)

type ErrorSeverity int

const (
	SeverityFatal ErrorSeverity = iota
	SeverityError
	SeverityWarning
)

type AppError struct {
	Err        error
	Message    string
	Category   string
	Severity   ErrorSeverity
	File       string
	Line       int
	Position   int
	Suggestion string
}

func (e *AppError) Error() string {
	if e.File != "" && e.Line > 0 {
		return fmt.Sprintf("[%s] %s at %s:%d: %v", e.Category, e.Message, e.File, e.Line, e.Err)
	}

	if e.Message != "" {
		return fmt.Sprintf("[%s] %s: %v", e.Category, e.Message, e.Err)
	}

	return fmt.Sprintf("[%s] %v", e.Category, e.Err)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func (e *AppError) IsFatal() bool {
	return e.Severity == SeverityFatal
}

func (e *AppError) WithSuggestion(suggestion string) *AppError {
	e.Suggestion = suggestion
	return e
}

func New(err error, message string) *AppError {
	return &AppError{
		Err:      err,
		Message:  message,
		Category: "General",
		Severity: SeverityError,
	}
}

func NewFatal(err error, message string) *AppError {
	return &AppError{
		Err:      err,
		Message:  message,
		Category: "General",
		Severity: SeverityFatal,
	}
}

func NewWarning(err error, message string) *AppError {
	return &AppError{
		Err:      err,
		Message:  message,
		Category: "General",
		Severity: SeverityWarning,
	}
}

func (e *AppError) WithCategory(category string) *AppError {
	e.Category = category
	return e
}

func (e *AppError) WithFile(file string, line, position int) *AppError {
	e.File = file
	e.Line = line
	e.Position = position
	return e
}

func (e *AppError) Is(target error) bool {
	return errors.Is(e.Err, target)
}

func HandleError(err error) bool {
	if err == nil {
		return false
	}

	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Severity == SeverityFatal
	}

	return false
}

func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}

	var appErr *AppError
	if errors.As(err, &appErr) {
		if message != "" {
			appErr.Message = fmt.Sprintf("%s: %s", message, appErr.Message)
		}
		return appErr
	}

	return New(err, message)
}

func WrapIfNotNil(err error, message string) error {
	if err == nil {
		return nil
	}
	return Wrap(err, message)
}
