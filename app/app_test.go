package app_test

import (
	"context"
	"testing"

	"github.com/eberle1080/go-textual/app"
	"github.com/eberle1080/go-textual/geometry"
	"github.com/eberle1080/go-textual/msg"
	"github.com/eberle1080/go-textual/screen"
	"github.com/eberle1080/go-textual/strip"
	"github.com/eberle1080/go-textual/widget"
)

// quitScreen sends a QuitMsg on mount so the event loop exits immediately.
type quitScreen struct {
	screen.BaseScreen
}

func (s *quitScreen) Render(region geometry.Region) []strip.Strip {
	return make([]strip.Strip, region.Height)
}

func (s *quitScreen) OnMount(_ context.Context) msg.Cmd {
	return func(_ context.Context) msg.Msg {
		return msg.QuitMsg{}
	}
}

func (s *quitScreen) WidgetChildren() []widget.Widget { return nil }

func TestEventLoop_QuitMsg(t *testing.T) {
	a := app.New(app.WithHeadless(80, 24))
	err := a.Run(context.Background(), &quitScreen{})
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
}

// panicScreen sends a synthetic PanicMsg from OnMount, then a QuitMsg.
type panicScreen struct {
	screen.BaseScreen
}

func (s *panicScreen) Render(region geometry.Region) []strip.Strip {
	return make([]strip.Strip, region.Height)
}

func (s *panicScreen) OnMount(_ context.Context) msg.Cmd {
	// Return a Cmd that panics; the event loop should recover and continue.
	return func(_ context.Context) msg.Msg {
		panic("test panic")
	}
}

func (s *panicScreen) Update(_ context.Context, m msg.Msg) msg.Cmd {
	switch m.(type) {
	case msg.PanicMsg:
		// Panic was caught; now quit.
		return func(_ context.Context) msg.Msg {
			return msg.QuitMsg{}
		}
	}
	return nil
}

func (s *panicScreen) WidgetChildren() []widget.Widget { return nil }

func TestEventLoop_PanicRecovery(t *testing.T) {
	a := app.New(app.WithHeadless(80, 24))
	err := a.Run(context.Background(), &panicScreen{})
	if err != nil {
		t.Fatalf("Run returned error after panic recovery: %v", err)
	}
}
