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

// RichLog is a scrollable log of text lines. Lines are appended and the widget
// always shows the most recent lines that fit in the region.
type RichLog struct {
	widget.BaseWidget
	lines    []logLine
	maxLines int
}

type logLine struct {
	text  string
	style rich.Style
}

// NewRichLog creates a RichLog with a maximum retained line count.
func NewRichLog(maxLines int) *RichLog {
	if maxLines <= 0 {
		maxLines = 1000
	}
	return &RichLog{
		BaseWidget: *widget.NewBaseWidget(
			widget.WithDOMOptions(dom.WithCSSTypeName("RichLog", "Widget")),
		),
		maxLines: maxLines,
	}
}

// Write appends a plain text line.
func (r *RichLog) Write(line string) {
	r.WriteStyled(line, rich.NewStyle())
}

// WriteStyled appends a styled line.
func (r *RichLog) WriteStyled(line string, style rich.Style) {
	r.lines = append(r.lines, logLine{text: line, style: style})
	if len(r.lines) > r.maxLines {
		r.lines = r.lines[len(r.lines)-r.maxLines:]
	}
	r.MarkDirty()
}

// Clear removes all lines.
func (r *RichLog) Clear() {
	r.lines = r.lines[:0]
	r.MarkDirty()
}

// Update is a no-op for RichLog.
func (r *RichLog) Update(_ context.Context, _ msg.Msg) msg.Cmd { return nil }

// Render draws the last N lines that fit in the region.
func (r *RichLog) Render(region geometry.Region) []strip.Strip {
	strips := make([]strip.Strip, region.Height)
	if region.Height == 0 {
		return strips
	}

	// Show the most recent lines that fit.
	start := len(r.lines) - region.Height
	if start < 0 {
		start = 0
	}

	visible := r.lines[start:]
	for i, line := range visible {
		if i >= region.Height {
			break
		}
		seg := rich.Segment{Text: line.text, Style: line.style}
		s := strip.New(rich.Segments{seg})
		s = s.TextAlign(region.Width, css.AlignHorizontal("left"))
		strips[i] = s
	}
	// Fill remaining lines.
	for i := len(visible); i < region.Height; i++ {
		strips[i] = strip.New(nil)
	}
	return strips
}
