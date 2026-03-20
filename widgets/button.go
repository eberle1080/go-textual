package widgets

import (
	"context"

	rich "github.com/eberle1080/go-rich"

	"github.com/eberle1080/go-textual/css"
	"github.com/eberle1080/go-textual/dom"
	"github.com/eberle1080/go-textual/geometry"
	"github.com/eberle1080/go-textual/keys"
	"github.com/eberle1080/go-textual/msg"
	"github.com/eberle1080/go-textual/strip"
	"github.com/eberle1080/go-textual/widget"
)

// ButtonPressedMsg is sent when the button is activated.
type ButtonPressedMsg struct {
	msg.BaseMsg
	Button *Button
}

// Button is a focusable widget that can be activated with Enter or Space.
type Button struct {
	widget.BaseWidget
	label   string
	pressed bool
}

// NewButton creates a Button with the given label.
func NewButton(label string) *Button {
	b := &Button{
		BaseWidget: *widget.NewBaseWidget(
			widget.WithDOMOptions(dom.WithCSSTypeName("Button", "Widget")),
			widget.WithCanFocus(true),
		),
		label: label,
	}
	return b
}

// Label returns the button label.
func (b *Button) Label() string { return b.label }

// Update handles key events to activate the button.
func (b *Button) Update(_ context.Context, m msg.Msg) msg.Cmd {
	switch v := m.(type) {
	case msg.KeyMsg:
		switch v.Key {
		case keys.Enter, " ":
			b.pressed = true
			b.MarkDirty()
			btn := b
			return func(ctx context.Context) msg.Msg {
				return ButtonPressedMsg{Button: btn}
			}
		}
	case msg.MouseDownMsg:
		if v.Button == msg.MouseButtonLeft {
			b.pressed = true
			b.MarkDirty()
			btn := b
			return func(ctx context.Context) msg.Msg {
				return ButtonPressedMsg{Button: btn}
			}
		}
	case msg.MouseUpMsg:
		b.pressed = false
		b.MarkDirty()
	case msg.FocusMsg:
		b.MarkDirty()
	case msg.BlurMsg:
		b.pressed = false
		b.MarkDirty()
	}
	return nil
}

// Render draws the button as [ label ] with a focused highlight.
func (b *Button) Render(region geometry.Region) []strip.Strip {
	strips := make([]strip.Strip, region.Height)
	if region.Height == 0 {
		return strips
	}

	style := rich.NewStyle()
	if b.pressed {
		style = style.Bold()
	}

	text := "[ " + b.label + " ]"
	seg := rich.Segment{Text: text, Style: style}

	s := strip.New(rich.Segments{seg})
	s = s.TextAlign(region.Width, css.AlignHorizontal("center"))
	strips[0] = s

	for i := 1; i < region.Height; i++ {
		strips[i] = strip.New(nil)
	}
	return strips
}
