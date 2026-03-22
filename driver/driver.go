package driver

import (
	"github.com/eberle1080/go-textual/geometry"
	"github.com/eberle1080/go-textual/msg"
)

// EventSink is the minimal interface the driver uses to deliver messages to
// the application. It replaces the old AppRef/MessagePump interface.
type EventSink interface {
	Send(m msg.Msg)
}

// Driver is the interface implemented by all terminal driver implementations.
type Driver interface {
	// Write sends raw text to the terminal output.
	Write(data string)
	// Flush flushes any buffered output.
	Flush()
	// StartApplicationMode enters the alternate screen and raw mode and starts
	// reading input.
	StartApplicationMode()
	// StopApplicationMode leaves the alternate screen, restores terminal state,
	// and stops reading input.
	StopApplicationMode()
	// DisableInput stops the input reader without restoring the terminal.
	DisableInput()
	// Close releases all resources held by the driver.
	Close()
	// IsHeadless reports whether the driver runs without a real terminal.
	IsHeadless() bool
	// IsInline reports whether the driver operates in inline (non-fullscreen) mode.
	IsInline() bool
	// CanSuspend reports whether the driver supports SIGTSTP suspension.
	CanSuspend() bool
	// SuspendApplicationMode temporarily leaves application mode (e.g. for SIGTSTP).
	SuspendApplicationMode()
	// ResumeApplicationMode re-enters application mode after a suspension.
	ResumeApplicationMode()
	// OpenURL opens a URL in the default browser or terminal handler.
	OpenURL(url string)
	// SetCursorOrigin sets an (x, y) offset applied to all mouse events.
	SetCursorOrigin(x, y int)
	// ClearCursorOrigin removes any previously set cursor origin offset.
	ClearCursorOrigin()
}

// DriverOption is a functional option for NewBaseDriver.
type DriverOption func(*BaseDriver)

// WithDebug enables verbose parser debugging.
func WithDebug(debug bool) DriverOption {
	return func(b *BaseDriver) { b.debug = debug }
}

// WithMouse enables mouse reporting.
func WithMouse(mouse bool) DriverOption {
	return func(b *BaseDriver) { b.mouse = mouse }
}

// WithSize overrides the terminal size reported to the app.
func WithSize(w, h int) DriverOption {
	size := [2]int{w, h}
	return func(b *BaseDriver) { b.size = &size }
}

// BaseDriver holds state shared by all platform driver implementations.
type BaseDriver struct {
	sink         EventSink
	debug        bool
	mouse        bool
	size         *[2]int
	downButtons  []msg.MouseButton
	lastMoveX    float64
	lastMoveY    float64
	cursorOrigin *geometry.Offset
}

// NewBaseDriver constructs a BaseDriver.
func NewBaseDriver(sink EventSink, opts ...DriverOption) *BaseDriver {
	b := &BaseDriver{sink: sink, mouse: true}
	for _, opt := range opts {
		opt(b)
	}
	return b
}

// SetCursorOrigin sets the coordinate offset applied to mouse events.
func (b *BaseDriver) SetCursorOrigin(x, y int) {
	o := geometry.Offset{X: x, Y: y}
	b.cursorOrigin = &o
}

// ClearCursorOrigin removes the cursor origin offset.
func (b *BaseDriver) ClearCursorOrigin() {
	b.cursorOrigin = nil
}

// ProcessMsg applies per-driver transformations (coordinate offset,
// stale-button synthesis) to a message and then sends it to the sink.
func (b *BaseDriver) ProcessMsg(m msg.Msg) {
	switch ev := m.(type) {
	case msg.MouseDownMsg:
		if b.cursorOrigin != nil {
			ev.X -= float64(b.cursorOrigin.X)
			ev.Y -= float64(b.cursorOrigin.Y)
			ev.ScreenX -= float64(b.cursorOrigin.X)
			ev.ScreenY -= float64(b.cursorOrigin.Y)
		}
		b.downButtons = append(b.downButtons, ev.Button)
		b.sink.Send(ev)
		return
	case msg.MouseUpMsg:
		if b.cursorOrigin != nil {
			ev.X -= float64(b.cursorOrigin.X)
			ev.Y -= float64(b.cursorOrigin.Y)
			ev.ScreenX -= float64(b.cursorOrigin.X)
			ev.ScreenY -= float64(b.cursorOrigin.Y)
		}
		b.downButtons = removeButton(b.downButtons, ev.Button)
		b.sink.Send(ev)
		return
	case msg.MouseMoveMsg:
		if b.cursorOrigin != nil {
			ev.X -= float64(b.cursorOrigin.X)
			ev.Y -= float64(b.cursorOrigin.Y)
			ev.ScreenX -= float64(b.cursorOrigin.X)
			ev.ScreenY -= float64(b.cursorOrigin.Y)
		}
		// Synthesize MouseUp for stale held buttons.
		if ev.Button == 0 && len(b.downButtons) > 0 {
			for _, btn := range b.downButtons {
				synth := msg.NewMouseUp(nil, ev.X, ev.Y, ev.ScreenX, ev.ScreenY, 0, 0, btn, false, false, false)
				b.sink.Send(synth)
			}
			b.downButtons = b.downButtons[:0]
		}
		b.lastMoveX = ev.X
		b.lastMoveY = ev.Y
		b.sink.Send(ev)
		return
	}
	b.sink.Send(m)
}

// Send delivers a message directly to the sink (no coordinate transformation).
func (b *BaseDriver) Send(m msg.Msg) {
	b.sink.Send(m)
}

func removeButton(buttons []msg.MouseButton, btn msg.MouseButton) []msg.MouseButton {
	for i, b := range buttons {
		if b == btn {
			return append(buttons[:i:i], buttons[i+1:]...)
		}
	}
	return buttons
}
