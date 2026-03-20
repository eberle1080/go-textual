package strip

import (
	rich "github.com/eberle1080/go-rich"
)

// StripRenderable implements [rich.Renderable] for a slice of Strips.
// Each Strip renders as one line of ANSI output.
type StripRenderable struct {
	Strips    []Strip
	ColorMode rich.ColorMode
}

// NewStripRenderable creates a StripRenderable from the given strips.
func NewStripRenderable(strips []Strip, colorMode rich.ColorMode) *StripRenderable {
	return &StripRenderable{Strips: strips, ColorMode: colorMode}
}

// Render implements [rich.Renderable]. Each strip is rendered as one line
// of ANSI escape sequences joined by newlines.
func (r *StripRenderable) Render(console *rich.Console, width int) rich.Segments {
	var result rich.Segments
	for i, s := range r.Strips {
		line := s.AdjustCellLength(width, rich.Style{})
		segs := line.Segments()
		result = append(result, segs...)
		if i < len(r.Strips)-1 {
			result = append(result, rich.Segment{Text: "\n"})
		}
	}
	return result
}
