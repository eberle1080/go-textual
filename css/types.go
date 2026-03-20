package css

import "github.com/eberle1080/go-textual/color"

// Display controls how a widget is rendered.
type Display string

// Visibility controls whether a widget is visible or hidden.
type Visibility string

// AlignHorizontal controls horizontal alignment.
type AlignHorizontal string

// AlignVertical controls vertical alignment.
type AlignVertical string

// ScrollbarGutter controls scrollbar gutter behaviour.
type ScrollbarGutter string

// BoxSizing controls how width and height are calculated.
type BoxSizing string

// Overflow controls what happens when content overflows.
type Overflow string

// EdgeType is the style of a border or outline edge.
type EdgeType string

// TextAlign controls text alignment.
type TextAlign string

// Constrain controls positional constraining.
type Constrain string

// Overlay controls overlay behaviour.
type Overlay string

// Position controls how offset is applied.
type Position string

// PointerShape is the mouse pointer shape.
type PointerShape string

// TextWrap controls text wrapping.
type TextWrap string

// TextOverflow controls text overflow behaviour.
type TextOverflow string

// Expand controls expansion behaviour.
type Expand string

// ScrollbarVisibility controls scrollbar visibility.
type ScrollbarVisibility string

// DockEdge is the edge a widget is docked to.
type DockEdge string

// EdgeStyle is a border/outline edge style: a type name and a colour.
type EdgeStyle struct {
	Type  EdgeType
	Color color.Color
}

// Specificity3 is a three-part CSS specificity value [id, class, type].
type Specificity3 [3]int

// Specificity6 is a six-part specificity used internally to resolve rule
// priority: [isUser, isImportant, id, class, type, tieBreaker].
type Specificity6 [6]int

// CSSLocation records the origin of a piece of CSS.
// Path is the file path (or an identifier); Variable is the class variable
// name when CSS is inlined in a widget definition, e.g. "Widget.DEFAULT_CSS".
type CSSLocation struct {
	Path     string
	Variable string
}

// AddSpecificity returns the element-wise sum of two Specificity3 values.
func AddSpecificity(a, b Specificity3) Specificity3 {
	return Specificity3{a[0] + b[0], a[1] + b[1], a[2] + b[2]}
}
