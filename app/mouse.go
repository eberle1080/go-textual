package app

import (
	"github.com/eberle1080/go-textual/msg"
	"github.com/eberle1080/go-textual/widget"
)

// widgetAt returns the deepest widget in the tree rooted at w whose recorded
// region contains screen coordinate (x, y). Children are checked before
// parents so overlapping children win.
func widgetAt(w widget.Widget, x, y int) widget.Widget {
	r := w.Region()
	if r.Width == 0 && r.Height == 0 {
		// Region not set — widget was never rendered, skip.
		return nil
	}
	if x < r.X || x >= r.X+r.Width || y < r.Y || y >= r.Y+r.Height {
		return nil
	}
	for _, child := range w.WidgetChildren() {
		if hit := widgetAt(child, x, y); hit != nil {
			return hit
		}
	}
	return w
}

// translateMouseMsg returns a copy of m with X/Y translated to be relative to
// the widget's region top-left corner. ScreenX/ScreenY are left unchanged.
func translateMouseMsg(m msg.Msg, w widget.Widget) msg.Msg {
	r := w.Region()
	dx := float64(r.X)
	dy := float64(r.Y)
	switch ev := m.(type) {
	case msg.MouseDownMsg:
		ev.X -= dx
		ev.Y -= dy
		return ev
	case msg.MouseUpMsg:
		ev.X -= dx
		ev.Y -= dy
		return ev
	case msg.MouseMoveMsg:
		ev.X -= dx
		ev.Y -= dy
		return ev
	case msg.MouseScrollUpMsg:
		ev.X -= dx
		ev.Y -= dy
		return ev
	case msg.MouseScrollDownMsg:
		ev.X -= dx
		ev.Y -= dy
		return ev
	case msg.MouseScrollLeftMsg:
		ev.X -= dx
		ev.Y -= dy
		return ev
	case msg.MouseScrollRightMsg:
		ev.X -= dx
		ev.Y -= dy
		return ev
	}
	return m
}

// mouseScreenXY extracts the screen-space integer coordinates from any mouse
// message. Returns (0, 0, false) if m is not a mouse message.
func mouseScreenXY(m msg.Msg) (x, y int, ok bool) {
	type screener interface{ ScreenOffset() interface{ X(); Y() } }
	switch ev := m.(type) {
	case msg.MouseDownMsg:
		return int(ev.ScreenX), int(ev.ScreenY), true
	case msg.MouseUpMsg:
		return int(ev.ScreenX), int(ev.ScreenY), true
	case msg.MouseMoveMsg:
		return int(ev.ScreenX), int(ev.ScreenY), true
	case msg.MouseScrollUpMsg:
		return int(ev.ScreenX), int(ev.ScreenY), true
	case msg.MouseScrollDownMsg:
		return int(ev.ScreenX), int(ev.ScreenY), true
	case msg.MouseScrollLeftMsg:
		return int(ev.ScreenX), int(ev.ScreenY), true
	case msg.MouseScrollRightMsg:
		return int(ev.ScreenX), int(ev.ScreenY), true
	}
	return 0, 0, false
}
