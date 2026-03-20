package layout

import (
	"math/big"
	"testing"

	"github.com/eberle1080/go-textual/geometry"
)

func TestWidgetPlacement_ResetOrigin(t *testing.T) {
	wp := WidgetPlacement{
		Region: geometry.Region{X: 10, Y: 20, Width: 30, Height: 40},
	}
	wp.ResetOrigin()
	if wp.Region.X != 0 || wp.Region.Y != 0 {
		t.Fatalf("expected (0,0), got (%d,%d)", wp.Region.X, wp.Region.Y)
	}
	if wp.Region.Width != 30 || wp.Region.Height != 40 {
		t.Fatalf("size should be preserved, got %dx%d", wp.Region.Width, wp.Region.Height)
	}
}

func TestWidgetPlacement_ProcessOffset_Normal(t *testing.T) {
	wp := WidgetPlacement{
		Region: geometry.Region{X: 5, Y: 5, Width: 10, Height: 10},
		Offset: geometry.Offset{X: 2, Y: 3},
	}
	constrainRegion := geometry.Region{X: 10, Y: 10, Width: 100, Height: 100}
	processed := wp.ProcessOffset(constrainRegion, geometry.Offset{})
	if processed.Region.X != 17 { // 5 + 10 + 2
		t.Fatalf("expected X=17, got %d", processed.Region.X)
	}
	if processed.Region.Y != 18 { // 5 + 10 + 3
		t.Fatalf("expected Y=18, got %d", processed.Region.Y)
	}
}

func TestWidgetPlacement_ProcessOffset_Absolute(t *testing.T) {
	wp := WidgetPlacement{
		Region:   geometry.Region{X: 5, Y: 5, Width: 10, Height: 10},
		Absolute: true,
	}
	processed := wp.ProcessOffset(geometry.Region{}, geometry.Offset{X: 42, Y: 7})
	if processed.Region.X != 42 || processed.Region.Y != 7 {
		t.Fatalf("expected absolute position (42,7), got (%d,%d)",
			processed.Region.X, processed.Region.Y)
	}
}

func TestTranslatePlacements(t *testing.T) {
	placements := []WidgetPlacement{
		{Region: geometry.Region{X: 0, Y: 0, Width: 10, Height: 10}},
		{Region: geometry.Region{X: 10, Y: 0, Width: 10, Height: 10}},
	}
	translated := TranslatePlacements(placements, geometry.Offset{X: 5, Y: 3})
	if translated[0].Region.X != 5 || translated[0].Region.Y != 3 {
		t.Fatalf("expected (5,3), got (%d,%d)", translated[0].Region.X, translated[0].Region.Y)
	}
	if translated[1].Region.X != 15 {
		t.Fatalf("expected X=15, got %d", translated[1].Region.X)
	}
}

func TestTranslatePlacements_ZeroOffset(t *testing.T) {
	placements := []WidgetPlacement{
		{Region: geometry.Region{X: 5, Y: 5, Width: 10, Height: 10}},
	}
	result := TranslatePlacements(placements, geometry.Offset{})
	if result[0].Region.X != 5 {
		t.Fatal("zero offset should not change placements")
	}
}

func TestGetBounds(t *testing.T) {
	placements := []WidgetPlacement{
		{Region: geometry.Region{X: 0, Y: 0, Width: 10, Height: 5}},
		{Region: geometry.Region{X: 20, Y: 10, Width: 15, Height: 8}},
	}
	bounds := GetBounds(placements)
	if bounds.X != 0 || bounds.Y != 0 {
		t.Fatalf("expected origin (0,0), got (%d,%d)", bounds.X, bounds.Y)
	}
	if bounds.Width != 35 { // max x2=35
		t.Fatalf("expected Width=35, got %d", bounds.Width)
	}
	if bounds.Height != 18 { // max y2=18
		t.Fatalf("expected Height=18, got %d", bounds.Height)
	}
}

func TestGetBounds_Empty(t *testing.T) {
	bounds := GetBounds(nil)
	if bounds.Width != 0 || bounds.Height != 0 {
		t.Fatal("empty placements should yield zero bounds")
	}
}

func TestBoxModel_Fields(t *testing.T) {
	bm := BoxModel{
		Width:  new(big.Rat).SetInt64(80),
		Height: new(big.Rat).SetInt64(24),
		Margin: geometry.Spacing{Top: 1, Bottom: 1},
	}
	wf, _ := bm.Width.Float64()
	hf, _ := bm.Height.Float64()
	if int(wf) != 80 || int(hf) != 24 {
		t.Fatalf("expected 80x24, got %gx%g", wf, hf)
	}
}

func TestGetLayout_Known(t *testing.T) {
	for _, name := range []string{"vertical", "horizontal", "grid", "stream", ""} {
		l, err := GetLayout(name)
		if err != nil {
			t.Fatalf("GetLayout(%q) returned error: %v", name, err)
		}
		if l == nil {
			t.Fatalf("GetLayout(%q) returned nil layout", name)
		}
	}
}

func TestGetLayout_Unknown(t *testing.T) {
	_, err := GetLayout("magic")
	if err == nil {
		t.Fatal("expected error for unknown layout")
	}
	if mle, ok := err.(*MissingLayoutError); !ok || mle.Name != "magic" {
		t.Fatalf("unexpected error type or name: %v", err)
	}
}
