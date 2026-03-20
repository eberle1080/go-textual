// Calculator is a simple 4-function go-textual application.
// Use number keys, +, -, *, / operators. Enter or = evaluates. Backspace clears.
// Press q or Ctrl+C to quit.
package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"

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

type CalcScreen struct {
	screen.BaseScreen
	display string // current input/result shown on screen
	left    float64
	op      string
	hasLeft bool
	error   bool
}

func NewCalcScreen() *CalcScreen {
	return &CalcScreen{display: "0"}
}

func (s *CalcScreen) WidgetChildren() []widget.Widget { return nil }

func (s *CalcScreen) evaluate() {
	if !s.hasLeft || s.op == "" {
		return
	}
	right, err := strconv.ParseFloat(s.display, 64)
	if err != nil {
		s.display = "Error"
		s.error = true
		return
	}
	var result float64
	switch s.op {
	case "+":
		result = s.left + right
	case "-":
		result = s.left - right
	case "*":
		result = s.left * right
	case "/":
		if right == 0 {
			s.display = "Div/0"
			s.error = true
			s.hasLeft = false
			s.op = ""
			return
		}
		result = s.left / right
	}
	// Format result: avoid trailing zeros for integers
	if result == float64(int64(result)) {
		s.display = fmt.Sprintf("%.0f", result)
	} else {
		s.display = strconv.FormatFloat(result, 'f', 10, 64)
		s.display = strings.TrimRight(s.display, "0")
	}
	s.left = result
	s.op = ""
	s.hasLeft = true
	s.error = false
}

func (s *CalcScreen) Update(_ context.Context, m msg.Msg) msg.Cmd {
	switch v := m.(type) {
	case msg.KeyMsg:
		switch v.Key {
		case "q", "Q", keys.Escape, "ctrl+c":
			return func(_ context.Context) msg.Msg { return msg.QuitMsg{} }
		case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
			if s.error || s.display == "0" {
				s.display = v.Key
				s.error = false
			} else if s.op != "" && !s.hasLeft {
				// Waiting for right operand — start fresh
				s.display = v.Key
			} else {
				if len(s.display) < 15 {
					s.display += v.Key
				}
			}
			s.MarkDirty()
		case ".":
			if s.error {
				s.display = "0."
				s.error = false
			} else if !strings.Contains(s.display, ".") {
				s.display += "."
			}
			s.MarkDirty()
		case "+", "-", "*", "/":
			if !s.error {
				val, err := strconv.ParseFloat(s.display, 64)
				if err == nil {
					if s.hasLeft && s.op != "" {
						s.evaluate()
					} else {
						s.left = val
						s.hasLeft = true
					}
					s.op = v.Key
					s.display = "0"
				}
			}
			s.MarkDirty()
		case keys.Enter, "=":
			if !s.error {
				s.evaluate()
			}
			s.MarkDirty()
		case keys.Backspace, "c", "C":
			s.display = "0"
			s.left = 0
			s.op = ""
			s.hasLeft = false
			s.error = false
			s.MarkDirty()
		}
	}
	return nil
}

func (s *CalcScreen) Render(region geometry.Region) []strip.Strip {
	strips := make([]strip.Strip, region.Height)
	if region.Height == 0 {
		return strips
	}

	row := 0

	// Title
	titleStyle := rich.NewStyle().Bold().Foreground(rich.ANSIColor(rich.Cyan))
	titleStrip := strip.New(rich.Segments{{Text: "  Calculator", Style: titleStyle}})
	titleStrip = titleStrip.TextAlign(region.Width, css.AlignHorizontal("left"))
	strips[row] = titleStrip
	row++

	// Display
	if row < region.Height {
		displayStyle := rich.NewStyle().Bold()
		if s.error {
			displayStyle = displayStyle.Foreground(rich.ANSIColor(rich.Red))
		} else {
			displayStyle = displayStyle.Foreground(rich.ANSIColor(rich.White))
		}
		expr := s.display
		if s.hasLeft && s.op != "" {
			expr = fmt.Sprintf("%.10g %s %s", s.left, s.op, s.display)
		}
		displayStrip := strip.New(rich.Segments{{Text: "  " + expr, Style: displayStyle}})
		displayStrip = displayStrip.TextAlign(region.Width, css.AlignHorizontal("right"))
		strips[row] = displayStrip
		row++
	}

	// Separator
	if row < region.Height {
		sepStyle := rich.NewStyle().Foreground(rich.ANSIColor(rich.BrightBlack))
		sepStrip := strip.New(rich.Segments{{Text: strings.Repeat("─", region.Width), Style: sepStyle}})
		strips[row] = sepStrip
		row++
	}

	// Button grid (4 rows of 4)
	buttonRows := []string{
		"7  8  9  /",
		"4  5  6  *",
		"1  2  3  -",
		"0  .  =  +",
	}
	for _, btnRow := range buttonRows {
		if row >= region.Height {
			break
		}
		btnStyle := rich.NewStyle()
		opStyle := rich.NewStyle().Foreground(rich.ANSIColor(rich.Yellow)).Bold()
		var segs rich.Segments
		for i, part := range strings.Split(btnRow, "  ") {
			if i > 0 {
				segs = append(segs, rich.Segment{Text: "  ", Style: rich.NewStyle()})
			}
			st := btnStyle
			p := strings.TrimSpace(part)
			if p == "+" || p == "-" || p == "*" || p == "/" || p == "=" {
				st = opStyle
			}
			segs = append(segs, rich.Segment{Text: "[ " + p + " ]", Style: st})
		}
		s := strip.New(segs)
		s = s.TextAlign(region.Width, css.AlignHorizontal("center"))
		strips[row] = s
		row++
	}

	// Help
	if row < region.Height {
		helpStyle := rich.NewStyle().Foreground(rich.ANSIColor(rich.BrightBlack))
		helpStrip := strip.New(rich.Segments{{Text: "  c=clear  q=quit", Style: helpStyle}})
		strips[row] = helpStrip
		row++
	}

	for ; row < region.Height; row++ {
		strips[row] = strip.New(nil)
	}
	return strips
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo}))
	a := app.New(app.WithLogger(logger))
	if err := a.Run(context.Background(), NewCalcScreen()); err != nil {
		logger.Error("app exited with error", "err", err)
		os.Exit(1)
	}
}
