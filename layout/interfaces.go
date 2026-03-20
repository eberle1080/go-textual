package layout

import (
	"math/big"

	"github.com/eberle1080/go-textual/dom"
	"github.com/eberle1080/go-textual/geometry"
)

// Layoutable is the interface that widgets must implement so that the layout
// engine can position them without importing the widget package.
type Layoutable interface {
	dom.Node

	// GetBoxModel returns the resolved box model (width, height, margin) for
	// this widget given the container and viewport sizes and optional fraction
	// units.
	GetBoxModel(
		container, viewport geometry.Size,
		widthFrac, heightFrac *big.Rat,
		constrainWidth, greedy bool,
	) BoxModel

	// GetContentWidth returns the intrinsic content width for widgets that
	// don't have an explicit CSS width.
	GetContentWidth(container, viewport geometry.Size) int

	// GetContentHeight returns the intrinsic content height for a widget
	// rendered at the given width.
	GetContentHeight(container, viewport geometry.Size, width int) int

	// Expand reports whether the widget expands to fill available space.
	Expand() bool
	// Shrink reports whether the widget shrinks to fit content.
	Shrink() bool

	// SortOrder returns the widget's sort order (creation order).
	SortOrder() int

	// AbsoluteOffset returns a fixed absolute offset for the widget, or nil
	// if the widget is positioned by the layout engine.
	AbsoluteOffset() *geometry.Offset

	// Layer returns the CSS layer name for this widget.
	Layer() string

	// PreLayout is called before layout begins so the widget can clear
	// layout-dependent caches.
	PreLayout(l Layout)

	// ProcessLayout receives the computed placements for this widget's
	// children and may modify them.
	ProcessLayout(placements []WidgetPlacement) []WidgetPlacement

	// ViewportSize returns the size of the terminal viewport.
	ViewportSize() geometry.Size

	// AppSize returns the size of the application root widget.
	AppSize() geometry.Size

	// ScrollableContentRegion returns the region within the widget that is
	// scrollable content.
	ScrollableContentRegion() geometry.Region
}

// Layout is the interface implemented by concrete layout algorithms.
type Layout interface {
	// Name returns the CSS layout name (e.g. "vertical").
	Name() string

	// Arrange positions children within size and returns their placements.
	Arrange(parent Layoutable, children []Layoutable, size geometry.Size, greedy bool) []WidgetPlacement

	// GetContentWidth returns the intrinsic width of a widget's content.
	GetContentWidth(widget Layoutable, container, viewport geometry.Size) int

	// GetContentHeight returns the intrinsic height of a widget's content.
	GetContentHeight(widget Layoutable, container, viewport geometry.Size, width int) int
}
