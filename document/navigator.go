package document

import (
	"sort"
	"unicode/utf8"

	"github.com/eberle1080/go-textual/geometry"
)

// DocumentNavigator provides cursor movement within a [WrappedDocument].
// It caches the last X offset so vertical navigation preserves the column.
type DocumentNavigator struct {
	doc         *WrappedDocument
	LastXOffset int
}

// NewDocumentNavigator creates a DocumentNavigator for the given wrapped document.
func NewDocumentNavigator(doc *WrappedDocument) *DocumentNavigator {
	return &DocumentNavigator{doc: doc}
}

// lineRuneCount returns the rune count of document line row.
func (n *DocumentNavigator) lineRuneCount(row int) int {
	return utf8.RuneCountInString(n.doc.doc.GetLine(row))
}

// GetLocationLeft returns the location one rune to the left of loc.
// Wraps to the end of the previous line when at column 0.
func (n *DocumentNavigator) GetLocationLeft(loc Location) Location {
	if loc.Col > 0 {
		return Location{Row: loc.Row, Col: loc.Col - 1}
	}
	if loc.Row > 0 {
		prevRow := loc.Row - 1
		return Location{Row: prevRow, Col: n.lineRuneCount(prevRow)}
	}
	return loc
}

// GetLocationRight returns the location one rune to the right of loc.
// Wraps to the start of the next line when at end of line.
func (n *DocumentNavigator) GetLocationRight(loc Location) Location {
	lineLen := n.lineRuneCount(loc.Row)
	if loc.Col < lineLen {
		return Location{Row: loc.Row, Col: loc.Col + 1}
	}
	if loc.Row+1 < n.doc.doc.LineCount() {
		return Location{Row: loc.Row + 1, Col: 0}
	}
	return loc
}

// GetLocationAbove returns the location one visual row above loc.
func (n *DocumentNavigator) GetLocationAbove(loc Location) Location {
	offset := n.doc.LocationToOffset(loc)
	if offset.Y == 0 {
		return loc
	}
	targetX := n.LastXOffset
	if offset.X > targetX {
		targetX = offset.X
		n.LastXOffset = targetX
	}
	return n.doc.OffsetToLocation(geometry.Offset{X: targetX, Y: offset.Y - 1})
}

// GetLocationBelow returns the location one visual row below loc.
func (n *DocumentNavigator) GetLocationBelow(loc Location) Location {
	offset := n.doc.LocationToOffset(loc)
	totalLines := n.doc.TotalVisualLines()
	if offset.Y >= totalLines-1 {
		return loc
	}
	targetX := n.LastXOffset
	if offset.X > targetX {
		targetX = offset.X
		n.LastXOffset = targetX
	}
	return n.doc.OffsetToLocation(geometry.Offset{X: targetX, Y: offset.Y + 1})
}

// GetLocationHome returns the start of the visual line containing loc.
func (n *DocumentNavigator) GetLocationHome(loc Location) Location {
	offset := n.doc.LocationToOffset(loc)
	n.LastXOffset = 0
	return n.doc.OffsetToLocation(geometry.Offset{X: 0, Y: offset.Y})
}

// GetLocationEnd returns the end of the visual line containing loc.
func (n *DocumentNavigator) GetLocationEnd(loc Location) Location {
	offset := n.doc.LocationToOffset(loc)
	row, sectionIndex := n.doc.VisualLineToDocumentLine(offset.Y)
	if row < 0 {
		return loc
	}
	sections := n.doc.GetSections(row)
	if sectionIndex >= len(sections) {
		return loc
	}
	sectionLen := len([]rune(sections[sectionIndex]))
	offsets := n.doc.GetOffsets(row)
	var sectionStart int
	if sectionIndex > 0 && sectionIndex-1 < len(offsets) {
		sectionStart = offsets[sectionIndex-1]
	}
	col := sectionStart + sectionLen
	n.LastXOffset = sectionLen
	return Location{Row: row, Col: col}
}

// GetLocationAtYOffset returns the location at the given visual row offset
// from loc, preserving the X offset as closely as possible.
func (n *DocumentNavigator) GetLocationAtYOffset(loc Location, yOffset int) Location {
	offset := n.doc.LocationToOffset(loc)
	targetY := offset.Y + yOffset
	totalLines := n.doc.TotalVisualLines()
	if targetY < 0 {
		targetY = 0
	}
	if targetY >= totalLines {
		targetY = totalLines - 1
	}
	return n.doc.OffsetToLocation(geometry.Offset{X: n.LastXOffset, Y: targetY})
}

// IsStartOfDocumentLine reports whether loc is at column 0.
func (n *DocumentNavigator) IsStartOfDocumentLine(loc Location) bool {
	return loc.Col == 0
}

// IsEndOfDocumentLine reports whether loc is at the end of its document line.
func (n *DocumentNavigator) IsEndOfDocumentLine(loc Location) bool {
	return loc.Col >= n.lineRuneCount(loc.Row)
}

// IsFirstWrappedLine reports whether loc is in the first visual section of
// its document line.
func (n *DocumentNavigator) IsFirstWrappedLine(loc Location) bool {
	offsets := n.doc.GetOffsets(loc.Row)
	if len(offsets) == 0 {
		return true
	}
	return loc.Col < offsets[0]
}

// IsLastWrappedLine reports whether loc is in the last visual section of its
// document line.
func (n *DocumentNavigator) IsLastWrappedLine(loc Location) bool {
	offsets := n.doc.GetOffsets(loc.Row)
	if len(offsets) == 0 {
		return true
	}
	return loc.Col >= offsets[len(offsets)-1]
}

// ClampReachable clamps loc to a valid position within the document.
func (n *DocumentNavigator) ClampReachable(loc Location) Location {
	lineCount := n.doc.doc.LineCount()
	if loc.Row < 0 {
		loc.Row = 0
	}
	if loc.Row >= lineCount {
		loc.Row = lineCount - 1
	}
	lineLen := n.lineRuneCount(loc.Row)
	if loc.Col < 0 {
		loc.Col = 0
	}
	if loc.Col > lineLen {
		loc.Col = lineLen
	}
	return loc
}

// SectionIndexForColumn returns the visual section index for the given column
// within document line row.
func (n *DocumentNavigator) SectionIndexForColumn(row, col int) int {
	offsets := n.doc.GetOffsets(row)
	return sort.SearchInts(offsets, col)
}
