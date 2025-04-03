package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// Layout errors
var (
	ErrInvalidLayoutType = errors.New("invalid layout type")
)

// Container layout in Sway
type LayoutType string

// Layout type constants
const (
	LayoutHorizontal LayoutType = "splith"
	LayoutVertical   LayoutType = "splitv"
	LayoutTabbed     LayoutType = "tabbed"
	LayoutStacking   LayoutType = "stacking"
)

// Aliases of layouts
var layoutAliases = map[string]LayoutType{
	"splith":   LayoutHorizontal,
	"splitv":   LayoutVertical,
	"tabbed":   LayoutTabbed,
	"stacking": LayoutStacking,

	// Aliases
	"horizontal": LayoutHorizontal,
	"h":          LayoutHorizontal,
	"vertical":   LayoutVertical,
	"v":          LayoutVertical,
	"stack":      LayoutStacking,
	"s":          LayoutStacking,
	"tab":        LayoutTabbed,
	"t":          LayoutTabbed,
}

// String representation of the layout type
func (l LayoutType) String() string {
	return string(l)
}

// Sway command for this layout
func (l LayoutType) Command() string {
	return fmt.Sprintf("layout %s", l)
}

// Sway command for splitting in this layout
func (l LayoutType) SplitCommand() string {
	return fmt.Sprintf("split %s", l)
}

// Dimension to use for resizing with this layout
func (l LayoutType) ResizeDimension() string {
	switch l {
	case LayoutHorizontal, LayoutTabbed:
		return "width"
	case LayoutVertical, LayoutStacking:
		return "height"
	default:
		return "width"
	}
}

// Parses a string into a LayoutType
func ParseLayoutType(s string) (LayoutType, error) {
	if s == "" {
		return "", ErrInvalidLayoutType
	}

	s = strings.ToLower(strings.TrimSpace(s))
	if layout, ok := layoutAliases[s]; ok {
		return layout, nil
	}
	return "", ErrInvalidLayoutType
}

// json.Marshaler interface
func (l LayoutType) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(l))
}

// json.Unmarshaler interface
func (l *LayoutType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	layout, err := ParseLayoutType(s)
	if err != nil {
		return err
	}

	*l = layout
	return nil
}

func (l LayoutType) IsValid() bool {
	switch l {
	case LayoutHorizontal, LayoutVertical, LayoutTabbed, LayoutStacking:
		return true
	default:
		return false
	}
}
