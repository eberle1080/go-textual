// Todo is a go-textual application demonstrating the Input widget and list management.
// Type a task and press Enter to add it. Up/Down to navigate, d to delete selected.
// Press q or Ctrl+C to quit.
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

type TodoScreen struct {
	screen.BaseScreen
	input   *widgets.Input
	tasks   []todoTask
	cursor  int
	focused bool // true = list focused, false = input focused
}

type todoTask struct {
	text string
	done bool
}

func NewTodoScreen() *TodoScreen {
	s := &TodoScreen{
		input: widgets.NewInput("New task…"),
	}
	return s
}

func (s *TodoScreen) WidgetChildren() []widget.Widget {
	return []widget.Widget{s.input}
}

func (s *TodoScreen) OnMount(_ context.Context) msg.Cmd {
	// Input starts focused
	return nil
}

func (s *TodoScreen) addTask(text string) {
	if text == "" {
		return
	}
	s.tasks = append(s.tasks, todoTask{text: text})
	s.cursor = len(s.tasks) - 1
	s.MarkDirty()
}

func (s *TodoScreen) deleteSelected() {
	if len(s.tasks) == 0 {
		return
	}
	s.tasks = append(s.tasks[:s.cursor], s.tasks[s.cursor+1:]...)
	if s.cursor >= len(s.tasks) && s.cursor > 0 {
		s.cursor--
	}
	s.MarkDirty()
}

func (s *TodoScreen) toggleSelected() {
	if s.cursor < len(s.tasks) {
		s.tasks[s.cursor].done = !s.tasks[s.cursor].done
		s.MarkDirty()
	}
}

func (s *TodoScreen) Update(ctx context.Context, m msg.Msg) msg.Cmd {
	switch v := m.(type) {
	case msg.KeyMsg:
		switch v.Key {
		case "q", "Q", "ctrl+c":
			return func(_ context.Context) msg.Msg { return msg.QuitMsg{} }
		case keys.Escape:
			if s.focused {
				s.focused = false
				s.MarkDirty()
				return nil
			}
			return func(_ context.Context) msg.Msg { return msg.QuitMsg{} }
		case keys.Tab:
			s.focused = !s.focused
			s.MarkDirty()
			return nil
		}

		if s.focused {
			// List navigation
			switch v.Key {
			case keys.Up:
				if s.cursor > 0 {
					s.cursor--
					s.MarkDirty()
				}
			case keys.Down:
				if s.cursor < len(s.tasks)-1 {
					s.cursor++
					s.MarkDirty()
				}
			case "d", "D", keys.Delete:
				s.deleteSelected()
			case " ", keys.Enter:
				s.toggleSelected()
			}
		} else {
			// Input mode — forward to the input widget
			return s.input.Update(ctx, m)
		}

	case widgets.InputSubmittedMsg:
		s.addTask(v.Value)
		s.input.Clear()
		s.focused = true

	case widgets.InputChangedMsg:
		// nothing needed
	}
	return nil
}

func (s *TodoScreen) Render(region geometry.Region) []strip.Strip {
	strips := make([]strip.Strip, region.Height)
	if region.Height == 0 {
		return strips
	}

	row := 0

	// Title
	titleStyle := rich.NewStyle().Bold().Foreground(rich.ANSIColor(rich.Cyan))
	titleStrip := strip.New(rich.Segments{{Text: "  Todo List", Style: titleStyle}})
	titleStrip = titleStrip.TextAlign(region.Width, css.AlignHorizontal("left"))
	strips[row] = titleStrip
	row++

	// Input field
	if row < region.Height {
		inputLabel := "  Add: "
		labelStyle := rich.NewStyle().Foreground(rich.ANSIColor(rich.Yellow))
		inputRegion := geometry.Region{X: region.X + len(inputLabel), Y: region.Y + row, Width: region.Width - len(inputLabel), Height: 1}
		inputStrips := widget.RenderChild(s.input, inputRegion)
		inputSegs := rich.Segments{{Text: inputLabel, Style: labelStyle}}
		if len(inputStrips) > 0 {
			// Combine label + input
			labelStrip := strip.New(inputSegs)
			if len(inputStrips) > 0 {
				combined := strip.Join([]strip.Strip{labelStrip, inputStrips[0]})
				strips[row] = combined
			} else {
				strips[row] = labelStrip
			}
		}
		row++
	}

	// Separator
	if row < region.Height {
		sepStyle := rich.NewStyle().Foreground(rich.ANSIColor(rich.BrightBlack))
		sepStrip := strip.New(rich.Segments{{Text: "  " + repeatStr("─", region.Width-2), Style: sepStyle}})
		strips[row] = sepStrip
		row++
	}

	// Task list
	listHeight := region.Height - row - 1 // leave 1 for help
	if listHeight < 0 {
		listHeight = 0
	}

	// Scroll to keep cursor visible
	start := 0
	if s.cursor >= listHeight && listHeight > 0 {
		start = s.cursor - listHeight + 1
	}

	for i := 0; i < listHeight && row < region.Height-1; i++ {
		idx := start + i
		if idx >= len(s.tasks) {
			strips[row] = strip.New(nil)
			row++
			continue
		}
		task := s.tasks[idx]

		var prefix string
		if task.done {
			prefix = "  [x] "
		} else {
			prefix = "  [ ] "
		}

		var style rich.Style
		if idx == s.cursor && s.focused {
			style = rich.NewStyle().Reverse()
		} else if task.done {
			style = rich.NewStyle().Foreground(rich.ANSIColor(rich.BrightBlack)).Strikethrough()
		} else {
			style = rich.NewStyle()
		}

		text := fmt.Sprintf("%s%s", prefix, task.text)
		seg := rich.Segment{Text: text, Style: style}
		s := strip.New(rich.Segments{seg})
		s = s.TextAlign(region.Width, css.AlignHorizontal("left"))
		strips[row] = s
		row++
	}

	// Empty state
	if len(s.tasks) == 0 && row < region.Height-1 {
		emptyStyle := rich.NewStyle().Foreground(rich.ANSIColor(rich.BrightBlack)).Italic()
		emptyStrip := strip.New(rich.Segments{{Text: "  No tasks yet. Type above and press Enter.", Style: emptyStyle}})
		strips[row] = emptyStrip
		row++
	}

	// Help line (last row)
	if region.Height > 0 {
		helpStyle := rich.NewStyle().Foreground(rich.ANSIColor(rich.BrightBlack))
		mode := "INPUT"
		if s.focused {
			mode = "LIST "
		}
		help := fmt.Sprintf("  [%s] Tab=switch  Space/Enter=toggle  d=delete  q=quit", mode)
		helpStrip := strip.New(rich.Segments{{Text: help, Style: helpStyle}})
		strips[region.Height-1] = helpStrip
	}

	for ; row < region.Height-1; row++ {
		strips[row] = strip.New(nil)
	}

	return strips
}

func repeatStr(s string, n int) string {
	result := make([]byte, 0, len(s)*n)
	for range n {
		result = append(result, s...)
	}
	return string(result)
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo}))
	a := app.New(app.WithLogger(logger))
	if err := a.Run(context.Background(), NewTodoScreen()); err != nil {
		logger.Error("app exited with error", "err", err)
		os.Exit(1)
	}
}
