package css

import "testing"

// mockNode implements SelectorNode and SelectorMatchNode for testing.
type mockNode struct {
	typeNames []string
	classes   map[string]bool
	id        string
	pseudos   map[string]bool
	path      []*mockNode
}

func (n *mockNode) CSSTypeNames() []string      { return n.typeNames }
func (n *mockNode) CSSClasses() map[string]bool { return n.classes }
func (n *mockNode) NodeID() string              { return n.id }
func (n *mockNode) HasClass(name string) bool   { return n.classes[name] }
func (n *mockNode) HasPseudoClasses(set map[string]bool) bool {
	for pc := range set {
		if !n.pseudos[pc] {
			return false
		}
	}
	return true
}

func (n *mockNode) CSSPathNodes() []SelectorNode {
	result := make([]SelectorNode, len(n.path))
	for i, p := range n.path {
		result[i] = p
	}
	return result
}

func newMockNode(typeNames []string, classes map[string]bool, id string) *mockNode {
	if classes == nil {
		classes = make(map[string]bool)
	}
	return &mockNode{typeNames: typeNames, classes: classes, id: id, pseudos: make(map[string]bool)}
}

func TestCheckSelectorsTypeMatch(t *testing.T) {
	node := newMockNode([]string{"Label"}, nil, "")
	node.path = []*mockNode{node}

	sel := Selector{
		Name:        "Label",
		Type:        SelectorType_Type,
		Combinator:  CombinatorDescendant,
		Specificity: Specificity3{0, 0, 1},
		Advance:     1,
	}
	if !CheckSelectors([]Selector{sel}, node.CSSPathNodes()) {
		t.Error("expected Label selector to match Label node")
	}
}

func TestCheckSelectorsTypeMismatch(t *testing.T) {
	node := newMockNode([]string{"Button"}, nil, "")
	node.path = []*mockNode{node}

	sel := Selector{
		Name:        "Label",
		Type:        SelectorType_Type,
		Combinator:  CombinatorDescendant,
		Specificity: Specificity3{0, 0, 1},
		Advance:     1,
	}
	if CheckSelectors([]Selector{sel}, node.CSSPathNodes()) {
		t.Error("expected Label selector to NOT match Button node")
	}
}

func TestCheckSelectorsClassMatch(t *testing.T) {
	node := newMockNode([]string{"Widget"}, map[string]bool{"active": true}, "")
	node.path = []*mockNode{node}

	sel := Selector{
		Name:        "active",
		Type:        SelectorClass,
		Combinator:  CombinatorDescendant,
		Specificity: Specificity3{0, 1, 0},
		Advance:     1,
	}
	if !CheckSelectors([]Selector{sel}, node.CSSPathNodes()) {
		t.Error("expected .active to match node with class 'active'")
	}
}

func TestCheckSelectorsIDMatch(t *testing.T) {
	node := newMockNode([]string{"Widget"}, nil, "main")
	node.path = []*mockNode{node}

	sel := Selector{
		Name:        "main",
		Type:        SelectorID,
		Combinator:  CombinatorDescendant,
		Specificity: Specificity3{1, 0, 0},
		Advance:     1,
	}
	if !CheckSelectors([]Selector{sel}, node.CSSPathNodes()) {
		t.Error("expected #main to match node with id 'main'")
	}
}

func TestCheckSelectorsDescendant(t *testing.T) {
	parent := newMockNode([]string{"Screen"}, nil, "")
	child := newMockNode([]string{"Label"}, nil, "")
	child.path = []*mockNode{parent, child}

	selectors := []Selector{
		{Name: "Screen", Type: SelectorType_Type, Combinator: CombinatorDescendant, Specificity: Specificity3{0, 0, 1}, Advance: 1},
		{Name: "Label", Type: SelectorType_Type, Combinator: CombinatorDescendant, Specificity: Specificity3{0, 0, 1}, Advance: 1},
	}
	if !CheckSelectors(selectors, child.CSSPathNodes()) {
		t.Error("expected 'Screen Label' to match Label inside Screen")
	}
}

func TestCheckSelectorsChild(t *testing.T) {
	parent := newMockNode([]string{"Screen"}, nil, "")
	middle := newMockNode([]string{"Container"}, nil, "")
	child := newMockNode([]string{"Label"}, nil, "")
	child.path = []*mockNode{parent, middle, child}

	// "Screen > Label" should NOT match because Label is not a direct child of Screen
	selectors := []Selector{
		{Name: "Screen", Type: SelectorType_Type, Combinator: CombinatorDescendant, Specificity: Specificity3{0, 0, 1}, Advance: 1},
		{Name: "Label", Type: SelectorType_Type, Combinator: CombinatorChild, Specificity: Specificity3{0, 0, 1}, Advance: 1},
	}
	if CheckSelectors(selectors, child.CSSPathNodes()) {
		t.Error("expected 'Screen > Label' to NOT match Label inside Screen/Container/Label")
	}
}

func TestCheckSelectorsEmptySelectors(t *testing.T) {
	node := newMockNode([]string{"Label"}, nil, "")
	node.path = []*mockNode{node}

	if CheckSelectors([]Selector{}, node.CSSPathNodes()) {
		t.Error("empty selectors should not match")
	}
}

func TestMatchFunction(t *testing.T) {
	node := newMockNode([]string{"Label"}, nil, "")
	node.path = []*mockNode{node}

	ss := SelectorSet{
		Selectors: []Selector{
			{Name: "Label", Type: SelectorType_Type, Combinator: CombinatorDescendant, Specificity: Specificity3{0, 0, 1}, Advance: 1},
		},
	}
	ss.TotalSpecificity()

	if !Match([]SelectorSet{ss}, node) {
		t.Error("Match should return true for matching node")
	}
}
