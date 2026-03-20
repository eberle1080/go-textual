package widgets

import (
	"context"
	"os"
	"path/filepath"
	"sort"

	rich "github.com/eberle1080/go-rich"

	"github.com/eberle1080/go-textual/css"
	"github.com/eberle1080/go-textual/dom"
	"github.com/eberle1080/go-textual/geometry"
	"github.com/eberle1080/go-textual/keys"
	"github.com/eberle1080/go-textual/msg"
	"github.com/eberle1080/go-textual/strip"
	"github.com/eberle1080/go-textual/widget"
)

// DirectorySelectedMsg is sent when the user selects a file or directory.
type DirectorySelectedMsg struct {
	msg.BaseMsg
	Path  string
	IsDir bool
}

type dirEntry struct {
	name  string
	path  string
	isDir bool
}

// DirectoryTree displays the contents of a directory as a navigable list.
type DirectoryTree struct {
	widget.BaseWidget
	root    string
	entries []dirEntry
	cursor  int
}

// NewDirectoryTree creates a DirectoryTree rooted at the given path.
func NewDirectoryTree(root string) *DirectoryTree {
	dt := &DirectoryTree{
		BaseWidget: *widget.NewBaseWidget(
			widget.WithDOMOptions(dom.WithCSSTypeName("DirectoryTree", "Widget")),
			widget.WithCanFocus(true),
		),
		root: root,
	}
	dt.reload()
	return dt
}

func (dt *DirectoryTree) reload() {
	entries, err := os.ReadDir(dt.root)
	if err != nil {
		dt.entries = nil
		return
	}

	// Directories first, then files, both sorted alphabetically.
	var dirs, files []dirEntry
	for _, e := range entries {
		entry := dirEntry{
			name:  e.Name(),
			path:  filepath.Join(dt.root, e.Name()),
			isDir: e.IsDir(),
		}
		if e.IsDir() {
			dirs = append(dirs, entry)
		} else {
			files = append(files, entry)
		}
	}
	sort.Slice(dirs, func(i, j int) bool { return dirs[i].name < dirs[j].name })
	sort.Slice(files, func(i, j int) bool { return files[i].name < files[j].name })

	dt.entries = append(dirs, files...)
	if dt.cursor >= len(dt.entries) {
		dt.cursor = len(dt.entries) - 1
	}
	if dt.cursor < 0 {
		dt.cursor = 0
	}
}

// SetRoot changes the root directory.
func (dt *DirectoryTree) SetRoot(path string) {
	dt.root = path
	dt.cursor = 0
	dt.reload()
	dt.MarkDirty()
}

// Root returns the current root directory.
func (dt *DirectoryTree) Root() string { return dt.root }

// Update handles navigation keys.
func (dt *DirectoryTree) Update(_ context.Context, m msg.Msg) msg.Cmd {
	switch v := m.(type) {
	case msg.FocusMsg:
		dt.MarkDirty()
	case msg.BlurMsg:
		dt.MarkDirty()
	case msg.MouseDownMsg:
		// Translate Y to an entry index accounting for scroll.
		region := dt.Region()
		scrollStart := 0
		if dt.cursor >= region.Height && region.Height > 0 {
			scrollStart = dt.cursor - region.Height + 1
		}
		idx := scrollStart + int(v.Y)
		if idx >= 0 && idx < len(dt.entries) {
			dt.cursor = idx
			dt.MarkDirty()
			if v.Button == msg.MouseButtonLeft {
				e := dt.entries[idx]
				if e.isDir {
					dt.SetRoot(e.path)
				}
				path := e.path
				isDir := e.isDir
				return func(_ context.Context) msg.Msg {
					return DirectorySelectedMsg{Path: path, IsDir: isDir}
				}
			}
		}
		return nil
	case msg.KeyMsg:
		switch v.Key {
		case keys.Up:
			if dt.cursor > 0 {
				dt.cursor--
				dt.MarkDirty()
			}
		case keys.Down:
			if dt.cursor < len(dt.entries)-1 {
				dt.cursor++
				dt.MarkDirty()
			}
		case keys.Enter:
			if dt.cursor < len(dt.entries) {
				e := dt.entries[dt.cursor]
				if e.isDir {
					dt.SetRoot(e.path)
				}
				path := e.path
				isDir := e.isDir
				return func(_ context.Context) msg.Msg {
					return DirectorySelectedMsg{Path: path, IsDir: isDir}
				}
			}
		}
	}
	return nil
}

// Render draws the directory listing.
func (dt *DirectoryTree) Render(region geometry.Region) []strip.Strip {
	strips := make([]strip.Strip, region.Height)
	if region.Height == 0 {
		return strips
	}

	// Scroll so cursor is always visible.
	start := 0
	if dt.cursor >= region.Height {
		start = dt.cursor - region.Height + 1
	}

	for row := 0; row < region.Height; row++ {
		idx := start + row
		if idx >= len(dt.entries) {
			strips[row] = strip.New(nil)
			continue
		}
		e := dt.entries[idx]

		prefix := "  "
		if e.isDir {
			prefix = "▶ "
		}
		text := prefix + e.name

		var style rich.Style
		if idx == dt.cursor {
			style = rich.NewStyle().Reverse()
		} else if e.isDir {
			style = rich.NewStyle().Foreground(rich.ANSIColor(rich.Blue)).Bold()
		} else {
			style = rich.NewStyle()
		}

		seg := rich.Segment{Text: text, Style: style}
		s := strip.New(rich.Segments{seg})
		s = s.TextAlign(region.Width, css.AlignHorizontal("left"))
		strips[row] = s
	}
	return strips
}
