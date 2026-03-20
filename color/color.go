package color

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	rich "github.com/eberle1080/go-rich"
)

// HSL represents a color in HSL (Hue, Saturation, Lightness) format.
type HSL struct {
	H float64 // Hue in range 0 to 1
	S float64 // Saturation in range 0 to 1
	L float64 // Lightness in range 0 to 1
}

// CSS returns the color in CSS hsl() format.
func (h HSL) CSS() string {
	return fmt.Sprintf("hsl(%s,%s%%,%s%%)",
		formatFloat(h.H*360),
		formatFloat(h.S*100),
		formatFloat(h.L*100),
	)
}

// formatFloat formats a float for CSS output: strip trailing zeros after decimal.
func formatFloat(f float64) string {
	s := strconv.FormatFloat(f, 'f', 1, 64)
	s = strings.TrimRight(s, "0")
	s = strings.TrimRight(s, ".")
	return s
}

// HSV represents a color in HSV (Hue, Saturation, Value) format.
type HSV struct {
	H float64 // Hue in range 0 to 1
	S float64 // Saturation in range 0 to 1
	V float64 // Value in range 0 to 1
}

// Lab represents a color in CIE-L*ab format.
type Lab struct {
	L float64 // Lightness in range 0 to 100
	A float64 // A axis in range -127 to 128
	B float64 // B axis in range -127 to 128
}

// Color represents an RGBA color.
// R, G, B are in the range 0-255; A is in the range 0.0-1.0.
// ANSI holds the ANSI color index (-1 = default, nil = not ANSI).
// Auto indicates the color is automatic (white or black for maximum contrast).
type Color struct {
	R    int
	G    int
	B    int
	A    float64
	ANSI *int
	Auto bool
}

// Package-level color constants.
var (
	White       = New(255, 255, 255)
	Black       = New(0, 0, 0)
	Transparent = NewWithAlpha(0, 0, 0, 0)
)

// New creates an opaque color with the given RGB values.
func New(r, g, b int) Color {
	return Color{R: r, G: g, B: b, A: 1.0}
}

// NewWithAlpha creates a color with the given RGB and alpha values.
func NewWithAlpha(r, g, b int, a float64) Color {
	return Color{R: r, G: g, B: b, A: a}
}

// Automatic creates an automatic color (selects white or black for contrast).
// alphaPercentage is 0-100; 100 means fully opaque.
func Automatic(alphaPercentage float64) Color {
	return Color{A: alphaPercentage / 100.0, Auto: true}
}

// IsTransparent reports whether the color is transparent (A==0 and not ANSI).
func (c Color) IsTransparent() bool {
	return c.A == 0 && c.ANSI == nil
}

// Clamped returns the color with all values restricted to their valid ranges.
func (c Color) Clamped() Color {
	r := clampInt(c.R, 0, 255)
	g := clampInt(c.G, 0, 255)
	b := clampInt(c.B, 0, 255)
	a := clampFloat(c.A, 0.0, 1.0)
	return Color{r, g, b, a, c.ANSI, c.Auto}
}

func clampInt(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func clampFloat(v, lo, hi float64) float64 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

// Inverse returns the inverse of this color (255-component for each RGB channel).
func (c Color) Inverse() Color {
	return Color{255 - c.R, 255 - c.G, 255 - c.B, c.A, nil, false}
}

// Normalized returns (R, G, B) each divided by 255.
func (c Color) Normalized() (float64, float64, float64) {
	return float64(c.R) / 255, float64(c.G) / 255, float64(c.B) / 255
}

// RGB returns the (R, G, B) components as integers.
func (c Color) RGB() (int, int, int) {
	return c.R, c.G, c.B
}

// HSL converts this color to HSL format.
func (c Color) HSL() HSL {
	r, g, b := c.Normalized()
	h, l, s := rgbToHLS(r, g, b)
	return HSL{h, s, l}
}

// HSV converts this color to HSV format.
func (c Color) HSV() HSV {
	r, g, b := c.Normalized()
	h, s, v := rgbToHSV(r, g, b)
	return HSV{h, s, v}
}

// Brightness returns the perceptual brightness (0=black, 1=white).
func (c Color) Brightness() float64 {
	r, g, b := c.Normalized()
	return (299*r + 587*g + 114*b) / 1000
}

// Hex returns the color as a CSS hex string (#RRGGBB or #RRGGBBAA).
func (c Color) Hex() string {
	cl := c.Clamped()
	if cl.ANSI != nil {
		if *cl.ANSI == -1 {
			return "ansi_default"
		}
		return "ansi_" + ANSIColors[*cl.ANSI]
	}
	if cl.A == 1 {
		return fmt.Sprintf("#%02X%02X%02X", cl.R, cl.G, cl.B)
	}
	return fmt.Sprintf("#%02X%02X%02X%02X", cl.R, cl.G, cl.B, int(cl.A*255))
}

// Hex6 returns the color as a 6-digit CSS hex string (alpha ignored).
func (c Color) Hex6() string {
	cl := c.Clamped()
	return fmt.Sprintf("#%02X%02X%02X", cl.R, cl.G, cl.B)
}

// CSS returns the color in CSS rgb(), rgba(), ansi_, or auto format.
func (c Color) CSS() string {
	if c.Auto {
		alphaPct := clampFloat(c.A, 0, 1) * 100.0
		if alphaPct == 100 {
			return "auto"
		}
		if math.Mod(alphaPct, 1) == 0 {
			return fmt.Sprintf("auto %d%%", int(alphaPct))
		}
		return fmt.Sprintf("auto %.1f%%", alphaPct)
	}
	if c.ANSI != nil {
		if *c.ANSI == -1 {
			return "ansi_default"
		}
		return "ansi_" + ANSIColors[*c.ANSI]
	}
	if c.A == 1 {
		return fmt.Sprintf("rgb(%d,%d,%d)", c.R, c.G, c.B)
	}
	return fmt.Sprintf("rgba(%d,%d,%d,%g)", c.R, c.G, c.B, c.A)
}

// Monochrome returns the luminance-weighted grayscale version of this color.
func (c Color) Monochrome() Color {
	gray := int(math.Round(float64(c.R)*0.2126 + float64(c.G)*0.7152 + float64(c.B)*0.0722))
	return Color{gray, gray, gray, c.A, nil, false}
}

// WithAlpha returns a copy of the color with the given alpha.
func (c Color) WithAlpha(a float64) Color {
	return Color{c.R, c.G, c.B, a, nil, false}
}

// MultiplyAlpha returns a copy with alpha multiplied by the given factor.
// Returns self unchanged if ANSI color.
func (c Color) MultiplyAlpha(a float64) Color {
	if c.ANSI != nil {
		return c
	}
	return Color{c.R, c.G, c.B, c.A * a, nil, c.Auto}
}

// Blend generates a new color between c and dest at the given factor (0=c, 1=dest).
// If alpha is non-nil, it overrides the blended alpha.
func (c Color) Blend(dest Color, factor float64, alpha *float64) Color {
	if dest.Auto {
		dest = c.ContrastText(dest.A)
	}
	if dest.ANSI != nil {
		return dest
	}
	if factor <= 0 {
		return c
	}
	if factor >= 1 {
		return dest
	}
	r1, g1, b1, a1 := float64(c.R), float64(c.G), float64(c.B), c.A
	r2, g2, b2, a2 := float64(dest.R), float64(dest.G), float64(dest.B), dest.A
	var newAlpha float64
	if alpha == nil {
		newAlpha = a1 + (a2-a1)*factor
	} else {
		newAlpha = *alpha
	}
	return Color{
		int(r1 + (r2-r1)*factor),
		int(g1 + (g2-g1)*factor),
		int(b1 + (b2-b1)*factor),
		newAlpha,
		nil,
		false,
	}
}

// Add blends c towards other using other's alpha, producing an opaque result.
func (c Color) Add(other Color) Color {
	a := 1.0
	return c.Blend(other, other.A, &a)
}

// Tint applies a tint (blend with alpha component) to the color.
func (c Color) Tint(other Color) Color {
	if c.ANSI != nil || other.ANSI != nil {
		return c
	}
	a2 := other.A
	return Color{
		int(float64(c.R) + float64(other.R-c.R)*a2),
		int(float64(c.G) + float64(other.G-c.G)*a2),
		int(float64(c.B) + float64(other.B-c.B)*a2),
		c.A,
		nil,
		false,
	}
}

// Darken reduces the luminance by amount (0-1) using Lab color space.
// If alpha is non-nil it overrides the result's alpha.
func (c Color) Darken(amount float64, alpha *float64) Color {
	lab := RGBToLab(c)
	lab.L -= amount * 100
	a := c.A
	if alpha != nil {
		a = *alpha
	}
	return LabToRGB(lab, a).Clamped()
}

// Lighten increases the luminance by amount (0-1).
func (c Color) Lighten(amount float64, alpha *float64) Color {
	return c.Darken(-amount, alpha)
}

// ContrastText returns white or black (with the given alpha) for maximum contrast.
func (c Color) ContrastText(alpha float64) Color {
	if c.Brightness() < 0.5 {
		return White.WithAlpha(alpha)
	}
	return Black.WithAlpha(alpha)
}

// ToRichRGB converts this Color to go-rich's RGBColor.
func (c Color) ToRichRGB() rich.RGBColor {
	cl := c.Clamped()
	return rich.RGBColor{R: uint8(cl.R), G: uint8(cl.G), B: uint8(cl.B)}
}

// FromRichRGB creates a Color from go-rich's RGBColor.
func FromRichRGB(rc rich.RGBColor) Color {
	return New(int(rc.R), int(rc.G), int(rc.B))
}

// FromHSL creates a color from HSL components.
func FromHSL(h, s, l float64) Color {
	r, g, b := hlsToRGB(h, l, s)
	return Color{int(r*255 + 0.5), int(g*255 + 0.5), int(b*255 + 0.5), 1.0, nil, false}
}

// FromHSV creates a color from HSV components.
func FromHSV(h, s, v float64) Color {
	r, g, b := hsvToRGB(h, s, v)
	return Color{int(r*255 + 0.5), int(g*255 + 0.5), int(b*255 + 0.5), 1.0, nil, false}
}

// ColorParseError is returned when a color string cannot be parsed.
type ColorParseError struct {
	Message        string
	SuggestedColor string
}

func (e *ColorParseError) Error() string { return e.Message }

// CSS color regex — matches hex, rgb, rgba, hsl, hsla forms.
// Groups: rgb3, rgb4, rgb6, rgb8, rgb, rgba, hsl, hsla
// All alternatives are wrapped in a single ^(?:...)$ so every branch is
// fully anchored; without the outer grouping only the first branch was
// anchored to the start of the string.
var reColor = regexp.MustCompile(`(?i)^(?:` +
	`#([0-9a-fA-F]{3})|` +
	`#([0-9a-fA-F]{4})|` +
	`#([0-9a-fA-F]{6})|` +
	`#([0-9a-fA-F]{8})|` +
	`rgb\(\s*([\d.]+\s*,\s*[\d.]+\s*,\s*[\d.]+)\s*\)|` +
	`rgba\(\s*([\d.]+\s*,\s*[\d.]+\s*,\s*[\d.]+\s*,\s*[\d.]+)\s*\)|` +
	`hsl\(\s*([\d.+-]+\s*,\s*\d+(?:\.\d+)?%\s*,\s*\d+(?:\.\d+)?%)\s*\)|` +
	`hsla\(\s*([\d.+-]+\s*,\s*\d+(?:\.\d+)?%\s*,\s*\d+(?:\.\d+)?%\s*,\s*[\d.]+)\s*\)` +
	`)$`)

// Parse parses a CSS color string into a Color.
// Supports: named colors, #RGB, #RGBA, #RRGGBB, #RRGGBBAA, rgb(), rgba(), hsl(), hsla(), ansi_*.
func Parse(s string) (Color, error) {
	if s == "ansi_default" {
		ansi := -1
		return Color{R: 0, G: 0, B: 0, A: 1.0, ANSI: &ansi}, nil
	}
	if strings.HasPrefix(s, "ansi_") {
		name := s[5:]
		for i, n := range ANSIColors {
			if n == name {
				if rgb, ok := ColorNameToRGB["ansi_"+name]; ok {
					ansi := i
					return Color{rgb[0], rgb[1], rgb[2], 1.0, &ansi, false}, nil
				}
			}
		}
	}
	if rgb, ok := ColorNameToRGB[strings.ToLower(s)]; ok {
		return Color{rgb[0], rgb[1], rgb[2], float64(rgb[3]) / 255.0, nil, false}, nil
	}

	m := reColor.FindStringSubmatch(s)
	if m == nil {
		return Color{}, &ColorParseError{Message: fmt.Sprintf("failed to parse %q as a color", s)}
	}
	// Groups: m[1]=rgb3 m[2]=rgb4 m[3]=rgb6 m[4]=rgb8 m[5]=rgb m[6]=rgba m[7]=hsl m[8]=hsla
	parseErr := func() error {
		return &ColorParseError{Message: fmt.Sprintf("failed to parse %q as a color", s)}
	}
	switch {
	case m[1] != "":
		hex := m[1]
		r, err1 := hexByte(hex[0:1] + hex[0:1])
		g, err2 := hexByte(hex[1:2] + hex[1:2])
		b, err3 := hexByte(hex[2:3] + hex[2:3])
		if err1 != nil || err2 != nil || err3 != nil {
			return Color{}, parseErr()
		}
		return Color{r, g, b, 1.0, nil, false}, nil
	case m[2] != "":
		hex := m[2]
		r, err1 := hexByte(hex[0:1] + hex[0:1])
		g, err2 := hexByte(hex[1:2] + hex[1:2])
		b, err3 := hexByte(hex[2:3] + hex[2:3])
		av, err4 := hexByte(hex[3:4] + hex[3:4])
		if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
			return Color{}, parseErr()
		}
		return Color{r, g, b, float64(av) / 255.0, nil, false}, nil
	case m[3] != "":
		hex := m[3]
		r, err1 := hexByte(hex[0:2])
		g, err2 := hexByte(hex[2:4])
		b, err3 := hexByte(hex[4:6])
		if err1 != nil || err2 != nil || err3 != nil {
			return Color{}, parseErr()
		}
		return Color{r, g, b, 1.0, nil, false}, nil
	case m[4] != "":
		hex := m[4]
		r, err1 := hexByte(hex[0:2])
		g, err2 := hexByte(hex[2:4])
		b, err3 := hexByte(hex[4:6])
		av, err4 := hexByte(hex[6:8])
		if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
			return Color{}, parseErr()
		}
		return Color{r, g, b, float64(av) / 255.0, nil, false}, nil
	case m[5] != "":
		parts := splitTrim(m[5], ",")
		rf, err1 := parseFloat(parts[0])
		gf, err2 := parseFloat(parts[1])
		bf, err3 := parseFloat(parts[2])
		if err1 != nil || err2 != nil || err3 != nil {
			return Color{}, parseErr()
		}
		r := clampInt(int(rf), 0, 255)
		g := clampInt(int(gf), 0, 255)
		b := clampInt(int(bf), 0, 255)
		return Color{r, g, b, 1.0, nil, false}, nil
	case m[6] != "":
		parts := splitTrim(m[6], ",")
		rf, err1 := parseFloat(parts[0])
		gf, err2 := parseFloat(parts[1])
		bf, err3 := parseFloat(parts[2])
		af, err4 := parseFloat(parts[3])
		if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
			return Color{}, parseErr()
		}
		r := clampInt(int(rf), 0, 255)
		g := clampInt(int(gf), 0, 255)
		b := clampInt(int(bf), 0, 255)
		a := clampFloat(af, 0.0, 1.0)
		return Color{r, g, b, a, nil, false}, nil
	case m[7] != "":
		parts := splitTrim(m[7], ",")
		hf, err := parseFloat(parts[0])
		if err != nil {
			return Color{}, parseErr()
		}
		h := math.Mod(hf, 360)
		if h < 0 {
			h += 360
		}
		h /= 360
		sv, err1 := percentageToFloat(parts[1])
		lv, err2 := percentageToFloat(parts[2])
		if err1 != nil || err2 != nil {
			return Color{}, parseErr()
		}
		return FromHSL(h, sv, lv), nil
	case m[8] != "":
		parts := splitTrim(m[8], ",")
		hf, err1 := parseFloat(parts[0])
		af, err2 := parseFloat(parts[3])
		if err1 != nil || err2 != nil {
			return Color{}, parseErr()
		}
		h := math.Mod(hf, 360)
		if h < 0 {
			h += 360
		}
		h /= 360
		sv, err3 := percentageToFloat(parts[1])
		lv, err4 := percentageToFloat(parts[2])
		if err3 != nil || err4 != nil {
			return Color{}, parseErr()
		}
		a := clampFloat(af, 0.0, 1.0)
		return FromHSL(h, sv, lv).WithAlpha(a), nil
	}
	return Color{}, &ColorParseError{Message: fmt.Sprintf("failed to parse %q as a color", s)}
}

func hexByte(s string) (int, error) {
	n, err := strconv.ParseInt(s, 16, 64)
	return int(n), err
}

func parseFloat(s string) (float64, error) {
	return strconv.ParseFloat(strings.TrimSpace(s), 64)
}

func percentageToFloat(s string) (float64, error) {
	s = strings.TrimSpace(s)
	s = strings.TrimSuffix(s, "%")
	f, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
	if err != nil {
		return 0, err
	}
	return clampFloat(f/100.0, 0.0, 1.0), nil
}

func splitTrim(s, sep string) []string {
	parts := strings.Split(s, sep)
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

// GradientStop is a position+color pair for use in a Gradient.
type GradientStop struct {
	Position float64
	Color    Color
}

// Gradient defines a color gradient with blended stops.
type Gradient struct {
	stops   []GradientStop
	quality int
	colors  []Color
}

// NewGradient creates a Gradient from explicit stops.
// Requires at least 2 stops; first must be at position 0, last at position 1.
func NewGradient(stops []GradientStop, quality int) (*Gradient, error) {
	if len(stops) < 2 {
		return nil, fmt.Errorf("at least 2 stops required")
	}
	// Sort stops by position
	sorted := make([]GradientStop, len(stops))
	copy(sorted, stops)
	sortStops(sorted)
	if sorted[0].Position != 0.0 {
		return nil, fmt.Errorf("first stop must be 0")
	}
	if sorted[len(sorted)-1].Position != 1.0 {
		return nil, fmt.Errorf("last stop must be 1")
	}
	return &Gradient{stops: sorted, quality: quality}, nil
}

// GradientFromColors creates a Gradient from evenly spaced colors.
func GradientFromColors(colors []Color, quality int) (*Gradient, error) {
	if len(colors) < 2 {
		return nil, fmt.Errorf("two or more colors required")
	}
	stops := make([]GradientStop, len(colors))
	for i, c := range colors {
		stops[i] = GradientStop{float64(i) / float64(len(colors)-1), c}
	}
	return NewGradient(stops, quality)
}

func sortStops(stops []GradientStop) {
	// Simple insertion sort (small slices)
	for i := 1; i < len(stops); i++ {
		for j := i; j > 0 && stops[j].Position < stops[j-1].Position; j-- {
			stops[j], stops[j-1] = stops[j-1], stops[j]
		}
	}
}

// Colors returns the pre-computed color list for the gradient.
func (g *Gradient) Colors() []Color {
	if g.colors != nil {
		return g.colors
	}
	quality := g.quality
	colors := make([]Color, quality)
	position := 0
	stop1 := g.stops[0]
	stop2 := g.stops[1]
	for step := 0; step < quality; step++ {
		t := float64(step) / float64(quality-1)
		for t > stop2.Position && position+2 < len(g.stops) {
			position++
			stop1 = g.stops[position]
			stop2 = g.stops[position+1]
		}
		factor := (t - stop1.Position) / (stop2.Position - stop1.Position)
		colors[step] = stop1.Color.Blend(stop2.Color, factor, nil)
	}
	g.colors = colors
	return colors
}

// GetColor returns the color at the given position (0-1).
func (g *Gradient) GetColor(position float64) Color {
	cs := g.Colors()
	if position <= 0 {
		return cs[0]
	}
	if position >= 1 {
		return cs[len(cs)-1]
	}
	colorPos := position * float64(g.quality-1)
	idx := int(colorPos)
	c1 := cs[idx]
	c2 := cs[idx+1]
	return c1.Blend(c2, math.Mod(colorPos, 1), nil)
}

// --- colorsys port ---

func rgbToHLS(r, g, b float64) (h, l, s float64) {
	maxC := math.Max(r, math.Max(g, b))
	minC := math.Min(r, math.Min(g, b))
	l = (minC + maxC) / 2.0
	if minC == maxC {
		return 0.0, l, 0.0
	}
	if l <= 0.5 {
		s = (maxC - minC) / (maxC + minC)
	} else {
		s = (maxC - minC) / (2.0 - maxC - minC)
	}
	rc := (maxC - r) / (maxC - minC)
	gc := (maxC - g) / (maxC - minC)
	bc := (maxC - b) / (maxC - minC)
	switch {
	case r == maxC:
		h = bc - gc
	case g == maxC:
		h = 2.0 + rc - bc
	default:
		h = 4.0 + gc - rc
	}
	h = math.Mod(h/6.0, 1.0)
	if h < 0 {
		h += 1.0
	}
	return h, l, s
}

func hlsToRGB(h, l, s float64) (r, g, b float64) {
	if s == 0 {
		return l, l, l
	}
	var m2 float64
	if l <= 0.5 {
		m2 = l * (1.0 + s)
	} else {
		m2 = l + s - l*s
	}
	m1 := 2.0*l - m2
	r = hlsValue(m1, m2, h+1.0/3.0)
	g = hlsValue(m1, m2, h)
	b = hlsValue(m1, m2, h-1.0/3.0)
	return
}

func hlsValue(m1, m2, hue float64) float64 {
	hue = math.Mod(hue, 1.0)
	if hue < 0 {
		hue += 1.0
	}
	switch {
	case hue < 1.0/6.0:
		return m1 + (m2-m1)*hue*6.0
	case hue < 0.5:
		return m2
	case hue < 2.0/3.0:
		return m1 + (m2-m1)*(2.0/3.0-hue)*6.0
	default:
		return m1
	}
}

func rgbToHSV(r, g, b float64) (h, s, v float64) {
	maxC := math.Max(r, math.Max(g, b))
	minC := math.Min(r, math.Min(g, b))
	v = maxC
	if minC == maxC {
		return 0.0, 0.0, v
	}
	s = (maxC - minC) / maxC
	rc := (maxC - r) / (maxC - minC)
	gc := (maxC - g) / (maxC - minC)
	bc := (maxC - b) / (maxC - minC)
	switch {
	case r == maxC:
		h = bc - gc
	case g == maxC:
		h = 2.0 + rc - bc
	default:
		h = 4.0 + gc - rc
	}
	h = math.Mod(h/6.0, 1.0)
	if h < 0 {
		h += 1.0
	}
	return h, s, v
}

func hsvToRGB(h, s, v float64) (r, g, b float64) {
	if s == 0 {
		return v, v, v
	}
	i := math.Floor(h * 6)
	f := h*6 - i
	p := v * (1.0 - s)
	q := v * (1.0 - s*f)
	t := v * (1.0 - s*(1.0-f))
	switch int(i) % 6 {
	case 0:
		return v, t, p
	case 1:
		return q, v, p
	case 2:
		return p, v, t
	case 3:
		return p, q, v
	case 4:
		return t, p, v
	default:
		return v, p, q
	}
}
