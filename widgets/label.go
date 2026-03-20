// Package widgets provides concrete widget implementations.
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

// Label displays a single line of text.
type Label struct {
	widget.BaseWidget
	text string
}

// NewLabel creates a Label with the given text.
func NewLabel(text string) *Label {
	l := &Label{
		BaseWidget: *widget.NewBaseWidget(
			widget.WithDOMOptions(dom.WithCSSTypeName("Label", "Widget")),
		),
		text: text,
	}
	return l
}

// SetText updates the label text and marks the widget dirty.
func (l *Label) SetText(text string) {
	l.text = text
	l.MarkDirty()
}

// Text returns the current label text.
func (l *Label) Text() string { return l.text }

// Update handles messages. Label has no interactive behaviour.
func (l *Label) Update(_ context.Context, _ msg.Msg) msg.Cmd { return nil }

// Render returns a single strip containing the label text.
func (l *Label) Render(region geometry.Region) []strip.Strip {
	strips := make([]strip.Strip, region.Height)
	if region.Height == 0 {
		return strips
	}
	seg := rich.Segment{Text: l.text, Style: rich.NewStyle()}
	s := strip.New(rich.Segments{seg})
	s = s.TextAlign(region.Width, css.AlignHorizontal("left"))
	strips[0] = s
	for i := 1; i < region.Height; i++ {
		strips[i] = strip.New(nil)
	}
	return strips
}
