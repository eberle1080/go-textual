// Dashboard is a go-textual application demonstrating multiple widgets in a layout.
// Shows simulated system metrics with progress bars and a log. Press q or Ctrl+C to quit.
package main

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
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
	"github.com/eberle1080/go-textual/widgets"
)

type metric struct {
	name    string
	value   float64
	total   float64
	bar     *widgets.ProgressBar
}

type DashboardScreen struct {
	screen.BaseScreen
	metrics []*metric
	log     *widgets.RichLog
	ticks   int
	rng     *rand.Rand
}

func NewDashboardScreen() *DashboardScreen {
	s := &DashboardScreen{
		log: widgets.NewRichLog(50),
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	s.metrics = []*metric{
		{name: "CPU", total: 100, bar: widgets.NewProgressBar(100)},
		{name: "MEM", total: 100, bar: widgets.NewProgressBar(100)},
		{name: "DSK", total: 100, bar: widgets.NewProgressBar(100)},
		{name: "NET", total: 100, bar: widgets.NewProgressBar(100)},
	}
	for _, m := range s.metrics {
		m.value = s.rng.Float64() * 60
		m.bar.SetProgress(m.value)
	}
	s.log.Write("Dashboard started.")
	return s
}

func (s *DashboardScreen) WidgetChildren() []widget.Widget {
	ws := make([]widget.Widget, 0, len(s.metrics)+1)
	for _, m := range s.metrics {
		ws = append(ws, m.bar)
	}
	ws = append(ws, s.log)
	return ws
}

func (s *DashboardScreen) OnMount(_ context.Context) msg.Cmd {
	return msg.Tick(500 * time.Millisecond)
}

func (s *DashboardScreen) Update(_ context.Context, m msg.Msg) msg.Cmd {
	switch v := m.(type) {
	case msg.TickMsg:
		_ = v
		s.ticks++
		s.updateMetrics()
		s.MarkDirty()
		return msg.Tick(500 * time.Millisecond)
	case msg.KeyMsg:
		switch v.Key {
		case "q", "Q", keys.Escape, "ctrl+c":
			return func(_ context.Context) msg.Msg { return msg.QuitMsg{} }
		}
	}
	return nil
}

func (s *DashboardScreen) updateMetrics() {
	for _, met := range s.metrics {
		// Random walk
		delta := (s.rng.Float64() - 0.4) * 15
		met.value += delta
		if met.value < 0 {
			met.value = 0
		}
		if met.value > met.total {
			met.value = met.total
		}
		met.bar.SetProgress(met.value)
	}

	if s.ticks%4 == 0 {
		cpu := s.metrics[0].value
		var logStyle rich.Style
		var logMsg string
		if cpu > 80 {
			logStyle = rich.NewStyle().Foreground(rich.ANSIColor(rich.Red))
			logMsg = fmt.Sprintf("[WARN] High CPU: %.0f%%", cpu)
		} else if cpu > 60 {
			logStyle = rich.NewStyle().Foreground(rich.ANSIColor(rich.Yellow))
			logMsg = fmt.Sprintf("[INFO] CPU moderate: %.0f%%", cpu)
		} else {
			logStyle = rich.NewStyle().Foreground(rich.ANSIColor(rich.Green))
			logMsg = fmt.Sprintf("[OK]   CPU normal: %.0f%%", cpu)
		}
		s.log.WriteStyled(logMsg, logStyle)
	}
}

func (s *DashboardScreen) Render(region geometry.Region) []strip.Strip {
	strips := make([]strip.Strip, region.Height)
	if region.Height == 0 {
		return strips
	}

	row := 0

	// Title
	titleStyle := rich.NewStyle().Bold().Foreground(rich.ANSIColor(rich.Cyan))
	title := fmt.Sprintf("  System Dashboard  (tick #%d)", s.ticks)
	titleStrip := strip.New(rich.Segments{{Text: title, Style: titleStyle}})
	titleStrip = titleStrip.TextAlign(region.Width, css.AlignHorizontal("left"))
	strips[row] = titleStrip
	row++

	// Separator
	if row < region.Height {
		sepStyle := rich.NewStyle().Foreground(rich.ANSIColor(rich.BrightBlack))
		sep := ""
		for range region.Width {
			sep += "─"
		}
		strips[row] = strip.New(rich.Segments{{Text: sep, Style: sepStyle}})
		row++
	}

	// Metrics rows
	for _, met := range s.metrics {
		if row >= region.Height {
			break
		}
		labelStyle := rich.NewStyle().Bold()
		label := fmt.Sprintf("  %-4s ", met.name)
		labelStrip := strip.New(rich.Segments{{Text: label, Style: labelStyle}})
		barRegion := geometry.Region{X: region.X + len(label), Y: region.Y + row, Width: region.Width - len(label), Height: 1}
		barStrips := widget.RenderChild(met.bar, barRegion)
		if len(barStrips) > 0 {
			combined := strip.Join([]strip.Strip{labelStrip, barStrips[0]})
			strips[row] = combined
		} else {
			strips[row] = labelStrip
		}
		row++
	}

	// Separator
	if row < region.Height {
		sepStyle := rich.NewStyle().Foreground(rich.ANSIColor(rich.BrightBlack))
		strips[row] = strip.New(rich.Segments{{Text: "  Log:", Style: sepStyle}})
		row++
	}

	// Log fills remaining space (minus help row)
	logHeight := region.Height - row - 1
	if logHeight < 0 {
		logHeight = 0
	}
	logRegion := geometry.Region{X: region.X, Y: region.Y + row, Width: region.Width, Height: logHeight}
	logStrips := widget.RenderChild(s.log, logRegion)
	for i, ls := range logStrips {
		if row+i >= region.Height-1 {
			break
		}
		strips[row+i] = ls
	}
	row += logHeight

	// Help
	if row < region.Height {
		helpStyle := rich.NewStyle().Foreground(rich.ANSIColor(rich.BrightBlack))
		strips[region.Height-1] = strip.New(rich.Segments{{Text: "  q=quit", Style: helpStyle}})
	}

	return strips
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo}))
	a := app.New(app.WithLogger(logger))
	if err := a.Run(context.Background(), NewDashboardScreen()); err != nil {
		logger.Error("app exited with error", "err", err)
		os.Exit(1)
	}
}
