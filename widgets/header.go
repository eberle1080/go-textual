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

// Header displays a styled title bar at the top of a screen.
type Header struct {
	widget.BaseWidget
	title    string
	subtitle string
	bg       rich.Color
	fg       rich.Color
}

// NewHeader creates a Header with the given title.
func NewHeader(title string) *Header {
	h := &Header{
		BaseWidget: *widget.NewBaseWidget(
			widget.WithDOMOptions(dom.WithCSSTypeName("Header", "Widget")),
		),
		title: title,
		bg:    rich.ANSIColor(rich.Blue),
		fg:    rich.ANSIColor(rich.White),
	}
	return h
}

// SetTitle updates the header title.
func (h *Header) SetTitle(title string) {
	h.title = title
	h.MarkDirty()
}

// SetSubtitle sets an optional subtitle shown to the right.
func (h *Header) SetSubtitle(subtitle string) {
	h.subtitle = subtitle
	h.MarkDirty()
}

// SetColors sets the background and foreground colors.
func (h *Header) SetColors(bg, fg rich.Color) {
	h.bg = bg
	h.fg = fg
	h.MarkDirty()
}

// Update is a no-op; Header is display-only.
func (h *Header) Update(_ context.Context, _ msg.Msg) msg.Cmd { return nil }

// Render draws the header bar.
func (h *Header) Render(region geometry.Region) []strip.Strip {
	strips := make([]strip.Strip, region.Height)
	if region.Height == 0 || region.Width == 0 {
		return strips
	}

	style := rich.NewStyle().Bold().Foreground(h.fg).Background(h.bg)

	var segs rich.Segments
	titleText := " " + h.title
	segs = append(segs, rich.Segment{Text: titleText, Style: style})

	if h.subtitle != "" {
		subtitleText := h.subtitle + " "
		segs = append(segs, rich.Segment{Text: subtitleText, Style: style.Dim()})
	}

	s := strip.New(segs)
	// Fill the rest of the line with the background color
	s = s.ExtendCellLength(region.Width, style)
	s = s.TextAlign(region.Width, css.AlignHorizontal("left"))
	strips[0] = s

	for i := 1; i < region.Height; i++ {
		strips[i] = strip.New(nil)
	}
	return strips
}
