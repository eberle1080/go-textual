package widgets

import (
	"context"

	rich "github.com/eberle1080/go-rich"

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
	value        []rune
	cursor       int
	placeholder  string
	scrollOffset int // first visible rune index; keeps cursor in view
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
	inp.scrollOffset = 0
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

// Render draws the input field with a scrolling window that keeps the cursor visible.
//
// Mirrors Textual's approach: render the full content with the cursor styled
// in-place (the character under the cursor is reversed; a space is appended when
// the cursor is at end-of-input), then crop the resulting Strip to the visible
// window [scrollOffset, scrollOffset+width). scrollOffset is adjusted here because
// region.Width is only known at render time.
func (inp *Input) Render(region geometry.Region) []strip.Strip {
	strips := make([]strip.Strip, region.Height)
	if region.Height == 0 {
		return strips
	}

	style := rich.NewStyle()
	cursorStyle := rich.NewStyle().Reverse()

	var s strip.Strip

	if len(inp.value) == 0 {
		placeholder := rich.Segments{{
			Text:  inp.placeholder,
			Style: style.Foreground(rich.ANSIColor(rich.BrightBlack)),
		}}
		s = strip.New(placeholder)
		s = s.AdjustCellLength(region.Width, style)
	} else {
		runes := inp.value
		cursor := inp.cursor

		// Build segments for the full content with cursor styled in-place.
		// When cursor is at end-of-value, append a space as the cursor block
		// (mirrors Textual's pad_right(1) + stylize at end).
		var segs rich.Segments
		before := string(runes[:cursor])
		if cursor < len(runes) {
			cursorChar := string(runes[cursor : cursor+1])
			after := string(runes[cursor+1:])
			segs = rich.Segments{
				{Text: before, Style: style},
				{Text: cursorChar, Style: cursorStyle},
				{Text: after, Style: style},
			}
		} else {
			segs = rich.Segments{
				{Text: before, Style: style},
				{Text: " ", Style: cursorStyle},
			}
		}

		// Adjust scrollOffset so cursor cell stays within [scrollOffset, scrollOffset+width).
		// cursor cell == inp.cursor because this codebase treats 1 rune = 1 cell.
		if inp.cursor < inp.scrollOffset {
			inp.scrollOffset = inp.cursor
		}
		if inp.cursor >= inp.scrollOffset+region.Width {
			inp.scrollOffset = inp.cursor - region.Width + 1
		}
		if inp.scrollOffset < 0 {
			inp.scrollOffset = 0
		}

		// Crop to the visible window, padding to exactly region.Width cells.
		s = strip.New(segs)
		s = s.CropExtend(inp.scrollOffset, inp.scrollOffset+region.Width, style)
	}

	strips[0] = s
	for i := 1; i < region.Height; i++ {
		strips[i] = strip.New(nil)
	}
	return strips
}
