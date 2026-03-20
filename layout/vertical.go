package layout

import (
	"github.com/eberle1080/go-textual/css"
	"github.com/eberle1080/go-textual/geometry"
	"github.com/eberle1080/go-textual/internal/resolve"
)

// VerticalLayout stacks children vertically, distributing height according to
// each child's CSS height scalar.
type VerticalLayout struct{}

// Name returns "vertical".
func (v VerticalLayout) Name() string { return "vertical" }

// Arrange positions children one below the other within size.
func (v VerticalLayout) Arrange(
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
	if g := parent.Styles().Gutter(); g.Top > 0 || g.Bottom > 0 {
		gutter = g.Top + g.Bottom
		if gutter > 1 {
			gutter = 1 // treat as 1-cell gutter between items
		}
	}

	// Build height scalar list from each child's CSS height.
	dims := make([]css.Scalar, len(children))
	for i, child := range children {
		if h := child.Styles().Height(); h != nil {
			dims[i] = *h
		} else {
			// auto: use content height
			dims[i] = css.Scalar{Value: 1, Unit: css.UnitAuto}
		}
	}

	// Resolve auto heights: for auto children, measure content height.
	widgets := make([]resolve.BoxModelable, len(children))
	for i, child := range children {
		widgets[i] = &boxModelAdapter{child}
	}

	slots := resolve.Resolve(dims, size.Height, gutter, containerSize, viewport, widgets)

	// For auto-height children, derive the slot length from the resolved outer
	// box height (from GetBoxModel) so that padding, borders, and box-sizing
	// adjustments are included, matching the outer-box approach used by
	// HorizontalLayout and GridLayout.
	for i, child := range children {
		if dims[i].IsAuto() && child.Display() {
			bm := child.GetBoxModel(containerSize, viewport, nil, nil, false, greedy)
			if bm.Height != nil {
				f, _ := bm.Height.Float64()
				slots[i].Length = int(f + 0.5)
			} else {
				slots[i].Length = child.GetContentHeight(containerSize, viewport, size.Width)
			}
		}
	}

	// Recompute offsets after auto resolution.
	offset := 0
	for i, child := range children {
		if !child.Display() {
			continue
		}
		slots[i].Offset = offset
		offset += slots[i].Length + gutter
	}

	// Compute total content height for vertical alignment.
	totalHeight := 0
	for i, child := range children {
		if child.Display() {
			totalHeight += slots[i].Length
		}
	}
	if gutter > 0 && totalHeight > 0 {
		displayed := 0
		for _, child := range children {
			if child.Display() {
				displayed++
			}
		}
		if displayed > 1 {
			totalHeight += gutter * (displayed - 1)
		}
	}

	// Vertical alignment offset (align-vertical: middle / bottom).
	yBase := 0
	if st := parent.Styles(); st != nil {
		switch st.AlignVertical() {
		case "middle":
			if size.Height > totalHeight {
				yBase = (size.Height - totalHeight) / 2
			}
		case "bottom":
			if size.Height > totalHeight {
				yBase = size.Height - totalHeight
			}
		}
	}

	placements := make([]WidgetPlacement, 0, len(children))
	for i, child := range children {
		if !child.Display() {
			continue
		}
		bm := child.GetBoxModel(containerSize, viewport, nil, nil, true, greedy)
		marginLeft := bm.Margin.Left
		marginRight := bm.Margin.Right

		// Determine width: respect explicit CSS width, otherwise fill container.
		x := marginLeft
		width := size.Width - marginLeft - marginRight
		if width < 0 {
			width = 0
		}
		if bm.Width != nil {
			if f, _ := bm.Width.Float64(); f > 0 {
				childWidth := int(f + 0.5)
				if childWidth < width {
					width = childWidth
					// Apply horizontal alignment from parent.
					if st := parent.Styles(); st != nil {
						switch st.AlignHorizontal() {
						case "center":
							x = (size.Width - childWidth) / 2
						case "right":
							x = size.Width - childWidth - marginRight
						}
					}
				}
			}
		}

		placements = append(placements, WidgetPlacement{
			Region: geometry.Region{
				X:      x,
				Y:      yBase + slots[i].Offset,
				Width:  width,
				Height: slots[i].Length,
			},
			Margin: bm.Margin,
			Widget: child,
			Order:  child.SortOrder(),
		})
	}
	return placements
}

// GetContentWidth returns the maximum content width of any child.
func (v VerticalLayout) GetContentWidth(
	widget Layoutable,
	container, viewport geometry.Size,
) int {
	maxW := 0
	for _, child := range widget.Children().Slice() {
		if lc, ok := child.(Layoutable); ok {
			w := lc.GetContentWidth(container, viewport)
			if w > maxW {
				maxW = w
			}
		}
	}
	return maxW
}

// GetContentHeight returns the sum of content heights of all children.
func (v VerticalLayout) GetContentHeight(
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
