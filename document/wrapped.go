package document

import (
	"sort"
	"strings"

	"github.com/eberle1080/go-textual/geometry"
)

// WrappedDocument wraps a [DocumentBase] at a given column width, providing a
// visual line model where one document line may span several visual sections.
type WrappedDocument struct {
	doc   DocumentBase
	width int

	// offsets[row] is the sorted list of column offsets where wrapping occurs
	// on document line row. An empty slice means the line fits on one section.
	offsets map[int][]int
}

// NewWrappedDocument creates a WrappedDocument for doc at the given width.
// If width <= 0, no wrapping is performed.
func NewWrappedDocument(doc DocumentBase, width int) *WrappedDocument {
	return &WrappedDocument{
		doc:     doc,
		width:   width,
		offsets: make(map[int][]int),
	}
}

// Wrap recalculates all wrap offsets for the current document contents.
func (w *WrappedDocument) Wrap() {
	w.offsets = make(map[int][]int, w.doc.LineCount())
	for row := 0; row < w.doc.LineCount(); row++ {
		w.offsets[row] = w.calculateOffsets(w.doc.GetLine(row))
	}
}

// SetWidth updates the wrap width and recalculates all offsets.
func (w *WrappedDocument) SetWidth(width int) {
	w.width = width
	w.Wrap()
}

// GetOffsets returns the wrap offsets for document line row.
// Each offset is the column at which a new visual section begins.
func (w *WrappedDocument) GetOffsets(row int) []int {
	if offs, ok := w.offsets[row]; ok {
		return offs
	}
	offs := w.calculateOffsets(w.doc.GetLine(row))
	w.offsets[row] = offs
	return offs
}

// GetSections returns the visual sections (sub-strings) for document line row.
func (w *WrappedDocument) GetSections(row int) []string {
	line := w.doc.GetLine(row)
	runes := []rune(line)
	offsets := w.GetOffsets(row)
	if len(offsets) == 0 {
		return []string{line}
	}
	var sections []string
	prev := 0
	for _, off := range offsets {
		end := off
		if end > len(runes) {
			end = len(runes)
		}
		sections = append(sections, string(runes[prev:end]))
		prev = end
	}
	sections = append(sections, string(runes[prev:]))
	return sections
}

// GetTargetDocumentColumn converts a visual offset within a section back to
// a document column.
func (w *WrappedDocument) GetTargetDocumentColumn(row, visualOffset, sectionIndex int) int {
	offsets := w.GetOffsets(row)
	var sectionStart int
	if sectionIndex > 0 && sectionIndex-1 < len(offsets) {
		sectionStart = offsets[sectionIndex-1]
	}
	return sectionStart + visualOffset
}

// LocationToOffset converts a document Location to a visual geometry.Offset
// (X = visual column within section, Y = visual row across all wrapped lines).
func (w *WrappedDocument) LocationToOffset(loc Location) geometry.Offset {
	visualRow := 0
	for row := 0; row < loc.Row && row < w.doc.LineCount(); row++ {
		visualRow += w.SectionCount(row)
	}

	offsets := w.GetOffsets(loc.Row)
	sectionIndex := sort.SearchInts(offsets, loc.Col)
	visualRow += sectionIndex

	var sectionStart int
	if sectionIndex > 0 && sectionIndex-1 < len(offsets) {
		sectionStart = offsets[sectionIndex-1]
	}
	visualCol := loc.Col - sectionStart

	return geometry.Offset{X: visualCol, Y: visualRow}
}

// OffsetToLocation converts a visual geometry.Offset back to a document Location.
func (w *WrappedDocument) OffsetToLocation(offset geometry.Offset) Location {
	visualRow := 0
	for row := 0; row < w.doc.LineCount(); row++ {
		sections := w.SectionCount(row)
		if visualRow+sections > offset.Y {
			sectionIndex := offset.Y - visualRow
			return w.sectionToLocation(row, sectionIndex, offset.X)
		}
		visualRow += sections
	}
	return w.doc.End()
}

// TotalVisualLines returns the total number of visual lines.
func (w *WrappedDocument) TotalVisualLines() int {
	total := 0
	for row := 0; row < w.doc.LineCount(); row++ {
		total += w.SectionCount(row)
	}
	return total
}

// SectionCount returns the number of visual sections (wrapped sub-lines) for
// document line row.
func (w *WrappedDocument) SectionCount(row int) int {
	return len(w.GetOffsets(row)) + 1
}

// calculateOffsets returns the wrap break column offsets for line.
func (w *WrappedDocument) calculateOffsets(line string) []int {
	if w.width <= 0 {
		return nil
	}
	runes := []rune(line)
	if len(runes) <= w.width {
		return nil
	}
	var offsets []int
	pos := 0
	for pos+w.width < len(runes) {
		pos += w.width
		offsets = append(offsets, pos)
	}
	return offsets
}

// sectionToLocation converts (row, sectionIndex, visualCol) to a Location.
func (w *WrappedDocument) sectionToLocation(row, sectionIndex, visualCol int) Location {
	offsets := w.GetOffsets(row)
	var sectionStart int
	if sectionIndex > 0 && sectionIndex-1 < len(offsets) {
		sectionStart = offsets[sectionIndex-1]
	}
	line := w.doc.GetLine(row)
	runes := []rune(line)
	col := sectionStart + visualCol
	if col > len(runes) {
		col = len(runes)
	}
	return Location{Row: row, Col: col}
}

// InvalidateRow clears the cached wrap offsets for the given row, causing
// them to be recalculated on next access.
func (w *WrappedDocument) InvalidateRow(row int) {
	delete(w.offsets, row)
}

// GetVisualLine returns the content of visual line y, applying the horizontal
// offset. If there is no content at y, an empty string is returned.
func (w *WrappedDocument) GetVisualLine(y, xOffset int) string {
	row, sectionIndex := w.VisualLineToDocumentLine(y)
	if row < 0 {
		return ""
	}
	sections := w.GetSections(row)
	if sectionIndex >= len(sections) {
		return ""
	}
	section := sections[sectionIndex]
	runes := []rune(section)
	if xOffset >= len(runes) {
		return ""
	}
	return strings.TrimRight(string(runes[xOffset:]), "")
}

// VisualLineToDocumentLine converts a visual line index y to (documentRow,
// sectionIndex). Returns (-1, 0) if y is out of range.
func (w *WrappedDocument) VisualLineToDocumentLine(y int) (row, sectionIndex int) {
	current := 0
	for row = 0; row < w.doc.LineCount(); row++ {
		count := w.SectionCount(row)
		if current+count > y {
			return row, y - current
		}
		current += count
	}
	return -1, 0
}

// GetTextAtVisualLine returns the text content of visual line y, split at the
// horizontal offset for rendering.
func (w *WrappedDocument) GetTextAtVisualLine(y int) string {
	row, sectionIndex := w.VisualLineToDocumentLine(y)
	if row < 0 {
		return ""
	}
	sections := w.GetSections(row)
	if sectionIndex >= len(sections) {
		return ""
	}
	return sections[sectionIndex]
}
