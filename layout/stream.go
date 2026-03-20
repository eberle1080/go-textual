package layout

import (
	"github.com/eberle1080/go-textual/geometry"
)

// StreamLayout is a simplified vertical layout where all widgets have full
// width and automatic height. It does not handle fr units or gutters; it is
// intended for stream-style content (logs, feeds, chat) where heights are
// content-driven.
type StreamLayout struct{}

// Name returns "stream".
func (s StreamLayout) Name() string { return "stream" }

// Arrange stacks children vertically at full width with content-driven heights.
func (s StreamLayout) Arrange(
	parent Layoutable,
	children []Layoutable,
	size geometry.Size,
	greedy bool,
) []WidgetPlacement {
	if len(children) == 0 {
		return nil
	}

	viewport := parent.ViewportSize()
	containerSize := size

	placements := make([]WidgetPlacement, 0, len(children))
	y := 0
	for _, child := range children {
		if !child.Display() {
			continue
		}
		h := child.GetContentHeight(containerSize, viewport, size.Width)
		if h == 0 {
			h = 1
		}
		placements = append(placements, WidgetPlacement{
			Region: geometry.Region{
				X:      0,
				Y:      y,
				Width:  size.Width,
				Height: h,
			},
			Widget: child,
			Order:  child.SortOrder(),
		})
		y += h
	}
	return placements
}

// GetContentWidth returns the container width (stream content is always full width).
func (s StreamLayout) GetContentWidth(
	widget Layoutable,
	container, viewport geometry.Size,
) int {
	return container.Width
}

// GetContentHeight returns the sum of all children's content heights.
func (s StreamLayout) GetContentHeight(
	widget Layoutable,
	container, viewport geometry.Size,
	width int,
) int {
	total := 0
	for _, child := range widget.Children().Slice() {
		if lc, ok := child.(Layoutable); ok {
			total += lc.GetContentHeight(container, viewport, width)
		}
	}
	return total
}
