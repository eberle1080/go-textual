package geometry

import (
	"fmt"
	"math"
)

// NullOffset is the zero-value Offset at the origin.
var NullOffset = Offset{}

// NullSize is the zero-value Size with no area.
var NullSize = Size{}

// NullRegion is the zero-value Region at the origin with no area.
var NullRegion = Region{}

// NullSpacing is the zero-value Spacing with no margins.
var NullSpacing = Spacing{}

// ClampInt restricts an integer value to the range [min, max].
// If min > max the arguments are treated as reversed.
func ClampInt(value, minimum, maximum int) int {
	if minimum > maximum {
		if value < maximum {
			return maximum
		}
		if value > minimum {
			return minimum
		}
		return value
	}
	if value < minimum {
		return minimum
	}
	if value > maximum {
		return maximum
	}
	return value
}

// ClampFloat restricts a float64 value to the range [min, max].
// If min > max the arguments are treated as reversed.
func ClampFloat(value, minimum, maximum float64) float64 {
	if minimum > maximum {
		if value < maximum {
			return maximum
		}
		if value > minimum {
			return minimum
		}
		return value
	}
	if value < minimum {
		return minimum
	}
	if value > maximum {
		return maximum
	}
	return value
}

// Offset is a cell coordinate defined by X and Y values.
// Offsets are typically relative to the top left of the terminal or container.
// X corresponds to the column; Y corresponds to the row.
type Offset struct {
	X int
	Y int
}

// IsOrigin reports whether the offset is at (0, 0).
func (o Offset) IsOrigin() bool {
	return o.X == 0 && o.Y == 0
}

// IsNonZero reports whether the offset is not at (0, 0).
func (o Offset) IsNonZero() bool {
	return o.X != 0 || o.Y != 0
}

// Clamped returns a copy of the offset with X and Y restricted to values >= 0.
func (o Offset) Clamped() Offset {
	x := o.X
	if x < 0 {
		x = 0
	}
	y := o.Y
	if y < 0 {
		y = 0
	}
	return Offset{x, y}
}

// Transpose returns (Y, X) — the offset with axes swapped.
func (o Offset) Transpose() (int, int) {
	return o.Y, o.X
}

// Add returns the sum of two offsets.
func (o Offset) Add(other Offset) Offset {
	return Offset{o.X + other.X, o.Y + other.Y}
}

// Sub returns the difference of two offsets.
func (o Offset) Sub(other Offset) Offset {
	return Offset{o.X - other.X, o.Y - other.Y}
}

// MulScalar multiplies both components by a scalar factor and truncates to int.
func (o Offset) MulScalar(factor float64) Offset {
	return Offset{int(float64(o.X) * factor), int(float64(o.Y) * factor)}
}

// MulXY multiplies X by fx and Y by fy, truncating to int.
func (o Offset) MulXY(fx, fy float64) Offset {
	return Offset{int(float64(o.X) * fx), int(float64(o.Y) * fy)}
}

// Neg returns the negation of the offset.
func (o Offset) Neg() Offset {
	return Offset{-o.X, -o.Y}
}

// Blend returns a new offset on the line between o and dest at the given factor.
// factor=0 returns o; factor=1 returns dest.
func (o Offset) Blend(dest Offset, factor float64) Offset {
	return Offset{
		int(float64(o.X) + float64(dest.X-o.X)*factor),
		int(float64(o.Y) + float64(dest.Y-o.Y)*factor),
	}
}

// DistanceTo returns the Euclidean distance to another offset.
func (o Offset) DistanceTo(other Offset) float64 {
	dx := float64(other.X - o.X)
	dy := float64(other.Y - o.Y)
	return math.Sqrt(dx*dx + dy*dy)
}

// Clamp restricts the offset to fit within a rectangle of [0, width-1] x [0, height-1].
func (o Offset) Clamp(width, height int) Offset {
	return Offset{
		ClampInt(o.X, 0, width-1),
		ClampInt(o.Y, 0, height-1),
	}
}

// Size represents the dimensions (width and height) of a rectangular region.
type Size struct {
	Width  int
	Height int
}

// IsNonZero reports whether the size has non-zero area.
func (s Size) IsNonZero() bool {
	return s.Width*s.Height != 0
}

// Area returns the area (Width * Height).
func (s Size) Area() int {
	return s.Width * s.Height
}

// Region returns a Region of this size at the origin.
func (s Size) Region() Region {
	return Region{0, 0, s.Width, s.Height}
}

// LineRange returns (0, s.Height), the range of row indices.
func (s Size) LineRange() (int, int) {
	return 0, s.Height
}

// WithWidth returns a new Size with only the width changed.
func (s Size) WithWidth(w int) Size {
	return Size{w, s.Height}
}

// WithHeight returns a new Size with only the height changed.
func (s Size) WithHeight(h int) Size {
	return Size{s.Width, h}
}

// Add returns the element-wise sum, clamping each dimension to 0.
func (s Size) Add(other Size) Size {
	w := s.Width + other.Width
	if w < 0 {
		w = 0
	}
	h := s.Height + other.Height
	if h < 0 {
		h = 0
	}
	return Size{w, h}
}

// Sub returns the element-wise difference, clamping each dimension to 0.
func (s Size) Sub(other Size) Size {
	w := s.Width - other.Width
	if w < 0 {
		w = 0
	}
	h := s.Height - other.Height
	if h < 0 {
		h = 0
	}
	return Size{w, h}
}

// Contains reports whether the point (x, y) falls within the size boundary.
func (s Size) Contains(x, y int) bool {
	return s.Width > x && x >= 0 && s.Height > y && y >= 0
}

// ContainsPoint reports whether the offset falls within the size boundary.
func (s Size) ContainsPoint(p Offset) bool {
	return s.Width > p.X && p.X >= 0 && s.Height > p.Y && p.Y >= 0
}

// ClampOffset clamps an offset to fit within the size boundary.
func (s Size) ClampOffset(o Offset) Offset {
	return o.Clamp(s.Width, s.Height)
}

// Region defines a rectangular area by position and dimensions.
//
//	(x, y)
//	  ┌────────────────────┐ ▲
//	  │                    │ │
//	  │                    │ height
//	  │                    │ │
//	  └────────────────────┘ ▼
//	  ◀─────── width ──────▶
type Region struct {
	X      int
	Y      int
	Width  int
	Height int
}

// RegionFromUnion returns the smallest Region that encloses all given regions.
// Returns an error if regions is empty.
func RegionFromUnion(regions []Region) (Region, error) {
	if len(regions) == 0 {
		return NullRegion, fmt.Errorf("at least one region expected")
	}
	minX := regions[0].X
	minY := regions[0].Y
	maxX := regions[0].X + regions[0].Width
	maxY := regions[0].Y + regions[0].Height
	for _, r := range regions[1:] {
		if r.X < minX {
			minX = r.X
		}
		if r.Y < minY {
			minY = r.Y
		}
		rx2 := r.X + r.Width
		ry2 := r.Y + r.Height
		if rx2 > maxX {
			maxX = rx2
		}
		if ry2 > maxY {
			maxY = ry2
		}
	}
	return Region{minX, minY, maxX - minX, maxY - minY}, nil
}

// RegionFromCorners creates a Region from top-left (x1,y1) and bottom-right (x2,y2) corners.
func RegionFromCorners(x1, y1, x2, y2 int) Region {
	return Region{x1, y1, x2 - x1, y2 - y1}
}

// RegionFromOffset creates a Region from an offset and size.
func RegionFromOffset(offset Offset, size Size) Region {
	return Region{offset.X, offset.Y, size.Width, size.Height}
}

// ScrollToVisible calculates the smallest offset needed to translate window so
// that region is visible inside it. If top is true, scrolls region to the top.
func ScrollToVisible(window, region Region, top bool) Offset {
	if !top && window.ContainsRegion(region) {
		return NullOffset
	}
	windowLeft, windowTop, windowRight, windowBottom := window.Corners()
	region = region.CropSize(window.Size())
	left, top_, right, bottom := region.Corners()
	deltaX, deltaY := 0, 0

	if !((windowRight > left && left >= windowLeft) && (windowRight > right && right >= windowLeft)) {
		a := left - windowLeft
		b := left - (windowRight - region.Width)
		if abs(a) < abs(b) {
			deltaX = a
		} else {
			deltaX = b
		}
	}

	if top {
		deltaY = top_ - windowTop
	} else if !((windowBottom > top_ && top_ >= windowTop) && (windowBottom > bottom && bottom >= windowTop)) {
		a := top_ - windowTop
		b := top_ - (windowBottom - region.Height)
		if abs(a) < abs(b) {
			deltaY = a
		} else {
			deltaY = b
		}
	}
	return Offset{deltaX, deltaY}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// IsNonZero reports whether the region has non-zero area.
func (r Region) IsNonZero() bool {
	return r.Width*r.Height > 0
}

// Right returns the maximum X value (exclusive).
func (r Region) Right() int { return r.X + r.Width }

// Bottom returns the maximum Y value (exclusive).
func (r Region) Bottom() int { return r.Y + r.Height }

// Area returns the area of the region.
func (r Region) Area() int { return r.Width * r.Height }

// Offset returns the top-left corner as an Offset.
func (r Region) Offset() Offset { return Offset{r.X, r.Y} }

// Center returns the center of the region as (float64, float64).
func (r Region) Center() (float64, float64) {
	return float64(r.X) + float64(r.Width)/2.0, float64(r.Y) + float64(r.Height)/2.0
}

// BottomLeft returns the bottom-left corner.
func (r Region) BottomLeft() Offset { return Offset{r.X, r.Y + r.Height} }

// TopRight returns the top-right corner.
func (r Region) TopRight() Offset { return Offset{r.X + r.Width, r.Y} }

// BottomRight returns the bottom-right corner (exclusive).
func (r Region) BottomRight() Offset { return Offset{r.X + r.Width, r.Y + r.Height} }

// BottomRightInclusive returns the bottom-right corner within the region boundaries.
func (r Region) BottomRightInclusive() Offset {
	return Offset{r.X + r.Width - 1, r.Y + r.Height - 1}
}

// Size returns the size of the region.
func (r Region) Size() Size { return Size{r.Width, r.Height} }

// Corners returns (x1, y1, x2, y2) — top-left and bottom-right corners.
func (r Region) Corners() (int, int, int, int) {
	return r.X, r.Y, r.X + r.Width, r.Y + r.Height
}

// ColumnSpan returns (start, end) column indices (end is exclusive).
func (r Region) ColumnSpan() (int, int) { return r.X, r.X + r.Width }

// LineSpan returns (start, end) row indices (end is exclusive).
func (r Region) LineSpan() (int, int) { return r.Y, r.Y + r.Height }

// ColumnRange returns (start, end) for the column range.
func (r Region) ColumnRange() (int, int) { return r.X, r.X + r.Width }

// LineRange returns (start, end) for the line range.
func (r Region) LineRange() (int, int) { return r.Y, r.Y + r.Height }

// ResetOffset returns a Region of the same size at (0, 0).
func (r Region) ResetOffset() Region { return Region{0, 0, r.Width, r.Height} }

// Add shifts the region by the given offset.
func (r Region) Add(o Offset) Region {
	return Region{r.X + o.X, r.Y + o.Y, r.Width, r.Height}
}

// Sub shifts the region by the negation of the given offset.
func (r Region) Sub(o Offset) Region {
	return Region{r.X - o.X, r.Y - o.Y, r.Width, r.Height}
}

// SpacingBetween returns the Spacing that, if subtracted from r, produces other.
func (r Region) SpacingBetween(other Region) Spacing {
	return Spacing{
		Top:    other.Y - r.Y,
		Right:  r.Right() - other.Right(),
		Bottom: r.Bottom() - other.Bottom(),
		Left:   other.X - r.X,
	}
}

// AtOffset returns a new Region with the same size at the given offset.
func (r Region) AtOffset(o Offset) Region {
	return Region{o.X, o.Y, r.Width, r.Height}
}

// CropSize returns a Region with the same offset but size no larger than s.
func (r Region) CropSize(s Size) Region {
	w := r.Width
	if s.Width < w {
		w = s.Width
	}
	h := r.Height
	if s.Height < h {
		h = s.Height
	}
	return Region{r.X, r.Y, w, h}
}

// Expand increases the size of the region by adding a border of the given size.
func (r Region) Expand(s Size) Region {
	return Region{
		r.X - s.Width,
		r.Y - s.Height,
		r.Width + s.Width*2,
		r.Height + s.Height*2,
	}
}

// Overlaps reports whether other shares any cells with r.
func (r Region) Overlaps(other Region) bool {
	x, y, x2, y2 := r.Corners()
	ox, oy, ox2, oy2 := other.Corners()
	xOverlap := (x2 > ox && ox >= x) || (x2 > ox2 && ox2 > x) || (ox < x && ox2 >= x2)
	yOverlap := (y2 > oy && oy >= y) || (y2 > oy2 && oy2 > y) || (oy < y && oy2 >= y2)
	return xOverlap && yOverlap
}

// Contains reports whether the point (x, y) falls within the region.
func (r Region) Contains(x, y int) bool {
	return r.X+r.Width > x && x >= r.X && r.Y+r.Height > y && y >= r.Y
}

// ContainsPoint reports whether the offset falls within the region.
func (r Region) ContainsPoint(p Offset) bool {
	x1, y1, x2, y2 := r.Corners()
	return x2 > p.X && p.X >= x1 && y2 > p.Y && p.Y >= y1
}

// ContainsRegion reports whether other is entirely within r.
func (r Region) ContainsRegion(other Region) bool {
	x1, y1, x2, y2 := r.Corners()
	ox, oy, ox2, oy2 := other.Corners()
	return ox >= x1 && ox <= x2 && oy >= y1 && oy <= y2 && ox2 >= x1 && ox2 <= x2 && oy2 >= y1 && oy2 <= y2
}

// Translate moves the region by the given offset.
func (r Region) Translate(o Offset) Region {
	return Region{r.X + o.X, r.Y + o.Y, r.Width, r.Height}
}

// Clip clips the region to fit within a bounding box of width x height.
func (r Region) Clip(width, height int) Region {
	x1, y1, x2, y2 := r.Corners()
	return RegionFromCorners(
		ClampInt(x1, 0, width),
		ClampInt(y1, 0, height),
		ClampInt(x2, 0, width),
		ClampInt(y2, 0, height),
	)
}

// Grow returns a new Region expanded by the given Spacing margin.
func (r Region) Grow(s Spacing) Region {
	if s == NullSpacing {
		return r
	}
	w := r.Width + s.Left + s.Right
	if w < 0 {
		w = 0
	}
	h := r.Height + s.Top + s.Bottom
	if h < 0 {
		h = 0
	}
	return Region{r.X - s.Left, r.Y - s.Top, w, h}
}

// Shrink returns a new Region reduced by the given Spacing margin.
func (r Region) Shrink(s Spacing) Region {
	if s == NullSpacing {
		return r
	}
	w := r.Width - (s.Left + s.Right)
	if w < 0 {
		w = 0
	}
	h := r.Height - (s.Top + s.Bottom)
	if h < 0 {
		h = 0
	}
	return Region{r.X + s.Left, r.Y + s.Top, w, h}
}

// Intersection returns the overlapping portion of r and other.
func (r Region) Intersection(other Region) Region {
	x1, y1, w1, h1 := r.X, r.Y, r.Width, r.Height
	cx1, cy1, w2, h2 := other.X, other.Y, other.Width, other.Height
	x2 := x1 + w1
	y2 := y1 + h1
	cx2 := cx1 + w2
	cy2 := cy1 + h2

	rx1 := cx1
	if x1 > cx2 {
		rx1 = cx2
	} else if x1 > cx1 {
		rx1 = x1
	}
	ry1 := cy1
	if y1 > cy2 {
		ry1 = cy2
	} else if y1 > cy1 {
		ry1 = y1
	}
	rx2 := cx2
	if x2 < cx1 {
		rx2 = cx1
	} else if x2 < cx2 {
		rx2 = x2
	}
	ry2 := cy2
	if y2 < cy1 {
		ry2 = cy1
	} else if y2 < cy2 {
		ry2 = y2
	}
	return Region{rx1, ry1, rx2 - rx1, ry2 - ry1}
}

// Union returns the smallest Region that contains both r and other.
func (r Region) Union(other Region) Region {
	x1, y1, x2, y2 := r.Corners()
	ox1, oy1, ox2, oy2 := other.Corners()
	minX := x1
	if ox1 < minX {
		minX = ox1
	}
	minY := y1
	if oy1 < minY {
		minY = oy1
	}
	maxX := x2
	if ox2 > maxX {
		maxX = ox2
	}
	maxY := y2
	if oy2 > maxY {
		maxY = oy2
	}
	return RegionFromCorners(minX, minY, maxX, maxY)
}

// Split divides r into 4 regions at (cutX, cutY).
// Negative cut values are measured from the opposite edge.
func (r Region) Split(cutX, cutY int) (Region, Region, Region, Region) {
	if cutX < 0 {
		cutX = r.Width + cutX
	}
	if cutY < 0 {
		cutY = r.Height + cutY
	}
	return Region{r.X, r.Y, cutX, cutY},
		Region{r.X + cutX, r.Y, r.Width - cutX, cutY},
		Region{r.X, r.Y + cutY, cutX, r.Height - cutY},
		Region{r.X + cutX, r.Y + cutY, r.Width - cutX, r.Height - cutY}
}

// SplitVertical divides r into two regions at the given x offset.
// Negative cut values are measured from the right edge.
func (r Region) SplitVertical(cut int) (Region, Region) {
	if cut < 0 {
		cut = r.Width + cut
	}
	return Region{r.X, r.Y, cut, r.Height},
		Region{r.X + cut, r.Y, r.Width - cut, r.Height}
}

// SplitHorizontal divides r into two regions at the given y offset.
// Negative cut values are measured from the bottom edge.
func (r Region) SplitHorizontal(cut int) (Region, Region) {
	if cut < 0 {
		cut = r.Height + cut
	}
	return Region{r.X, r.Y, r.Width, cut},
		Region{r.X, r.Y + cut, r.Width, r.Height - cut}
}

// TranslateInside translates r so it fits within container.
// xAxis and yAxis control which axes are adjusted.
func (r Region) TranslateInside(container Region, xAxis, yAxis bool) Region {
	x1, y1, w1, h1 := container.X, container.Y, container.Width, container.Height
	x2, y2, w2, h2 := r.X, r.Y, r.Width, r.Height
	nx, ny := x2, y2
	if xAxis {
		nx = clampInsideAxis(x2, x1, w1, w2)
	}
	if yAxis {
		ny = clampInsideAxis(y2, y1, h1, h2)
	}
	return Region{nx, ny, w2, h2}
}

func clampInsideAxis(pos, containerPos, containerSize, selfSize int) int {
	hi := containerPos + containerSize - selfSize
	if pos > hi {
		pos = hi
	}
	if pos < containerPos {
		pos = containerPos
	}
	return pos
}

// Inflect moves r around one or both axes.
// xAxis: +1 moves right, -1 moves left, 0 leaves unchanged.
// yAxis: +1 moves down, -1 moves up, 0 leaves unchanged.
// margin adds additional spacing between the regions (overlapping).
func (r Region) Inflect(xAxis, yAxis int, margin *Spacing) Region {
	inflectMargin := NullSpacing
	if margin != nil {
		inflectMargin = *margin
	}
	x, y, width, height := r.X, r.Y, r.Width, r.Height
	if xAxis != 0 {
		x += (width + inflectMargin.MaxWidth()) * xAxis
	}
	if yAxis != 0 {
		y += (height + inflectMargin.MaxHeight()) * yAxis
	}
	return Region{x, y, width, height}
}

// Constrain constrains r to fit within container using the specified methods per axis.
// constrainX and constrainY may each be "none", "inside", or "inflect".
func (r Region) Constrain(constrainX, constrainY string, margin Spacing, container Region) Region {
	marginRegion := r.Grow(margin)
	result := r

	compareSpan := func(spanStart, spanEnd, containerStart, containerEnd int) int {
		if spanStart >= containerStart && spanEnd <= containerEnd {
			return 0
		}
		if spanStart < containerStart {
			return -1
		}
		return 1
	}

	if constrainX == "inflect" || constrainY == "inflect" {
		var ix, iy int
		if constrainX == "inflect" {
			ix = -compareSpan(marginRegion.X, marginRegion.Right(), container.X, container.Right())
		}
		if constrainY == "inflect" {
			iy = -compareSpan(marginRegion.Y, marginRegion.Bottom(), container.Y, container.Bottom())
		}
		result = result.Inflect(ix, iy, &margin)
	}

	result = result.TranslateInside(
		container.Shrink(margin),
		constrainX != "none",
		constrainY != "none",
	)
	return result
}

// Spacing stores space around a widget (padding/border) as top, right, bottom, left.
type Spacing struct {
	Top    int
	Right  int
	Bottom int
	Left   int
}

// IsNonZero reports whether any spacing value is non-zero.
func (s Spacing) IsNonZero() bool {
	return s != NullSpacing
}

// Width returns the total horizontal spacing (left + right).
func (s Spacing) Width() int { return s.Left + s.Right }

// Height returns the total vertical spacing (top + bottom).
func (s Spacing) Height() int { return s.Top + s.Bottom }

// MaxWidth returns max(left, right) — overlap spacing for the X axis.
func (s Spacing) MaxWidth() int {
	if s.Left > s.Right {
		return s.Left
	}
	return s.Right
}

// MaxHeight returns max(top, bottom) — overlap spacing for the Y axis.
func (s Spacing) MaxHeight() int {
	if s.Top > s.Bottom {
		return s.Top
	}
	return s.Bottom
}

// TopLeft returns (left, top) as a pair.
func (s Spacing) TopLeft() (int, int) { return s.Left, s.Top }

// BottomRight returns (right, bottom) as a pair.
func (s Spacing) BottomRight() (int, int) { return s.Right, s.Bottom }

// Totals returns (left+right, top+bottom) — total horizontal and vertical spacing.
func (s Spacing) Totals() (int, int) { return s.Left + s.Right, s.Top + s.Bottom }

// CSS returns the spacing in CSS shorthand notation.
// One value if all sides equal; two values if top==bottom and right==left; four values otherwise.
func (s Spacing) CSS() string {
	if s.Top == s.Right && s.Right == s.Bottom && s.Bottom == s.Left {
		return itoa(s.Top)
	}
	if s.Top == s.Bottom && s.Right == s.Left {
		return itoa(s.Top) + " " + itoa(s.Right)
	}
	return itoa(s.Top) + " " + itoa(s.Right) + " " + itoa(s.Bottom) + " " + itoa(s.Left)
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := false
	if n < 0 {
		neg = true
		n = -n
	}
	buf := [20]byte{}
	pos := len(buf)
	for n > 0 {
		pos--
		buf[pos] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		pos--
		buf[pos] = '-'
	}
	return string(buf[pos:])
}

// UnpackSpacing creates a Spacing from CSS-style values: 1, 2, or 4 integers.
// 1 value: all sides equal.
// 2 values: (top/bottom, left/right).
// 4 values: (top, right, bottom, left).
func UnpackSpacing(values ...int) (Spacing, error) {
	switch len(values) {
	case 1:
		v := values[0]
		return Spacing{v, v, v, v}, nil
	case 2:
		return Spacing{values[0], values[1], values[0], values[1]}, nil
	case 4:
		return Spacing{values[0], values[1], values[2], values[3]}, nil
	default:
		return NullSpacing, fmt.Errorf("1, 2 or 4 integers required for spacing properties; %d given", len(values))
	}
}

// SpacingVertical returns Spacing with the given amount on top and bottom only.
func SpacingVertical(amount int) Spacing { return Spacing{amount, 0, amount, 0} }

// SpacingHorizontal returns Spacing with the given amount on left and right only.
func SpacingHorizontal(amount int) Spacing { return Spacing{0, amount, 0, amount} }

// SpacingAll returns Spacing with the given amount on all sides.
func SpacingAll(amount int) Spacing { return Spacing{amount, amount, amount, amount} }

// Add returns the element-wise sum of two Spacing values.
func (s Spacing) Add(other Spacing) Spacing {
	return Spacing{s.Top + other.Top, s.Right + other.Right, s.Bottom + other.Bottom, s.Left + other.Left}
}

// Sub returns the element-wise difference of two Spacing values.
func (s Spacing) Sub(other Spacing) Spacing {
	return Spacing{s.Top - other.Top, s.Right - other.Right, s.Bottom - other.Bottom, s.Left - other.Left}
}

// GrowMaximum returns a new Spacing where each side is the max of the two inputs.
func (s Spacing) GrowMaximum(other Spacing) Spacing {
	top := s.Top
	if other.Top > top {
		top = other.Top
	}
	right := s.Right
	if other.Right > right {
		right = other.Right
	}
	bottom := s.Bottom
	if other.Bottom > bottom {
		bottom = other.Bottom
	}
	left := s.Left
	if other.Left > left {
		left = other.Left
	}
	return Spacing{top, right, bottom, left}
}
