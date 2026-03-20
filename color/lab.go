package color

import "math"

// RGBToLab converts an RGB color to the CIE-L*ab format.
// Uses the standard RGB color space with a D65/2° standard illuminant.
// Conversion passes through the XYZ color space.
// Cf. http://www.easyrgb.com/en/math.php
func RGBToLab(c Color) Lab {
	r := float64(c.R) / 255.0
	g := float64(c.G) / 255.0
	b := float64(c.B) / 255.0

	if r > 0.04045 {
		r = math.Pow((r+0.055)/1.055, 2.4)
	} else {
		r = r / 12.92
	}
	if g > 0.04045 {
		g = math.Pow((g+0.055)/1.055, 2.4)
	} else {
		g = g / 12.92
	}
	if b > 0.04045 {
		b = math.Pow((b+0.055)/1.055, 2.4)
	} else {
		b = b / 12.92
	}

	x := (r*41.24 + g*35.76 + b*18.05) / 95.047
	y := (r*21.26 + g*71.52 + b*7.22) / 100
	z := (r*1.93 + g*11.92 + b*95.05) / 108.883

	off := 16.0 / 116.0
	if x > 0.008856 {
		x = math.Pow(x, 1.0/3.0)
	} else {
		x = 7.787*x + off
	}
	if y > 0.008856 {
		y = math.Pow(y, 1.0/3.0)
	} else {
		y = 7.787*y + off
	}
	if z > 0.008856 {
		z = math.Pow(z, 1.0/3.0)
	} else {
		z = 7.787*z + off
	}

	return Lab{116*y - 16, 500 * (x - y), 200 * (y - z)}
}

// LabToRGB converts a CIE-L*ab color to RGB with the given alpha.
// Uses the standard RGB color space with a D65/2° standard illuminant.
// Conversion passes through the XYZ color space.
// Cf. http://www.easyrgb.com/en/math.php
func LabToRGB(lab Lab, alpha float64) Color {
	y := (lab.L + 16) / 116
	x := lab.A/500 + y
	z := y - lab.B/200

	off := 16.0 / 116.0
	if y > 0.2068930344 {
		y = math.Pow(y, 3)
	} else {
		y = (y - off) / 7.787
	}
	if x > 0.2068930344 {
		x = 0.95047 * math.Pow(x, 3)
	} else {
		x = 0.122059 * (x - off)
	}
	if z > 0.2068930344 {
		z = 1.08883 * math.Pow(z, 3)
	} else {
		z = 0.139827 * (z - off)
	}

	r := x*3.2406 + y*-1.5372 + z*-0.4986
	g := x*-0.9689 + y*1.8758 + z*0.0415
	b := x*0.0557 + y*-0.2040 + z*1.0570

	if r > 0.0031308 {
		r = 1.055*math.Pow(r, 1.0/2.4) - 0.055
	} else {
		r = 12.92 * r
	}
	if g > 0.0031308 {
		g = 1.055*math.Pow(g, 1.0/2.4) - 0.055
	} else {
		g = 12.92 * g
	}
	if b > 0.0031308 {
		b = 1.055*math.Pow(b, 1.0/2.4) - 0.055
	} else {
		b = 12.92 * b
	}

	return Color{int(r * 255), int(g * 255), int(b * 255), alpha, nil, false}
}
