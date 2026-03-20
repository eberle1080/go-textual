// Stopwatch is a go-textual application demonstrating Tick-based timers.
// Press Space to start/stop, r to reset, q or Ctrl+C to quit.
package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	rich "github.com/eberle1080/go-rich"

	"github.com/eberle1080/go-textual/app"
	"github.com/eberle1080/go-textual/css"
	"github.com/eberle1080/go-textual/geometry"
	"github.com/eberle1080/go-textual/keys"
	"github.com/eberle1080/go-textual/msg"
	"github.com/eberle1080/go-textual/screen"
	"github.com/eberle1080/go-textual/strip"
	"github.com/eberle1080/go-textual/widget"
)

type StopwatchScreen struct {
	screen.BaseScreen
	elapsed time.Duration
	running bool
}

func NewStopwatchScreen() *StopwatchScreen {
	return &StopwatchScreen{}
}

const tickInterval = 100 * time.Millisecond

func (s *StopwatchScreen) WidgetChildren() []widget.Widget { return nil }

func (s *StopwatchScreen) OnMount(_ context.Context) msg.Cmd {
	return msg.Tick(100 * time.Millisecond)
}

func (s *StopwatchScreen) Update(_ context.Context, m msg.Msg) msg.Cmd {
	switch v := m.(type) {
	case msg.TickMsg:
		if s.running {
			s.elapsed += v.Duration
			s.MarkDirty()
		}
		return msg.Tick(tickInterval)
	case msg.KeyMsg:
		switch v.Key {
		case " ":
			s.running = !s.running
			s.MarkDirty()
		case "r", "R":
			s.running = false
			s.elapsed = 0
			s.MarkDirty()
		case "q", "Q", keys.Escape, "ctrl+c":
			return func(_ context.Context) msg.Msg { return msg.QuitMsg{} }
		}
	}
	return nil
}

func (s *StopwatchScreen) Render(region geometry.Region) []strip.Strip {
	strips := make([]strip.Strip, region.Height)
	if region.Height == 0 {
		return strips
	}

	// Title
	titleStyle := rich.NewStyle().Bold().Foreground(rich.ANSIColor(rich.Cyan))
	titleStrip := strip.New(rich.Segments{{Text: "  Stopwatch", Style: titleStyle}})
	titleStrip = titleStrip.TextAlign(region.Width, css.AlignHorizontal("left"))
	if region.Height > 0 {
		strips[0] = titleStrip
	}

	// Time display
	if region.Height > 1 {
		h := int(s.elapsed.Hours())
		m := int(s.elapsed.Minutes()) % 60
		sec := int(s.elapsed.Seconds()) % 60
		cs := int(s.elapsed.Milliseconds()/10) % 100

		timeStr := fmt.Sprintf("  %02d:%02d:%02d.%02d", h, m, sec, cs)
		timeStyle := rich.NewStyle().Bold()
		if s.running {
			timeStyle = timeStyle.Foreground(rich.ANSIColor(rich.Green))
		} else {
			timeStyle = timeStyle.Foreground(rich.ANSIColor(rich.White))
		}
		timeStrip := strip.New(rich.Segments{{Text: timeStr, Style: timeStyle}})
		timeStrip = timeStrip.TextAlign(region.Width, css.AlignHorizontal("center"))
		strips[1] = timeStrip
	}

	// Status
	if region.Height > 2 {
		status := "  [STOPPED]  Space=start  r=reset  q=quit"
		statusStyle := rich.NewStyle().Foreground(rich.ANSIColor(rich.BrightBlack))
		if s.running {
			status = "  [RUNNING]  Space=stop   r=reset  q=quit"
			statusStyle = rich.NewStyle().Foreground(rich.ANSIColor(rich.Yellow))
		}
		statusStrip := strip.New(rich.Segments{{Text: status, Style: statusStyle}})
		statusStrip = statusStrip.TextAlign(region.Width, css.AlignHorizontal("left"))
		strips[2] = statusStrip
	}

	for i := 3; i < region.Height; i++ {
		strips[i] = strip.New(nil)
	}
	return strips
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo}))
	a := app.New(app.WithLogger(logger))
	if err := a.Run(context.Background(), NewStopwatchScreen()); err != nil {
		logger.Error("app exited with error", "err", err)
		os.Exit(1)
	}
}
