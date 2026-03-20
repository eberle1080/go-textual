package widgets

import (
	"context"
	"strings"

	rich "github.com/eberle1080/go-rich"

	"github.com/eberle1080/go-textual/css"
	"github.com/eberle1080/go-textual/dom"
	"github.com/eberle1080/go-textual/geometry"
	"github.com/eberle1080/go-textual/keys"
	"github.com/eberle1080/go-textual/msg"
	"github.com/eberle1080/go-textual/strip"
	"github.com/eberle1080/go-textual/widget"
)

// TextArea displays and optionally edits multi-line text.
// When ReadOnly is true it acts as a scrollable text viewer.
type TextArea struct {
	widget.BaseWidget
	lines    []string
	scrollY  int
	cursorX  int
	cursorY  int
	readOnly bool
}

// NewTextArea creates a TextArea pre-filled with text.
// Pass readOnly=true for a plain scrollable viewer.
func NewTextArea(text string, readOnly bool) *TextArea {
	ta := &TextArea{
		BaseWidget: *widget.NewBaseWidget(
			widget.WithDOMOptions(dom.WithCSSTypeName("TextArea", "Widget")),
			widget.WithCanFocus(!readOnly),
		),
		readOnly: readOnly,
	}
	ta.SetText(text)
	return ta
}

// SetText replaces the content and resets the scroll/cursor position.
func (ta *TextArea) SetText(text string) {
	ta.lines = strings.Split(text, "\n")
	ta.scrollY = 0
	ta.cursorX = 0
	ta.cursorY = 0
	ta.MarkDirty()
}

// Text returns the current content as a single string.
func (ta *TextArea) Text() string { return strings.Join(ta.lines, "\n") }

// AppendLine adds a line at the end and scrolls to show it.
func (ta *TextArea) AppendLine(line string) {
	ta.lines = append(ta.lines, line)
	ta.MarkDirty()
}

func (ta *TextArea) currentLine() []rune {
	if ta.cursorY >= len(ta.lines) {
		return nil
	}
	return []rune(ta.lines[ta.cursorY])
}

func (ta *TextArea) setLine(y int, runes []rune) {
	for len(ta.lines) <= y {
		ta.lines = append(ta.lines, "")
	}
	ta.lines[y] = string(runes)
}

// Update handles scroll and (when editable) keyboard input.
func (ta *TextArea) Update(_ context.Context, m msg.Msg) msg.Cmd {
	switch v := m.(type) {
	case msg.FocusMsg:
		ta.MarkDirty()
	case msg.BlurMsg:
		ta.MarkDirty()
	case msg.KeyMsg:
		switch v.Key {
		case keys.Up:
			if ta.cursorY > 0 {
				ta.cursorY--
				ta.clampCursorX()
				ta.MarkDirty()
			}
		case keys.Down:
			if ta.cursorY < len(ta.lines)-1 {
				ta.cursorY++
				ta.clampCursorX()
				ta.MarkDirty()
			}
		case keys.PageUp:
			ta.scrollY -= 10
			ta.cursorY -= 10
			ta.clampScroll(0)
			ta.MarkDirty()
		case keys.PageDown:
			ta.scrollY += 10
			ta.cursorY += 10
			ta.clampScroll(0)
			ta.MarkDirty()
		default:
			if ta.readOnly {
				return nil
			}
			switch v.Key {
			case keys.Left:
				if ta.cursorX > 0 {
					ta.cursorX--
					ta.MarkDirty()
				} else if ta.cursorY > 0 {
					ta.cursorY--
					ta.cursorX = len([]rune(ta.lines[ta.cursorY]))
					ta.MarkDirty()
				}
			case keys.Right:
				line := ta.currentLine()
				if ta.cursorX < len(line) {
					ta.cursorX++
					ta.MarkDirty()
				} else if ta.cursorY < len(ta.lines)-1 {
					ta.cursorY++
					ta.cursorX = 0
					ta.MarkDirty()
				}
			case keys.Home:
				ta.cursorX = 0
				ta.MarkDirty()
			case keys.End:
				ta.cursorX = len(ta.currentLine())
				ta.MarkDirty()
			case keys.Enter:
				line := ta.currentLine()
				before := string(line[:ta.cursorX])
				after := string(line[ta.cursorX:])
				ta.lines[ta.cursorY] = before
				newLines := make([]string, len(ta.lines)+1)
				copy(newLines, ta.lines[:ta.cursorY+1])
				newLines[ta.cursorY+1] = after
				copy(newLines[ta.cursorY+2:], ta.lines[ta.cursorY+1:])
				ta.lines = newLines
				ta.cursorY++
				ta.cursorX = 0
				ta.MarkDirty()
			case keys.Backspace:
				line := ta.currentLine()
				if ta.cursorX > 0 {
					newLine := append(line[:ta.cursorX-1:ta.cursorX-1], line[ta.cursorX:]...)
					ta.setLine(ta.cursorY, newLine)
					ta.cursorX--
					ta.MarkDirty()
				} else if ta.cursorY > 0 {
					prev := []rune(ta.lines[ta.cursorY-1])
					cur := ta.currentLine()
					ta.cursorX = len(prev)
					merged := append(prev, cur...)
					ta.lines = append(ta.lines[:ta.cursorY-1], ta.lines[ta.cursorY:]...)
					ta.cursorY--
					ta.setLine(ta.cursorY, merged)
					ta.MarkDirty()
				}
			default:
				if v.IsPrintable() && v.Character != nil {
					ch := []rune(*v.Character)
					if len(ch) > 0 {
						line := ta.currentLine()
						newLine := make([]rune, len(line)+len(ch))
						copy(newLine[:ta.cursorX], line[:ta.cursorX])
						copy(newLine[ta.cursorX:], ch)
						copy(newLine[ta.cursorX+len(ch):], line[ta.cursorX:])
						ta.setLine(ta.cursorY, newLine)
						ta.cursorX += len(ch)
						ta.MarkDirty()
					}
				}
			}
		}
	}
	return nil
}

func (ta *TextArea) clampCursorX() {
	if ta.cursorY < len(ta.lines) {
		max := len([]rune(ta.lines[ta.cursorY]))
		if ta.cursorX > max {
			ta.cursorX = max
		}
	}
}

func (ta *TextArea) clampScroll(_ int) {
	if ta.scrollY < 0 {
		ta.scrollY = 0
	}
	if ta.cursorY < 0 {
		ta.cursorY = 0
	}
	if ta.cursorY >= len(ta.lines) && len(ta.lines) > 0 {
		ta.cursorY = len(ta.lines) - 1
	}
}

// Render draws the visible lines of text.
func (ta *TextArea) Render(region geometry.Region) []strip.Strip {
	strips := make([]strip.Strip, region.Height)
	if region.Height == 0 {
		return strips
	}

	// Keep cursor in view.
	if ta.cursorY < ta.scrollY {
		ta.scrollY = ta.cursorY
	}
	if ta.cursorY >= ta.scrollY+region.Height {
		ta.scrollY = ta.cursorY - region.Height + 1
	}

	defaultStyle := rich.NewStyle()
	cursorStyle := rich.NewStyle().Reverse()

	for row := range region.Height {
		lineIdx := ta.scrollY + row
		if lineIdx >= len(ta.lines) {
			strips[row] = strip.New(nil)
			continue
		}

		lineRunes := []rune(ta.lines[lineIdx])

		if !ta.readOnly && lineIdx == ta.cursorY {
			// Draw line with cursor inserted.
			before := string(lineRunes[:ta.cursorX])
			var cursor string
			if ta.cursorX < len(lineRunes) {
				cursor = string(lineRunes[ta.cursorX : ta.cursorX+1])
			} else {
				cursor = " "
			}
			after := ""
			if ta.cursorX+1 < len(lineRunes) {
				after = string(lineRunes[ta.cursorX+1:])
			}
			segs := rich.Segments{
				{Text: before, Style: defaultStyle},
				{Text: cursor, Style: cursorStyle},
				{Text: after, Style: defaultStyle},
			}
			s := strip.New(segs)
			s = s.TextAlign(region.Width, css.AlignHorizontal("left"))
			strips[row] = s
		} else {
			seg := rich.Segment{Text: string(lineRunes), Style: defaultStyle}
			s := strip.New(rich.Segments{seg})
			s = s.TextAlign(region.Width, css.AlignHorizontal("left"))
			strips[row] = s
		}
	}
	return strips
}
