package layout

import (
	"math/big"
	"testing"

	"github.com/eberle1080/go-textual/css"
	"github.com/eberle1080/go-textual/dom"
	"github.com/eberle1080/go-textual/geometry"
)

// arrangeTestWidget is a minimal Layoutable for arrange regression tests.
type arrangeTestWidget struct {
	*dom.DOMNode
	fixedWidth  int
	fixedHeight int
	// boxWidth and boxHeight represent the outer box (content + padding + border).
	// When non-zero, GetBoxModel returns these instead of fixedWidth/fixedHeight,
	// letting tests simulate bordered or padded widgets without a full widget.Widget.
	boxWidth  int
	boxHeight int
	sortOrd   int
}

func newArrangeWidget(typeName string) *arrangeTestWidget {
	return &arrangeTestWidget{
		DOMNode: dom.NewDOMNode(dom.WithCSSTypeName(typeName, "Widget", "DOMNode")),
		sortOrd: 1,
	}
}

func (w *arrangeTestWidget) GetBoxModel(
	container, viewport geometry.Size,
	widthFrac, heightFrac *big.Rat,
	constrainWidth, greedy bool,
) BoxModel {
	var bw, bh *big.Rat
	outerW := w.boxWidth
	if outerW == 0 {
		outerW = w.fixedWidth
	}
	if outerW > 0 {
		bw = new(big.Rat).SetInt64(int64(outerW))
	}
	outerH := w.boxHeight
	if outerH == 0 {
		outerH = w.fixedHeight
	}
	if outerH > 0 {
		bh = new(big.Rat).SetInt64(int64(outerH))
	}
	return BoxModel{Width: bw, Height: bh, Margin: geometry.Spacing{}}
}

func (w *arrangeTestWidget) GetContentWidth(container, viewport geometry.Size) int {
	if w.fixedWidth > 0 {
		return w.fixedWidth
	}
	return container.Width
}

func (w *arrangeTestWidget) GetContentHeight(container, viewport geometry.Size, width int) int {
	if w.fixedHeight > 0 {
		return w.fixedHeight
	}
	return 1
}

func (w *arrangeTestWidget) Expand() bool                                         { return true }
func (w *arrangeTestWidget) Shrink() bool                                         { return true }
func (w *arrangeTestWidget) SortOrder() int                                       { return w.sortOrd }
func (w *arrangeTestWidget) AbsoluteOffset() *geometry.Offset                     { return nil }
func (w *arrangeTestWidget) Layer() string                                        { return "" }
func (w *arrangeTestWidget) PreLayout(_ Layout)                                   {}
func (w *arrangeTestWidget) ProcessLayout(ps []WidgetPlacement) []WidgetPlacement { return ps }
func (w *arrangeTestWidget) ViewportSize() geometry.Size                          { return geometry.Size{Width: 80, Height: 24} }
func (w *arrangeTestWidget) AppSize() geometry.Size                               { return geometry.Size{Width: 80, Height: 24} }
func (w *arrangeTestWidget) ScrollableContentRegion() geometry.Region {
	return geometry.Region{Width: 80, Height: 24}
}

// Verify compile-time satisfaction of Layoutable.
var _ Layoutable = (*arrangeTestWidget)(nil)

// TestArrange_PlainChildren verifies that children with no dock or split rule
// reach the layout engine and receive non-zero placements.
func TestArrange_PlainChildren(t *testing.T) {
	parent := newArrangeWidget("Parent")
	parent.DOMNode.InlineStyles().SetLayout("vertical")
	c1 := newArrangeWidget("Child")
	c2 := newArrangeWidget("Child")

	children := []Layoutable{c1, c2}
	size := geometry.Size{Width: 80, Height: 24}

	result := Arrange(parent, children, size, size, false)
	if len(result.Placements) != 2 {
		t.Fatalf("expected 2 placements, got %d", len(result.Placements))
	}
	for i, p := range result.Placements {
		if p.Region.Width != 80 {
			t.Fatalf("placement[%d] width=%d, want 80", i, p.Region.Width)
		}
	}
}

// TestArrange_DockNoneNotDocked verifies that a widget whose Dock() returns
// "none" (the default) is NOT classified as docked and goes through layout.
func TestArrange_DockNoneNotDocked(t *testing.T) {
	parent := newArrangeWidget("Parent")
	parent.DOMNode.InlineStyles().SetLayout("vertical")
	plain := newArrangeWidget("Plain")
	// Leave dock unset; Dock() returns "none" — must not be classified docked.

	result := Arrange(parent, []Layoutable{plain}, geometry.Size{Width: 80, Height: 24}, geometry.Size{Width: 80, Height: 24}, false)
	if len(result.Placements) != 1 {
		t.Fatalf("expected 1 placement, got %d", len(result.Placements))
	}
	if result.Placements[0].Region.Width != 80 {
		t.Fatalf("plain child width=%d, want 80", result.Placements[0].Region.Width)
	}
}

// TestArrange_DockEdges verifies that docked children consume space from each
// edge and the main child gets the remaining area.
func TestArrange_DockEdges(t *testing.T) {
	for _, edge := range []string{"top", "bottom", "left", "right"} {
		t.Run(edge, func(t *testing.T) {
			parent := newArrangeWidget("Parent")
			parent.DOMNode.InlineStyles().SetLayout("vertical")

			docked := newArrangeWidget("DockedChild")
			docked.DOMNode.InlineStyles().SetDock(edge)
			docked.fixedHeight = 3
			docked.fixedWidth = 10

			main := newArrangeWidget("Main")

			result := Arrange(parent, []Layoutable{docked, main}, geometry.Size{Width: 80, Height: 24}, geometry.Size{Width: 80, Height: 24}, false)
			if len(result.Placements) != 2 {
				t.Fatalf("edge=%q: expected 2 placements, got %d", edge, len(result.Placements))
			}
			dockedP := result.Placements[0]
			if dockedP.Region.Width == 0 && dockedP.Region.Height == 0 {
				t.Fatalf("edge=%q: docked widget got zero-size placement", edge)
			}
		})
	}
}

// TestArrange_SplitEdges verifies that split children consume space from the
// named edge and the remaining area shrinks accordingly.
func TestArrange_SplitEdges(t *testing.T) {
	for _, edge := range []string{"left", "right", "top", "bottom"} {
		t.Run(edge, func(t *testing.T) {
			parent := newArrangeWidget("Parent")
			parent.DOMNode.InlineStyles().SetLayout("vertical")

			splitter := newArrangeWidget("Splitter")
			splitter.DOMNode.InlineStyles().SetSplit(edge)
			splitter.fixedWidth = 20
			splitter.fixedHeight = 6

			main := newArrangeWidget("Main")

			result := Arrange(parent, []Layoutable{splitter, main}, geometry.Size{Width: 80, Height: 24}, geometry.Size{Width: 80, Height: 24}, false)
			if len(result.Placements) != 2 {
				t.Fatalf("edge=%q: expected 2 placements, got %d", edge, len(result.Placements))
			}
			splitP := result.Placements[0]
			if splitP.Region.Width == 0 && splitP.Region.Height == 0 {
				t.Fatalf("edge=%q: split widget got zero-size placement", edge)
			}
		})
	}
}

// TestArrange_SplitNoneNotSplit verifies that split:none does not enter the
// split arrangement path.
func TestArrange_SplitNoneNotSplit(t *testing.T) {
	parent := newArrangeWidget("Parent")
	parent.DOMNode.InlineStyles().SetLayout("vertical")
	child := newArrangeWidget("Child")
	child.DOMNode.InlineStyles().SetSplit("none")

	result := Arrange(parent, []Layoutable{child}, geometry.Size{Width: 80, Height: 24}, geometry.Size{Width: 80, Height: 24}, false)
	if len(result.Placements) != 1 {
		t.Fatalf("expected 1 placement (split:none should be plain), got %d", len(result.Placements))
	}
	if result.Placements[0].Region.Width != 80 {
		t.Fatalf("split:none child should get full width 80, got %d", result.Placements[0].Region.Width)
	}
}

// TestArrange_GridWithSpansAndGutters verifies that grid layout honours
// grid_gutter_horizontal and column_span.
func TestArrange_GridWithSpansAndGutters(t *testing.T) {
	parent := newArrangeWidget("Parent")
	parent.DOMNode.InlineStyles().SetLayout("grid")
	parent.DOMNode.InlineStyles().SetGridSizeColumns(2)
	parent.DOMNode.InlineStyles().SetGridGutterHorizontal(2)

	c1 := newArrangeWidget("C1")
	c2 := newArrangeWidget("C2")
	c3 := newArrangeWidget("C3")
	c3.DOMNode.InlineStyles().SetColumnSpan(2)

	// 2 columns × 40 wide + 2 gutter = 82 total width
	size := geometry.Size{Width: 82, Height: 24}
	result := Arrange(parent, []Layoutable{c1, c2, c3}, size, size, false)
	if len(result.Placements) != 3 {
		t.Fatalf("expected 3 placements, got %d", len(result.Placements))
	}
	for _, p := range result.Placements {
		if p.Widget == c3 {
			if p.Region.Width != 82 {
				t.Fatalf("c3 span=2 width=%d, want 82", p.Region.Width)
			}
		}
	}
}

// TestArrange_HorizontalAutoWidth verifies that a child with width:auto in a
// HorizontalLayout receives its intrinsic content width rather than collapsing
// to zero (the UnitAuto → resolve.Resolve fallback).
func TestArrange_HorizontalAutoWidth(t *testing.T) {
	parent := newArrangeWidget("Parent")
	parent.DOMNode.InlineStyles().SetLayout("horizontal")

	child := newArrangeWidget("Child")
	child.fixedWidth = 25 // intrinsic content width
	child.DOMNode.InlineStyles().SetWidth(css.Scalar{Value: 1, Unit: css.UnitAuto})

	size := geometry.Size{Width: 80, Height: 24}
	result := Arrange(parent, []Layoutable{child}, size, size, false)

	if len(result.Placements) != 1 {
		t.Fatalf("expected 1 placement, got %d", len(result.Placements))
	}
	if result.Placements[0].Region.Width != 25 {
		t.Fatalf("expected auto width=25 (intrinsic), got %d", result.Placements[0].Region.Width)
	}
}

// TestArrange_HorizontalAutoWidth_OuterBox verifies that a width:auto child
// in a HorizontalLayout is allocated its outer box width (which includes
// borders and padding) rather than only raw content width.
func TestArrange_HorizontalAutoWidth_OuterBox(t *testing.T) {
	parent := newArrangeWidget("Parent")
	parent.DOMNode.InlineStyles().SetLayout("horizontal")

	child := newArrangeWidget("Child")
	child.fixedWidth = 25 // intrinsic content width
	child.boxWidth = 27   // outer box = content + 2 border cells
	child.DOMNode.InlineStyles().SetWidth(css.Scalar{Value: 1, Unit: css.UnitAuto})

	size := geometry.Size{Width: 80, Height: 24}
	result := Arrange(parent, []Layoutable{child}, size, size, false)

	if len(result.Placements) != 1 {
		t.Fatalf("expected 1 placement, got %d", len(result.Placements))
	}
	if result.Placements[0].Region.Width != 27 {
		t.Fatalf("expected auto width=27 (outer box with borders), got %d", result.Placements[0].Region.Width)
	}
}

// TestArrange_GridAutoTracks_OuterBox verifies that auto grid tracks are sized
// from the outer box width/height of children in that track, so bordered or
// padded widgets are fully accommodated.
func TestArrange_GridAutoTracks_OuterBox(t *testing.T) {
	parent := newArrangeWidget("Parent")
	parent.DOMNode.InlineStyles().SetLayout("grid")
	parent.DOMNode.InlineStyles().SetGridColumns([]css.Scalar{
		{Value: 1, Unit: css.UnitAuto},
	})
	parent.DOMNode.InlineStyles().SetGridRows([]css.Scalar{
		{Value: 1, Unit: css.UnitAuto},
	})

	c1 := newArrangeWidget("C1")
	c1.fixedWidth = 20 // intrinsic content width
	c1.fixedHeight = 5 // intrinsic content height
	c1.boxWidth = 22   // outer box = content + 2 border columns
	c1.boxHeight = 7   // outer box = content + 2 border rows

	size := geometry.Size{Width: 80, Height: 24}
	result := Arrange(parent, []Layoutable{c1}, size, size, false)

	if len(result.Placements) != 1 {
		t.Fatalf("expected 1 placement, got %d", len(result.Placements))
	}
	if result.Placements[0].Region.Width != 22 {
		t.Fatalf("expected auto column width=22 (outer box), got %d", result.Placements[0].Region.Width)
	}
	if result.Placements[0].Region.Height != 7 {
		t.Fatalf("expected auto row height=7 (outer box), got %d", result.Placements[0].Region.Height)
	}
}

// TestArrange_VerticalAutoHeight_OuterBox verifies that a height:auto child
// in a VerticalLayout is allocated its outer box height (including borders and
// padding) rather than only raw content height.
func TestArrange_VerticalAutoHeight_OuterBox(t *testing.T) {
	parent := newArrangeWidget("Parent")
	parent.DOMNode.InlineStyles().SetLayout("vertical")

	child := newArrangeWidget("Child")
	child.fixedHeight = 5 // intrinsic content height
	child.boxHeight = 7   // outer box = content + 2 border rows
	// No explicit CSS height set → VerticalLayout assigns UnitAuto.

	size := geometry.Size{Width: 80, Height: 24}
	result := Arrange(parent, []Layoutable{child}, size, size, false)

	if len(result.Placements) != 1 {
		t.Fatalf("expected 1 placement, got %d", len(result.Placements))
	}
	if result.Placements[0].Region.Height != 7 {
		t.Fatalf("expected auto height=7 (outer box with borders), got %d",
			result.Placements[0].Region.Height)
	}
}

// TestArrange_GridAutoTracks verifies that grid rows and columns declared as
// auto size to intrinsic child content rather than collapsing to zero.
func TestArrange_GridAutoTracks(t *testing.T) {
	parent := newArrangeWidget("Parent")
	parent.DOMNode.InlineStyles().SetLayout("grid")
	// 2 auto columns, 1 auto row.
	parent.DOMNode.InlineStyles().SetGridColumns([]css.Scalar{
		{Value: 1, Unit: css.UnitAuto},
		{Value: 1, Unit: css.UnitAuto},
	})
	parent.DOMNode.InlineStyles().SetGridRows([]css.Scalar{
		{Value: 1, Unit: css.UnitAuto},
	})

	c1 := newArrangeWidget("C1")
	c1.fixedWidth = 20
	c1.fixedHeight = 5
	c2 := newArrangeWidget("C2")
	c2.fixedWidth = 20
	c2.fixedHeight = 5

	size := geometry.Size{Width: 80, Height: 24}
	result := Arrange(parent, []Layoutable{c1, c2}, size, size, false)

	if len(result.Placements) != 2 {
		t.Fatalf("expected 2 placements, got %d", len(result.Placements))
	}
	for _, p := range result.Placements {
		if p.Region.Width == 0 {
			t.Fatalf("auto column should have non-zero width, got placement %+v", p.Region)
		}
		if p.Region.Height == 0 {
			t.Fatalf("auto row should have non-zero height, got placement %+v", p.Region)
		}
		if p.Region.Width != 20 {
			t.Fatalf("expected auto column width=20 (intrinsic), got %d", p.Region.Width)
		}
		if p.Region.Height != 5 {
			t.Fatalf("expected auto row height=5 (intrinsic), got %d", p.Region.Height)
		}
	}
}
