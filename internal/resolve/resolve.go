// Package resolve converts CSS scalar dimensions to concrete cell counts and
// box models. It is a port of textual/_resolve.py and uses math/big.Rat for
// fractional precision, matching the css.Scalar.Resolve contract.
package resolve

import (
	"math/big"

	"github.com/eberle1080/go-textual/css"
	"github.com/eberle1080/go-textual/geometry"
)

// Slot is the resolved offset and length for one dimension slot.
type Slot struct {
	Offset int
	Length int
}

// BoxModelable is the subset of a widget that the box-model resolver needs.
type BoxModelable interface {
	// Styles returns the merged CSS render styles for the widget.
	Styles() *css.RenderStyles
	// Display returns false if the widget is hidden (display:none).
	Display() bool
}

// ResolveFractionUnit calculates the value of 1fr given the total available
// space, the gutter size, and the full (non-compacted) list of widgets that
// correspond positionally to dimensions. It iterates dimensions and widgets
// in lockstep so that non-fraction entries correctly subtract from the
// remaining space even when they precede fractional ones.
//
// Returns nil if there are no visible fractional dimensions.
func ResolveFractionUnit(
	dimensions []css.Scalar,
	total, gutter int,
	size, viewport geometry.Size,
	widgets []BoxModelable,
) *big.Rat {
	// First pass: check whether any visible widget has a fractional dimension
	// and count visible widgets for gutter calculation.
	hasFr := false
	visibleCount := 0
	for i, dim := range dimensions {
		if i >= len(widgets) {
			break
		}
		if !widgets[i].Display() {
			continue
		}
		visibleCount++
		if dim.IsFraction() {
			hasFr = true
		}
	}
	if !hasFr {
		return nil
	}

	// Second pass: sum non-fractional space and total fraction units.
	remaining := new(big.Rat).SetInt64(int64(total))
	totalFractions := new(big.Rat)

	for i, dim := range dimensions {
		if i >= len(widgets) {
			break
		}
		if !widgets[i].Display() {
			continue
		}
		if dim.IsFraction() {
			totalFractions.Add(totalFractions, new(big.Rat).SetFloat64(dim.Value))
		} else if !dim.IsAuto() {
			resolved, err := dim.Resolve(size, viewport, nil)
			if err == nil {
				remaining.Sub(remaining, resolved)
			}
		}
	}

	// Account for gutters between visible widgets.
	if gutter > 0 && visibleCount > 1 {
		g := new(big.Rat).SetInt64(int64(gutter * (visibleCount - 1)))
		remaining.Sub(remaining, g)
	}

	if remaining.Sign() <= 0 || totalFractions.Sign() == 0 {
		return new(big.Rat)
	}

	return new(big.Rat).Quo(remaining, totalFractions)
}

// ResolveBoxModel resolves a single CSS scalar into a rational cell count,
// handling auto, fr, %, and explicit unit types.
//
// fractionUnit is the value of 1fr (pass nil to treat fr as 1).
func ResolveBoxModel(
	dim css.Scalar,
	fractionUnit *big.Rat,
	size, viewport geometry.Size,
) (*big.Rat, error) {
	return dim.Resolve(size, viewport, fractionUnit)
}

// Resolve converts a list of CSS scalar dimensions into concrete Slot values
// (offset, length) for layout. gutter is added between visible slots.
//
// This is the primary entry point used by layout engines such as VerticalLayout
// and HorizontalLayout.
func Resolve(
	dimensions []css.Scalar,
	total, gutter int,
	size, viewport geometry.Size,
	widgets []BoxModelable,
) []Slot {
	n := len(dimensions)
	if n == 0 {
		return nil
	}

	// Pass the full widgets slice so ResolveFractionUnit iterates in lockstep.
	fractionUnit := ResolveFractionUnit(dimensions, total, gutter, size, viewport, widgets)

	// Resolve each dimension to a rational, then round to int.
	lengths := make([]int, n)
	for i, dim := range dimensions {
		if i >= len(widgets) || !widgets[i].Display() {
			lengths[i] = 0
			continue
		}
		rat, err := dim.Resolve(size, viewport, fractionUnit)
		if err != nil {
			lengths[i] = 0
			continue
		}
		f, _ := rat.Float64()
		if f < 0 {
			f = 0
		}
		lengths[i] = int(f + 0.5) // round half-up
	}

	// Compute offsets with gutter.
	slots := make([]Slot, n)
	offset := 0
	for i, length := range lengths {
		if i < len(widgets) && !widgets[i].Display() {
			slots[i] = Slot{Offset: offset, Length: 0}
			continue
		}
		slots[i] = Slot{Offset: offset, Length: length}
		if length > 0 {
			offset += length + gutter
		}
	}

	return slots
}
