package msg

import "github.com/eberle1080/go-textual/geometry"

// MouseButton identifies a mouse button.
type MouseButton int

const (
	MouseButtonLeft   MouseButton = 0
	MouseButtonMiddle MouseButton = 1
	MouseButtonRight  MouseButton = 2
)

// mouseBase holds common fields for all mouse messages.
type mouseBase struct {
	BaseMsg
	Widget  any
	X       float64
	Y       float64
	ScreenX float64
	ScreenY float64
	DeltaX  int
	DeltaY  int
	Button  MouseButton
	Shift   bool
	Alt     bool
	Ctrl    bool
}

// Offset returns the widget-relative position.
func (m mouseBase) Offset() geometry.Offset {
	return geometry.Offset{X: int(m.X), Y: int(m.Y)}
}

// ScreenOffset returns the screen-relative position.
func (m mouseBase) ScreenOffset() geometry.Offset {
	return geometry.Offset{X: int(m.ScreenX), Y: int(m.ScreenY)}
}

func newMouseBase(widget any, x, y, screenX, screenY float64, deltaX, deltaY int, button MouseButton, shift, alt, ctrl bool) mouseBase {
	return mouseBase{
		Widget: widget, X: x, Y: y, ScreenX: screenX, ScreenY: screenY,
		DeltaX: deltaX, DeltaY: deltaY,
		Button: button, Shift: shift, Alt: alt, Ctrl: ctrl,
	}
}

// MouseDownMsg is sent when a mouse button is pressed.
type MouseDownMsg struct{ mouseBase }

// MouseUpMsg is sent when a mouse button is released.
type MouseUpMsg struct{ mouseBase }

// MouseMoveMsg is sent when the mouse moves.
type MouseMoveMsg struct{ mouseBase }

// MouseScrollUpMsg is sent on upward scroll.
type MouseScrollUpMsg struct{ mouseBase }

// MouseScrollDownMsg is sent on downward scroll.
type MouseScrollDownMsg struct{ mouseBase }

// MouseScrollLeftMsg is sent on leftward scroll.
type MouseScrollLeftMsg struct{ mouseBase }

// MouseScrollRightMsg is sent on rightward scroll.
type MouseScrollRightMsg struct{ mouseBase }

func NewMouseDown(widget any, x, y, screenX, screenY float64, deltaX, deltaY int, button MouseButton, shift, alt, ctrl bool) MouseDownMsg {
	return MouseDownMsg{newMouseBase(widget, x, y, screenX, screenY, deltaX, deltaY, button, shift, alt, ctrl)}
}

func NewMouseUp(widget any, x, y, screenX, screenY float64, deltaX, deltaY int, button MouseButton, shift, alt, ctrl bool) MouseUpMsg {
	return MouseUpMsg{newMouseBase(widget, x, y, screenX, screenY, deltaX, deltaY, button, shift, alt, ctrl)}
}

func NewMouseMove(widget any, x, y, screenX, screenY float64, deltaX, deltaY int, button MouseButton, shift, alt, ctrl bool) MouseMoveMsg {
	return MouseMoveMsg{newMouseBase(widget, x, y, screenX, screenY, deltaX, deltaY, button, shift, alt, ctrl)}
}

func NewMouseScrollUp(widget any, x, y, screenX, screenY float64, deltaX, deltaY int, button MouseButton, shift, alt, ctrl bool) MouseScrollUpMsg {
	return MouseScrollUpMsg{newMouseBase(widget, x, y, screenX, screenY, deltaX, deltaY, button, shift, alt, ctrl)}
}

func NewMouseScrollDown(widget any, x, y, screenX, screenY float64, deltaX, deltaY int, button MouseButton, shift, alt, ctrl bool) MouseScrollDownMsg {
	return MouseScrollDownMsg{newMouseBase(widget, x, y, screenX, screenY, deltaX, deltaY, button, shift, alt, ctrl)}
}

func NewMouseScrollLeft(widget any, x, y, screenX, screenY float64, deltaX, deltaY int, button MouseButton, shift, alt, ctrl bool) MouseScrollLeftMsg {
	return MouseScrollLeftMsg{newMouseBase(widget, x, y, screenX, screenY, deltaX, deltaY, button, shift, alt, ctrl)}
}

func NewMouseScrollRight(widget any, x, y, screenX, screenY float64, deltaX, deltaY int, button MouseButton, shift, alt, ctrl bool) MouseScrollRightMsg {
	return MouseScrollRightMsg{newMouseBase(widget, x, y, screenX, screenY, deltaX, deltaY, button, shift, alt, ctrl)}
}
