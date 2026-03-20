// Package color provides a Color type with RGBA support, CSS parsing, color
// blending, HSL/HSV conversion, and CIE-L*ab conversion utilities.
//
// # Quick start
//
//	c := color.New(255, 0, 0)                 // red
//	c2, _ := color.Parse("#00ff00")           // green from hex
//	blended := c.Blend(c2, 0.5, nil)          // midpoint color
//	darker := c.Darken(0.2, nil)              // 20% darker via Lab
package color
