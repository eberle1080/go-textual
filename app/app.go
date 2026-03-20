// Package app provides the top-level application: a single event loop
// goroutine that owns all widget state, with input/signal/timer goroutines
// communicating through a typed channel.
package app

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/eberle1080/go-textual/driver"
	"github.com/eberle1080/go-textual/geometry"
	"github.com/eberle1080/go-textual/msg"
	"github.com/eberle1080/go-textual/screen"
	"github.com/eberle1080/go-textual/widget"
)

// App is the top-level application. Create one with New, then call Run.
type App struct {
	ctx    context.Context
	cancel context.CancelFunc
	events chan msg.Msg

	drv     driver.Driver
	screen  screen.Screen
	size    geometry.Size
	focused widget.Widget

	logger   *slog.Logger
	quitting bool

	// mouseDownTarget is the widget that received the last MouseDown, used to
	// synthesize click events on MouseUp.
	mouseDownTarget widget.Widget

	// headless holds the headless driver config, if any.
	headless *headlessConfig
}

type headlessConfig struct {
	width, height int
}

// Option configures an App.
type Option func(*App)

// WithLogger sets the logger. Defaults to slog.Default().
func WithLogger(l *slog.Logger) Option {
	return func(a *App) { a.logger = l }
}

// WithHeadless configures the app to run without a real terminal, useful for
// testing.
func WithHeadless(width, height int) Option {
	return func(a *App) {
		a.headless = &headlessConfig{width: width, height: height}
	}
}

// New creates a new App.
func New(opts ...Option) *App {
	a := &App{
		events: make(chan msg.Msg, 256),
		logger: slog.Default(),
	}
	for _, opt := range opts {
		opt(a)
	}
	return a
}

// Send implements driver.EventSink, delivering a message to the event loop.
func (a *App) Send(m msg.Msg) {
	select {
	case a.events <- m:
	default:
		a.logger.Debug("event channel full, dropping message",
			slog.String("type", msgTypeName(m)))
	}
}

// Run mounts the screen and runs the event loop until the app exits.
// It returns nil on clean exit, or an error if something went wrong.
func (a *App) Run(ctx context.Context, s screen.Screen, opts ...RunOption) error {
	a.ctx, a.cancel = context.WithCancel(ctx)
	defer a.cancel()

	a.screen = s

	// Build driver.
	if a.headless != nil {
		a.drv = newHeadlessDriver(a, a.headless.width, a.headless.height)
	} else {
		a.drv = newPlatformDriver(a)
	}

	a.drv.StartApplicationMode()
	defer a.drv.Close()
	defer a.drv.StopApplicationMode()

	// Start signal handler goroutine.
	go a.runSignalHandler(a.ctx)

	// Mount the initial screen.
	if cmd := a.mountScreen(a.ctx, s); cmd != nil {
		go a.runCmd(a.ctx, cmd)
	}

	return a.runLoop()
}

// Exit cleanly terminates the event loop.
func (a *App) Exit() {
	a.Send(msg.QuitMsg{})
}

// PushScreen replaces the current screen with s and returns a Cmd that
// performs the swap inside the event loop.
func (a *App) PushScreen(s screen.Screen) msg.Cmd {
	return func(ctx context.Context) msg.Msg {
		return pushScreenMsg{screen: s}
	}
}

// RunOption configures Run behaviour.
type RunOption func(*App)

// mountScreen mounts s as the active screen, calling Compose and OnMount.
func (a *App) mountScreen(ctx context.Context, s screen.Screen) msg.Cmd {
	// Build child tree from Compose.
	if base, ok := s.(interface{ AddChild(widget.Widget) }); ok {
		for _, child := range s.Compose() {
			base.AddChild(child)
			a.dispatchMount(ctx, child)
		}
	}
	a.focused = widget.NextFocus(s, nil)

	return s.OnMount(ctx)
}

// dispatchMount sends MountMsg to w and all its descendants.
func (a *App) dispatchMount(ctx context.Context, w widget.Widget) {
	if cmd := w.Update(ctx, msg.MountMsg{}); cmd != nil {
		go a.runCmd(ctx, cmd)
	}
	for _, child := range w.WidgetChildren() {
		a.dispatchMount(ctx, child)
	}
}

// pushScreenMsg is internal — sent by PushScreen cmd.
type pushScreenMsg struct {
	msg.BaseMsg
	screen screen.Screen
}

// msgTypeName returns a short name for a message for logging.
func msgTypeName(m msg.Msg) string {
	if m == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%T", m)
}
