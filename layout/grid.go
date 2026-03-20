package layout

import (
	"github.com/eberle1080/go-textual/css"
	"github.com/eberle1080/go-textual/geometry"
	"github.com/eberle1080/go-textual/internal/layoutresolve"
	"github.com/eberle1080/go-textual/internal/resolve"
)

// GridLayout arranges children in a rectangular grid with optional column/row
// spans. Grid dimensions are specified via CSS grid-rows and grid-columns
// scalar lists.
type GridLayout struct {
	// MinColumnWidth and MaxColumnWidth clamp each column's resolved width.
	MinColumnWidth *int
	MaxColumnWidth *int
	// StretchHeight makes rows fill all available vertical space.
	StretchHeight bool
	// Regular forces all cells to be the same size.
	Regular bool
	// gridSize caches the last-used grid dimensions [columns, rows].
	gridSize *[2]int
}

// Name returns "grid".
func (g *GridLayout) Name() string { return "grid" }

// Arrange positions children in a grid.
func (g *GridLayout) Arrange(
	parent Layoutable,
	children []Layoutable,
	size geometry.Size,
	greedy bool,
) []WidgetPlacement {
	if len(children) == 0 {
		return nil
	}

	viewport := parent.ViewportSize()
	containerSize := size

	styles := parent.Styles()

	// Use grid-specific gutters when set; fall back to the general gutter.
	gutterH := 0
	gutterV := 0
	if v, ok := styles.GetRule("grid_gutter_horizontal"); ok {
		if n, ok := v.(int); ok {
			gutterH = n
		}
	} else {
		gutterH = styles.Gutter().Left + styles.Gutter().Right
	}
	if v, ok := styles.GetRule("grid_gutter_vertical"); ok {
		if n, ok := v.(int); ok {
			gutterV = n
		}
	} else {
		gutterV = styles.Gutter().Top + styles.Gutter().Bottom
	}

	// Parse grid-columns scalars (one per column). Default to 1-column grid.
	colScalars := g.parseGridScalars(styles, "grid_columns")
	rowScalars := g.parseGridScalars(styles, "grid_rows")

	// If no explicit grid specification, derive column count from CSS
	// grid_size_columns or auto.
	gridCols := len(colScalars)
	if gridCols == 0 {
		gridCols = g.getGridCols(styles, children)
		colScalars = makeFrScalars(gridCols)
	}

	gridRows := len(rowScalars)
	numCells := len(children)
	if gridRows == 0 {
		gridRows = (numCells + gridCols - 1) / gridCols
		if gridRows == 0 {
			gridRows = 1
		}
		rowScalars = makeFrScalars(gridRows)
	}

	// Pre-assign children to cells so per-track auto sizing only considers
	// the children that actually occupy each track, not all children.
	cellAssigns := preAssignCells(children, gridCols)

	// Substitute auto column tracks with the max resolved outer box width of
	// the children occupying that column. Using the outer box (from GetBoxModel)
	// rather than raw content ensures padding and borders are included.
	for i, sc := range colScalars {
		if sc.IsAuto() {
			maxW := 0
			for ci, child := range children {
				cp := cellAssigns[ci]
				if cp.col < 0 {
					continue // non-displayed
				}
				if cp.col <= i && i < cp.col+cp.colSpan {
					bm := child.GetBoxModel(containerSize, viewport, nil, nil, false, false)
					var w int
					if bm.Width != nil {
						f, _ := bm.Width.Float64()
						w = int(f + 0.5)
					} else {
						w = child.GetContentWidth(containerSize, viewport)
					}
					if w > maxW {
						maxW = w
					}
				}
			}
			colScalars[i] = css.Scalar{Value: float64(maxW), Unit: css.UnitCells}
		}
	}

	// Substitute auto row tracks with the max resolved outer box height of
	// the children occupying that row.
	for i, sc := range rowScalars {
		if sc.IsAuto() {
			maxH := 0
			for ci, child := range children {
				cp := cellAssigns[ci]
				if cp.col < 0 {
					continue // non-displayed
				}
				if cp.row <= i && i < cp.row+cp.rowSpan {
					bm := child.GetBoxModel(containerSize, viewport, nil, nil, false, false)
					var h int
					if bm.Height != nil {
						f, _ := bm.Height.Float64()
						h = int(f + 0.5)
					} else {
						h = child.GetContentHeight(containerSize, viewport, containerSize.Width)
					}
					if h > maxH {
						maxH = h
					}
				}
			}
			rowScalars[i] = css.Scalar{Value: float64(maxH), Unit: css.UnitCells}
		}
	}

	// Resolve column widths.
	colWidgets := makeAdapters(children, gridCols)
	colSlots := resolve.Resolve(colScalars, size.Width, gutterH, containerSize, viewport, colWidgets)

	// Apply column width limits.
	if g.MinColumnWidth != nil || g.MaxColumnWidth != nil {
		for i := range colSlots {
			if g.MinColumnWidth != nil && colSlots[i].Length < *g.MinColumnWidth {
				colSlots[i].Length = *g.MinColumnWidth
			}
			if g.MaxColumnWidth != nil && colSlots[i].Length > *g.MaxColumnWidth {
				colSlots[i].Length = *g.MaxColumnWidth
			}
		}
	}

	// Resolve row heights.
	rowWidgets := makeAdapters(children, gridRows)
	rowSlots := resolve.Resolve(rowScalars, size.Height, gutterV, containerSize, viewport, rowWidgets)

	// Build a sparse occupancy grid large enough for all children plus any
	// row_span overflows. We allocate extra rows defensively.
	totalGridRows := gridRows + numCells // worst case: all 1-col spans
	occupied := make([][]bool, totalGridRows)
	for i := range occupied {
		occupied[i] = make([]bool, gridCols)
	}

	// Span-aware placement: for each child find the next unoccupied cell that
	// fits its (col_span × row_span) footprint, then mark those cells.
	placements := make([]WidgetPlacement, 0, numCells)
	curRow, curCol := 0, 0

	for _, child := range children {
		if !child.Display() {
			continue
		}

		// Read per-child spans.
		colSpan := 1
		rowSpan := 1
		if v, ok := child.Styles().GetRule("column_span"); ok {
			if n, ok := v.(int); ok && n > 1 {
				colSpan = n
			}
		}
		if v, ok := child.Styles().GetRule("row_span"); ok {
			if n, ok := v.(int); ok && n > 1 {
				rowSpan = n
			}
		}
		if colSpan > gridCols {
			colSpan = gridCols
		}

		// Advance curRow/curCol until we find a cell where the span fits.
	search:
		for {
			// Wrap columns.
			if curCol+colSpan > gridCols {
				curCol = 0
				curRow++
			}
			// Grow occupancy grid if needed.
			for curRow+rowSpan > len(occupied) {
				occupied = append(occupied, make([]bool, gridCols))
			}
			// Check whether the footprint is clear.
			for r := 0; r < rowSpan; r++ {
				for c := 0; c < colSpan; c++ {
					if occupied[curRow+r][curCol+c] {
						curCol++
						continue search
					}
				}
			}
			break
		}

		// Mark the footprint as occupied.
		for r := 0; r < rowSpan; r++ {
			for c := 0; c < colSpan; c++ {
				occupied[curRow+r][curCol+c] = true
			}
		}

		// Compute the pixel region spanned by this child.
		// Extend col/row slot arrays if the span reaches beyond what was resolved.
		for len(colSlots) < curCol+colSpan {
			colSlots = append(colSlots, resolve.Slot{})
		}
		for len(rowSlots) < curRow+rowSpan {
			rowSlots = append(rowSlots, resolve.Slot{})
		}

		x := colSlots[curCol].Offset
		y := rowSlots[curRow].Offset
		w := 0
		h := 0
		for c := 0; c < colSpan; c++ {
			w += colSlots[curCol+c].Length
			if c > 0 {
				w += gutterH
			}
		}
		for r := 0; r < rowSpan; r++ {
			h += rowSlots[curRow+r].Length
			if r > 0 {
				h += gutterV
			}
		}

		placements = append(placements, WidgetPlacement{
			Region: geometry.Region{X: x, Y: y, Width: w, Height: h},
			Widget: child,
			Order:  child.SortOrder(),
		})

		// Advance past this child's columns.
		curCol += colSpan
		if curCol >= gridCols {
			curCol = 0
			curRow++
		}
	}
	return placements
}

// GetContentWidth returns the sum of auto column widths.
func (g *GridLayout) GetContentWidth(
	widget Layoutable,
	container, viewport geometry.Size,
) int {
	styles := widget.Styles()
	colScalars := g.parseGridScalars(styles, "grid_columns")
	if len(colScalars) == 0 {
		return container.Width
	}
	total := 0
	for _, sc := range colScalars {
		if sc.IsCells() {
			total += int(sc.Value)
		}
	}
	return total
}

// GetContentHeight returns the sum of auto row heights.
func (g *GridLayout) GetContentHeight(
	widget Layoutable,
	container, viewport geometry.Size,
	width int,
) int {
	styles := widget.Styles()
	rowScalars := g.parseGridScalars(styles, "grid_rows")
	if len(rowScalars) == 0 {
		return container.Height
	}
	total := 0
	for _, sc := range rowScalars {
		if sc.IsCells() {
			total += int(sc.Value)
		}
	}
	return total
}

// parseGridScalars extracts a slice of css.Scalar from the named style rule.
// The rule value is expected to be a []css.Scalar; returns nil if absent.
func (g *GridLayout) parseGridScalars(styles interface{ GetRule(string) (any, bool) }, rule string) []css.Scalar {
	val, ok := styles.GetRule(rule)
	if !ok {
		return nil
	}
	if scalars, ok := val.([]css.Scalar); ok {
		return scalars
	}
	return nil
}

// getGridCols reads grid_size_columns from the parent's styles, or defaults
// to ceil(sqrt(len(children))).
func (g *GridLayout) getGridCols(styles interface{ GetRule(string) (any, bool) }, children []Layoutable) int {
	val, ok := styles.GetRule("grid_size_columns")
	if ok {
		if n, ok := val.(int); ok && n > 0 {
			return n
		}
	}
	// Default: auto columns.
	n := len(children)
	if n == 0 {
		return 1
	}
	cols := 1
	for cols*cols < n {
		cols++
	}
	return cols
}

// makeFrScalars creates n equal 1fr scalars.
func makeFrScalars(n int) []css.Scalar {
	scalars := make([]css.Scalar, n)
	for i := range scalars {
		scalars[i] = css.Scalar{Value: 1, Unit: css.UnitFraction}
	}
	return scalars
}

// makeAdapters wraps the first n children as resolve.BoxModelable.
func makeAdapters(children []Layoutable, n int) []resolve.BoxModelable {
	count := n
	if count > len(children) {
		count = len(children)
	}
	result := make([]resolve.BoxModelable, count)
	for i := 0; i < count; i++ {
		result[i] = &boxModelAdapter{children[i]}
	}
	return result
}

// layoutResolveEdge adapts a css.Scalar to layoutresolve.Edge.
type layoutResolveEdge struct {
	scalar css.Scalar
	min    int
}

func (e *layoutResolveEdge) Size() *int {
	if e.scalar.IsCells() {
		v := int(e.scalar.Value)
		return &v
	}
	return nil
}

func (e *layoutResolveEdge) Fraction() int {
	if e.scalar.IsFraction() {
		return int(e.scalar.Value)
	}
	return 0
}

func (e *layoutResolveEdge) MinSize() int { return e.min }

// makeEdges converts css.Scalar slice to layoutresolve.Edge slice.
func makeEdges(scalars []css.Scalar, minSize int) []layoutresolve.Edge {
	edges := make([]layoutresolve.Edge, len(scalars))
	for i, sc := range scalars {
		edges[i] = &layoutResolveEdge{scalar: sc, min: minSize}
	}
	return edges
}

// boxModelAdapter wraps a Layoutable for use with the internal resolve package.
type boxModelAdapter struct {
	w Layoutable
}

func (a *boxModelAdapter) Styles() *css.RenderStyles { return a.w.Styles() }
func (a *boxModelAdapter) Display() bool             { return a.w.Display() }

// Prevent unused import.
var _ = layoutresolve.Resolve
var _ = makeEdges

// cellPos records the grid position and span of a single child as determined
// by the span-aware cell assignment algorithm.
type cellPos struct {
	col, row         int
	colSpan, rowSpan int
}

// preAssignCells runs the span-aware cell assignment for children and returns
// a cellPos for each child. Non-displayed children get col == -1.
// This mirrors the main Arrange loop's placement logic so per-track auto
// sizing can consult only the children that occupy a given track.
func preAssignCells(children []Layoutable, gridCols int) []cellPos {
	positions := make([]cellPos, len(children))
	occupied := make([][]bool, len(children)+gridCols)
	for i := range occupied {
		occupied[i] = make([]bool, gridCols)
	}
	curRow, curCol := 0, 0
	for ci, child := range children {
		if !child.Display() {
			positions[ci] = cellPos{col: -1, row: -1}
			continue
		}
		cs, rs := 1, 1
		if v, ok := child.Styles().GetRule("column_span"); ok {
			if n, ok := v.(int); ok && n > 1 {
				cs = n
			}
		}
		if v, ok := child.Styles().GetRule("row_span"); ok {
			if n, ok := v.(int); ok && n > 1 {
				rs = n
			}
		}
		if cs > gridCols {
			cs = gridCols
		}
	search:
		for {
			if curCol+cs > gridCols {
				curCol = 0
				curRow++
			}
			for curRow+rs > len(occupied) {
				occupied = append(occupied, make([]bool, gridCols))
			}
			for r := 0; r < rs; r++ {
				for c := 0; c < cs; c++ {
					if occupied[curRow+r][curCol+c] {
						curCol++
						continue search
					}
				}
			}
			break
		}
		for r := 0; r < rs; r++ {
			for c := 0; c < cs; c++ {
				occupied[curRow+r][curCol+c] = true
			}
		}
		positions[ci] = cellPos{col: curCol, row: curRow, colSpan: cs, rowSpan: rs}
		curCol += cs
		if curCol >= gridCols {
			curCol = 0
			curRow++
		}
	}
	return positions
}
