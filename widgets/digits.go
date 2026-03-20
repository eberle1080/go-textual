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

// Digits displays a string of characters in a large 3-row tall font using
// 7-segment-inspired ASCII art. Useful for counters, clocks, and scoreboards.
type Digits struct {
	widget.BaseWidget
	value string
	style rich.Style
}

// digitRows[ch] holds the 3 rows for a character. Characters not in the map
// fall back to a 3-space wide blank.
var digitRows = map[rune][3]string{
	'0': {" _ ", "| |", "|_|"},
	'1': {"   ", "  |", "  |"},
	'2': {" _ ", " _|", "|_ "},
	'3': {" _ ", " _|", " _|"},
	'4': {"   ", "|_|", "  |"},
	'5': {" _ ", "|_ ", " _|"},
	'6': {" _ ", "|_ ", "|_|"},
	'7': {" _ ", "  |", "  |"},
	'8': {" _ ", "|_|", "|_|"},
	'9': {" _ ", "|_|", " _|"},
	':': {"   ", " . ", " . "},
	'.': {"   ", "   ", " . "},
	'-': {"   ", " _ ", "   "},
	'+': {"   ", " + ", "   "},
	' ': {"   ", "   ", "   "},
}

// NewDigits creates a Digits widget displaying the given value string.
func NewDigits(value string) *Digits {
	d := &Digits{
		BaseWidget: *widget.NewBaseWidget(
			widget.WithDOMOptions(dom.WithCSSTypeName("Digits", "Widget")),
		),
		value: value,
		style: rich.NewStyle().Bold(),
	}
	return d
}

// SetValue updates the displayed value.
func (d *Digits) SetValue(v string) {
	d.value = v
	d.MarkDirty()
}

// Value returns the current value.
func (d *Digits) Value() string { return d.value }

// SetStyle sets the text style for the digit characters.
func (d *Digits) SetStyle(s rich.Style) {
	d.style = s
	d.MarkDirty()
}

// Update is a no-op; Digits is display-only.
func (d *Digits) Update(_ context.Context, _ msg.Msg) msg.Cmd { return nil }

// Render draws the digit characters across 3 rows.
func (d *Digits) Render(region geometry.Region) []strip.Strip {
	strips := make([]strip.Strip, region.Height)
	if region.Height < 3 || region.Width == 0 {
		return strips
	}

	runes := []rune(d.value)
	rows := [3]strings.Builder{}

	for _, r := range runes {
		glyphs, ok := digitRows[r]
		if !ok {
			glyphs = digitRows[' ']
		}
		for i := range 3 {
			rows[i].WriteString(glyphs[i])
		}
	}

	for i := range 3 {
		text := rows[i].String()
		seg := rich.Segment{Text: text, Style: d.style}
		s := strip.New(rich.Segments{seg})
		s = s.TextAlign(region.Width, css.AlignHorizontal("center"))
		strips[i] = s
	}
	for i := 3; i < region.Height; i++ {
		strips[i] = strip.New(nil)
	}
	return strips
}
