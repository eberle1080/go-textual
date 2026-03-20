// Package dom implements the DOM tree: a parent-child node hierarchy with
// CSS selector matching. It has no message-pump or reactive dependencies.
package dom

import "github.com/eberle1080/go-textual/css"

// Node is the interface implemented by all DOM nodes (including Widget).
type Node interface {
	// CSS selector matching.
	NodeID() string
	CSSTypeNames() []string
	CSSClasses() map[string]bool
	HasClass(name string) bool
	HasAllClasses(names ...string) bool
	HasPseudoClasses(set map[string]bool) bool
	CSSPathNodes() []css.SelectorNode

	// Identification.
	NodeName() string
	CSSTypeName() string

	// Tree structure.
	Parent() Node
	Children() *NodeList

	// Style access.
	Styles() *css.RenderStyles
	CSSStyles() *css.Styles
	InlineStyles() *css.Styles

	// Visibility.
	Display() bool
	Visible() bool
}
