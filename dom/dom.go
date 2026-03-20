package dom

import "github.com/eberle1080/go-textual/css"

// DOMNode is the base implementation of Node. Widget embeds DOMNode.
//
// DOMNode is not safe for concurrent mutation; it is expected to be accessed
// only from the event loop goroutine.
type DOMNode struct {
	name string
	id   string

	classes       map[string]bool
	pseudoClasses map[string]bool

	nodes  *NodeList
	parent Node

	cssStyles    *css.Styles
	inlineStyles *css.Styles
	renderStyles *css.RenderStyles

	cssTypeName  string
	cssTypeNames []string
}

// DOMNodeOption is a functional option for NewDOMNode.
type DOMNodeOption func(*DOMNode)

// WithName sets the optional display name.
func WithName(name string) DOMNodeOption {
	return func(d *DOMNode) { d.name = name }
}

// WithID sets the DOM id (without the # prefix).
func WithID(id string) DOMNodeOption {
	return func(d *DOMNode) { d.id = id }
}

// WithClasses sets the initial CSS classes.
func WithClasses(classes ...string) DOMNodeOption {
	return func(d *DOMNode) {
		for _, c := range classes {
			d.classes[c] = true
		}
	}
}

// WithCSSTypeName sets the CSS type name and hierarchy.
func WithCSSTypeName(typeName string, hierarchy ...string) DOMNodeOption {
	return func(d *DOMNode) {
		d.cssTypeName = typeName
		d.cssTypeNames = append([]string{typeName}, hierarchy...)
	}
}

// NewDOMNode constructs a DOMNode with the given functional options.
func NewDOMNode(opts ...DOMNodeOption) *DOMNode {
	d := &DOMNode{
		classes:       make(map[string]bool),
		pseudoClasses: make(map[string]bool),
		nodes:         NewNodeList(),
		cssStyles:     css.NewStyles(),
		inlineStyles:  css.NewStyles(),
	}
	d.renderStyles = css.NewRenderStyles(d.cssStyles, d.inlineStyles)
	for _, opt := range opts {
		opt(d)
	}
	return d
}

// ── dom.Node implementation ───────────────────────────────────────────────

func (d *DOMNode) NodeID() string          { return d.id }
func (d *DOMNode) NodeName() string        { return d.name }
func (d *DOMNode) CSSTypeName() string     { return d.cssTypeName }
func (d *DOMNode) CSSTypeNames() []string  { return d.cssTypeNames }
func (d *DOMNode) CSSClasses() map[string]bool { return d.classes }

func (d *DOMNode) HasClass(name string) bool { return d.classes[name] }

func (d *DOMNode) HasAllClasses(names ...string) bool {
	for _, name := range names {
		if !d.classes[name] {
			return false
		}
	}
	return true
}

func (d *DOMNode) HasPseudoClasses(set map[string]bool) bool {
	for pc := range set {
		if !d.pseudoClasses[pc] {
			return false
		}
	}
	return true
}

func (d *DOMNode) CSSPathNodes() []css.SelectorNode {
	var path []css.SelectorNode
	var current Node = d
	for current != nil {
		path = append(path, current)
		current = current.Parent()
	}
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return path
}

func (d *DOMNode) Parent() Node      { return d.parent }
func (d *DOMNode) Children() *NodeList {
	if d.nodes == nil {
		d.nodes = NewNodeList()
	}
	return d.nodes
}

func (d *DOMNode) Styles() *css.RenderStyles { return d.renderStyles }
func (d *DOMNode) CSSStyles() *css.Styles    { return d.cssStyles }
func (d *DOMNode) InlineStyles() *css.Styles { return d.inlineStyles }

func (d *DOMNode) Display() bool {
	return d.renderStyles.Display() != css.Display("none")
}

func (d *DOMNode) Visible() bool {
	if d.renderStyles.Display() == css.Display("none") {
		return false
	}
	if d.renderStyles.Visibility() == css.Visibility("hidden") {
		return false
	}
	if d.parent != nil {
		return d.parent.Visible()
	}
	return true
}

// ── Tree manipulation ─────────────────────────────────────────────────────

// SetParent sets the parent node. Called by the framework during mounting.
func (d *DOMNode) SetParent(parent Node) { d.parent = parent }

// ── CSS class management ──────────────────────────────────────────────────

func (d *DOMNode) AddClass(names ...string) {
	for _, name := range names {
		d.classes[name] = true
	}
}

func (d *DOMNode) RemoveClass(names ...string) {
	for _, name := range names {
		delete(d.classes, name)
	}
}

func (d *DOMNode) ToggleClass(names ...string) {
	for _, name := range names {
		if d.classes[name] {
			delete(d.classes, name)
		} else {
			d.classes[name] = true
		}
	}
}

func (d *DOMNode) SetClass(add bool, names ...string) {
	if add {
		d.AddClass(names...)
	} else {
		d.RemoveClass(names...)
	}
}

func (d *DOMNode) HasPseudoClass(name string) bool { return d.pseudoClasses[name] }

func (d *DOMNode) SetPseudoClass(name string, active bool) {
	if active {
		d.pseudoClasses[name] = true
	} else {
		delete(d.pseudoClasses, name)
	}
}

// ── Style management ──────────────────────────────────────────────────────

// SetStyles parses cssStr and merges the declarations into inline styles.
func (d *DOMNode) SetStyles(cssStr string) error {
	loc := css.CSSLocation{Variable: "inline"}
	styles, err := css.ParseDeclarations(cssStr, loc)
	if err != nil {
		return err
	}
	d.inlineStyles.Merge(styles)
	return nil
}

// ApplyCSSStyles replaces the base CSS styles with the given styles, then
// recomputes the render styles. Called by the stylesheet resolver.
func (d *DOMNode) ApplyCSSStyles(s *css.Styles) {
	d.cssStyles = s
	d.renderStyles = css.NewRenderStyles(d.cssStyles, d.inlineStyles)
}

// ResetStyles clears the base CSS styles.
func (d *DOMNode) ResetStyles() {
	d.cssStyles = css.NewStyles()
	d.renderStyles = css.NewRenderStyles(d.cssStyles, d.inlineStyles)
}

// ── Query ─────────────────────────────────────────────────────────────────

// Query returns a DOMQuery for the given CSS selector rooted at this node.
func (d *DOMNode) Query(selector string) *DOMQuery {
	return newQuery(d, selector, true)
}

// QueryOne finds exactly one matching node.
func (d *DOMNode) QueryOne(selector string) (Node, error) {
	results, err := d.Query(selector).Results()
	if err != nil {
		return nil, err
	}
	switch len(results) {
	case 0:
		return nil, &NoMatchesError{Selector: selector}
	case 1:
		return results[0], nil
	default:
		return nil, &TooManyMatchesError{Selector: selector, Count: len(results)}
	}
}
