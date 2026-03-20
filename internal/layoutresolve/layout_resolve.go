// Package layoutresolve divides a fixed total space among a set of edges,
// each of which may specify an explicit size, a fractional share, or a
// minimum size. It is a port of textual/_layout_resolve.py.
package layoutresolve

// Edge describes one slot in a linear layout (column or row).
type Edge interface {
	// Size returns a pointer to an explicit cell count, or nil if not set.
	Size() *int
	// Fraction returns the fractional weight (≥0). Zero means not fractional.
	Fraction() int
	// MinSize returns the minimum number of cells this edge must occupy.
	MinSize() int
}

// Resolve divides total cells among edges.
//
// Edges with an explicit Size() receive that many cells (clamped to their
// MinSize). The remaining space is divided among fractional edges in
// proportion to their Fraction() weight. Edges with no size and fraction=0
// receive their MinSize.
//
// The returned slice has the same length as edges.
func Resolve(total int, edges []Edge) []int {
	n := len(edges)
	result := make([]int, n)

	remaining := total
	totalFractions := 0

	// First pass: allocate explicit sizes and count fractions.
	for i, e := range edges {
		if s := e.Size(); s != nil {
			sz := *s
			if sz < e.MinSize() {
				sz = e.MinSize()
			}
			result[i] = sz
			remaining -= sz
		} else if e.Fraction() > 0 {
			totalFractions += e.Fraction()
		} else {
			// No size, no fraction → use MinSize.
			result[i] = e.MinSize()
			remaining -= e.MinSize()
		}
	}

	if remaining < 0 {
		remaining = 0
	}

	// Second pass: allocate fractional edges.
	if totalFractions > 0 {
		allocated := 0
		fractionEdges := make([]int, 0, n) // indices of fraction edges
		for i, e := range edges {
			if e.Size() == nil && e.Fraction() > 0 {
				fractionEdges = append(fractionEdges, i)
			}
		}
		for k, i := range fractionEdges {
			e := edges[i]
			var share int
			if k == len(fractionEdges)-1 {
				// Last fraction gets the rest to avoid rounding gaps.
				share = remaining - allocated
			} else {
				share = (remaining * e.Fraction()) / totalFractions
			}
			if share < e.MinSize() {
				share = e.MinSize()
			}
			result[i] = share
			allocated += share
		}
	}

	return result
}
