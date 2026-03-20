package widgets

import (
	"context"
	"strings"

	rich "github.com/eberle1080/go-rich"

	"github.com/eberle1080/go-textual/css"
	"github.com/eberle1080/go-textual/dom"
	"github.com/eberle1080/go-textual/geometry"
	"github.com/eberle1080/go-textual/msg"
	"github.com/eberle1080/go-textual/strip"
	"github.com/eberle1080/go-textual/widget"
)

// Static displays multiple lines of static text.
type Static struct {
	widget.BaseWidget
	lines []string
}

// NewStatic creates a Static widget with the given text content.
// Newlines in content are split into separate lines.
func NewStatic(content string) *Static {
	s := &Static{
		BaseWidget: *widget.NewBaseWidget(
			widget.WithDOMOptions(dom.WithCSSTypeName("Static", "Widget")),
		),
		lines: strings.Split(content, "\n"),
	}
	return s
}

// SetContent updates the content and marks the widget dirty.
func (s *Static) SetContent(content string) {
	s.lines = strings.Split(content, "\n")
	s.MarkDirty()
}

// Update is a no-op for Static.
func (s *Static) Update(_ context.Context, _ msg.Msg) msg.Cmd { return nil }

// Render returns one strip per line of content.
func (s *Static) Render(region geometry.Region) []strip.Strip {
	strips := make([]strip.Strip, region.Height)
	for i := 0; i < region.Height; i++ {
		if i < len(s.lines) {
			seg := rich.Segment{Text: s.lines[i], Style: rich.NewStyle()}
			st := strip.New(rich.Segments{seg})
			strips[i] = st.TextAlign(region.Width, css.AlignHorizontal("left"))
		} else {
			strips[i] = strip.New(nil)
		}
	}
	return strips
}
