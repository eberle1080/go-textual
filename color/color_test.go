package color

import (
	"math"
	"testing"

	rich "github.com/eberle1080/go-rich"
)

func approxEq(a, b, tol float64) bool {
	return math.Abs(a-b) <= tol
}

func TestColorNew(t *testing.T) {
	c := New(10, 20, 30)
	if c.R != 10 || c.G != 20 || c.B != 30 || c.A != 1.0 {
		t.Errorf("New = %+v", c)
	}
}

func TestColorNormalized(t *testing.T) {
	r, g, b := Color{255, 128, 64, 1, nil, false}.Normalized()
	if !approxEq(r, 1.0, 0.001) {
		t.Errorf("R normalized = %f", r)
	}
	if !approxEq(g, 128.0/255.0, 0.001) {
		t.Errorf("G normalized = %f", g)
	}
	if !approxEq(b, 64.0/255.0, 0.001) {
		t.Errorf("B normalized = %f", b)
	}
}

func TestColorClamped(t *testing.T) {
	c := Color{300, 100, -20, 1.5, nil, false}.Clamped()
	want := Color{255, 100, 0, 1.0, nil, false}
	if c != want {
		t.Errorf("Clamped = %+v, want %+v", c, want)
	}
}

func TestColorCSS(t *testing.T) {
	_tmp1 := Color{10, 20, 30, 1.0, nil, false}
	if _tmp1.CSS() != "rgb(10,20,30)" {
		t.Error("CSS rgb failed")
	}
	_tmp2 := Color{10, 20, 30, 0.5, nil, false}
	if _tmp2.CSS() != "rgba(10,20,30,0.5)" {
		t.Error("CSS rgba failed")
	}
	ansi1 := 1
	if (Color{0, 0, 0, 0, &ansi1, false}).CSS() != "ansi_red" {
		t.Error("CSS ansi_red failed")
	}
	autoC := Color{10, 20, 30, 0.5, nil, true}.CSS()
	if autoC != "auto 50%" {
		t.Errorf("CSS auto 50%% = %q", autoC)
	}
	auto2 := Automatic(70.5)
	if auto2.CSS() != "auto 70.5%" {
		t.Errorf("Automatic(70.5).CSS() = %q", auto2.CSS())
	}
}

func TestColorMonochrome(t *testing.T) {
	tests := []struct {
		in   Color
		want Color
	}{
		{Color{10, 20, 30, 1, nil, false}, Color{19, 19, 19, 1, nil, false}},
		{Color{10, 20, 30, 0.5, nil, false}, Color{19, 19, 19, 0.5, nil, false}},
		{Color{255, 255, 255, 1, nil, false}, Color{255, 255, 255, 1, nil, false}},
		{Color{0, 0, 0, 1, nil, false}, Color{0, 0, 0, 1, nil, false}},
	}
	for _, tt := range tests {
		got := tt.in.Monochrome()
		if got != tt.want {
			t.Errorf("Monochrome(%+v) = %+v, want %+v", tt.in, got, tt.want)
		}
	}
}

func TestColorRGB(t *testing.T) {
	r, g, b := Color{10, 20, 30, 0.55, nil, false}.RGB()
	if r != 10 || g != 20 || b != 30 {
		t.Errorf("RGB = (%d,%d,%d)", r, g, b)
	}
}

func TestColorHSL(t *testing.T) {
	red := Color{200, 20, 32, 1, nil, false}
	hsl := red.HSL()
	if !approxEq(hsl.H, 0.9888888888888889, 0.001) {
		t.Errorf("HSL.H = %f", hsl.H)
	}
	if !approxEq(hsl.S, 0.818181818181818, 0.001) {
		t.Errorf("HSL.S = %f", hsl.S)
	}
	if !approxEq(hsl.L, 0.43137254901960786, 0.001) {
		t.Errorf("HSL.L = %f", hsl.L)
	}
	if hsl.CSS() != "hsl(356,81.8%,43.1%)" {
		t.Errorf("HSL.CSS() = %q", hsl.CSS())
	}
	// Round-trip
	back := FromHSL(hsl.H, hsl.S, hsl.L)
	r1, g1, b1 := red.Normalized()
	r2, g2, b2 := back.Normalized()
	if !approxEq(r1, r2, 0.01) || !approxEq(g1, g2, 0.01) || !approxEq(b1, b2, 0.01) {
		t.Errorf("HSL round-trip failed: got normalized (%f,%f,%f)", r2, g2, b2)
	}
}

func TestColorHSV(t *testing.T) {
	red := Color{200, 20, 32, 1, nil, false}
	hsv := red.HSV()
	if !approxEq(hsv.H, 0.9888888888888889, 0.001) {
		t.Errorf("HSV.H = %f", hsv.H)
	}
	if !approxEq(hsv.S, 0.8999999999999999, 0.001) {
		t.Errorf("HSV.S = %f", hsv.S)
	}
	if !approxEq(hsv.V, 0.7843137254901961, 0.001) {
		t.Errorf("HSV.V = %f", hsv.V)
	}
	back := FromHSV(hsv.H, hsv.S, hsv.V)
	r1, g1, b1 := red.Normalized()
	r2, g2, b2 := back.Normalized()
	if !approxEq(r1, r2, 0.01) || !approxEq(g1, g2, 0.01) || !approxEq(b1, b2, 0.01) {
		t.Errorf("HSV round-trip failed")
	}
}

func TestColorBrightness(t *testing.T) {
	_tmp3 := Color{255, 255, 255, 1, nil, false}
	if _tmp3.Brightness() != 1 {
		t.Error("white brightness should be 1")
	}
	_tmp4 := Color{0, 0, 0, 1, nil, false}
	if _tmp4.Brightness() != 0 {
		t.Error("black brightness should be 0")
	}
	if !approxEq(Color{127, 127, 127, 1, nil, false}.Brightness(), 0.49803921568627446, 0.001) {
		t.Error("gray brightness wrong")
	}
}

func TestColorHex(t *testing.T) {
	_tmp5 := Color{255, 0, 127, 1, nil, false}
	if _tmp5.Hex() != "#FF007F" {
		t.Error("Hex opaque failed")
	}
	_tmp6 := Color{255, 0, 127, 0.5, nil, false}
	if _tmp6.Hex() != "#FF007F7F" {
		t.Error("Hex with alpha failed")
	}
}

func TestColorHex6(t *testing.T) {
	_tmp7 := Color{0, 0, 0, 1, nil, false}
	if _tmp7.Hex6() != "#000000" {
		t.Error("Hex6 black failed")
	}
	_tmp8 := Color{255, 255, 255, 0.25, nil, false}
	if _tmp8.Hex6() != "#FFFFFF" {
		t.Error("Hex6 alpha ignored failed")
	}
}

func TestColorWithAlpha(t *testing.T) {
	c := Color{255, 50, 100, 1, nil, false}.WithAlpha(0.25)
	want := Color{255, 50, 100, 0.25, nil, false}
	if c != want {
		t.Errorf("WithAlpha = %+v", c)
	}
}

func TestMultiplyAlpha(t *testing.T) {
	c1 := Color{100, 100, 100, 1, nil, false}.MultiplyAlpha(0.5)
	if !approxEq(c1.A, 0.5, 0.001) {
		t.Errorf("MultiplyAlpha(0.5) = %f", c1.A)
	}
	c2 := Color{100, 100, 100, 0.5, nil, false}.MultiplyAlpha(0.5)
	if !approxEq(c2.A, 0.25, 0.001) {
		t.Errorf("MultiplyAlpha(0.5) from 0.5 = %f", c2.A)
	}
}

func TestColorBlend(t *testing.T) {
	black := Color{0, 0, 0, 1, nil, false}
	white := Color{255, 255, 255, 1, nil, false}
	got0 := black.Blend(white, 0, nil)
	if got0 != black {
		t.Error("Blend(0) should return self")
	}
	got1 := black.Blend(white, 1, nil)
	if got1 != white {
		t.Error("Blend(1) should return dest")
	}
	got05 := black.Blend(white, 0.5, nil)
	if got05.R != 127 || got05.G != 127 || got05.B != 127 {
		t.Errorf("Blend(0.5) = %+v", got05)
	}
}

func TestColorParse(t *testing.T) {
	tests := []struct {
		text    string
		r, g, b int
		a       float64
	}{
		{"#000000", 0, 0, 0, 1.0},
		{"#ffffff", 255, 255, 255, 1.0},
		{"#FFFFFF", 255, 255, 255, 1.0},
		{"#fab", 255, 170, 187, 1.0},
		{"#fab0", 255, 170, 187, 0.0},
		{"#020304ff", 2, 3, 4, 1.0},
		{"#02030400", 2, 3, 4, 0.0},
		{"rgb(0,0,0)", 0, 0, 0, 1.0},
		{"rgb(255,255,255)", 255, 255, 255, 1.0},
		{"rgba(255,255,255,1)", 255, 255, 255, 1.0},
		{"rgb(2,3,4)", 2, 3, 4, 1.0},
		{"rgba(2,3,4,1.0)", 2, 3, 4, 1.0},
		{"hsl(45,25%,25%)", 80, 72, 48, 1.0},
		{"hsla(45,25%,25%,0.35)", 80, 72, 48, 0.35},
	}
	for _, tt := range tests {
		got, err := Parse(tt.text)
		if err != nil {
			t.Errorf("Parse(%q) error: %v", tt.text, err)
			continue
		}
		if got.R != tt.r || got.G != tt.g || got.B != tt.b || !approxEq(got.A, tt.a, 0.001) {
			t.Errorf("Parse(%q) = (%d,%d,%d,%.3f), want (%d,%d,%d,%.3f)",
				tt.text, got.R, got.G, got.B, got.A, tt.r, tt.g, tt.b, tt.a)
		}
	}
}

func TestColorParseClamp(t *testing.T) {
	got, err := Parse("rgb(300, 300, 300)")
	if err != nil {
		t.Fatal(err)
	}
	if got.R != 255 || got.G != 255 || got.B != 255 {
		t.Errorf("clamped rgb = %+v", got)
	}
}

func TestColorParseHSLNegativeDegrees(t *testing.T) {
	c1, _ := Parse("hsl(-90, 50%, 50%)")
	c2, _ := Parse("hsl(270, 50%, 50%)")
	if c1.R != c2.R || c1.G != c2.G || c1.B != c2.B {
		t.Errorf("hsl(-90) != hsl(270): %+v vs %+v", c1, c2)
	}
}

func TestColorParseNamedColors(t *testing.T) {
	red, err := Parse("red")
	if err != nil {
		t.Fatal(err)
	}
	if red.R != 255 || red.G != 0 || red.B != 0 {
		t.Errorf("red = %+v", red)
	}
	transparent, _ := Parse("transparent")
	if !transparent.IsTransparent() {
		t.Error("transparent should be transparent")
	}
}

func TestColorParseError(t *testing.T) {
	_, err := Parse("notacolor!!!")
	if err == nil {
		t.Error("expected parse error")
	}
}

func TestColorParseRejectsGarbagePrefix(t *testing.T) {
	// Regex must be fully anchored — a valid-looking suffix inside garbage must
	// not match.
	bad := []string{
		"x#ff0000",
		"junk rgb(1,2,3)",
		"  #abc",
		"#abcXYZ",
		"rgb(1.2.3,0,0)",
		"rgba(0,0,0,1.2.3)",
	}
	for _, s := range bad {
		_, err := Parse(s)
		if err == nil {
			t.Errorf("Parse(%q) should have returned an error", s)
		}
	}
}

func TestColorParseRejectsMalformedHSLPercentages(t *testing.T) {
	// Malformed percentage tokens must be rejected, not silently coerced.
	bad := []string{
		"hsl(0,1.2.3%,50%)",
		"hsl(0,50%,1.2.3%)",
		"hsla(0,1.2.3%,50%,1)",
		"hsla(0,50%,1.2.3%,1)",
	}
	for _, s := range bad {
		_, err := Parse(s)
		if err == nil {
			t.Errorf("Parse(%q) should have returned an error", s)
		}
	}
}

func TestColorParseANSIDefault(t *testing.T) {
	c, err := Parse("ansi_default")
	if err != nil {
		t.Fatalf("Parse(\"ansi_default\") error: %v", err)
	}
	if c.A != 1.0 {
		t.Errorf("ansi_default A = %f, want 1.0", c.A)
	}
	if c.ANSI == nil || *c.ANSI != -1 {
		t.Errorf("ansi_default ANSI = %v, want -1", c.ANSI)
	}
	// Must be opaque like all other ANSI colors produced by Parse.
	ansiRed, _ := Parse("ansi_red")
	if c.A != ansiRed.A {
		t.Errorf("ansi_default A (%f) != ansi_red A (%f)", c.A, ansiRed.A)
	}
}

func TestColorDarken(t *testing.T) {
	c := Color{200, 210, 220, 1, nil, false}
	if c.Darken(1, nil).R > 5 || c.Darken(1, nil).G > 5 || c.Darken(1, nil).B > 5 {
		t.Error("Darken(1) should approach black")
	}
	if c.Darken(-1, nil).R < 250 {
		t.Error("Darken(-1) should approach white")
	}
	d := c.Darken(0.1, nil)
	if d.R != 172 || d.G != 182 || d.B != 192 {
		t.Errorf("Darken(0.1) = (%d,%d,%d), want (172,182,192)", d.R, d.G, d.B)
	}
}

func TestColorLighten(t *testing.T) {
	c := Color{200, 210, 220, 1, nil, false}
	if c.Lighten(1, nil).R < 250 {
		t.Error("Lighten(1) should approach white")
	}
	l := c.Lighten(0.1, nil)
	if l.R != 228 || l.G != 238 || l.B != 248 {
		t.Errorf("Lighten(0.1) = (%d,%d,%d), want (228,238,248)", l.R, l.G, l.B)
	}
}

func TestRGBToLab(t *testing.T) {
	data := []struct {
		r, g, b  int
		L, a, b_ float64
	}{
		{10, 23, 73, 10.245, 15.913, -32.672},
		{200, 34, 123, 45.438, 67.750, -8.008},
		{0, 0, 0, 0, 0, 0},
		{255, 255, 255, 100, 0, 0},
	}
	for _, tt := range data {
		rgb := Color{tt.r, tt.g, tt.b, 1, nil, false}
		lab := RGBToLab(rgb)
		if !approxEq(lab.L, tt.L, 0.1) {
			t.Errorf("RGBToLab(%d,%d,%d) L = %f, want %f", tt.r, tt.g, tt.b, lab.L, tt.L)
		}
		if !approxEq(lab.A, tt.a, 0.1) {
			t.Errorf("RGBToLab(%d,%d,%d) a = %f, want %f", tt.r, tt.g, tt.b, lab.A, tt.a)
		}
		if !approxEq(lab.B, tt.b_, 0.1) {
			t.Errorf("RGBToLab(%d,%d,%d) b = %f, want %f", tt.r, tt.g, tt.b, lab.B, tt.b_)
		}
	}
}

func TestLabToRGB(t *testing.T) {
	data := []struct {
		r, g, b  int
		L, a, b_ float64
	}{
		{10, 23, 73, 10.245, 15.913, -32.672},
		{0, 0, 0, 0, 0, 0},
		{255, 255, 255, 100, 0, 0},
	}
	for _, tt := range data {
		lab := Lab{tt.L, tt.a, tt.b_}
		rgb := LabToRGB(lab, 1.0)
		if math.Abs(float64(rgb.R-tt.r)) > 1 || math.Abs(float64(rgb.G-tt.g)) > 1 || math.Abs(float64(rgb.B-tt.b)) > 1 {
			t.Errorf("LabToRGB(%v) = (%d,%d,%d), want (%d,%d,%d)",
				lab, rgb.R, rgb.G, rgb.B, tt.r, tt.g, tt.b)
		}
	}
}

func TestRGBLabRoundtrip(t *testing.T) {
	for r := 0; r <= 255; r += 32 {
		for g := 0; g <= 255; g += 32 {
			for b := 0; b <= 255; b += 32 {
				c := Color{r, g, b, 1, nil, false}
				back := LabToRGB(RGBToLab(c), 1.0)
				if math.Abs(float64(back.R-r)) > 1 || math.Abs(float64(back.G-g)) > 1 || math.Abs(float64(back.B-b)) > 1 {
					t.Errorf("roundtrip (%d,%d,%d) -> (%d,%d,%d)", r, g, b, back.R, back.G, back.B)
				}
			}
		}
	}
}

func TestColorInverse(t *testing.T) {
	c := Color{55, 0, 255, 0.1, nil, false}.Inverse()
	if c.R != 200 || c.G != 255 || c.B != 0 || !approxEq(c.A, 0.1, 0.001) {
		t.Errorf("Inverse = %+v", c)
	}
}

func TestIsTransparent(t *testing.T) {
	_tmp9 := Color{0, 0, 0, 0, nil, false}
	if !_tmp9.IsTransparent() {
		t.Error("alpha=0 should be transparent")
	}
	_tmp10 := Color{20, 20, 30, 0.01, nil, false}
	if _tmp10.IsTransparent() {
		t.Error("alpha=0.01 should not be transparent")
	}
	ansi1 := 1
	if (Color{20, 20, 30, 0, &ansi1, false}).IsTransparent() {
		t.Error("ANSI color should not be transparent")
	}
}

func TestColorTint(t *testing.T) {
	tests := []struct {
		base, tint, want Color
	}{
		{
			Color{0, 0, 0, 1, nil, false},
			Color{10, 20, 30, 1, nil, false},
			Color{10, 20, 30, 1, nil, false},
		},
		{
			Color{0, 0, 0, 0.5, nil, false},
			Color{255, 255, 255, 0.5, nil, false},
			Color{127, 127, 127, 0.5, nil, false},
		},
	}
	for _, tt := range tests {
		got := tt.base.Tint(tt.tint)
		if got.R != tt.want.R || got.G != tt.want.G || got.B != tt.want.B || !approxEq(got.A, tt.want.A, 0.001) {
			t.Errorf("Tint = %+v, want %+v", got, tt.want)
		}
	}
}

func TestColorAdd(t *testing.T) {
	c1 := Color{1, 2, 3, 1, nil, false}
	c2 := Color{20, 30, 40, 1, nil, false}
	got := c1.Add(c2)
	if got.R != 20 || got.G != 30 || got.B != 40 {
		t.Errorf("Add full alpha = %+v", got)
	}
}

func TestGradientErrors(t *testing.T) {
	_, err := NewGradient(nil, 50)
	if err == nil {
		t.Error("NewGradient(nil) should error")
	}
	_, err = NewGradient([]GradientStop{{0.1, New(255, 0, 0)}, {1, New(0, 0, 255)}}, 50)
	if err == nil {
		t.Error("gradient not starting at 0 should error")
	}
	_, err = NewGradient([]GradientStop{{0, New(255, 0, 0)}, {0.8, New(0, 0, 255)}}, 50)
	if err == nil {
		t.Error("gradient not ending at 1 should error")
	}
	_, err = GradientFromColors([]Color{New(255, 0, 0)}, 50)
	if err == nil {
		t.Error("GradientFromColors with one color should error")
	}
}

func TestGradient(t *testing.T) {
	g, err := NewGradient([]GradientStop{
		{0, Color{255, 0, 0, 1, nil, false}},
		{0.5, Color{0, 0, 255, 1, nil, false}},
		{1, Color{0, 255, 0, 1, nil, false}},
	}, 11)
	if err != nil {
		t.Fatal(err)
	}
	start := g.GetColor(-1)
	if start.R != 255 || start.G != 0 || start.B != 0 {
		t.Errorf("GetColor(-1) = %+v", start)
	}
	end := g.GetColor(1)
	if end.G != 255 || end.R != 0 || end.B != 0 {
		t.Errorf("GetColor(1) = %+v", end)
	}
	mid := g.GetColor(0.5)
	if mid.B != 255 {
		t.Errorf("GetColor(0.5) = %+v, want blue", mid)
	}
}

func TestContrastText(t *testing.T) {
	light := Color{200, 200, 200, 1, nil, false}
	dark := Color{20, 20, 20, 1, nil, false}
	lt := light.ContrastText(0.95)
	if lt.R != 0 {
		t.Errorf("light background contrast text should be black, got %+v", lt)
	}
	dt := dark.ContrastText(0.95)
	if dt.R != 255 {
		t.Errorf("dark background contrast text should be white, got %+v", dt)
	}
}

func TestToRichRGB(t *testing.T) {
	c := Color{10, 20, 30, 1, nil, false}
	got := c.ToRichRGB()
	want := rich.RGBColor{R: 10, G: 20, B: 30}
	if got != want {
		t.Errorf("ToRichRGB = %+v, want %+v", got, want)
	}
}

func TestFromRichRGB(t *testing.T) {
	rc := rich.RGBColor{R: 10, G: 20, B: 30}
	got := FromRichRGB(rc)
	if got.R != 10 || got.G != 20 || got.B != 30 || got.A != 1.0 {
		t.Errorf("FromRichRGB = %+v", got)
	}
}
