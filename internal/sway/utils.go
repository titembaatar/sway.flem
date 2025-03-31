package sway

import (
	"strings"
)

func StringReplace(s, old, new string) string {
	return strings.Replace(s, old, new, -1)
}

func MatchApp(id, targetId int64) bool {
	if id == 0 || targetId == 0 {
		return false
	}

	return id == targetId
}

func ParseSize(size string) (width, height string, err error) {
	parts := strings.Fields(size)

	if len(parts) == 0 {
		return "", "", nil
	}

	if len(parts) == 1 {
		return parts[0], parts[0], nil
	}

	return parts[0], parts[1], nil
}
