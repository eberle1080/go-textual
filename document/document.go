package document

import (
	"strings"
	"unicode/utf8"

	"github.com/eberle1080/go-textual/geometry"
)

// Newline represents the line-ending style of a document.
type Newline string

const (
	// NewlineLF uses Unix-style LF (\n) line endings.
	NewlineLF Newline = "\n"
	// NewlineCRLF uses Windows-style CRLF (\r\n) line endings.
	NewlineCRLF Newline = "\r\n"
	// NewlineCR uses old Mac-style CR (\r) line endings.
	NewlineCR Newline = "\r"
)

// Location represents a row/column position in a document.
// Row and Col are zero-based.
type Location struct {
	Row, Col int
}

// Less reports whether l comes before other in document order.
func (l Location) Less(other Location) bool {
	if l.Row != other.Row {
		return l.Row < other.Row
	}
	return l.Col < other.Col
}

// Equal reports whether l and other are the same position.
func (l Location) Equal(other Location) bool {
	return l.Row == other.Row && l.Col == other.Col
}

// Clamped returns the location with Row and Col floored at zero.
func (l Location) Clamped() Location {
	if l.Row < 0 {
		l.Row = 0
	}
	if l.Col < 0 {
		l.Col = 0
	}
	return l
}

// Selection represents a range within a document. When Start == End the
// selection is a cursor (empty range).
type Selection struct {
	Start Location
	End   Location
}

// Cursor returns a Selection representing a cursor (empty selection) at loc.
func Cursor(loc Location) Selection {
	return Selection{Start: loc, End: loc}
}

// IsEmpty reports whether the selection contains no characters.
func (s Selection) IsEmpty() bool {
	return s.Start.Equal(s.End)
}

// ContainsLine reports whether row y falls within the selection's row range.
func (s Selection) ContainsLine(y int) bool {
	start, end := s.Start, s.End
	if end.Less(start) {
		start, end = end, start
	}
	return y >= start.Row && y <= end.Row
}

// Normalized returns a Selection where Start <= End.
func (s Selection) Normalized() Selection {
	if s.End.Less(s.Start) {
		return Selection{Start: s.End, End: s.Start}
	}
	return s
}

// EditResult is returned by [DocumentBase.ReplaceRange].
type EditResult struct {
	// EndLocation is the position immediately after the inserted text.
	EndLocation Location
	// ReplacedText is the text that was overwritten.
	ReplacedText string
}

// DocumentBase is the minimal interface that document implementations must satisfy.
type DocumentBase interface {
	ReplaceRange(start, end Location, text string) EditResult
	Text() string
	Newline() Newline
	Lines() []string
	GetLine(row int) string
	GetTextRange(start, end Location) string
	GetSize() geometry.Size
	LineCount() int
	Start() Location
	End() Location
}

// Document is the standard in-memory document.
// It stores text as a slice of lines (without line terminators) and a detected
// newline style. All column indices are Unicode codepoint (rune) offsets.
type Document struct {
	lines   []string
	newline Newline
}

// NewDocument creates a Document from text, auto-detecting the newline style.
func NewDocument(text string) *Document {
	nl := detectNewlineStyle(text)
	return &Document{
		lines:   splitLines(text, nl),
		newline: nl,
	}
}

// detectNewlineStyle returns the dominant newline style found in text.
func detectNewlineStyle(text string) Newline {
	if strings.Contains(text, "\r\n") {
		return NewlineCRLF
	}
	if strings.Contains(text, "\r") {
		return NewlineCR
	}
	return NewlineLF
}

// splitLines splits text into lines using the given newline, omitting the
// terminators. An empty string returns a single empty line.
func splitLines(text string, nl Newline) []string {
	if text == "" {
		return []string{""}
	}
	sep := string(nl)
	// Normalise CRLF to LF before splitting for simpler code, then switch back.
	lines := strings.Split(text, sep)
	// If the text ends with a newline, strings.Split leaves an empty trailing
	// element which is the correct "last empty line" behaviour.
	return lines
}

// Text returns the full document text with the original newline style.
func (d *Document) Text() string {
	return strings.Join(d.lines, string(d.newline))
}

// Newline returns the document's line-ending style.
func (d *Document) Newline() Newline { return d.newline }

// Lines returns a copy of the document's line slice.
func (d *Document) Lines() []string {
	result := make([]string, len(d.lines))
	copy(result, d.lines)
	return result
}

// GetLine returns the text of line row (0-based), or "" if out of range.
func (d *Document) GetLine(row int) string {
	if row < 0 || row >= len(d.lines) {
		return ""
	}
	return d.lines[row]
}

// LineCount returns the number of lines in the document.
func (d *Document) LineCount() int { return len(d.lines) }

// Start returns the first location in the document (always 0,0).
func (d *Document) Start() Location { return Location{} }

// End returns the location one past the last character in the document.
func (d *Document) End() Location {
	if len(d.lines) == 0 {
		return Location{}
	}
	lastRow := len(d.lines) - 1
	return Location{Row: lastRow, Col: utf8.RuneCountInString(d.lines[lastRow])}
}

// GetSize returns the document dimensions: Width = max column width,
// Height = line count.
func (d *Document) GetSize() geometry.Size {
	maxCol := 0
	for _, line := range d.lines {
		n := utf8.RuneCountInString(line)
		if n > maxCol {
			maxCol = n
		}
	}
	return geometry.Size{Width: maxCol, Height: len(d.lines)}
}

// GetTextRange returns the text between start and end. If start > end they
// are swapped. Returns "" for an empty or degenerate range.
func (d *Document) GetTextRange(start, end Location) string {
	if start.Equal(end) {
		return ""
	}
	if end.Less(start) {
		start, end = end, start
	}

	if start.Row == end.Row {
		runes := []rune(d.GetLine(start.Row))
		sc := clampCol(start.Col, len(runes))
		ec := clampCol(end.Col, len(runes))
		return string(runes[sc:ec])
	}

	var sb strings.Builder
	// First line from start.Col to end.
	firstRunes := []rune(d.GetLine(start.Row))
	sc := clampCol(start.Col, len(firstRunes))
	sb.WriteString(string(firstRunes[sc:]))
	sb.WriteString(string(d.newline))

	// Middle lines.
	for row := start.Row + 1; row < end.Row; row++ {
		sb.WriteString(d.GetLine(row))
		sb.WriteString(string(d.newline))
	}

	// Last line up to end.Col.
	lastRunes := []rune(d.GetLine(end.Row))
	ec := clampCol(end.Col, len(lastRunes))
	sb.WriteString(string(lastRunes[:ec]))
	return sb.String()
}

// ReplaceRange replaces the text between start and end with text, returning an
// EditResult containing the new end position and the replaced text.
// If end < start they are swapped. The replacement is clamped to the document.
func (d *Document) ReplaceRange(start, end Location, text string) EditResult {
	if end.Less(start) {
		start, end = end, start
	}
	// Clamp to document bounds.
	docEnd := d.End()
	if start.Less(d.Start()) {
		start = d.Start()
	}
	if docEnd.Less(end) {
		end = docEnd
	}

	replaced := d.GetTextRange(start, end)
	insertLines := splitInsertLines(text, d.newline)

	// Build the replacement result.
	var result []string
	result = append(result, d.lines[:start.Row]...)

	if start.Row < len(d.lines) {
		startLineRunes := []rune(d.lines[start.Row])
		sc := clampCol(start.Col, len(startLineRunes))
		prefix := string(startLineRunes[:sc])

		var endLineRunes []rune
		if end.Row < len(d.lines) {
			endLineRunes = []rune(d.lines[end.Row])
		}
		ec := clampCol(end.Col, len(endLineRunes))
		suffix := string(endLineRunes[ec:])

		switch len(insertLines) {
		case 0:
			result = append(result, prefix+suffix)
		case 1:
			result = append(result, prefix+insertLines[0]+suffix)
		default:
			result = append(result, prefix+insertLines[0])
			result = append(result, insertLines[1:len(insertLines)-1]...)
			result = append(result, insertLines[len(insertLines)-1]+suffix)
		}
	} else {
		result = append(result, insertLines...)
	}

	// Append lines after end.Row.
	if end.Row+1 < len(d.lines) {
		result = append(result, d.lines[end.Row+1:]...)
	}

	if len(result) == 0 {
		result = []string{""}
	}

	d.lines = result

	// Compute the end location after insertion.
	var endLoc Location
	switch len(insertLines) {
	case 0:
		endLoc = start
	case 1:
		endLoc = Location{
			Row: start.Row,
			Col: start.Col + utf8.RuneCountInString(insertLines[0]),
		}
	default:
		endLoc = Location{
			Row: start.Row + len(insertLines) - 1,
			Col: utf8.RuneCountInString(insertLines[len(insertLines)-1]),
		}
	}

	return EditResult{
		EndLocation:  endLoc,
		ReplacedText: replaced,
	}
}

// splitInsertLines splits text on the document's newline style. An empty
// string returns a single-element slice containing "".
func splitInsertLines(text string, nl Newline) []string {
	if text == "" {
		return []string{""}
	}
	return strings.Split(text, string(nl))
}

// clampCol clamps col to the range [0, max].
func clampCol(col, max int) int {
	if col < 0 {
		return 0
	}
	if col > max {
		return max
	}
	return col
}

// GetIndexFromLocation converts a Location to a linear character index.
func (d *Document) GetIndexFromLocation(loc Location) int {
	idx := 0
	for row := 0; row < loc.Row && row < len(d.lines); row++ {
		idx += utf8.RuneCountInString(d.lines[row]) + 1 // +1 for newline
	}
	if loc.Row < len(d.lines) {
		line := d.lines[loc.Row]
		col := clampCol(loc.Col, utf8.RuneCountInString(line))
		idx += col
	}
	return idx
}

// GetLocationFromIndex converts a linear character index to a Location.
func (d *Document) GetLocationFromIndex(idx int) Location {
	for row, line := range d.lines {
		lineLen := utf8.RuneCountInString(line) + 1 // +1 for newline
		if idx < lineLen {
			return Location{Row: row, Col: idx}
		}
		idx -= lineLen
	}
	return d.End()
}
