// Package strip provides the Strip type — an immutable horizontal line of
// Rich segments — and associated rendering utilities.
//
// A Strip is the fundamental unit used by the Widget rendering pipeline.
// Widgets render themselves as slices of Strips (one per screen row), and
// the compositor assembles these into a complete frame.
//
// # Strip
//
// [Strip] wraps a [github.com/eberle1080/go-rich.Segments] slice and provides
// operations for cropping, padding, styling, aligning, and joining strips.
// Operations return new Strips rather than mutating the receiver; common
// results are memoized using [github.com/eberle1080/go-textual/internal/cache].
//
// # Helpers
//
// [LinePad] pads a raw Segments line with blank segments on the left/right.
// [GetLineLength] computes the total cell width of a Segments line.
// [IndexToCellPosition] converts a character index to a cell position for
// multi-byte character handling.
package strip
