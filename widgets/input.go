package widgets

import (
	"context"
	"unicode/utf8"

	rich "github.com/eberle1080/go-rich"

	"github.com/eberle1080/go-textual/css"
	"github.com/eberle1080/go-textual/dom"
	"github.com/eberle1080/go-textual/geometry"
	"github.com/eberle1080/go-textual/keys"
	"github.com/eberle1080/go-textual/msg"
	"github.com/eberle1080/go-textual/strip"
	"github.com/eberle1080/go-textual/widget"
)

// InputChangedMsg is sent when the input value changes.
type InputChangedMsg struct {
	msg.BaseMsg
	Input *Input
	Value string
}

// InputSubmittedMsg is sent when the user presses Enter.
type InputSubmittedMsg struct {
	msg.BaseMsg
	Input *Input
	Value string
}

// Input is a single-line text input widget.
type Input struct {
	widget.BaseWidget
	value       []rune
	cursor      int
	placeholder string
}

// NewInput creates an Input widget.
func NewInput(placeholder string) *Input {
	inp := &Input{
		BaseWidget: *widget.NewBaseWidget(
			widget.WithDOMOptions(dom.WithCSSTypeName("Input", "Widget")),
			widget.WithCanFocus(true),
		),
		placeholder: placeholder,
	}
	return inp
}

// Value returns the current input value.
func (inp *Input) Value() string { return string(inp.value) }

// SetValue sets the input value and moves the cursor to the end.
func (inp *Input) SetValue(v string) {
	inp.value = []rune(v)
	inp.cursor = len(inp.value)
	inp.MarkDirty()
}

// Clear resets the input to empty.
func (inp *Input) Clear() {
	inp.value = inp.value[:0]
	inp.cursor = 0
	inp.MarkDirty()
}

// Update handles keyboard input.
func (inp *Input) Update(_ context.Context, m msg.Msg) msg.Cmd {
	switch v := m.(type) {
	case msg.FocusMsg:
		inp.MarkDirty()
	case msg.BlurMsg:
		inp.MarkDirty()
	case msg.KeyMsg:
		switch v.Key {
		case keys.Enter:
			val := string(inp.value)
			i := inp
			return func(_ context.Context) msg.Msg {
				return InputSubmittedMsg{Input: i, Value: val}
			}
		case keys.Backspace:
			if inp.cursor > 0 {
				inp.value = append(inp.value[:inp.cursor-1], inp.value[inp.cursor:]...)
				inp.cursor--
				inp.MarkDirty()
				val := string(inp.value)
				i := inp
				return func(_ context.Context) msg.Msg {
					return InputChangedMsg{Input: i, Value: val}
				}
			}
		case keys.Delete:
			if inp.cursor < len(inp.value) {
				inp.value = append(inp.value[:inp.cursor], inp.value[inp.cursor+1:]...)
				inp.MarkDirty()
				val := string(inp.value)
				i := inp
				return func(_ context.Context) msg.Msg {
					return InputChangedMsg{Input: i, Value: val}
				}
			}
		case keys.Left:
			if inp.cursor > 0 {
				inp.cursor--
				inp.MarkDirty()
			}
		case keys.Right:
			if inp.cursor < len(inp.value) {
				inp.cursor++
				inp.MarkDirty()
			}
		case keys.Home:
			inp.cursor = 0
			inp.MarkDirty()
		case keys.End:
			inp.cursor = len(inp.value)
			inp.MarkDirty()
		default:
			if v.IsPrintable() && v.Character != nil {
				ch := []rune(*v.Character)
				if len(ch) > 0 {
					// Insert character at cursor position
					newVal := make([]rune, len(inp.value)+len(ch))
					copy(newVal[:inp.cursor], inp.value[:inp.cursor])
					copy(newVal[inp.cursor:], ch)
					copy(newVal[inp.cursor+len(ch):], inp.value[inp.cursor:])
					inp.value = newVal
					inp.cursor += len(ch)
					inp.MarkDirty()
					val := string(inp.value)
					i := inp
					return func(_ context.Context) msg.Msg {
						return InputChangedMsg{Input: i, Value: val}
					}
				}
			}
		}
	}
	return nil
}

// Render draws the input field.
func (inp *Input) Render(region geometry.Region) []strip.Strip {
	strips := make([]strip.Strip, region.Height)
	if region.Height == 0 {
		return strips
	}

	var text string
	style := rich.NewStyle()

	if len(inp.value) == 0 {
		text = inp.placeholder
		style = style.Foreground(rich.ANSIColor(rich.BrightBlack))
	} else {
		text = string(inp.value)
	}

	// Insert cursor marker if focused (simple approach: show "|" at cursor)
	_ = utf8.RuneCountInString(text) // ensure valid UTF-8

	var segs rich.Segments
	if len(inp.value) == 0 {
		segs = rich.Segments{{Text: text, Style: style}}
	} else {
		runes := []rune(text)
		before := string(runes[:inp.cursor])
		after := string(runes[inp.cursor:])
		segs = rich.Segments{
			{Text: before, Style: style},
			{Text: "|", Style: rich.NewStyle().Reverse()},
			{Text: after, Style: style},
		}
	}

	s := strip.New(segs)
	s = s.TextAlign(region.Width, css.AlignHorizontal("left"))
	strips[0] = s
	for i := 1; i < region.Height; i++ {
		strips[i] = strip.New(nil)
	}
	return strips
}
