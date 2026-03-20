package widgets

import (
	"context"
	"fmt"

	rich "github.com/eberle1080/go-rich"

	"github.com/eberle1080/go-textual/css"
	"github.com/eberle1080/go-textual/dom"
	"github.com/eberle1080/go-textual/geometry"
	"github.com/eberle1080/go-textual/keys"
	"github.com/eberle1080/go-textual/msg"
	"github.com/eberle1080/go-textual/strip"
	"github.com/eberle1080/go-textual/widget"
)

// DataTableSelectedMsg is sent when the user activates a row with Enter.
type DataTableSelectedMsg struct {
	msg.BaseMsg
	Row  int
	Cols []string
}

// TableColumn defines a column in a DataTable.
type TableColumn struct {
	Header string
	Width  int // 0 = auto (share remaining space equally)
}

// DataTable displays tabular data with column headers and keyboard navigation.
type DataTable struct {
	widget.BaseWidget
	columns []TableColumn
	rows    [][]string
	cursor  int
	scrollY int
}

// NewDataTable creates a DataTable with the given columns.
func NewDataTable(columns ...TableColumn) *DataTable {
	dt := &DataTable{
		BaseWidget: *widget.NewBaseWidget(
			widget.WithDOMOptions(dom.WithCSSTypeName("DataTable", "Widget")),
			widget.WithCanFocus(true),
		),
		columns: columns,
	}
	return dt
}

// AddRow appends a row. Values beyond the column count are ignored; missing
// values are rendered as empty strings.
func (dt *DataTable) AddRow(values ...string) {
	row := make([]string, len(dt.columns))
	copy(row, values)
	dt.rows = append(dt.rows, row)
	dt.MarkDirty()
}

// ClearRows removes all data rows.
func (dt *DataTable) ClearRows() {
	dt.rows = dt.rows[:0]
	dt.cursor = 0
	dt.scrollY = 0
	dt.MarkDirty()
}

// Cursor returns the selected row index.
func (dt *DataTable) Cursor() int { return dt.cursor }

// Update handles keyboard navigation.
func (dt *DataTable) Update(_ context.Context, m msg.Msg) msg.Cmd {
	switch v := m.(type) {
	case msg.FocusMsg:
		dt.MarkDirty()
	case msg.BlurMsg:
		dt.MarkDirty()
	case msg.MouseDownMsg:
		// Row 0 is the header; data rows start at row 1.
		rowClicked := int(v.Y) - 1
		dataIdx := dt.scrollY + rowClicked
		if rowClicked >= 0 && dataIdx < len(dt.rows) {
			dt.cursor = dataIdx
			dt.MarkDirty()
			if v.Button == msg.MouseButtonLeft {
				row := dt.rows[dataIdx]
				return func(_ context.Context) msg.Msg {
					return DataTableSelectedMsg{Row: dataIdx, Cols: row}
				}
			}
		}
		return nil
	case msg.KeyMsg:
		switch v.Key {
		case keys.Up:
			if dt.cursor > 0 {
				dt.cursor--
				dt.MarkDirty()
			}
		case keys.Down:
			if dt.cursor < len(dt.rows)-1 {
				dt.cursor++
				dt.MarkDirty()
			}
		case keys.Home:
			dt.cursor = 0
			dt.MarkDirty()
		case keys.End:
			if len(dt.rows) > 0 {
				dt.cursor = len(dt.rows) - 1
			}
			dt.MarkDirty()
		case keys.PageUp:
			dt.cursor -= 10
			if dt.cursor < 0 {
				dt.cursor = 0
			}
			dt.MarkDirty()
		case keys.PageDown:
			dt.cursor += 10
			if dt.cursor >= len(dt.rows) && len(dt.rows) > 0 {
				dt.cursor = len(dt.rows) - 1
			}
			dt.MarkDirty()
		case keys.Enter:
			if dt.cursor < len(dt.rows) {
				idx := dt.cursor
				row := dt.rows[idx]
				return func(_ context.Context) msg.Msg {
					return DataTableSelectedMsg{Row: idx, Cols: row}
				}
			}
		}
	}
	return nil
}

// colWidths computes effective column widths for a given total width.
func (dt *DataTable) colWidths(totalWidth int) []int {
	widths := make([]int, len(dt.columns))
	fixed := 0
	autoCount := 0
	for i, col := range dt.columns {
		if col.Width > 0 {
			widths[i] = col.Width
			fixed += col.Width + 1 // +1 for separator
		} else {
			autoCount++
		}
	}
	remaining := totalWidth - fixed
	if autoCount > 0 && remaining > 0 {
		share := remaining / autoCount
		for i, col := range dt.columns {
			if col.Width == 0 {
				widths[i] = share
			}
		}
	}
	return widths
}

// Render draws the header and visible rows.
func (dt *DataTable) Render(region geometry.Region) []strip.Strip {
	strips := make([]strip.Strip, region.Height)
	if region.Height == 0 || region.Width == 0 {
		return strips
	}

	widths := dt.colWidths(region.Width)
	row := 0

	// Header row
	headerStyle := rich.NewStyle().Bold().Underline()
	var headerSegs rich.Segments
	for i, col := range dt.columns {
		cell := truncatePad(col.Header, widths[i])
		headerSegs = append(headerSegs, rich.Segment{Text: cell, Style: headerStyle})
		if i < len(dt.columns)-1 {
			headerSegs = append(headerSegs, rich.Segment{Text: " ", Style: headerStyle})
		}
	}
	hdr := strip.New(headerSegs)
	hdr = hdr.TextAlign(region.Width, css.AlignHorizontal("left"))
	strips[row] = hdr
	row++

	// Scroll to keep cursor visible
	dataHeight := region.Height - 1
	if dt.cursor < dt.scrollY {
		dt.scrollY = dt.cursor
	}
	if dt.cursor >= dt.scrollY+dataHeight && dataHeight > 0 {
		dt.scrollY = dt.cursor - dataHeight + 1
	}

	// Data rows
	for ; row < region.Height; row++ {
		dataIdx := dt.scrollY + (row - 1)
		if dataIdx >= len(dt.rows) {
			strips[row] = strip.New(nil)
			continue
		}
		dataRow := dt.rows[dataIdx]

		var style rich.Style
		if dataIdx == dt.cursor {
			style = rich.NewStyle().Reverse()
		} else if dataIdx%2 == 0 {
			style = rich.NewStyle()
		} else {
			style = rich.NewStyle().Dim()
		}

		var segs rich.Segments
		for i := range dt.columns {
			val := ""
			if i < len(dataRow) {
				val = dataRow[i]
			}
			cell := truncatePad(val, widths[i])
			segs = append(segs, rich.Segment{Text: cell, Style: style})
			if i < len(dt.columns)-1 {
				segs = append(segs, rich.Segment{Text: " ", Style: style})
			}
		}
		s := strip.New(segs)
		s = s.TextAlign(region.Width, css.AlignHorizontal("left"))
		strips[row] = s
	}
	return strips
}

// truncatePad truncates or pads a string to exactly width runes.
func truncatePad(s string, width int) string {
	runes := []rune(s)
	if len(runes) > width {
		if width > 1 {
			return string(runes[:width-1]) + "…"
		}
		return string(runes[:width])
	}
	return fmt.Sprintf("%-*s", width, s)
}
