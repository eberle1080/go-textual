// Counter is a simple go-textual application demonstrating stateful widgets.
// Press + or = to increment, - to decrement, r to reset, q or Ctrl+C to quit.
package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	rich "github.com/eberle1080/go-rich"

	"github.com/eberle1080/go-textual/app"
	"github.com/eberle1080/go-textual/css"
	"github.com/eberle1080/go-textual/geometry"
	"github.com/eberle1080/go-textual/keys"
	"github.com/eberle1080/go-textual/msg"
	"github.com/eberle1080/go-textual/screen"
	"github.com/eberle1080/go-textual/strip"
	"github.com/eberle1080/go-textual/widget"
	"github.com/eberle1080/go-textual/widgets"
)

type CounterScreen struct {
	screen.BaseScreen
	count      int
	countLabel *widgets.Label
	helpLabel  *widgets.Label
	incBtn     *widgets.Button
	decBtn     *widgets.Button
	resetBtn   *widgets.Button
}

func NewCounterScreen() *CounterScreen {
	s := &CounterScreen{
		countLabel: widgets.NewLabel("0"),
		helpLabel:  widgets.NewLabel("  +/= increment   -/_ decrement   r reset   q quit"),
		incBtn:     widgets.NewButton("+"),
		decBtn:     widgets.NewButton("-"),
		resetBtn:   widgets.NewButton("Reset"),
	}
	s.updateLabel()
	return s
}

func (s *CounterScreen) updateLabel() {
	color := rich.ANSIColor(rich.Green)
	if s.count < 0 {
		color = rich.ANSIColor(rich.Red)
	} else if s.count == 0 {
		color = rich.ANSIColor(rich.White)
	}
	_ = color
	s.countLabel.SetText(fmt.Sprintf("Count: %d", s.count))
}

func (s *CounterScreen) WidgetChildren() []widget.Widget {
	return []widget.Widget{s.countLabel, s.helpLabel, s.incBtn, s.decBtn, s.resetBtn}
}

func (s *CounterScreen) Update(_ context.Context, m msg.Msg) msg.Cmd {
	switch v := m.(type) {
	case msg.KeyMsg:
		switch v.Key {
		case "+", "=":
			s.count++
			s.updateLabel()
		case "-", "_":
			s.count--
			s.updateLabel()
		case "r", "R":
			s.count = 0
			s.updateLabel()
		case "q", "Q", keys.Escape, "ctrl+c":
			return func(_ context.Context) msg.Msg { return msg.QuitMsg{} }
		}
	case widgets.ButtonPressedMsg:
		switch v.Button {
		case s.incBtn:
			s.count++
			s.updateLabel()
		case s.decBtn:
			s.count--
			s.updateLabel()
		case s.resetBtn:
			s.count = 0
			s.updateLabel()
		}
	}
	return nil
}

func (s *CounterScreen) Render(region geometry.Region) []strip.Strip {
	strips := make([]strip.Strip, region.Height)
	if region.Height < 5 {
		return strips
	}

	// Title bar
	title := "  Counter"
	titleStyle := rich.NewStyle().Bold().Foreground(rich.ANSIColor(rich.Cyan))
	titleStrip := strip.New(rich.Segments{{Text: title, Style: titleStyle}})
	titleStrip = titleStrip.TextAlign(region.Width, css.AlignHorizontal("left"))
	strips[0] = titleStrip

	// Separator
	sep := strip.New(rich.Segments{{Text: "  " + repeat("─", region.Width-2), Style: rich.NewStyle().Foreground(rich.ANSIColor(rich.BrightBlack))}})
	strips[1] = sep

	// Count value (centered, large-ish)
	countStyle := rich.NewStyle().Bold()
	if s.count > 0 {
		countStyle = countStyle.Foreground(rich.ANSIColor(rich.Green))
	} else if s.count < 0 {
		countStyle = countStyle.Foreground(rich.ANSIColor(rich.Red))
	}
	countText := fmt.Sprintf("  %d", s.count)
	countStrip := strip.New(rich.Segments{{Text: countText, Style: countStyle}})
	countStrip = countStrip.TextAlign(region.Width, css.AlignHorizontal("center"))
	strips[2] = countStrip

	// Buttons row
	if region.Width >= 30 {
		btnWidth := region.Width / 3
		decStrips := widget.RenderChild(s.decBtn, geometry.Region{X: region.X, Y: region.Y + 3, Width: btnWidth, Height: 1})
		resetStrips := widget.RenderChild(s.resetBtn, geometry.Region{X: region.X + btnWidth, Y: region.Y + 3, Width: btnWidth, Height: 1})
		incStrips := widget.RenderChild(s.incBtn, geometry.Region{X: region.X + 2*btnWidth, Y: region.Y + 3, Width: region.Width - 2*btnWidth, Height: 1})
		if len(decStrips) > 0 && len(resetStrips) > 0 && len(incStrips) > 0 {
			combined := strip.Join([]strip.Strip{decStrips[0], resetStrips[0], incStrips[0]})
			strips[3] = combined
		}
	}

	// Help line
	if region.Height > 4 {
		helpStrips := widget.RenderChild(s.helpLabel, geometry.Region{X: region.X, Y: region.Y + 4, Width: region.Width, Height: 1})
		if len(helpStrips) > 0 {
			strips[4] = helpStrips[0]
		}
	}

	for i := 5; i < region.Height; i++ {
		strips[i] = strip.New(nil)
	}
	return strips
}

func repeat(s string, n int) string {
	result := ""
	for i := 0; i < n; i++ {
		result += s
	}
	return result
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo}))
	a := app.New(app.WithLogger(logger))
	if err := a.Run(context.Background(), NewCounterScreen()); err != nil {
		logger.Error("app exited with error", "err", err)
		os.Exit(1)
	}
}
