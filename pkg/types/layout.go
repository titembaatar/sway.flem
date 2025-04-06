package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrInvalidLayoutType = errors.New("invalid layout type")
)

type LayoutType string

const (
	LayoutHorizontal LayoutType = "splith"
	LayoutVertical   LayoutType = "splitv"
	LayoutTabbed     LayoutType = "tabbed"
	LayoutStacking   LayoutType = "stacking"
)

var layoutAliases = map[string]LayoutType{
	"splith":   LayoutHorizontal,
	"splitv":   LayoutVertical,
	"tabbed":   LayoutTabbed,
	"stacking": LayoutStacking,

	"horizontal": LayoutHorizontal,
	"h":          LayoutHorizontal,
	"vertical":   LayoutVertical,
	"v":          LayoutVertical,
	"stack":      LayoutStacking,
	"s":          LayoutStacking,
	"tab":        LayoutTabbed,
	"t":          LayoutTabbed,
}

func (l LayoutType) String() string {
	return string(l)
}

func (l LayoutType) Command() string {
	return fmt.Sprintf("layout %s", l)
}

func (l LayoutType) SplitCommand() string {
	switch l {
	case LayoutHorizontal:
		return "split h"
	case LayoutVertical:
		return "split v"
	case LayoutTabbed:
		return "split h"
	case LayoutStacking:
		return "split v"
	default:
		return fmt.Sprintf("split h") // Default to horizontal
	}
}

func (l LayoutType) Orientation() string {
	switch l {
	case LayoutHorizontal, LayoutTabbed:
		return "width"
	case LayoutVertical, LayoutStacking:
		return "height"
	default:
		return "width"
	}
}

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

func (l LayoutType) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(l))
}

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
