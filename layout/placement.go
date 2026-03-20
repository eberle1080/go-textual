package layout

import (
	"math/big"

	"github.com/eberle1080/go-textual/geometry"
	"github.com/eberle1080/go-textual/internal/spatial"
)

// BoxModel holds the resolved width, height, and margin for a widget.
type BoxModel struct {
	Width  *big.Rat
	Height *big.Rat
	Margin geometry.Spacing
}

// WidgetPlacement describes where one widget is positioned within its parent's
// coordinate space.
type WidgetPlacement struct {
	Region   geometry.Region
	Offset   geometry.Offset
	Margin   geometry.Spacing
	Widget   Layoutable
	Order    int
	Fixed    bool
	Overlay  bool
	Absolute bool
}

// ResetOrigin resets the placement's region to start at (0,0) while
// preserving its width and height.
func (wp *WidgetPlacement) ResetOrigin() {
	wp.Region = geometry.Region{
		X:      0,
		Y:      0,
		Width:  wp.Region.Width,
		Height: wp.Region.Height,
	}
}

// ProcessOffset adjusts the placement's region by the given offset, optionally
// constraining it within constrainRegion. absoluteOffset overrides the position
// entirely if the placement is flagged as absolute.
func (wp WidgetPlacement) ProcessOffset(
	constrainRegion geometry.Region,
	absoluteOffset geometry.Offset,
) WidgetPlacement {
	if wp.Absolute {
		wp.Region.X = absoluteOffset.X
		wp.Region.Y = absoluteOffset.Y
		return wp
	}
	wp.Region.X += constrainRegion.X + wp.Offset.X
	wp.Region.Y += constrainRegion.Y + wp.Offset.Y
	return wp
}

// TranslatePlacements returns a copy of placements with every region shifted
// by offset.
func TranslatePlacements(placements []WidgetPlacement, offset geometry.Offset) []WidgetPlacement {
	if offset.X == 0 && offset.Y == 0 {
		return placements
	}
	result := make([]WidgetPlacement, len(placements))
	for i, p := range placements {
		p.Region.X += offset.X
		p.Region.Y += offset.Y
		result[i] = p
	}
	return result
}

// GetBounds returns the smallest region that contains all non-overlay
// placements.
func GetBounds(placements []WidgetPlacement) geometry.Region {
	if len(placements) == 0 {
		return geometry.Region{}
	}
	minX, minY := int(^uint(0)>>1), int(^uint(0)>>1)
	maxX, maxY := -minX-1, -minY-1
	any := false
	for _, p := range placements {
		if p.Overlay {
			continue
		}
		any = true
		if p.Region.X < minX {
			minX = p.Region.X
		}
		if p.Region.Y < minY {
			minY = p.Region.Y
		}
		x2 := p.Region.X + p.Region.Width
		y2 := p.Region.Y + p.Region.Height
		if x2 > maxX {
			maxX = x2
		}
		if y2 > maxY {
			maxY = y2
		}
	}
	if !any {
		return geometry.Region{}
	}
	return geometry.Region{X: minX, Y: minY, Width: maxX - minX, Height: maxY - minY}
}

// DockArrangeResult holds all widget placements produced by a single arrange
// call, plus cached spatial-map and scroll-spacing data.
type DockArrangeResult struct {
	Placements    []WidgetPlacement
	Widgets       map[Layoutable]bool
	ScrollSpacing geometry.Spacing

	spatialMap *spatial.SpatialMap[WidgetPlacement]
}

// NewDockArrangeResult constructs a DockArrangeResult from placements.
func NewDockArrangeResult(placements []WidgetPlacement, scrollSpacing geometry.Spacing) DockArrangeResult {
	widgets := make(map[Layoutable]bool, len(placements))
	for _, p := range placements {
		if p.Widget != nil {
			widgets[p.Widget] = true
		}
	}
	return DockArrangeResult{
		Placements:    placements,
		Widgets:       widgets,
		ScrollSpacing: scrollSpacing,
	}
}

// TotalRegion returns the bounding region of all non-overlay placements.
func (r *DockArrangeResult) TotalRegion() geometry.Region {
	return GetBounds(r.Placements)
}

// SpatialMap lazily builds and returns a spatial.SpatialMap over the
// placements so that visible placements can be queried efficiently.
func (r *DockArrangeResult) SpatialMap() *spatial.SpatialMap[WidgetPlacement] {
	if r.spatialMap != nil {
		return r.spatialMap
	}
	total := r.TotalRegion()
	// Choose a cell size that gives roughly 10×10 grid cells.
	cellW := 1
	cellH := 1
	if total.Width > 10 {
		cellW = total.Width / 10
	}
	if total.Height > 10 {
		cellH = total.Height / 10
	}
	sm := spatial.New[WidgetPlacement](total, cellW, cellH)
	entries := make([]spatial.Entry[WidgetPlacement], len(r.Placements))
	for i, p := range r.Placements {
		entries[i] = spatial.Entry[WidgetPlacement]{
			Region: p.Region,
			Value:  p,
			Fixed:  p.Fixed || p.Overlay,
		}
	}
	sm.Insert(entries)
	r.spatialMap = sm
	return r.spatialMap
}

// GetVisiblePlacements returns placements whose regions overlap with region.
func (r *DockArrangeResult) GetVisiblePlacements(region geometry.Region) []WidgetPlacement {
	return r.SpatialMap().GetValuesInRegion(region)
}
