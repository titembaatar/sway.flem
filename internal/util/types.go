package util

type Result struct {
	Success bool
	Error   string
	Data    any
}

type Size struct {
	Width  string
	Height string
}

type ParsedSize struct {
	Width       int
	Height      int
	WidthInPct  bool
	HeightInPct bool
}

type Position struct {
	X string
	Y string
}

type ParsedPosition struct {
	X      int
	Y      int
	XInPct bool
	YInPct bool
}
