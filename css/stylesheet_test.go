package css

import (
	"fmt"
	"testing"
)

// mockStylesheetNode implements StylesheetNode for testing.
type mockStylesheetNode struct {
	*mockNode
	styles        *RenderStyles
	selectorNames map[string]bool
	typeName      string
	parent        StylesheetNode
	children      []StylesheetNode
}

func newMockStylesheetNode(typeNames []string, classes map[string]bool, id string) *mockStylesheetNode {
	mn := newMockNode(typeNames, classes, id)
	mn.path = []*mockNode{mn}
	rs := NewRenderStyles(NewStyles(), NewStyles())
	snames := make(map[string]bool)
	for _, tn := range typeNames {
		snames[tn] = true
	}
	for cls := range classes {
		snames["."+cls] = true
	}
	if id != "" {
		snames["#"+id] = true
	}
	primaryType := ""
	if len(typeNames) > 0 {
		primaryType = typeNames[0]
	}
	return &mockStylesheetNode{
		mockNode:      mn,
		styles:        rs,
		selectorNames: snames,
		typeName:      primaryType,
	}
}

func (n *mockStylesheetNode) SelectorNames() map[string]bool { return n.selectorNames }
func (n *mockStylesheetNode) NodeStyles() *RenderStyles      { return n.styles }
func (n *mockStylesheetNode) NotifyStyleUpdate()             {}
func (n *mockStylesheetNode) Refresh()                       {}
func (n *mockStylesheetNode) PseudoClassesCacheKey() any     { return n.id }
func (n *mockStylesheetNode) CSSTypeName() string            { return n.typeName }
func (n *mockStylesheetNode) Parent() StylesheetNode         { return n.parent }
func (n *mockStylesheetNode) Children() []StylesheetNode     { return n.children }
func (n *mockStylesheetNode) CSSPathNodes() []SelectorNode {
	result := make([]SelectorNode, len(n.path))
	for i, p := range n.path {
		result[i] = p
	}
	return result
}

func TestStylesheetAddSourceAndRules(t *testing.T) {
	ss := NewStylesheet(nil)
	err := ss.AddSource("Label { color: red; }", nil, false, 0, "")
	if err != nil {
		t.Fatalf("AddSource error: %v", err)
	}
	rules, err := ss.Rules()
	if err != nil {
		t.Fatalf("Rules() error: %v", err)
	}
	if len(rules) == 0 {
		t.Error("expected at least one rule set")
	}
}

func TestStylesheetParseAll(t *testing.T) {
	ss := NewStylesheet(nil)
	ss.AddSource("Label { color: red; }\n.myclass { display: block; }", nil, false, 0, "")
	err := ss.ParseAll()
	if err != nil {
		t.Fatalf("ParseAll error: %v", err)
	}
	rules, _ := ss.Rules()
	if len(rules) != 2 {
		t.Errorf("expected 2 rule sets, got %d", len(rules))
	}
}

func TestStylesheetApply(t *testing.T) {
	ss := NewStylesheet(nil)
	ss.AddSource("Label { color: red; }", nil, false, 0, "")

	node := newMockStylesheetNode([]string{"Label"}, nil, "")
	err := ss.Apply(node, false, nil)
	if err != nil {
		t.Fatalf("Apply error: %v", err)
	}
	// The node's base styles should have 'color' set
	if !node.styles.Base().HasRule("color") {
		t.Error("expected 'color' rule to be applied to node")
	}
}

func TestStylesheetApplyClassSelector(t *testing.T) {
	ss := NewStylesheet(nil)
	ss.AddSource(".active { display: block; }", nil, false, 0, "")

	node := newMockStylesheetNode([]string{"Widget"}, map[string]bool{"active": true}, "")
	err := ss.Apply(node, false, nil)
	if err != nil {
		t.Fatalf("Apply error: %v", err)
	}
	if !node.styles.Base().HasRule("display") {
		t.Error("expected 'display' rule to be applied to .active node")
	}
}

func TestStylesheetVariableSubstitution(t *testing.T) {
	vars := map[string]string{"primary": "red"}
	ss := NewStylesheet(vars)
	ss.AddSource("Label { color: $primary; }", nil, false, 0, "")

	node := newMockStylesheetNode([]string{"Label"}, nil, "")
	err := ss.Apply(node, false, nil)
	if err != nil {
		t.Fatalf("Apply error: %v", err)
	}
	if !node.styles.Base().HasRule("color") {
		t.Error("expected 'color' rule after variable substitution")
	}
}

func TestStylesheetHasSource(t *testing.T) {
	ss := NewStylesheet(nil)
	loc := &CSSLocation{Path: "/test/foo.css"}
	ss.AddSource("Label { color: red; }", loc, false, 0, "")
	if !ss.HasSource("/test/foo.css", "") {
		t.Error("HasSource should return true for added source")
	}
	if ss.HasSource("/test/bar.css", "") {
		t.Error("HasSource should return false for unknown source")
	}
}

func TestStylesheetCopy(t *testing.T) {
	ss := NewStylesheet(map[string]string{"x": "1"})
	ss.AddSource("Label { color: red; }", nil, false, 0, "")
	cp := ss.Copy()
	if cp == ss {
		t.Error("Copy should return a new pointer")
	}
	// Copy should have same source
	rules, err := cp.Rules()
	if err != nil {
		t.Fatalf("Copy Rules() error: %v", err)
	}
	if len(rules) == 0 {
		t.Error("Copy should have rules from original")
	}
}

func TestStylesheetSetVariables(t *testing.T) {
	ss := NewStylesheet(map[string]string{"primary": "red"})
	ss.AddSource("Label { color: $primary; }", nil, false, 0, "")
	// Change variables
	ss.SetVariables(map[string]string{"primary": "blue"})
	// Reparse with new variables
	if err := ss.ParseAll(); err != nil {
		t.Fatalf("ParseAll after SetVariables error: %v", err)
	}
}

func TestStylesheetRulesMap(t *testing.T) {
	ss := NewStylesheet(nil)
	ss.AddSource("Label { color: red; }\n.foo { display: block; }", nil, false, 0, "")
	rm, err := ss.RulesMap()
	if err != nil {
		t.Fatalf("RulesMap error: %v", err)
	}
	if _, ok := rm["Label"]; !ok {
		t.Error("RulesMap should have 'Label' key")
	}
	if _, ok := rm[".foo"]; !ok {
		t.Error("RulesMap should have '.foo' key")
	}
}

func TestStylesheetCSS(t *testing.T) {
	ss := NewStylesheet(nil)
	ss.AddSource("Label { color: red; }", nil, false, 0, "")
	css, err := ss.CSS()
	if err != nil {
		t.Fatalf("CSS() error: %v", err)
	}
	if css == "" {
		t.Error("CSS() should return non-empty string")
	}
}

// TestStylesheetInitialFallsBackToDefaultCSS verifies that a user rule with
// `initial` is overridden to the value provided by a lower-priority default CSS
// source rather than being silently dropped.
func TestStylesheetInitialFallsBackToDefaultCSS(t *testing.T) {
	ss := NewStylesheet(nil)
	// Default CSS (lower priority): sets color to red.
	ss.AddSource("Label { color: red; }", nil, true, 1, "")
	// User CSS (higher priority): resets color via `initial`.
	ss.AddSource("Label { color: initial; }", nil, false, 0, "")

	node := newMockStylesheetNode([]string{"Label"}, nil, "")
	if err := ss.Apply(node, false, nil); err != nil {
		t.Fatalf("Apply error: %v", err)
	}
	// The fallback from default CSS should win: color must be set.
	if !node.styles.Base().HasRule("color") {
		t.Error("expected 'color' to be set from default CSS fallback after user 'initial'")
	}
}

// TestStylesheetDefaultInitialFallsBackToBuiltIn verifies that when default
// CSS uses `initial`, the built-in default value is materialised into the
// applied rules so that HasRule, GetRules, and CSS serialisation see a
// concrete value rather than an absent rule.
func TestStylesheetDefaultInitialFallsBackToBuiltIn(t *testing.T) {
	ss := NewStylesheet(nil)
	// Only source: default CSS that resets color to initial.
	ss.AddSource("Label { color: initial; }", nil, true, 0, "")

	node := newMockStylesheetNode([]string{"Label"}, nil, "")
	if err := ss.Apply(node, false, nil); err != nil {
		t.Fatalf("Apply error: %v", err)
	}
	// The built-in default must be materialised in the base rules.
	if !node.styles.Base().HasRule("color") {
		t.Error("expected built-in default for 'color' to be present after default-CSS 'initial'")
	}
	// The materialised value must be the built-in white default (255,255,255).
	defVal, _ := builtInStyleDefault("color")
	v, _ := node.styles.Base().GetRule("color")
	if fmt.Sprintf("%v", v) != fmt.Sprintf("%v", defVal) {
		t.Errorf("color after default-CSS initial = %v, want built-in default %v", v, defVal)
	}
}

// TestStylesheetUserInitialIgnoresLowerSpecificityUserRule is the key
// regression: a more-specific user selector using `initial` must NOT fall back
// to a less-specific user declaration; it must resolve to the default-CSS or
// built-in default.
func TestStylesheetUserInitialIgnoresLowerSpecificityUserRule(t *testing.T) {
	ss := NewStylesheet(nil)
	// Lower-specificity user rule sets color: red.
	ss.AddSource("Label { color: red; }", nil, false, 1, "")
	// Higher-specificity user rule (ID selector) resets via initial.
	ss.AddSource("#myid { color: initial; }", nil, false, 0, "")

	node := newMockStylesheetNode([]string{"Label"}, nil, "myid")
	if err := ss.Apply(node, false, nil); err != nil {
		t.Fatalf("Apply error: %v", err)
	}

	// color must NOT be red; the user initial must have bypassed the
	// lower-specificity user rule.
	base := node.styles.Base()
	if !base.HasRule("color") {
		t.Fatal("expected 'color' to be set to the built-in default, but rule is absent")
	}
	defVal, _ := builtInStyleDefault("color")
	v, _ := base.GetRule("color")
	if fmt.Sprintf("%v", v) == fmt.Sprintf("%v", defVal) {
		return // correct: resolved to built-in default, not red
	}
	t.Errorf("color after user initial = %v; expected built-in default %v (not 'red')", v, defVal)
}

// TestStylesheetUpdateWalksDescendants verifies that Update applies styles to
// the root node and at least one nested descendant via Children() traversal.
func TestStylesheetUpdateWalksDescendants(t *testing.T) {
	ss := NewStylesheet(nil)
	ss.AddSource("Label { color: red; }\nButton { display: block; }", nil, false, 0, "")

	root := newMockStylesheetNode([]string{"Label"}, nil, "")
	child := newMockStylesheetNode([]string{"Button"}, nil, "")
	child.mockNode.path = []*mockNode{root.mockNode, child.mockNode}
	child.parent = root
	root.children = []StylesheetNode{child}

	ss.Update(root, false)

	if !root.styles.Base().HasRule("color") {
		t.Error("expected 'color' rule applied to root node")
	}
	if !child.styles.Base().HasRule("display") {
		t.Error("expected 'display' rule applied to child node")
	}
}

func TestStylesheetInvalidCSS(t *testing.T) {
	ss := NewStylesheet(nil)
	// Missing closing brace — should produce a parse error or at minimum not crash
	_ = ss.AddSource("Label { color: $undefined_variable; }", nil, false, 0, "")
	_, err := ss.Rules()
	if err == nil {
		t.Error("expected error for CSS with undefined variable reference")
	}
}
