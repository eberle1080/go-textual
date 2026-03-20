package app

import (
	rich "github.com/eberle1080/go-rich"

	"github.com/eberle1080/go-textual/geometry"
	"github.com/eberle1080/go-textual/strip"
	"github.com/eberle1080/go-textual/widget"
)

// renderIfDirty re-renders the screen if any widget is dirty.
func (a *App) renderIfDirty() {
	if a.screen == nil {
		return
	}
	if !isDirtyRecursive(a.screen) {
		return
	}

	region := geometry.Region{
		X:      0,
		Y:      0,
		Width:  a.size.Width,
		Height: a.size.Height,
	}

	// Record the screen's region so mouse hit-testing can find it.
	if rr, ok := a.screen.(interface{ SetRegion(geometry.Region) }); ok {
		rr.SetRegion(region)
	}
	strips := a.screen.Render(region)
	a.flushStrips(strips, region)
	clearDirtyRecursive(a.screen)
}

// flushStrips writes rendered strips to the terminal.
func (a *App) flushStrips(strips []strip.Strip, region geometry.Region) {
	if a.drv == nil || len(strips) == 0 {
		return
	}

	colorMode := rich.ColorMode256

	var out string
	for i, s := range strips {
		out += moveCursor(region.X, region.Y+i)
		out += s.Render(colorMode)
	}

	a.drv.Write(out)
}

// isDirtyRecursive returns true if w or any descendant is dirty.
func isDirtyRecursive(w widget.Widget) bool {
	if w.IsDirty() {
		return true
	}
	for _, child := range w.WidgetChildren() {
		if isDirtyRecursive(child) {
			return true
		}
	}
	return false
}

// clearDirtyRecursive clears dirty flags on all widgets.
func clearDirtyRecursive(w widget.Widget) {
	if base, ok := w.(interface{ ClearDirty() }); ok {
		base.ClearDirty()
	}
	for _, child := range w.WidgetChildren() {
		clearDirtyRecursive(child)
	}
}

// moveCursor returns the ANSI escape for moving the cursor to (x, y) (0-based).
func moveCursor(x, y int) string {
	return "\x1b[" + itoa(y+1) + ";" + itoa(x+1) + "H"
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	buf := [20]byte{}
	i := len(buf)
	neg := n < 0
	if neg {
		n = -n
	}
	for n > 0 {
		i--
		buf[i] = byte(n%10) + '0'
		n /= 10
	}
	if neg {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}
