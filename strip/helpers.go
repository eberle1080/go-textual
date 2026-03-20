package strip

import (
	"strings"
	"unicode/utf8"

	rich "github.com/eberle1080/go-rich"
)

// GetLineLength returns the total number of Unicode code points (runes) across
// all segments. For terminal rendering, each code point is treated as one cell
// (double-width characters are not handled here).
func GetLineLength(segments rich.Segments) int {
	total := 0
	for _, seg := range segments {
		total += utf8.RuneCountInString(seg.Text)
	}
	return total
}

// IndexToCellPosition converts a character (rune) index within the plain text
// of segments to the corresponding cell (visual column) position. Returns the
// total cell length if index exceeds the segment length.
func IndexToCellPosition(segments rich.Segments, index int) int {
	pos := 0
	remaining := index
	for _, seg := range segments {
		runes := []rune(seg.Text)
		if remaining <= len(runes) {
			pos += remaining
			return pos
		}
		pos += len(runes)
		remaining -= len(runes)
	}
	return pos
}

// LinePad adds left and right blank padding (filled with spaces in the given
// style) to a Segments line. Returns the padded Segments.
func LinePad(segments rich.Segments, left, right int, style rich.Style) rich.Segments {
	var result rich.Segments
	if left > 0 {
		result = append(result, rich.Segment{
			Text:  strings.Repeat(" ", left),
			Style: style,
		})
	}
	result = append(result, segments...)
	if right > 0 {
		result = append(result, rich.Segment{
			Text:  strings.Repeat(" ", right),
			Style: style,
		})
	}
	return result
}
