// FileBrowser is a go-textual application demonstrating the DirectoryTree widget.
// Use Up/Down to navigate, Enter to open a directory. Press q or Ctrl+C to quit.
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

type FileBrowserScreen struct {
	screen.BaseScreen
	tree     *widgets.DirectoryTree
	selected string
}

func NewFileBrowserScreen(root string) *FileBrowserScreen {
	return &FileBrowserScreen{
		tree: widgets.NewDirectoryTree(root),
	}
}

func (s *FileBrowserScreen) WidgetChildren() []widget.Widget {
	return []widget.Widget{s.tree}
}

func (s *FileBrowserScreen) Update(ctx context.Context, m msg.Msg) msg.Cmd {
	switch v := m.(type) {
	case msg.KeyMsg:
		switch v.Key {
		case "q", "Q", keys.Escape, "ctrl+c":
			return func(_ context.Context) msg.Msg { return msg.QuitMsg{} }
		case keys.Backspace:
			// Navigate to parent directory
			root := s.tree.Root()
			parent := parentDir(root)
			if parent != root {
				s.tree.SetRoot(parent)
				s.MarkDirty()
			}
			return nil
		}
		// Forward navigation keys to the tree
		return s.tree.Update(ctx, m)
	case widgets.DirectorySelectedMsg:
		s.selected = v.Path
		if !v.IsDir {
			s.MarkDirty()
		}
	}
	return nil
}

func parentDir(path string) string {
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' || path[i] == '\\' {
			if i == 0 {
				return path[:1]
			}
			return path[:i]
		}
	}
	return path
}

func (s *FileBrowserScreen) Render(region geometry.Region) []strip.Strip {
	strips := make([]strip.Strip, region.Height)
	if region.Height == 0 {
		return strips
	}

	row := 0

	// Title bar
	titleStyle := rich.NewStyle().Bold().Foreground(rich.ANSIColor(rich.Cyan))
	title := fmt.Sprintf("  File Browser: %s", s.tree.Root())
	titleStrip := strip.New(rich.Segments{{Text: title, Style: titleStyle}})
	titleStrip = titleStrip.TextAlign(region.Width, css.AlignHorizontal("left"))
	strips[row] = titleStrip
	row++

	// Separator
	if row < region.Height {
		sepStyle := rich.NewStyle().Foreground(rich.ANSIColor(rich.BrightBlack))
		sep := make([]byte, region.Width)
		for i := range sep {
			sep[i] = '-'
		}
		sepStrip := strip.New(rich.Segments{{Text: string(sep), Style: sepStyle}})
		strips[row] = sepStrip
		row++
	}

	// Directory tree fills the middle
	treeHeight := region.Height - row - 2 // leave 2 rows at bottom
	if treeHeight < 0 {
		treeHeight = 0
	}
	treeRegion := geometry.Region{X: 0, Y: row, Width: region.Width, Height: treeHeight}
	treeStrips := widget.RenderChild(s.tree, treeRegion)
	for i, ts := range treeStrips {
		if row+i >= region.Height {
			break
		}
		strips[row+i] = ts
	}
	row += treeHeight

	// Selected file info
	if row < region.Height {
		infoStyle := rich.NewStyle().Foreground(rich.ANSIColor(rich.Yellow))
		infoText := "  No file selected"
		if s.selected != "" {
			infoText = fmt.Sprintf("  Selected: %s", s.selected)
		}
		infoStrip := strip.New(rich.Segments{{Text: infoText, Style: infoStyle}})
		infoStrip = infoStrip.TextAlign(region.Width, css.AlignHorizontal("left"))
		strips[row] = infoStrip
		row++
	}

	// Help
	if row < region.Height {
		helpStyle := rich.NewStyle().Foreground(rich.ANSIColor(rich.BrightBlack))
		helpStrip := strip.New(rich.Segments{{Text: "  ↑↓=navigate  Enter=open  Backspace=parent  q=quit", Style: helpStyle}})
		strips[row] = helpStrip
		row++
	}

	for ; row < region.Height; row++ {
		strips[row] = strip.New(nil)
	}
	return strips
}

func main() {
	root, err := os.Getwd()
	if err != nil {
		root = "/"
	}
	if len(os.Args) > 1 {
		root = os.Args[1]
	}

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo}))
	a := app.New(app.WithLogger(logger))
	if err := a.Run(context.Background(), NewFileBrowserScreen(root)); err != nil {
		logger.Error("app exited with error", "err", err)
		os.Exit(1)
	}
}
