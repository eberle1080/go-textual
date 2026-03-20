package widgets

import rich "github.com/eberle1080/go-rich"

// TabPane represents a single tab with a title and content widget.
// It is a data structure, not a Widget itself — TabbedContent manages rendering.
type TabPane struct {
	// Title is the text shown on the tab button.
	Title string
	// Style overrides the default title style.
	Style rich.Style
	// ID is an optional application-assigned identifier.
	ID string
}

// NewTabPane creates a TabPane with the given title.
func NewTabPane(title string) TabPane {
	return TabPane{Title: title}
}

// WithID returns a copy with the given ID set.
func (tp TabPane) WithID(id string) TabPane {
	tp.ID = id
	return tp
}

// WithStyle returns a copy with the given title style.
func (tp TabPane) WithStyle(style rich.Style) TabPane {
	tp.Style = style
	return tp
}
