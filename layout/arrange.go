package layout

import (
	"github.com/eberle1080/go-textual/geometry"
	"github.com/eberle1080/go-textual/internal/partition"
)

// Arrange positions a widget's children within the given size and returns a
// DockArrangeResult. If optimal is true the layout may use content-size
// heuristics to produce a more accurate arrangement.
func Arrange(
	widget Layoutable,
	children []Layoutable,
	size, viewport geometry.Size,
	optimal bool,
) DockArrangeResult {
	if len(children) == 0 {
		return NewDockArrangeResult(nil, geometry.Spacing{})
	}

	// Separate docked widgets from non-docked.
	// Only treat a widget as docked when dock is a real edge value; the
	// CSS default "none" (returned by Dock() when unset) must not classify
	// ordinary children as docked.
	nonDocked, docked := partition.Partition(func(w Layoutable) bool {
		d := w.Styles().Dock()
		return d == "top" || d == "bottom" || d == "left" || d == "right"
	}, children)

	// Separate split widgets from regular.
	// HasRule("split") is true when the rule is explicitly set; we also
	// exclude the explicit "none" value so that split:none widgets go
	// through the normal layout path.
	nonSplit, split := partition.Partition(func(w Layoutable) bool {
		v, ok := w.Styles().GetRule("split")
		if !ok {
			return false
		}
		s, ok := v.(string)
		return ok && s != "" && s != "none"
	}, nonDocked)

	// Start with the full region available to the layout.
	mainRegion := geometry.Region{X: 0, Y: 0, Width: size.Width, Height: size.Height}

	var allPlacements []WidgetPlacement
	scrollSpacing := geometry.Spacing{}

	// Arrange docked widgets (they consume space from the edges).
	if len(docked) > 0 {
		dockedPlacements, remaining := arrangeDockWidgets(docked, mainRegion, viewport)
		allPlacements = append(allPlacements, dockedPlacements...)

		// Compute scroll spacing from what the docked widgets claimed.
		scrollSpacing = geometry.Spacing{
			Top:    remaining.Y - mainRegion.Y,
			Right:  (mainRegion.X + mainRegion.Width) - (remaining.X + remaining.Width),
			Bottom: (mainRegion.Y + mainRegion.Height) - (remaining.Y + remaining.Height),
			Left:   remaining.X - mainRegion.X,
		}
		mainRegion = remaining
	}

	// Arrange split widgets.
	if len(split) > 0 {
		splitPlacements, remaining := arrangeSplitWidgets(split, mainRegion, viewport)
		allPlacements = append(allPlacements, splitPlacements...)
		mainRegion = remaining
	}

	// Arrange the main (non-docked, non-split) widgets using the widget's layout.
	if len(nonSplit) > 0 {
		layout := widget.Styles().Layout()
		l, err := GetLayout(layout)
		if err != nil {
			l, _ = GetLayout("vertical")
		}

		widget.PreLayout(l)
		mainSize := geometry.Size{Width: mainRegion.Width, Height: mainRegion.Height}
		mainPlacements := l.Arrange(widget, nonSplit, mainSize, false)

		// Translate placements to the main region's origin.
		offset := geometry.Offset{X: mainRegion.X, Y: mainRegion.Y}
		mainPlacements = TranslatePlacements(mainPlacements, offset)
		allPlacements = append(allPlacements, mainPlacements...)
	}

	// Process absolute-offset overrides.
	processed := make([]WidgetPlacement, len(allPlacements))
	for i, p := range allPlacements {
		if p.Widget != nil {
			if abs := p.Widget.AbsoluteOffset(); abs != nil {
				p.Absolute = true
				p = p.ProcessOffset(mainRegion, *abs)
			}
		}
		processed[i] = p
	}

	result := NewDockArrangeResult(processed, scrollSpacing)
	return result
}

// arrangeDockWidgets consumes space from the edges of region for each docked
// widget and returns the placements and the remaining region.
func arrangeDockWidgets(
	docked []Layoutable,
	region geometry.Region,
	viewport geometry.Size,
) ([]WidgetPlacement, geometry.Region) {
	var placements []WidgetPlacement
	remaining := region

	for _, w := range docked {
		dock := w.Styles().Dock()
		containerSize := geometry.Size{Width: remaining.Width, Height: remaining.Height}
		bm := w.GetBoxModel(containerSize, viewport, nil, nil, true, false)

		var widgetRegion geometry.Region
		wWidth := 0
		wHeight := 0
		if bm.Width != nil {
			f, _ := bm.Width.Float64()
			wWidth = int(f + 0.5)
		}
		if bm.Height != nil {
			f, _ := bm.Height.Float64()
			wHeight = int(f + 0.5)
		}

		switch dock {
		case "top":
			if wHeight == 0 {
				wHeight = 1
			}
			widgetRegion = geometry.Region{
				X:      remaining.X,
				Y:      remaining.Y,
				Width:  remaining.Width,
				Height: wHeight,
			}
			remaining.Y += wHeight
			remaining.Height -= wHeight
		case "bottom":
			if wHeight == 0 {
				wHeight = 1
			}
			widgetRegion = geometry.Region{
				X:      remaining.X,
				Y:      remaining.Y + remaining.Height - wHeight,
				Width:  remaining.Width,
				Height: wHeight,
			}
			remaining.Height -= wHeight
		case "left":
			if wWidth == 0 {
				wWidth = 1
			}
			widgetRegion = geometry.Region{
				X:      remaining.X,
				Y:      remaining.Y,
				Width:  wWidth,
				Height: remaining.Height,
			}
			remaining.X += wWidth
			remaining.Width -= wWidth
		case "right":
			if wWidth == 0 {
				wWidth = 1
			}
			widgetRegion = geometry.Region{
				X:      remaining.X + remaining.Width - wWidth,
				Y:      remaining.Y,
				Width:  wWidth,
				Height: remaining.Height,
			}
			remaining.Width -= wWidth
		}

		placements = append(placements, WidgetPlacement{
			Region: widgetRegion,
			Widget: w,
			Fixed:  true,
		})
	}
	return placements, remaining
}

// arrangeSplitWidgets positions split-off (sidebar/panel) widgets beside the
// main region, honouring the CSS split edge ("left", "right", "top",
// "bottom"). The default edge when unset or "none" is "right".
func arrangeSplitWidgets(
	split []Layoutable,
	region geometry.Region,
	viewport geometry.Size,
) ([]WidgetPlacement, geometry.Region) {
	var placements []WidgetPlacement
	remaining := region

	for _, w := range split {
		containerSize := geometry.Size{Width: remaining.Width, Height: remaining.Height}
		bm := w.GetBoxModel(containerSize, viewport, nil, nil, true, false)

		// Determine the split edge from CSS.
		edge := "right"
		if v, ok := w.Styles().GetRule("split"); ok {
			if s, ok := v.(string); ok && s != "" && s != "none" {
				edge = s
			}
		}

		var widgetRegion geometry.Region

		switch edge {
		case "left":
			wWidth := 0
			if bm.Width != nil {
				f, _ := bm.Width.Float64()
				wWidth = int(f + 0.5)
			}
			if wWidth == 0 {
				wWidth = remaining.Width / 4
			}
			widgetRegion = geometry.Region{
				X:      remaining.X,
				Y:      remaining.Y,
				Width:  wWidth,
				Height: remaining.Height,
			}
			remaining.X += wWidth
			remaining.Width -= wWidth

		case "top":
			wHeight := 0
			if bm.Height != nil {
				f, _ := bm.Height.Float64()
				wHeight = int(f + 0.5)
			}
			if wHeight == 0 {
				wHeight = remaining.Height / 4
			}
			widgetRegion = geometry.Region{
				X:      remaining.X,
				Y:      remaining.Y,
				Width:  remaining.Width,
				Height: wHeight,
			}
			remaining.Y += wHeight
			remaining.Height -= wHeight

		case "bottom":
			wHeight := 0
			if bm.Height != nil {
				f, _ := bm.Height.Float64()
				wHeight = int(f + 0.5)
			}
			if wHeight == 0 {
				wHeight = remaining.Height / 4
			}
			widgetRegion = geometry.Region{
				X:      remaining.X,
				Y:      remaining.Y + remaining.Height - wHeight,
				Width:  remaining.Width,
				Height: wHeight,
			}
			remaining.Height -= wHeight

		default: // "right"
			wWidth := 0
			if bm.Width != nil {
				f, _ := bm.Width.Float64()
				wWidth = int(f + 0.5)
			}
			if wWidth == 0 {
				wWidth = remaining.Width / 4
			}
			widgetRegion = geometry.Region{
				X:      remaining.X + remaining.Width - wWidth,
				Y:      remaining.Y,
				Width:  wWidth,
				Height: remaining.Height,
			}
			remaining.Width -= wWidth
		}

		placements = append(placements, WidgetPlacement{
			Region: widgetRegion,
			Widget: w,
		})
	}
	return placements, remaining
}

// buildLayers groups widgets by their CSS layer name.
func buildLayers(widgets []Layoutable) map[string][]Layoutable {
	layers := make(map[string][]Layoutable)
	for _, w := range widgets {
		layer := w.Layer()
		layers[layer] = append(layers[layer], w)
	}
	return layers
}
