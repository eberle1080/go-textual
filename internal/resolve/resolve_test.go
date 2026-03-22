package resolve

import (
	"testing"

	"github.com/eberle1080/go-textual/css"
	"github.com/eberle1080/go-textual/geometry"
)

// mockWidget satisfies BoxModelable.
type mockWidget struct {
	display bool
}

func (m *mockWidget) Styles() *css.RenderStyles {
	return css.NewRenderStyles(css.NewStyles(), css.NewStyles())
}
func (m *mockWidget) Display() bool { return m.display }

func TestResolve_ExplicitCells(t *testing.T) {
	dims := []css.Scalar{
		{Value: 10, Unit: css.UnitCells},
		{Value: 20, Unit: css.UnitCells},
		{Value: 30, Unit: css.UnitCells},
	}
	widgets := []BoxModelable{
		&mockWidget{display: true},
		&mockWidget{display: true},
		&mockWidget{display: true},
	}
	size := geometry.Size{Width: 100, Height: 100}
	viewport := geometry.Size{Width: 100, Height: 100}

	slots := Resolve(dims, 100, 0, size, viewport, widgets)
	if len(slots) != 3 {
		t.Fatalf("expected 3 slots, got %d", len(slots))
	}
	if slots[0].Offset != 0 || slots[0].Length != 10 {
		t.Fatalf("slot[0] wrong: %+v", slots[0])
	}
	if slots[1].Offset != 10 || slots[1].Length != 20 {
		t.Fatalf("slot[1] wrong: %+v", slots[1])
	}
	if slots[2].Offset != 30 || slots[2].Length != 30 {
		t.Fatalf("slot[2] wrong: %+v", slots[2])
	}
}

func TestResolve_WithGutter(t *testing.T) {
	dims := []css.Scalar{
		{Value: 10, Unit: css.UnitCells},
		{Value: 10, Unit: css.UnitCells},
	}
	widgets := []BoxModelable{
		&mockWidget{display: true},
		&mockWidget{display: true},
	}
	size := geometry.Size{Width: 100, Height: 100}
	viewport := geometry.Size{Width: 100, Height: 100}

	slots := Resolve(dims, 100, 2, size, viewport, widgets)
	if slots[1].Offset != 12 { // 10 + 2 gutter
		t.Fatalf("expected slot[1].Offset=12, got %d", slots[1].Offset)
	}
}

func TestResolve_MixedCellsAndFr(t *testing.T) {
	// [20cells, 1fr] with total=100 → fr gets 80
	dims := []css.Scalar{
		{Value: 20, Unit: css.UnitCells},
		{Value: 1, Unit: css.UnitFraction},
	}
	widgets := []BoxModelable{
		&mockWidget{display: true},
		&mockWidget{display: true},
	}
	size := geometry.Size{Width: 100, Height: 100}
	viewport := geometry.Size{Width: 100, Height: 100}

	slots := Resolve(dims, 100, 0, size, viewport, widgets)
	if slots[0].Length != 20 {
		t.Fatalf("expected slot[0].Length=20, got %d", slots[0].Length)
	}
	if slots[1].Length != 80 {
		t.Fatalf("expected slot[1].Length=80, got %d", slots[1].Length)
	}
	if slots[1].Offset != 20 {
		t.Fatalf("expected slot[1].Offset=20, got %d", slots[1].Offset)
	}
}

func TestResolve_MixedCellsAndFrWithGutter(t *testing.T) {
	// [20cells, 1fr] with total=100, gutter=2 → fr gets 100-20-2=78
	dims := []css.Scalar{
		{Value: 20, Unit: css.UnitCells},
		{Value: 1, Unit: css.UnitFraction},
	}
	widgets := []BoxModelable{
		&mockWidget{display: true},
		&mockWidget{display: true},
	}
	size := geometry.Size{Width: 100, Height: 100}
	viewport := geometry.Size{Width: 100, Height: 100}

	slots := Resolve(dims, 100, 2, size, viewport, widgets)
	if slots[0].Length != 20 {
		t.Fatalf("expected slot[0].Length=20, got %d", slots[0].Length)
	}
	if slots[1].Length != 78 {
		t.Fatalf("expected slot[1].Length=78, got %d", slots[1].Length)
	}
}

func TestResolve_HiddenWidgetWithFr(t *testing.T) {
	// [20cells (hidden), 1fr] — hidden widget should not consume space; fr gets all 100
	dims := []css.Scalar{
		{Value: 20, Unit: css.UnitCells},
		{Value: 1, Unit: css.UnitFraction},
	}
	widgets := []BoxModelable{
		&mockWidget{display: false}, // hidden
		&mockWidget{display: true},
	}
	size := geometry.Size{Width: 100, Height: 100}
	viewport := geometry.Size{Width: 100, Height: 100}

	slots := Resolve(dims, 100, 0, size, viewport, widgets)
	if slots[0].Length != 0 {
		t.Fatalf("expected hidden slot[0].Length=0, got %d", slots[0].Length)
	}
	if slots[1].Length != 100 {
		t.Fatalf("expected slot[1].Length=100 (gets all space), got %d", slots[1].Length)
	}
}

func TestResolve_HiddenWidget(t *testing.T) {
	dims := []css.Scalar{
		{Value: 10, Unit: css.UnitCells},
		{Value: 10, Unit: css.UnitCells},
	}
	widgets := []BoxModelable{
		&mockWidget{display: false}, // hidden
		&mockWidget{display: true},
	}
	size := geometry.Size{Width: 100, Height: 100}
	viewport := geometry.Size{Width: 100, Height: 100}

	slots := Resolve(dims, 100, 0, size, viewport, widgets)
	if slots[0].Length != 0 {
		t.Fatalf("expected hidden widget to have Length=0, got %d", slots[0].Length)
	}
}
