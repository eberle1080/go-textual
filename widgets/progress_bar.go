package widgets

import (
	"context"
	"fmt"
	"strings"

	rich "github.com/eberle1080/go-rich"

	"github.com/eberle1080/go-textual/css"
	"github.com/eberle1080/go-textual/dom"
	"github.com/eberle1080/go-textual/geometry"
	"github.com/eberle1080/go-textual/msg"
	"github.com/eberle1080/go-textual/strip"
	"github.com/eberle1080/go-textual/widget"
)

// ProgressBar displays a horizontal progress bar.
type ProgressBar struct {
	widget.BaseWidget
	total   float64
	current float64
	label   string
}

// NewProgressBar creates a ProgressBar with the given total.
func NewProgressBar(total float64) *ProgressBar {
	p := &ProgressBar{
		BaseWidget: *widget.NewBaseWidget(
			widget.WithDOMOptions(dom.WithCSSTypeName("ProgressBar", "Widget")),
		),
		total: total,
	}
	return p
}

// SetProgress sets the current progress value.
func (p *ProgressBar) SetProgress(v float64) {
	if v < 0 {
		v = 0
	}
	if v > p.total {
		v = p.total
	}
	p.current = v
	p.MarkDirty()
}

// SetLabel sets an optional text label shown next to the bar.
func (p *ProgressBar) SetLabel(label string) {
	p.label = label
	p.MarkDirty()
}

// Progress returns the current progress value.
func (p *ProgressBar) Progress() float64 { return p.current }

// Update is a no-op for ProgressBar.
func (p *ProgressBar) Update(_ context.Context, _ msg.Msg) msg.Cmd { return nil }

// Render draws the progress bar.
func (p *ProgressBar) Render(region geometry.Region) []strip.Strip {
	strips := make([]strip.Strip, region.Height)
	if region.Height == 0 || region.Width == 0 {
		return strips
	}

	pct := 0.0
	if p.total > 0 {
		pct = p.current / p.total
	}
	pctStr := fmt.Sprintf(" %.0f%%", pct*100)

	// Reserve space for percentage text and optional label
	suffix := pctStr
	if p.label != "" {
		suffix = " " + p.label + pctStr
	}
	barWidth := region.Width - len(suffix)
	if barWidth < 4 {
		barWidth = 4
		suffix = ""
	}

	filled := int(float64(barWidth) * pct)
	if filled > barWidth-2 {
		filled = barWidth - 2
	}

	bar := "[" + strings.Repeat("#", filled) + strings.Repeat("-", barWidth-filled-2) + "]"

	filledStyle := rich.NewStyle().Foreground(rich.ANSIColor(rich.Green))
	emptyStyle := rich.NewStyle().Foreground(rich.ANSIColor(rich.BrightBlack))

	segs := rich.Segments{
		{Text: bar[:1+filled], Style: filledStyle},
		{Text: bar[1+filled:], Style: emptyStyle},
		{Text: suffix, Style: rich.NewStyle()},
	}

	s := strip.New(segs)
	s = s.TextAlign(region.Width, css.AlignHorizontal("left"))
	strips[0] = s
	for i := 1; i < region.Height; i++ {
		strips[i] = strip.New(nil)
	}
	return strips
}
