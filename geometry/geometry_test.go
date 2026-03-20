package geometry

import (
	"math"
	"testing"
)

func TestClampInt(t *testing.T) {
	tests := []struct {
		value, min, max, want int
	}{
		{5, 0, 10, 5},
		{-1, 0, 10, 0},
		{11, 0, 10, 10},
		{0, 0, 10, 0},
		{10, 0, 10, 10},
		// reversed min/max
		{5, 10, 0, 5},
		{-1, 10, 0, 0},
		{11, 10, 0, 10},
		{0, 10, 0, 0},
		{10, 10, 0, 10},
	}
	for _, tt := range tests {
		got := ClampInt(tt.value, tt.min, tt.max)
		if got != tt.want {
			t.Errorf("ClampInt(%d, %d, %d) = %d, want %d", tt.value, tt.min, tt.max, got, tt.want)
		}
	}
}

func TestOffsetIsOrigin(t *testing.T) {
	_tmp1 := Offset{0, 0}
	if !_tmp1.IsOrigin() {
		t.Error("Offset{0,0}.IsOrigin() should be true")
	}
	_tmp2 := Offset{1, 0}
	if _tmp2.IsOrigin() {
		t.Error("Offset{1,0}.IsOrigin() should be false")
	}
}

func TestOffsetIsNonZero(t *testing.T) {
	_tmp3 := Offset{0, 0}
	if _tmp3.IsNonZero() {
		t.Error("Offset{0,0}.IsNonZero() should be false")
	}
	_tmp4 := Offset{1, 0}
	if !_tmp4.IsNonZero() {
		t.Error("Offset{1,0}.IsNonZero() should be true")
	}
	_tmp5 := Offset{0, -1}
	if !_tmp5.IsNonZero() {
		t.Error("Offset{0,-1}.IsNonZero() should be true")
	}
}

func TestOffsetClamped(t *testing.T) {
	tests := []struct {
		in, want Offset
	}{
		{Offset{-10, 0}, Offset{0, 0}},
		{Offset{-10, -5}, Offset{0, 0}},
		{Offset{5, -5}, Offset{5, 0}},
		{Offset{5, 10}, Offset{5, 10}},
	}
	for _, tt := range tests {
		if got := tt.in.Clamped(); got != tt.want {
			t.Errorf("Offset%v.Clamped() = %v, want %v", tt.in, got, tt.want)
		}
	}
}

func TestOffsetTranspose(t *testing.T) {
	y, x := Offset{1, 2}.Transpose()
	if y != 2 || x != 1 {
		t.Errorf("Transpose() = (%d, %d), want (2, 1)", y, x)
	}
	y, x = Offset{5, 10}.Transpose()
	if y != 10 || x != 5 {
		t.Errorf("Transpose() = (%d, %d), want (10, 5)", y, x)
	}
}

func TestOffsetAdd(t *testing.T) {
	got := Offset{1, 1}.Add(Offset{2, 2})
	if got != (Offset{3, 3}) {
		t.Errorf("Add = %v, want {3,3}", got)
	}
	got = Offset{1, 2}.Add(Offset{3, 4})
	if got != (Offset{4, 6}) {
		t.Errorf("Add = %v, want {4,6}", got)
	}
}

func TestOffsetSub(t *testing.T) {
	got := Offset{1, 1}.Sub(Offset{2, 2})
	if got != (Offset{-1, -1}) {
		t.Errorf("Sub = %v, want {-1,-1}", got)
	}
	got = Offset{3, 4}.Sub(Offset{2, 1})
	if got != (Offset{1, 3}) {
		t.Errorf("Sub = %v, want {1,3}", got)
	}
}

func TestOffsetNeg(t *testing.T) {
	_tmp6 := Offset{2, -3}
	if _tmp6.Neg() != (Offset{-2, 3}) {
		t.Error("Neg failed")
	}
}

func TestOffsetMulScalar(t *testing.T) {
	tests := []struct {
		in     Offset
		factor float64
		want   Offset
	}{
		{Offset{2, 1}, 2, Offset{4, 2}},
		{Offset{2, 1}, -2, Offset{-4, -2}},
		{Offset{2, 1}, 0, Offset{0, 0}},
	}
	for _, tt := range tests {
		got := tt.in.MulScalar(tt.factor)
		if got != tt.want {
			t.Errorf("MulScalar(%v, %v) = %v, want %v", tt.in, tt.factor, got, tt.want)
		}
	}
}

func TestOffsetBlend(t *testing.T) {
	a := Offset{1, 2}
	b := Offset{3, 4}
	if a.Blend(b, 0) != a {
		t.Error("Blend(0) should return self")
	}
	if a.Blend(b, 1) != b {
		t.Error("Blend(1) should return dest")
	}
	if a.Blend(b, 0.5) != (Offset{2, 3}) {
		t.Error("Blend(0.5) should return midpoint")
	}
}

func TestOffsetDistanceTo(t *testing.T) {
	_tmp7 := Offset{20, 30}
	if _tmp7.DistanceTo(Offset{20, 30}) != 0 {
		t.Error("distance to self should be 0")
	}
	_tmp8 := Offset{0, 0}
	if _tmp8.DistanceTo(Offset{1, 0}) != 1.0 {
		t.Error("distance to adjacent should be 1")
	}
	got := Offset{2, 1}.DistanceTo(Offset{5, 5})
	if math.Abs(got-5.0) > 0.001 {
		t.Errorf("distance = %f, want 5.0", got)
	}
}

func TestOffsetClamp(t *testing.T) {
	tests := []struct {
		in          Offset
		width, height int
		want        Offset
	}{
		{Offset{1, 2}, 3, 3, Offset{1, 2}},
		{Offset{3, 2}, 3, 3, Offset{2, 2}},
		{Offset{-3, 2}, 3, 3, Offset{0, 2}},
		{Offset{5, 4}, 3, 3, Offset{2, 2}},
	}
	for _, tt := range tests {
		got := tt.in.Clamp(tt.width, tt.height)
		if got != tt.want {
			t.Errorf("Offset%v.Clamp(%d,%d) = %v, want %v", tt.in, tt.width, tt.height, got, tt.want)
		}
	}
}

func TestSizeRegion(t *testing.T) {
	got := Size{30, 40}.Region()
	want := Region{0, 0, 30, 40}
	if got != want {
		t.Errorf("Region() = %v, want %v", got, want)
	}
}

func TestSizeContains(t *testing.T) {
	s := Size{10, 10}
	tests := []struct {
		x, y int
		want bool
	}{
		{5, 5, true}, {9, 9, true}, {0, 0, true},
		{10, 9, false}, {9, 10, false}, {-1, 0, false}, {0, -1, false},
	}
	for _, tt := range tests {
		if got := s.Contains(tt.x, tt.y); got != tt.want {
			t.Errorf("Contains(%d,%d) = %v, want %v", tt.x, tt.y, got, tt.want)
		}
	}
}

func TestSizeContainsPoint(t *testing.T) {
	s := Size{10, 10}
	if !s.ContainsPoint(Offset{5, 5}) {
		t.Error("ContainsPoint(5,5) should be true")
	}
	if s.ContainsPoint(Offset{10, 9}) {
		t.Error("ContainsPoint(10,9) should be false")
	}
}

func TestSizeIsNonZero(t *testing.T) {
	_tmp9 := Size{1, 1}
	if !_tmp9.IsNonZero() {
		t.Error("1x1 should be non-zero")
	}
	_tmp10 := Size{0, 1}
	if _tmp10.IsNonZero() {
		t.Error("0x1 should be zero")
	}
	_tmp11 := Size{1, 0}
	if _tmp11.IsNonZero() {
		t.Error("1x0 should be zero")
	}
}

func TestSizeArea(t *testing.T) {
	tests := []struct{ w, h, want int }{{0, 0, 0}, {1, 0, 0}, {1, 1, 1}, {4, 5, 20}}
	for _, tt := range tests {
		_tmp1 := Size{tt.w, tt.h}
		if got := _tmp1.Area(); got != tt.want {
			t.Errorf("Size{%d,%d}.Area() = %d, want %d", tt.w, tt.h, got, tt.want)
		}
	}
}

func TestSizeLineRange(t *testing.T) {
	start, end := Size{20, 0}.LineRange()
	if start != 0 || end != 0 {
		t.Errorf("LineRange of 20x0 = (%d,%d), want (0,0)", start, end)
	}
	start, end = Size{0, 20}.LineRange()
	if start != 0 || end != 20 {
		t.Errorf("LineRange of 0x20 = (%d,%d), want (0,20)", start, end)
	}
}

func TestSizeAddSub(t *testing.T) {
	got := Size{5, 10}.Add(Size{2, 3})
	if got != (Size{7, 13}) {
		t.Errorf("Add = %v, want {7,13}", got)
	}
	got = Size{5, 10}.Sub(Size{2, 3})
	if got != (Size{3, 7}) {
		t.Errorf("Sub = %v, want {3,7}", got)
	}
	// clamp to 0
	got = Size{1, 1}.Sub(Size{5, 5})
	if got != (Size{0, 0}) {
		t.Errorf("Sub clamped = %v, want {0,0}", got)
	}
}

func TestSizeWithWidthHeight(t *testing.T) {
	_tmp12 := Size{1, 2}
	if _tmp12.WithHeight(10) != (Size{1, 10}) {
		t.Error("WithHeight failed")
	}
	_tmp13 := Size{1, 2}
	if _tmp13.WithWidth(10) != (Size{10, 2}) {
		t.Error("WithWidth failed")
	}
}

func TestSizeClampOffset(t *testing.T) {
	tests := []struct {
		s    Size
		in   Offset
		want Offset
	}{
		{Size{3, 3}, Offset{1, 2}, Offset{1, 2}},
		{Size{3, 3}, Offset{3, 2}, Offset{2, 2}},
		{Size{3, 3}, Offset{-3, 2}, Offset{0, 2}},
		{Size{3, 3}, Offset{5, 4}, Offset{2, 2}},
	}
	for _, tt := range tests {
		if got := tt.s.ClampOffset(tt.in); got != tt.want {
			t.Errorf("ClampOffset(%v) = %v, want %v", tt.in, got, tt.want)
		}
	}
}

func TestRegionNull(t *testing.T) {
	_tmp14 := Region{}
	if _tmp14 != (Region{0, 0, 0, 0}) {
		t.Error("null region mismatch")
	}
	_tmp15 := Region{}
	if _tmp15.IsNonZero() {
		t.Error("null region should not be non-zero")
	}
}

func TestRegionFromUnion(t *testing.T) {
	_, err := RegionFromUnion(nil)
	if err == nil {
		t.Error("expected error for empty regions")
	}
	regions := []Region{{10, 20, 30, 40}, {15, 25, 5, 5}, {30, 25, 20, 10}}
	got, err := RegionFromUnion(regions)
	if err != nil {
		t.Fatal(err)
	}
	if got != (Region{10, 20, 40, 40}) {
		t.Errorf("FromUnion = %v, want {10,20,40,40}", got)
	}
}

func TestRegionFromOffset(t *testing.T) {
	got := RegionFromOffset(Offset{3, 4}, Size{5, 6})
	if got != (Region{3, 4, 5, 6}) {
		t.Errorf("FromOffset = %v, want {3,4,5,6}", got)
	}
}

func TestRegionArea(t *testing.T) {
	_tmp16 := Region{3, 4, 0, 0}
	if _tmp16.Area() != 0 {
		t.Error("area of 0-size region should be 0")
	}
	_tmp17 := Region{3, 4, 5, 6}
	if _tmp17.Area() != 30 {
		t.Error("area should be 30")
	}
}

func TestRegionSize(t *testing.T) {
	_tmp18 := Region{3, 4, 5, 6}
	if _tmp18.Size() != (Size{5, 6}) {
		t.Error("Size failed")
	}
}

func TestRegionOffset(t *testing.T) {
	_tmp19 := Region{1, 2, 3, 4}
	if _tmp19.Offset() != (Offset{1, 2}) {
		t.Error("Offset failed")
	}
}

func TestRegionBottomLeft(t *testing.T) {
	_tmp20 := Region{1, 2, 3, 4}
	if _tmp20.BottomLeft() != (Offset{1, 6}) {
		t.Error("BottomLeft failed")
	}
}

func TestRegionTopRight(t *testing.T) {
	_tmp21 := Region{1, 2, 3, 4}
	if _tmp21.TopRight() != (Offset{4, 2}) {
		t.Error("TopRight failed")
	}
}

func TestRegionBottomRight(t *testing.T) {
	_tmp22 := Region{1, 2, 3, 4}
	if _tmp22.BottomRight() != (Offset{4, 6}) {
		t.Error("BottomRight failed")
	}
	_tmp23 := Region{1, 2, 3, 4}
	if _tmp23.BottomRightInclusive() != (Offset{3, 5}) {
		t.Error("BottomRightInclusive failed")
	}
}

func TestRegionAddSub(t *testing.T) {
	_tmp24 := Region{1, 2, 3, 4}
	if _tmp24.Add(Offset{10, 20}) != (Region{11, 22, 3, 4}) {
		t.Error("Add failed")
	}
	_tmp25 := Region{11, 22, 3, 4}
	if _tmp25.Sub(Offset{10, 20}) != (Region{1, 2, 3, 4}) {
		t.Error("Sub failed")
	}
}

func TestRegionAtOffset(t *testing.T) {
	got := Region{10, 10, 30, 40}.AtOffset(Offset{0, 0})
	if got != (Region{0, 0, 30, 40}) {
		t.Errorf("AtOffset(0,0) = %v", got)
	}
	got = Region{10, 10, 30, 40}.AtOffset(Offset{-15, 30})
	if got != (Region{-15, 30, 30, 40}) {
		t.Errorf("AtOffset(-15,30) = %v", got)
	}
}

func TestRegionCropSize(t *testing.T) {
	got := Region{10, 20, 100, 200}.CropSize(Size{50, 40})
	if got != (Region{10, 20, 50, 40}) {
		t.Errorf("CropSize = %v", got)
	}
	got = Region{10, 20, 100, 200}.CropSize(Size{500, 40})
	if got != (Region{10, 20, 100, 40}) {
		t.Errorf("CropSize = %v", got)
	}
}

func TestRegionOverlaps(t *testing.T) {
	_tmp26 := Region{10, 10, 30, 20}
	if !_tmp26.Overlaps(Region{0, 0, 20, 20}) {
		t.Error("should overlap")
	}
	_tmp27 := Region{10, 10, 5, 5}
	if _tmp27.Overlaps(Region{15, 15, 20, 20}) {
		t.Error("should not overlap")
	}
	_tmp28 := Region{10, 10, 5, 5}
	if _tmp28.Overlaps(Region{0, 0, 50, 10}) {
		t.Error("should not overlap")
	}
	_tmp29 := Region{10, 10, 5, 5}
	if !_tmp29.Overlaps(Region{0, 0, 50, 11}) {
		t.Error("should overlap")
	}
}

func TestRegionContains(t *testing.T) {
	r := Region{10, 10, 20, 30}
	if !r.Contains(10, 10) {
		t.Error("should contain (10,10)")
	}
	if !r.Contains(29, 39) {
		t.Error("should contain (29,39)")
	}
	if r.Contains(30, 40) {
		t.Error("should not contain (30,40)")
	}
}

func TestRegionContainsPoint(t *testing.T) {
	r := Region{10, 10, 20, 30}
	if !r.ContainsPoint(Offset{10, 10}) {
		t.Error("should contain (10,10)")
	}
	if r.ContainsPoint(Offset{30, 40}) {
		t.Error("should not contain (30,40)")
	}
}

func TestRegionContainsRegion(t *testing.T) {
	r := Region{10, 10, 20, 30}
	if !r.ContainsRegion(Region{10, 10, 5, 5}) {
		t.Error("should contain")
	}
	if r.ContainsRegion(Region{10, 9, 5, 5}) {
		t.Error("should not contain (y too small)")
	}
	if !r.ContainsRegion(Region{10, 10, 20, 30}) {
		t.Error("should contain itself")
	}
	if r.ContainsRegion(Region{10, 10, 21, 30}) {
		t.Error("should not contain (too wide)")
	}
}

func TestRegionTranslate(t *testing.T) {
	_tmp30 := Region{1, 2, 3, 4}
	if _tmp30.Translate(Offset{10, 20}) != (Region{11, 22, 3, 4}) {
		t.Error("Translate failed")
	}
}

func TestRegionClip(t *testing.T) {
	got := Region{10, 10, 20, 30}.Clip(20, 25)
	if got != (Region{10, 10, 10, 15}) {
		t.Errorf("Clip = %v, want {10,10,10,15}", got)
	}
}

func TestRegionShrink(t *testing.T) {
	margin := Spacing{1, 2, 3, 4}
	region := Region{10, 10, 50, 50}
	got := region.Shrink(margin)
	want := Region{14, 11, 44, 46}
	if got != want {
		t.Errorf("Shrink = %v, want %v", got, want)
	}
}

func TestRegionGrow(t *testing.T) {
	margin := Spacing{1, 2, 3, 4}
	region := Region{10, 10, 50, 50}
	got := region.Grow(margin)
	want := Region{6, 9, 56, 54}
	if got != want {
		t.Errorf("Grow = %v, want %v", got, want)
	}
}

func TestRegionIntersection(t *testing.T) {
	got := Region{0, 0, 100, 50}.Intersection(Region{10, 10, 10, 10})
	if got != (Region{10, 10, 10, 10}) {
		t.Errorf("Intersection = %v", got)
	}
	got = Region{10, 10, 30, 20}.Intersection(Region{20, 15, 60, 40})
	if got != (Region{20, 15, 20, 15}) {
		t.Errorf("Intersection = %v", got)
	}
	_tmp31 := Region{10, 10, 20, 30}
	if _tmp31.Intersection(Region{50, 50, 100, 200}).IsNonZero() {
		t.Error("non-overlapping intersection should be empty")
	}
}

func TestRegionUnion(t *testing.T) {
	got := Region{5, 5, 10, 10}.Union(Region{20, 30, 10, 5})
	if got != (Region{5, 5, 25, 30}) {
		t.Errorf("Union = %v, want {5,5,25,30}", got)
	}
}

func TestRegionColumnLineSpan(t *testing.T) {
	r := Region{5, 10, 20, 30}
	cs, ce := r.ColumnSpan()
	if cs != 5 || ce != 25 {
		t.Errorf("ColumnSpan = (%d,%d), want (5,25)", cs, ce)
	}
	ls, le := r.LineSpan()
	if ls != 10 || le != 40 {
		t.Errorf("LineSpan = (%d,%d), want (10,40)", ls, le)
	}
}

func TestRegionResetOffset(t *testing.T) {
	_tmp32 := Region{5, 10, 20, 30}
	if _tmp32.ResetOffset() != (Region{0, 0, 20, 30}) {
		t.Error("ResetOffset failed")
	}
}

func TestRegionExpand(t *testing.T) {
	got := Region{50, 10, 10, 5}.Expand(Size{2, 3})
	if got != (Region{48, 7, 14, 11}) {
		t.Errorf("Expand = %v, want {48,7,14,11}", got)
	}
}

func TestRegionSplit(t *testing.T) {
	a, b, c, d := Region{10, 5, 22, 15}.Split(10, 5)
	if a != (Region{10, 5, 10, 5}) || b != (Region{20, 5, 12, 5}) ||
		c != (Region{10, 10, 10, 10}) || d != (Region{20, 10, 12, 10}) {
		t.Errorf("Split = %v %v %v %v", a, b, c, d)
	}
}

func TestRegionSplitNegative(t *testing.T) {
	a, b, c, d := Region{10, 5, 22, 15}.Split(-1, -1)
	if a != (Region{10, 5, 21, 14}) || b != (Region{31, 5, 1, 14}) ||
		c != (Region{10, 19, 21, 1}) || d != (Region{31, 19, 1, 1}) {
		t.Errorf("Split(-1,-1) = %v %v %v %v", a, b, c, d)
	}
}

func TestRegionSplitVertical(t *testing.T) {
	a, b := Region{10, 5, 22, 15}.SplitVertical(10)
	if a != (Region{10, 5, 10, 15}) || b != (Region{20, 5, 12, 15}) {
		t.Errorf("SplitVertical = %v %v", a, b)
	}
	a, b = Region{10, 5, 22, 15}.SplitVertical(-1)
	if a != (Region{10, 5, 21, 15}) || b != (Region{31, 5, 1, 15}) {
		t.Errorf("SplitVertical(-1) = %v %v", a, b)
	}
}

func TestRegionSplitHorizontal(t *testing.T) {
	a, b := Region{10, 5, 22, 15}.SplitHorizontal(5)
	if a != (Region{10, 5, 22, 5}) || b != (Region{10, 10, 22, 10}) {
		t.Errorf("SplitHorizontal = %v %v", a, b)
	}
	a, b = Region{10, 5, 22, 15}.SplitHorizontal(-1)
	if a != (Region{10, 5, 22, 14}) || b != (Region{10, 19, 22, 1}) {
		t.Errorf("SplitHorizontal(-1) = %v %v", a, b)
	}
}

func TestRegionTranslateInside(t *testing.T) {
	got := Region{10, 20, 10, 20}.TranslateInside(Region{0, 0, 30, 25}, true, true)
	if got != (Region{10, 5, 10, 20}) {
		t.Errorf("TranslateInside = %v, want {10,5,10,20}", got)
	}
	got = Region{10, 10, 20, 5}.TranslateInside(Region{0, 0, 100, 100}, true, true)
	if got != (Region{10, 10, 20, 5}) {
		t.Errorf("TranslateInside already inside = %v", got)
	}
}

func TestRegionInflect(t *testing.T) {
	_tmp33 := Region{0, 0, 1, 1}
	if _tmp33.Inflect(1, 1, nil) != (Region{1, 1, 1, 1}) {
		t.Error("Inflect default failed")
	}
	s := SpacingAll(1)
	_tmp34 := Region{0, 0, 1, 1}
	if _tmp34.Inflect(1, 1, &s) != (Region{2, 2, 1, 1}) {
		t.Error("Inflect with margin failed")
	}
	m := Spacing{2, 2, 2, 2}
	_tmp35 := Region{10, 10, 30, 20}
	if _tmp35.Inflect(1, 1, &m) != (Region{42, 32, 30, 20}) {
		t.Error("Inflect positive both failed")
	}
	got := Region{10, 10, 30, 20}.Inflect(1, -1, &m)
	if got != (Region{42, -12, 30, 20}) {
		t.Errorf("Inflect y=-1 = %v, want {42,-12,30,20}", got)
	}
	got = Region{10, 10, 30, 20}.Inflect(-1, 1, &m)
	if got != (Region{-22, 32, 30, 20}) {
		t.Errorf("Inflect x=-1 = %v, want {-22,32,30,20}", got)
	}
}

func TestRegionScrollToVisible(t *testing.T) {
	tests := []struct {
		window, region Region
		want           Offset
	}{
		{Region{0, 0, 200, 100}, Region{0, 0, 200, 100}, Offset{0, 0}},
		{Region{0, 0, 200, 100}, Region{0, -100, 10, 10}, Offset{0, -100}},
		{Region{10, 15, 20, 10}, Region{0, 0, 50, 50}, Offset{-10, -15}},
	}
	for _, tt := range tests {
		got := ScrollToVisible(tt.window, tt.region, false)
		if got != tt.want {
			t.Errorf("ScrollToVisible = %v, want %v", got, tt.want)
		}
	}
}

func TestRegionSpacingBetween(t *testing.T) {
	tests := []struct {
		r1, r2   Region
		expected Spacing
	}{
		{Region{0, 0, 100, 80}, Region{0, 0, 100, 80}, Spacing{0, 0, 0, 0}},
		{Region{0, 0, 100, 80}, Region{10, 10, 10, 10}, Spacing{10, 80, 60, 10}},
	}
	for _, tt := range tests {
		got := tt.r1.SpacingBetween(tt.r2)
		if got != tt.expected {
			t.Errorf("SpacingBetween = %v, want %v", got, tt.expected)
		}
		if tt.r1.Shrink(got) != tt.r2 {
			t.Error("Shrink(SpacingBetween) should recover r2")
		}
	}
}

func TestRegionConstrain(t *testing.T) {
	tests := []struct {
		cx, cy    string
		margin    Spacing
		r, c, exp Region
	}{
		{"none", "none", SpacingAll(0), Region{0, 0, 10, 10}, Region{0, 0, 100, 100}, Region{0, 0, 10, 10}},
		{"inside", "inside", SpacingAll(1), Region{-5, -5, 10, 10}, Region{0, 0, 100, 100}, Region{1, 1, 10, 10}},
		{"inside", "inside", SpacingAll(1), Region{95, 95, 10, 10}, Region{0, 0, 100, 100}, Region{89, 89, 10, 10}},
		{"inside", "inflect", SpacingAll(1), Region{-5, -5, 10, 10}, Region{0, 0, 100, 100}, Region{1, 6, 10, 10}},
	}
	for i, tt := range tests {
		got := tt.r.Constrain(tt.cx, tt.cy, tt.margin, tt.c)
		if got != tt.exp {
			t.Errorf("test %d: Constrain = %v, want %v", i, got, tt.exp)
		}
	}
}

func TestSpacingIsNonZero(t *testing.T) {
	_tmp36 := Spacing{1, 0, 0, 0}
	if !_tmp36.IsNonZero() {
		t.Error("{1,0,0,0} should be non-zero")
	}
	_tmp37 := Spacing{}
	if _tmp37.IsNonZero() {
		t.Error("{0,0,0,0} should be zero")
	}
}

func TestSpacingWidthHeight(t *testing.T) {
	s := Spacing{2, 3, 4, 5}
	if s.Width() != 8 {
		t.Errorf("Width = %d, want 8", s.Width())
	}
	if s.Height() != 6 {
		t.Errorf("Height = %d, want 6", s.Height())
	}
}

func TestSpacingTopLeftBottomRight(t *testing.T) {
	s := Spacing{2, 3, 4, 5}
	l, top := s.TopLeft()
	if l != 5 || top != 2 {
		t.Errorf("TopLeft = (%d,%d), want (5,2)", l, top)
	}
	r, bot := s.BottomRight()
	if r != 3 || bot != 4 {
		t.Errorf("BottomRight = (%d,%d), want (3,4)", r, bot)
	}
}

func TestSpacingTotals(t *testing.T) {
	h, v := Spacing{2, 3, 4, 5}.Totals()
	if h != 8 || v != 6 {
		t.Errorf("Totals = (%d,%d), want (8,6)", h, v)
	}
}

func TestSpacingCSS(t *testing.T) {
	_tmp38 := Spacing{1, 1, 1, 1}
	if _tmp38.CSS() != "1" {
		t.Error("CSS all-equal failed")
	}
	_tmp39 := Spacing{1, 2, 1, 2}
	if _tmp39.CSS() != "1 2" {
		t.Error("CSS pair failed")
	}
	_tmp40 := Spacing{1, 2, 3, 4}
	if _tmp40.CSS() != "1 2 3 4" {
		t.Error("CSS four-value failed")
	}
}

func TestUnpackSpacing(t *testing.T) {
	s, _ := UnpackSpacing(1)
	if s != (Spacing{1, 1, 1, 1}) {
		t.Error("Unpack(1) failed")
	}
	s, _ = UnpackSpacing(1, 2)
	if s != (Spacing{1, 2, 1, 2}) {
		t.Error("Unpack(1,2) failed")
	}
	s, _ = UnpackSpacing(1, 2, 3, 4)
	if s != (Spacing{1, 2, 3, 4}) {
		t.Error("Unpack(1,2,3,4) failed")
	}
	_, err := UnpackSpacing()
	if err == nil {
		t.Error("Unpack() should error")
	}
	_, err = UnpackSpacing(1, 2, 3)
	if err == nil {
		t.Error("Unpack(3 values) should error")
	}
}

func TestSpacingAddSub(t *testing.T) {
	got := Spacing{1, 2, 3, 4}.Add(Spacing{5, 6, 7, 8})
	if got != (Spacing{6, 8, 10, 12}) {
		t.Errorf("Add = %v", got)
	}
	got = Spacing{1, 2, 3, 4}.Sub(Spacing{5, 6, 7, 8})
	if got != (Spacing{-4, -4, -4, -4}) {
		t.Errorf("Sub = %v", got)
	}
}

func TestSpacingConvenienceConstructors(t *testing.T) {
	if SpacingVertical(2) != (Spacing{2, 0, 2, 0}) {
		t.Error("SpacingVertical failed")
	}
	if SpacingHorizontal(2) != (Spacing{0, 2, 0, 2}) {
		t.Error("SpacingHorizontal failed")
	}
	if SpacingAll(2) != (Spacing{2, 2, 2, 2}) {
		t.Error("SpacingAll failed")
	}
}

func TestSpacingGrowMaximum(t *testing.T) {
	a := Spacing{1, 6, 3, 2}
	b := Spacing{4, 2, 1, 5}
	got := a.GrowMaximum(b)
	if got != (Spacing{4, 6, 3, 5}) {
		t.Errorf("GrowMaximum = %v, want {4,6,3,5}", got)
	}
}
