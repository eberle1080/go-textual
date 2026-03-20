package layout

import (
	"github.com/eberle1080/go-textual/css"
	"github.com/eberle1080/go-textual/geometry"
	"github.com/eberle1080/go-textual/internal/resolve"
)

// HorizontalLayout arranges children side by side, distributing width
// according to each child's CSS width scalar.
type HorizontalLayout struct{}

// Name returns "horizontal".
func (h HorizontalLayout) Name() string { return "horizontal" }

// Arrange positions children left to right within size.
func (h HorizontalLayout) Arrange(
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

	gutter := 0
	if g := parent.Styles().Gutter(); g.Left > 0 || g.Right > 0 {
		gutter = g.Left + g.Right
		if gutter > 1 {
			gutter = 1
		}
	}

	// Build width scalar list from each child's CSS width.
	dims := make([]css.Scalar, len(children))
	for i, child := range children {
		if w := child.Styles().Width(); w != nil {
			dims[i] = *w
		} else {
			dims[i] = css.Scalar{Value: 1, Unit: css.UnitFraction} // default 1fr
		}
	}

	widgets := make([]resolve.BoxModelable, len(children))
	for i, child := range children {
		widgets[i] = &boxModelAdapter{child}
	}

	slots := resolve.Resolve(dims, size.Width, gutter, containerSize, viewport, widgets)

	// For auto-width children, substitute the resolved outer box width from
	// GetBoxModel (which includes padding, borders, and box-sizing adjustments)
	// rather than the raw content width, so bordered or padded auto widgets
	// receive the full allocation they need.
	for i, child := range children {
		if dims[i].IsAuto() && child.Display() {
			bm := child.GetBoxModel(containerSize, viewport, nil, nil, false, greedy)
			if bm.Width != nil {
				f, _ := bm.Width.Float64()
				slots[i].Length = int(f + 0.5)
			} else {
				slots[i].Length = child.GetContentWidth(containerSize, viewport)
			}
		}
	}

	// Recompute offsets after any auto substitution.
	offset := 0
	for i, child := range children {
		if !child.Display() {
			continue
		}
		slots[i].Offset = offset
		offset += slots[i].Length + gutter
	}

	placements := make([]WidgetPlacement, 0, len(children))
	for i, child := range children {
		if !child.Display() {
			continue
		}
		bm := child.GetBoxModel(containerSize, viewport, nil, nil, false, greedy)
		marginTop := bm.Margin.Top
		marginBottom := bm.Margin.Bottom
		y := marginTop
		height := size.Height - marginTop - marginBottom
		if height < 0 {
			height = 0
		}

		placements = append(placements, WidgetPlacement{
			Region: geometry.Region{
				X:      slots[i].Offset,
				Y:      y,
				Width:  slots[i].Length,
				Height: height,
			},
			Margin: bm.Margin,
			Widget: child,
			Order:  child.SortOrder(),
		})
	}
	return placements
}

// GetContentWidth returns the sum of content widths of all children.
func (h HorizontalLayout) GetContentWidth(
	widget Layoutable,
	container, viewport geometry.Size,
) int {
	total := 0
	for _, child := range widget.Children().Slice() {
		if lc, ok := child.(Layoutable); ok {
			total += lc.GetContentWidth(container, viewport)
		}
	}
	return total
}

// GetContentHeight returns the maximum content height of any child.
func (h HorizontalLayout) GetContentHeight(
	widget Layoutable,
	container, viewport geometry.Size,
	width int,
) int {
	maxH := 0
	for _, child := range widget.Children().Slice() {
		if lc, ok := child.(Layoutable); ok {
			ch := lc.GetContentHeight(container, viewport, width)
			if ch > maxH {
				maxH = ch
			}
		}
	}
	return maxH
}
