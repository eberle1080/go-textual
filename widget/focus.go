package widget

// FocusableWidgets performs a DFS traversal of the widget tree and returns
// all widgets that report CanFocus() == true, in tree order.
func FocusableWidgets(root Widget) []Widget {
	var result []Widget
	var walk func(w Widget)
	walk = func(w Widget) {
		if w.CanFocus() {
			result = append(result, w)
		}
		for _, child := range w.WidgetChildren() {
			walk(child)
		}
	}
	walk(root)
	return result
}

// NextFocus returns the next focusable widget after current in tree order.
// If current is nil or not found, returns the first focusable widget.
// Returns nil if there are no focusable widgets.
func NextFocus(root Widget, current Widget) Widget {
	focusable := FocusableWidgets(root)
	if len(focusable) == 0 {
		return nil
	}
	for i, w := range focusable {
		if w == current {
			return focusable[(i+1)%len(focusable)]
		}
	}
	return focusable[0]
}

// PrevFocus returns the previous focusable widget before current in tree order.
// If current is nil or not found, returns the last focusable widget.
// Returns nil if there are no focusable widgets.
func PrevFocus(root Widget, current Widget) Widget {
	focusable := FocusableWidgets(root)
	if len(focusable) == 0 {
		return nil
	}
	for i, w := range focusable {
		if w == current {
			return focusable[(i-1+len(focusable))%len(focusable)]
		}
	}
	return focusable[len(focusable)-1]
}
