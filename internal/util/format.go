package util

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func ParseSize(size string) (string, string, error) {
	parts := strings.Fields(size)

	if len(parts) == 0 {
		return "", "", fmt.Errorf("empty size string")
	}

	if len(parts) == 1 {
		// Use same value for both dimensions
		return parts[0], parts[0], nil
	}

	if len(parts) >= 2 {
		return parts[0], parts[1], nil
	}

	return "", "", fmt.Errorf("invalid size format: %s", size)
}

func ParseDimension(dim string) (int, bool, error) {
	// Check for percentage points (ppt)
	if strings.HasSuffix(dim, "ppt") {
		val, err := strconv.Atoi(strings.TrimSuffix(dim, "ppt"))
		if err != nil {
			return 0, false, fmt.Errorf("invalid ppt value: %s", dim)
		}
		return val, true, nil
	}

	// Check for pixels (px) or just a number
	value := dim
	if strings.HasSuffix(dim, "px") {
		value = strings.TrimSuffix(dim, "px")
	}

	val, err := strconv.Atoi(value)
	if err != nil {
		return 0, false, fmt.Errorf("invalid pixel value: %s", dim)
	}

	return val, false, nil
}

func FormatPosition(position string) string {
	specialPositions := map[string]string{
		"center":  "position center",
		"middle":  "position center",
		"top":     "position 0 0",
		"bottom":  "position 0 999999",
		"left":    "position 0 center",
		"right":   "position 999999 center",
		"pointer": "position cursor",
		"cursor":  "position cursor",
		"mouse":   "position cursor",
	}

	if formatted, ok := specialPositions[strings.ToLower(position)]; ok {
		return formatted
	}

	return fmt.Sprintf("position %s", position)
}

func IsRegularExpression(s string) bool {
	// Look for common regex characters
	regexChars := []string{"*", "+", "?", "^", "$", ".", "[", "]", "(", ")", "{", "}"}

	for _, char := range regexChars {
		if strings.Contains(s, char) {
			return true
		}
	}

	return false
}

func MatchWithRegex(pattern, str string) (bool, error) {
	if !IsRegularExpression(pattern) {
		return pattern == str, nil
	}

	regex, err := regexp.Compile(pattern)
	if err != nil {
		return false, fmt.Errorf("invalid regex pattern: %v", err)
	}

	return regex.MatchString(str), nil
}

func IsValidLayout(layout string) bool {
	validLayouts := map[string]bool{
		"splith":   true,
		"splitv":   true,
		"stacking": true,
		"tabbed":   true,
	}

	return validLayouts[layout]
}
