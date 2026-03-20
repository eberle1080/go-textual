package strip

import (
	"strings"

	rich "github.com/eberle1080/go-rich"

	"github.com/eberle1080/go-textual/css"
	"github.com/eberle1080/go-textual/internal/cache"
)

// cropKey is the cache key for Crop operations.
type cropKey struct{ start, end int }

// styleKey is the cache key for ApplyStyle operations.
type styleKey = rich.Style

// Strip is an immutable horizontal line of Rich segments.
// All mutating operations (Crop, ApplyStyle, etc.) return a new Strip.
//
// Frequently used operations are memoized in small FIFO caches on the Strip.
// Because Strip is a value type the caches are pointers; nil caches are
// lazily initialised on first use.
type Strip struct {
	segments   rich.Segments
	cellLength int     // -1 means not yet computed
	cropCache  *cache.FIFOCache[cropKey, Strip]
	styleCache *cache.FIFOCache[styleKey, Strip]
}

// New creates a Strip from a Segments slice.
// If cellLength is provided it is used as the pre-computed cell length;
// otherwise it is computed on first access.
func New(segments rich.Segments, cellLength ...int) Strip {
	s := Strip{
		segments:   segments,
		cellLength: -1,
	}
	if len(cellLength) > 0 {
		s.cellLength = cellLength[0]
	}
	return s
}

// Blank creates a Strip of pure whitespace of the given cell length.
func Blank(cellLength int, style rich.Style) Strip {
	seg := rich.Segment{Text: strings.Repeat(" ", cellLength), Style: style}
	return Strip{
		segments:   rich.Segments{seg},
		cellLength: cellLength,
	}
}

// FromLines converts a slice of Segments (one per line) into a slice of Strips.
func FromLines(lines []rich.Segments, cellLength ...int) []Strip {
	result := make([]Strip, len(lines))
	for i, segs := range lines {
		if len(cellLength) > 0 {
			result[i] = New(segs, cellLength[0])
		} else {
			result[i] = New(segs)
		}
	}
	return result
}

// CellLength returns the total cell width of the strip, computing it lazily.
func (s Strip) CellLength() int {
	if s.cellLength >= 0 {
		return s.cellLength
	}
	s.cellLength = GetLineLength(s.segments)
	return s.cellLength
}

// Text returns the plain text content of the strip (no styling).
func (s Strip) Text() string {
	var b strings.Builder
	for _, seg := range s.segments {
		b.WriteString(seg.Text)
	}
	return b.String()
}

// Segments returns the underlying segment slice.
func (s Strip) Segments() rich.Segments { return s.segments }

// Len returns the number of segments in the strip.
func (s Strip) Len() int { return len(s.segments) }

// Equal reports whether two strips have identical text and style content.
func (s Strip) Equal(other Strip) bool {
	if len(s.segments) != len(other.segments) {
		return false
	}
	for i, seg := range s.segments {
		o := other.segments[i]
		if seg.Text != o.Text || seg.Style != o.Style {
			return false
		}
	}
	return true
}

// Crop returns a new Strip containing only the cells in [start, end).
// Both start and end are clamped to valid cell positions.
func (s Strip) Crop(start, end int) Strip {
	length := s.CellLength()
	if start < 0 {
		start = 0
	}
	if end > length {
		end = length
	}
	if start >= end {
		return New(nil, 0)
	}
	if start == 0 && end == length {
		return s
	}

	// Check cache.
	key := cropKey{start, end}
	if s.cropCache != nil {
		if cached, ok := s.cropCache.Get(key); ok {
			return cached
		}
	}

	result := cropSegments(s.segments, start, end)
	strip := New(result, end-start)

	if s.cropCache == nil {
		s.cropCache = cache.NewFIFO[cropKey, Strip](8)
	}
	s.cropCache.Set(key, strip)
	return strip
}

// CropExtend crops to [start, end) then pads with spaces in style so the
// result has exactly (end-start) cells.
func (s Strip) CropExtend(start, end int, style rich.Style) Strip {
	cropped := s.Crop(start, end)
	return cropped.ExtendCellLength(end-start, style)
}

// Divide splits the strip at each cell position in cuts.
// cuts must be in ascending order. Returns len(cuts)+1 strips.
func (s Strip) Divide(cuts []int) []Strip {
	if len(cuts) == 0 {
		return []Strip{s}
	}
	result := make([]Strip, 0, len(cuts)+1)
	prev := 0
	for _, cut := range cuts {
		result = append(result, s.Crop(prev, cut))
		prev = cut
	}
	result = append(result, s.Crop(prev, s.CellLength()))
	return result
}

// ApplyStyle returns a new Strip with the given style applied to all segments.
func (s Strip) ApplyStyle(style rich.Style) Strip {
	if s.styleCache != nil {
		if cached, ok := s.styleCache.Get(style); ok {
			return cached
		}
	}

	result := make(rich.Segments, len(s.segments))
	for i, seg := range s.segments {
		result[i] = rich.Segment{Text: seg.Text, Style: style}
	}
	strip := New(result, s.cellLength)

	if s.styleCache == nil {
		s.styleCache = cache.NewFIFO[styleKey, Strip](4)
	}
	s.styleCache.Set(style, strip)
	return strip
}

// AdjustCellLength pads or truncates the strip to exactly cellLength cells.
// Padding uses spaces in the given style.
func (s Strip) AdjustCellLength(cellLength int, style rich.Style) Strip {
	current := s.CellLength()
	if current == cellLength {
		return s
	}
	if current > cellLength {
		return s.Crop(0, cellLength)
	}
	return s.ExtendCellLength(cellLength, style)
}

// ExtendCellLength extends the strip to at least cellLength cells by appending
// spaces in style. Never truncates.
func (s Strip) ExtendCellLength(cellLength int, style rich.Style) Strip {
	current := s.CellLength()
	if current >= cellLength {
		return s
	}
	pad := cellLength - current
	segs := make(rich.Segments, len(s.segments)+1)
	copy(segs, s.segments)
	segs[len(segs)-1] = rich.Segment{
		Text:  strings.Repeat(" ", pad),
		Style: style,
	}
	return New(segs, cellLength)
}

// CropPad crops to [left, left+cellLength) then pads right to cellLength.
func (s Strip) CropPad(cellLength, left, right int, style rich.Style) Strip {
	cropped := s.Crop(left, left+cellLength)
	return cropped.AdjustCellLength(cellLength, style)
}

// Splice returns a new Strip with columns [x, x+src.CellLength()) replaced by
// src. Columns outside that range are taken from s. This is equivalent to
// s.Crop(0,x) + src + s.Crop(x+src.CellLength(), s.CellLength()).
func (s Strip) Splice(x int, src Strip) Strip {
	total := s.CellLength()
	srcLen := src.CellLength()
	if x <= 0 && srcLen >= total {
		return src
	}
	var segs rich.Segments
	if x > 0 {
		left := s.Crop(0, x)
		segs = append(segs, left.Segments()...)
	}
	segs = append(segs, src.Segments()...)
	end := x + srcLen
	if end < total {
		right := s.Crop(end, total)
		segs = append(segs, right.Segments()...)
	}
	return New(segs)
}

// Simplify merges adjacent segments that share the same style.
func (s Strip) Simplify() Strip {
	if len(s.segments) <= 1 {
		return s
	}
	result := make(rich.Segments, 0, len(s.segments))
	result = append(result, s.segments[0])
	for _, seg := range s.segments[1:] {
		last := &result[len(result)-1]
		if last.Style == seg.Style {
			last.Text += seg.Text
		} else {
			result = append(result, seg)
		}
	}
	return New(result, s.cellLength)
}

// Join concatenates a slice of Strips into a single Strip.
func Join(strips []Strip) Strip {
	if len(strips) == 0 {
		return New(nil, 0)
	}
	total := 0
	var segs rich.Segments
	for _, s := range strips {
		segs = append(segs, s.segments...)
		if s.cellLength >= 0 {
			total += s.cellLength
		} else {
			total += GetLineLength(s.segments)
		}
	}
	return New(segs, total)
}

// Align aligns a slice of strips to a target width and height using the given
// horizontal and vertical alignment. Missing rows are filled with blank strips.
func Align(
	strips []Strip,
	style rich.Style,
	width int,
	height *int,
	horizontal css.AlignHorizontal,
	vertical css.AlignVertical,
) []Strip {
	targetHeight := len(strips)
	if height != nil {
		targetHeight = *height
	}

	// Pad/truncate to targetHeight.
	result := make([]Strip, targetHeight)
	copyCount := len(strips)
	if copyCount > targetHeight {
		copyCount = targetHeight
	}

	for i := 0; i < copyCount; i++ {
		result[i] = strips[i]
	}
	for i := copyCount; i < targetHeight; i++ {
		result[i] = Blank(width, style)
	}

	// Apply vertical alignment for unused rows.
	switch vertical {
	case "bottom":
		// Shift strips down.
		shift := targetHeight - copyCount
		if shift > 0 {
			for i := targetHeight - 1; i >= 0; i-- {
				if i-shift >= 0 {
					result[i] = result[i-shift]
				} else {
					result[i] = Blank(width, style)
				}
			}
		}
	case "middle":
		shift := (targetHeight - copyCount) / 2
		if shift > 0 {
			for i := targetHeight - 1; i >= 0; i-- {
				if i-shift >= 0 && i-shift < copyCount {
					result[i] = result[i-shift]
				} else {
					result[i] = Blank(width, style)
				}
			}
		}
	}

	// Apply horizontal alignment to each strip.
	for i, s := range result {
		result[i] = s.TextAlign(width, horizontal)
	}
	return result
}

// TextAlign adjusts a strip's horizontal alignment within the given width.
func (s Strip) TextAlign(width int, align css.AlignHorizontal) Strip {
	current := s.CellLength()
	if current >= width {
		return s.Crop(0, width)
	}
	pad := width - current
	style := rich.Style{}
	switch align {
	case "right":
		return s.CropPad(width, 0, pad, style)
	case "center":
		left := pad / 2
		right := pad - left
		segs := LinePad(s.segments, left, right, style)
		return New(segs, width)
	default: // "left"
		return s.ExtendCellLength(width, style)
	}
}

// Render converts the strip to an ANSI escape sequence string.
func (s Strip) Render(colorMode rich.ColorMode) string {
	return s.segments.ToANSI(colorMode)
}

// --- internal helpers ---

// cropSegments crops a Segments slice to [start, end) cell positions.
func cropSegments(segs rich.Segments, start, end int) rich.Segments {
	if len(segs) == 0 || start >= end {
		return nil
	}
	var result rich.Segments
	cellPos := 0
	for _, seg := range segs {
		runes := []rune(seg.Text)
		segEnd := cellPos + len(runes)
		if segEnd <= start {
			cellPos = segEnd
			continue
		}
		if cellPos >= end {
			break
		}
		// This segment overlaps [start, end).
		segStart := start - cellPos
		if segStart < 0 {
			segStart = 0
		}
		segCropEnd := end - cellPos
		if segCropEnd > len(runes) {
			segCropEnd = len(runes)
		}
		text := string(runes[segStart:segCropEnd])
		if text != "" {
			result = append(result, rich.Segment{Text: text, Style: seg.Style})
		}
		cellPos = segEnd
	}
	return result
}

