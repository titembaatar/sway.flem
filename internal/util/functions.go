package util

import (
	"log"
	"strings"
)

func NewResult(success bool, err error, data any) Result {
	if err == nil {
		return Result{
			Success: success,
			Data:    data,
		}
	}

	return Result{
		Success: success,
		Error:   err.Error(),
		Data:    data,
	}
}

func NewSuccessResult(data any) Result {
	return Result{
		Success: true,
		Data:    data,
	}
}

func NewErrorResult(err error) Result {
	if err == nil {
		return Result{Success: true}
	}

	return Result{
		Success: false,
		Error:   err.Error(),
	}
}

func LogVerbose(verbose bool, format string, args ...any) {
	if verbose {
		log.Printf(format, args...)
	}
}

func StripPrefix(s, prefix string) string {
	return strings.TrimPrefix(s, prefix)
}

func MatchAppNames(running, config string) bool {
	runningLower := strings.ToLower(running)
	configLower := strings.ToLower(config)
	return runningLower == configLower
}
