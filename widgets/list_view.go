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

// ListViewSelectedMsg is sent when the user activates an item with Enter.
type ListViewSelectedMsg struct {
	msg.BaseMsg
	Index int
	Item  ListItem
}

// ListView displays a scrollable, focusable list of ListItems.
type ListView struct {
	widget.BaseWidget
	items  []ListItem
	cursor int
}

// NewListView creates an empty ListView.
func NewListView() *ListView {
	lv := &ListView{
		BaseWidget: *widget.NewBaseWidget(
			widget.WithDOMOptions(dom.WithCSSTypeName("ListView", "Widget")),
			widget.WithCanFocus(true),
		),
	}
	return lv
}

// SetItems replaces all items and resets the cursor.
func (lv *ListView) SetItems(items []ListItem) {
	lv.items = items
	lv.cursor = 0
	lv.MarkDirty()
}

// AppendItem adds a single item.
func (lv *ListView) AppendItem(item ListItem) {
	lv.items = append(lv.items, item)
	lv.MarkDirty()
}

// Items returns the current items.
func (lv *ListView) Items() []ListItem { return lv.items }

// Cursor returns the current cursor index.
func (lv *ListView) Cursor() int { return lv.cursor }

// SetCursor moves the cursor to index i (clamped to valid range).
func (lv *ListView) SetCursor(i int) {
	if i < 0 {
		i = 0
	}
	if i >= len(lv.items) && len(lv.items) > 0 {
		i = len(lv.items) - 1
	}
	lv.cursor = i
	lv.MarkDirty()
}

// Update handles keyboard navigation.
func (lv *ListView) Update(_ context.Context, m msg.Msg) msg.Cmd {
	switch v := m.(type) {
	case msg.FocusMsg:
		lv.MarkDirty()
	case msg.BlurMsg:
		lv.MarkDirty()
	case msg.MouseDownMsg:
		// Translate Y to an item index.
		idx := int(v.Y)
		// Account for scroll offset.
		visibleStart := 0
		region := lv.Region()
		if lv.cursor >= region.Height && region.Height > 0 {
			visibleStart = lv.cursor - region.Height + 1
		}
		idx = visibleStart + idx
		if idx >= 0 && idx < len(lv.items) {
			lv.cursor = idx
			lv.MarkDirty()
			if v.Button == msg.MouseButtonLeft {
				item := lv.items[idx]
				return func(_ context.Context) msg.Msg {
					return ListViewSelectedMsg{Index: idx, Item: item}
				}
			}
		}
		return nil
	case msg.KeyMsg:
		switch v.Key {
		case keys.Up:
			if lv.cursor > 0 {
				lv.cursor--
				lv.MarkDirty()
			}
		case keys.Down:
			if lv.cursor < len(lv.items)-1 {
				lv.cursor++
				lv.MarkDirty()
			}
		case keys.Home:
			if lv.cursor != 0 {
				lv.cursor = 0
				lv.MarkDirty()
			}
		case keys.End:
			if last := len(lv.items) - 1; lv.cursor != last && last >= 0 {
				lv.cursor = last
				lv.MarkDirty()
			}
		case keys.PageUp:
			lv.cursor -= 10
			if lv.cursor < 0 {
				lv.cursor = 0
			}
			lv.MarkDirty()
		case keys.PageDown:
			lv.cursor += 10
			if lv.cursor >= len(lv.items) && len(lv.items) > 0 {
				lv.cursor = len(lv.items) - 1
			}
			lv.MarkDirty()
		case keys.Enter, " ":
			if lv.cursor < len(lv.items) {
				idx := lv.cursor
				item := lv.items[idx]
				return func(_ context.Context) msg.Msg {
					return ListViewSelectedMsg{Index: idx, Item: item}
				}
			}
		}
	}
	return nil
}

// Render draws the visible items.
func (lv *ListView) Render(region geometry.Region) []strip.Strip {
	strips := make([]strip.Strip, region.Height)
	if region.Height == 0 {
		return strips
	}

	// Scroll window: keep cursor visible.
	start := 0
	if lv.cursor >= region.Height {
		start = lv.cursor - region.Height + 1
	}

	for row := range region.Height {
		idx := start + row
		if idx >= len(lv.items) {
			strips[row] = strip.New(nil)
			continue
		}
		item := lv.items[idx]

		prefix := item.Prefix
		if prefix != "" {
			prefix += " "
		}
		text := prefix + item.Text

		var style rich.Style
		if idx == lv.cursor {
			style = rich.NewStyle().Reverse()
		} else {
			style = item.Style
		}

		seg := rich.Segment{Text: text, Style: style}
		s := strip.New(rich.Segments{seg})
		s = s.TextAlign(region.Width, css.AlignHorizontal("left"))
		strips[row] = s
	}
	return strips
}
