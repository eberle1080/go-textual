package widgets

import (
	"context"

	rich "github.com/eberle1080/go-rich"

	"github.com/eberle1080/go-textual/css"
	"github.com/eberle1080/go-textual/dom"
	"github.com/eberle1080/go-textual/geometry"
	"github.com/eberle1080/go-textual/msg"
	"github.com/eberle1080/go-textual/strip"
	"github.com/eberle1080/go-textual/widget"
)

// Sparkline renders a compact bar chart from a series of float64 values.
// Each data point is rendered as a single character column using Unicode
// block elements (▁▂▃▄▅▆▇█).
type Sparkline struct {
	widget.BaseWidget
	data  []float64
	min   float64
	max   float64
	color rich.Color
}

var sparkBars = []rune{' ', '▁', '▂', '▃', '▄', '▅', '▆', '▇', '█'}

// NewSparkline creates a Sparkline widget.
func NewSparkline() *Sparkline {
	s := &Sparkline{
		BaseWidget: *widget.NewBaseWidget(
			widget.WithDOMOptions(dom.WithCSSTypeName("Sparkline", "Widget")),
		),
		color: rich.ANSIColor(rich.Green),
	}
	return s
}

// SetData replaces the data series. min/max are computed automatically.
func (s *Sparkline) SetData(data []float64) {
	s.data = data
	s.min = 0
	s.max = 0
	for i, v := range data {
		if i == 0 || v < s.min {
			s.min = v
		}
		if i == 0 || v > s.max {
			s.max = v
		}
	}
	s.MarkDirty()
}

// SetRange sets explicit min/max instead of auto-scaling.
func (s *Sparkline) SetRange(min, max float64) {
	s.min = min
	s.max = max
	s.MarkDirty()
}

// SetColor sets the bar color.
func (s *Sparkline) SetColor(c rich.Color) {
	s.color = c
	s.MarkDirty()
}

// AppendValue adds a single value, keeping the series length at most maxLen.
func (s *Sparkline) AppendValue(v float64, maxLen int) {
	s.data = append(s.data, v)
	if maxLen > 0 && len(s.data) > maxLen {
		s.data = s.data[len(s.data)-maxLen:]
	}
	// Recompute range.
	s.SetData(s.data)
}

// Update is a no-op; Sparkline is display-only.
func (s *Sparkline) Update(_ context.Context, _ msg.Msg) msg.Cmd { return nil }

// Render draws the sparkline.
func (s *Sparkline) Render(region geometry.Region) []strip.Strip {
	strips := make([]strip.Strip, region.Height)
	if region.Height == 0 || region.Width == 0 || len(s.data) == 0 {
		return strips
	}

	style := rich.NewStyle().Foreground(s.color)
	rng := s.max - s.min

	// Use only as many data points as fit in the width.
	data := s.data
	if len(data) > region.Width {
		data = data[len(data)-region.Width:]
	}

	// Build one bar-character per data point.
	var sb []rune
	for _, v := range data {
		var idx int
		if rng > 0 {
			frac := (v - s.min) / rng
			idx = int(frac * float64(len(sparkBars)-1))
			idx = max(0, min(idx, len(sparkBars)-1))
		}
		sb = append(sb, sparkBars[idx])
	}

	text := string(sb)
	seg := rich.Segment{Text: text, Style: style}
	st := strip.New(rich.Segments{seg})
	st = st.TextAlign(region.Width, css.AlignHorizontal("left"))
	strips[0] = st

	for i := 1; i < region.Height; i++ {
		strips[i] = strip.New(nil)
	}
	return strips
}
