package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

// Size errors
var (
	ErrInvalidSizeFormat = errors.New("invalid size format: must be a number, optionally followed by 'ppt' or 'px'")
)

// Unit for container sizes
type SizeUnit string

// Size units
const (
	UnitPercent SizeUnit = "ppt" // Percentage points (default in Sway)
	UnitPixels  SizeUnit = "px"  // Pixels
)

// Container size
type Size struct {
	Value int      // Numeric value
	Unit  SizeUnit // Unit (ppt or px)
}

// Regular expression for parsing sizes
var sizeRegex = regexp.MustCompile(`^(\d+)(ppt|px)?$`)

// String into a Size
func ParseSize(s string) (Size, error) {
	if s == "" {
		return Size{}, nil
	}

	if !sizeRegex.MatchString(s) {
		return Size{}, ErrInvalidSizeFormat
	}

	matches := sizeRegex.FindStringSubmatch(s)
	if len(matches) < 2 {
		return Size{}, ErrInvalidSizeFormat
	}

	value, err := strconv.Atoi(matches[1])
	if err != nil {
		return Size{}, ErrInvalidSizeFormat
	}

	if value < 0 {
		return Size{}, ErrInvalidSizeFormat
	}

	unit := UnitPercent // Default to percentage points (ppt)
	if len(matches) > 2 && matches[2] != "" {
		unit = SizeUnit(matches[2])
	}

	return Size{Value: value, Unit: unit}, nil
}

// String representation of the size
func (s Size) String() string {
	if s.Value == 0 {
		return ""
	}
	if s.Unit == "" {
		return fmt.Sprintf("%d", s.Value)
	}
	return fmt.Sprintf("%d%s", s.Value, s.Unit)
}

// json.Marshaler interface
func (s Size) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

// json.Unmarshaler interface
func (s *Size) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	size, err := ParseSize(str)
	if err != nil {
		return err
	}

	*s = size
	return nil
}

func (s Size) IsEmpty() bool {
	return s.Value == 0
}

func (s Size) IsValid() bool {
	if s.IsEmpty() {
		return true // Empty size is valid
	}

	return s.Value > 0 && (s.Unit == UnitPercent || s.Unit == UnitPixels)
}
