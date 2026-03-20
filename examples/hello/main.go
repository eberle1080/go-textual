// Hello is a minimal go-textual application that demonstrates the event loop
// and widget rendering. Press Ctrl+C or Q to exit.
package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/eberle1080/go-textual/app"
	"github.com/eberle1080/go-textual/geometry"
	"github.com/eberle1080/go-textual/keys"
	"github.com/eberle1080/go-textual/msg"
	"github.com/eberle1080/go-textual/screen"
	"github.com/eberle1080/go-textual/strip"
	"github.com/eberle1080/go-textual/widget"
	"github.com/eberle1080/go-textual/widgets"
)

// HelloScreen is the root screen for the hello world app.
type HelloScreen struct {
	screen.BaseScreen
	label  *widgets.Label
	button *widgets.Button
}

// NewHelloScreen constructs the hello screen.
func NewHelloScreen() *HelloScreen {
	s := &HelloScreen{
		label:  widgets.NewLabel("Hello, world!"),
		button: widgets.NewButton("Quit"),
	}
	return s
}

// Compose returns the initial widget tree.
func (s *HelloScreen) Compose() []widget.Widget {
	return []widget.Widget{s.label, s.button}
}

// OnMount sets up the initial focus.
func (s *HelloScreen) OnMount(_ context.Context) msg.Cmd {
	return nil
}

// Update handles keyboard events.
func (s *HelloScreen) Update(_ context.Context, m msg.Msg) msg.Cmd {
	switch v := m.(type) {
	case msg.KeyMsg:
		switch v.Key {
		case "q", "Q", keys.Escape, "ctrl+c":
			return func(ctx context.Context) msg.Msg {
				return msg.QuitMsg{}
			}
		}
	case widgets.ButtonPressedMsg:
		if v.Button == s.button {
			return func(ctx context.Context) msg.Msg {
				return msg.QuitMsg{}
			}
		}
	}
	return nil
}

// Render draws the screen.
func (s *HelloScreen) Render(region geometry.Region) []strip.Strip {
	strips := make([]strip.Strip, region.Height)
	if region.Height == 0 {
		return strips
	}

	// Render label in top half, button in bottom half.
	half := region.Height / 2
	labelRegion := geometry.Region{X: region.X, Y: region.Y, Width: region.Width, Height: half}
	buttonRegion := geometry.Region{X: region.X, Y: region.Y + half, Width: region.Width, Height: region.Height - half}

	labelStrips := widget.RenderChild(s.label, labelRegion)
	buttonStrips := widget.RenderChild(s.button, buttonRegion)

	copy(strips[:half], labelStrips)
	copy(strips[half:], buttonStrips)
	return strips
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	a := app.New(app.WithLogger(logger))
	if err := a.Run(context.Background(), NewHelloScreen()); err != nil {
		logger.Error("app exited with error", "err", err)
		os.Exit(1)
	}
}
