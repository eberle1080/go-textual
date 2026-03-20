package widgets

import (
	"context"
	"strings"

	rich "github.com/eberle1080/go-rich"

	"github.com/eberle1080/go-textual/css"
	"github.com/eberle1080/go-textual/dom"
	"github.com/eberle1080/go-textual/geometry"
	"github.com/eberle1080/go-textual/keys"
	"github.com/eberle1080/go-textual/msg"
	"github.com/eberle1080/go-textual/strip"
	"github.com/eberle1080/go-textual/widget"
)

// TabChangedMsg is sent when the active tab changes.
type TabChangedMsg struct {
	msg.BaseMsg
	Index int
	Pane  TabPane
}

// TabbedContent displays a row of tab buttons and renders the active tab's
// content widget below it. Tab switching is done with Left/Right or Tab/Shift+Tab.
type TabbedContent struct {
	widget.BaseWidget
	panes    []TabPane
	contents []widget.Widget
	active   int
}

// NewTabbedContent creates a TabbedContent widget.
func NewTabbedContent() *TabbedContent {
	tc := &TabbedContent{
		BaseWidget: *widget.NewBaseWidget(
			widget.WithDOMOptions(dom.WithCSSTypeName("TabbedContent", "Widget")),
			widget.WithCanFocus(true),
		),
	}
	return tc
}

// AddTab adds a tab with the given pane definition and content widget.
func (tc *TabbedContent) AddTab(pane TabPane, content widget.Widget) {
	tc.panes = append(tc.panes, pane)
	tc.contents = append(tc.contents, content)
	tc.MarkDirty()
}

// SetActive switches to the tab at index i.
func (tc *TabbedContent) SetActive(i int) {
	if i < 0 || i >= len(tc.panes) {
		return
	}
	tc.active = i
	tc.MarkDirty()
}

// Active returns the index of the currently active tab.
func (tc *TabbedContent) Active() int { return tc.active }

// ActiveContent returns the content widget of the active tab, or nil.
func (tc *TabbedContent) ActiveContent() widget.Widget {
	if tc.active < len(tc.contents) {
		return tc.contents[tc.active]
	}
	return nil
}

// WidgetChildren returns only the active content widget so the dirty/render
// walk only touches what's visible.
func (tc *TabbedContent) WidgetChildren() []widget.Widget {
	if tc.active < len(tc.contents) {
		return []widget.Widget{tc.contents[tc.active]}
	}
	return nil
}

// Update handles tab switching and forwards key events to the active content.
func (tc *TabbedContent) Update(ctx context.Context, m msg.Msg) msg.Cmd {
	switch v := m.(type) {
	case msg.FocusMsg:
		tc.MarkDirty()
	case msg.BlurMsg:
		tc.MarkDirty()
	case msg.KeyMsg:
		switch v.Key {
		case keys.Left, keys.BackTab:
			if tc.active > 0 {
				tc.active--
				tc.MarkDirty()
				idx := tc.active
				pane := tc.panes[idx]
				return func(_ context.Context) msg.Msg {
					return TabChangedMsg{Index: idx, Pane: pane}
				}
			}
			return nil
		case keys.Right, keys.Tab:
			if tc.active < len(tc.panes)-1 {
				tc.active++
				tc.MarkDirty()
				idx := tc.active
				pane := tc.panes[idx]
				return func(_ context.Context) msg.Msg {
					return TabChangedMsg{Index: idx, Pane: pane}
				}
			}
			return nil
		}
		// Forward other keys to active content.
		if tc.active < len(tc.contents) {
			return tc.contents[tc.active].Update(ctx, m)
		}
	}
	return nil
}

// Render draws the tab bar and the active content below it.
func (tc *TabbedContent) Render(region geometry.Region) []strip.Strip {
	strips := make([]strip.Strip, region.Height)
	if region.Height == 0 || len(tc.panes) == 0 {
		return strips
	}

	activeStyle := rich.NewStyle().Bold().Reverse()
	inactiveStyle := rich.NewStyle().Foreground(rich.ANSIColor(rich.BrightBlack))
	sepStyle := rich.NewStyle().Foreground(rich.ANSIColor(rich.BrightBlack))

	// Tab bar (row 0)
	var tabSegs rich.Segments
	for i, pane := range tc.panes {
		var st rich.Style
		if i == tc.active {
			st = activeStyle
		} else if (pane.Style != rich.Style{}) {
			st = pane.Style
		} else {
			st = inactiveStyle
		}
		tabSegs = append(tabSegs, rich.Segment{Text: " " + pane.Title + " ", Style: st})
		if i < len(tc.panes)-1 {
			tabSegs = append(tabSegs, rich.Segment{Text: "│", Style: sepStyle})
		}
	}
	tabBar := strip.New(tabSegs)
	tabBar = tabBar.TextAlign(region.Width, css.AlignHorizontal("left"))
	strips[0] = tabBar

	// Separator (row 1)
	if region.Height > 1 {
		var sb strings.Builder
		for range region.Width {
			sb.WriteString("─")
		}
		sepLine := sb.String()
		strips[1] = strip.New(rich.Segments{{Text: sepLine, Style: sepStyle}})
	}

	// Content (rows 2+)
	contentHeight := region.Height - 2
	if contentHeight > 0 && tc.active < len(tc.contents) {
		contentRegion := geometry.Region{X: region.X, Y: region.Y + 2, Width: region.Width, Height: contentHeight}
		contentStrips := widget.RenderChild(tc.contents[tc.active], contentRegion)
		for i, s := range contentStrips {
			if 2+i >= region.Height {
				break
			}
			strips[2+i] = s
		}
	}

	return strips
}
