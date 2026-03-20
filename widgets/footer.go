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

// FooterBinding is a key+description pair shown in the footer.
type FooterBinding struct {
	Key  string
	Desc string
}

// Footer displays a row of key bindings at the bottom of a screen.
type Footer struct {
	widget.BaseWidget
	bindings []FooterBinding
	bg       rich.Color
	keyFg    rich.Color
	descFg   rich.Color
}

// NewFooter creates a Footer with the given key bindings.
func NewFooter(bindings ...FooterBinding) *Footer {
	f := &Footer{
		BaseWidget: *widget.NewBaseWidget(
			widget.WithDOMOptions(dom.WithCSSTypeName("Footer", "Widget")),
		),
		bindings: bindings,
		bg:       rich.ANSIColor(rich.BrightBlack),
		keyFg:    rich.ANSIColor(rich.White),
		descFg:   rich.ANSIColor(rich.BrightBlack),
	}
	return f
}

// SetBindings replaces all displayed bindings.
func (f *Footer) SetBindings(bindings ...FooterBinding) {
	f.bindings = bindings
	f.MarkDirty()
}

// Update is a no-op; Footer is display-only.
func (f *Footer) Update(_ context.Context, _ msg.Msg) msg.Cmd { return nil }

// Render draws the footer bar.
func (f *Footer) Render(region geometry.Region) []strip.Strip {
	strips := make([]strip.Strip, region.Height)
	if region.Height == 0 || region.Width == 0 {
		return strips
	}

	keyStyle := rich.NewStyle().Bold().Foreground(f.keyFg).Background(f.bg)
	descStyle := rich.NewStyle().Foreground(f.descFg).Background(f.bg)
	sepStyle := rich.NewStyle().Foreground(f.descFg).Background(f.bg)

	var segs rich.Segments
	for i, b := range f.bindings {
		if i > 0 {
			segs = append(segs, rich.Segment{Text: "  ", Style: sepStyle})
		}
		segs = append(segs,
			rich.Segment{Text: " " + strings.ToUpper(b.Key) + " ", Style: keyStyle},
			rich.Segment{Text: " " + b.Desc, Style: descStyle},
		)
	}

	bgStyle := rich.NewStyle().Background(f.bg)
	s := strip.New(segs)
	s = s.ExtendCellLength(region.Width, bgStyle)
	s = s.TextAlign(region.Width, css.AlignHorizontal("left"))
	strips[0] = s

	for i := 1; i < region.Height; i++ {
		strips[i] = strip.New(nil)
	}
	return strips
}
