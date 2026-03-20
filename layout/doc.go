// Package layout provides the layout engine for go-textual.
//
// # Interfaces
//
// [Layoutable] is the interface implemented by all widgets. Layout code
// operates exclusively on Layoutable, so the layout package never imports
// the widget package (preventing an import cycle).
//
// [Layout] is implemented by the concrete layout algorithms. The framework
// calls [Layout.Arrange] to obtain a [DockArrangeResult] that describes where
// each child widget should be placed.
//
// # Layouts
//
// Four layouts are provided:
//
//   - [VerticalLayout] — stacks children vertically
//   - [HorizontalLayout] — arranges children horizontally
//   - [GridLayout] — arranges children in a grid with column/row spans
//   - [StreamLayout] — simplified vertical layout with auto heights
//
// Use [GetLayout] to obtain a named layout by string.
//
// # Placement
//
// [WidgetPlacement] describes where one widget is positioned within its
// parent's coordinate space. [DockArrangeResult] groups all placements for
// a single arrange call and provides spatial-map-based querying via
// [DockArrangeResult.GetVisiblePlacements].
package layout
