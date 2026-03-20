package widget

// Tree is a helper for building widget trees.
// Usage:
//
//	tree := widget.Tree(root,
//	    widget.Tree(header),
//	    widget.Tree(body, widget.Tree(button)),
//	)
type Tree struct {
	Widget   Widget
	Children []*Tree
}

// BuildTree attaches all children in a Tree to their parent widgets.
// It returns the root widget.
func BuildTree(t *Tree) Widget {
	if t == nil {
		return nil
	}
	if base, ok := t.Widget.(interface{ AddChild(Widget) }); ok {
		for _, child := range t.Children {
			if child != nil {
				BuildTree(child)
				base.AddChild(child.Widget)
			}
		}
	}
	return t.Widget
}

// NewTree constructs a Tree node.
func NewTree(w Widget, children ...*Tree) *Tree {
	return &Tree{Widget: w, Children: children}
}
