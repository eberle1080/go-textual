package widgets

import (
	rich "github.com/eberle1080/go-rich"
)

// ListItem is a single item in a ListView.
type ListItem struct {
	// Text is the main label shown in the list.
	Text string
	// Prefix is an optional icon/symbol prepended to the text.
	Prefix string
	// Style overrides the default text style when non-zero.
	Style rich.Style
	// Data is an arbitrary value that the application can attach to the item.
	Data any
}

// NewListItem creates a ListItem with the given text.
func NewListItem(text string) ListItem {
	return ListItem{Text: text}
}

// WithPrefix returns a copy of the item with the given prefix.
func (li ListItem) WithPrefix(prefix string) ListItem {
	li.Prefix = prefix
	return li
}

// WithStyle returns a copy of the item with the given style.
func (li ListItem) WithStyle(style rich.Style) ListItem {
	li.Style = style
	return li
}

// WithData returns a copy of the item with the given data.
func (li ListItem) WithData(data any) ListItem {
	li.Data = data
	return li
}
