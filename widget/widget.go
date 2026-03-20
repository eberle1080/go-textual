// Package widget defines the Widget interface and BaseWidget implementation.
// Widgets are pure structs; all mutation happens in Update(), which is called
// exclusively from the event loop goroutine. No locks are needed.
package widget

import (
	"context"

	"github.com/eberle1080/go-textual/css"
	"github.com/eberle1080/go-textual/dom"
	"github.com/eberle1080/go-textual/geometry"
	"github.com/eberle1080/go-textual/msg"
	"github.com/eberle1080/go-textual/strip"
)

// Widget is the interface implemented by all widgets.
type Widget interface {
	dom.Node

	// Update is called by the event loop with each incoming message. It
	// returns a Cmd for any async work, or nil.
	Update(ctx context.Context, m msg.Msg) msg.Cmd

	// Render returns the visual representation of the widget for the given
	// region.
	Render(region geometry.Region) []strip.Strip

	// WidgetChildren returns the widget's direct children as Widgets.
	WidgetChildren() []Widget

	// CanFocus reports whether the widget can receive keyboard focus.
	CanFocus() bool

	// IsDirty reports whether the widget needs re-rendering.
	IsDirty() bool

	// MarkDirty marks the widget as needing re-rendering.
	MarkDirty()

	// Region returns the screen region this widget currently occupies.
	// It is set by RenderChild and used for mouse hit-testing.
	Region() geometry.Region
}

// RenderChild renders w into r, recording r as w's region for mouse
// hit-testing, then returns the resulting strips. Use this instead of calling
// w.Render(r) directly whenever the widget should be reachable by mouse events.
func RenderChild(w Widget, r geometry.Region) []strip.Strip {
	if rr, ok := w.(interface{ SetRegion(geometry.Region) }); ok {
		rr.SetRegion(r)
	}
	return w.Render(r)
}

// BaseWidget provides the default implementation of Widget.
// Concrete widgets should embed *BaseWidget and override Update/Render.
type BaseWidget struct {
	dom.DOMNode

	dirty    bool
	region   geometry.Region
	canFocus bool
	children []Widget
}

// BaseWidgetOption is a functional option for NewBaseWidget.
type BaseWidgetOption func(*BaseWidget)

// WithDOMOptions passes options through to the underlying DOMNode.
func WithDOMOptions(opts ...dom.DOMNodeOption) BaseWidgetOption {
	return func(b *BaseWidget) {
		for _, opt := range opts {
			opt(&b.DOMNode)
		}
	}
}

// WithCanFocus sets the focusability of the widget.
func WithCanFocus(v bool) BaseWidgetOption {
	return func(b *BaseWidget) { b.canFocus = v }
}

// NewBaseWidget constructs a BaseWidget.
func NewBaseWidget(opts ...BaseWidgetOption) *BaseWidget {
	b := &BaseWidget{
		DOMNode: *dom.NewDOMNode(),
		dirty:   true,
	}
	for _, opt := range opts {
		opt(b)
	}
	return b
}

// Update is the default no-op implementation.
func (b *BaseWidget) Update(_ context.Context, _ msg.Msg) msg.Cmd { return nil }

// Render returns an empty strip list. Concrete widgets override this.
func (b *BaseWidget) Render(region geometry.Region) []strip.Strip {
	return make([]strip.Strip, region.Height)
}

// WidgetChildren returns the widget's children.
func (b *BaseWidget) WidgetChildren() []Widget { return b.children }

// AddChild adds a child widget.
func (b *BaseWidget) AddChild(w Widget) {
	b.children = append(b.children, w)
	b.DOMNode.Children().Append(w)
	b.MarkDirty()
}

// RemoveChild removes a child widget.
func (b *BaseWidget) RemoveChild(w Widget) {
	for i, c := range b.children {
		if c == w {
			b.children = append(b.children[:i], b.children[i+1:]...)
			break
		}
	}
	b.DOMNode.Children().Remove(w)
	b.MarkDirty()
}

// CanFocus reports whether the widget accepts keyboard focus.
func (b *BaseWidget) CanFocus() bool { return b.canFocus }

// IsDirty reports whether the widget needs re-rendering.
func (b *BaseWidget) IsDirty() bool { return b.dirty }

// MarkDirty marks the widget as needing re-rendering.
func (b *BaseWidget) MarkDirty() { b.dirty = true }

// ClearDirty clears the dirty flag after rendering.
func (b *BaseWidget) ClearDirty() { b.dirty = false }

// Region returns the region assigned by the compositor after layout.
func (b *BaseWidget) Region() geometry.Region { return b.region }

// SetRegion records the region assigned by the compositor.
func (b *BaseWidget) SetRegion(r geometry.Region) { b.region = r }

// ApplyStyles updates the base CSS styles and refreshes the render styles.
func (b *BaseWidget) ApplyStyles(s *css.Styles) {
	b.DOMNode.ApplyCSSStyles(s)
	b.MarkDirty()
}
