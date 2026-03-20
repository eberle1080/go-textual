package css

import "strings"

// SelectorType is the type of a CSS selector.
type SelectorType int

const (
	SelectorUniversal SelectorType = iota + 1 // * operator
	SelectorType_Type                          // CSS type selector e.g. Label
	SelectorClass                              // CSS class selector e.g. .loaded
	SelectorID                                 // CSS ID selector e.g. #main
	SelectorNested                             // Nesting placeholder &
)

// CombinatorType is the type of combinator between two selectors.
type CombinatorType int

const (
	CombinatorSame       CombinatorType = iota + 1 // Selectors are combined
	CombinatorDescendant                            // Descendant combinator (space)
	CombinatorChild                                 // Child combinator >
)

// SelectorNode is the interface implemented by DOM nodes for selector matching.
type SelectorNode interface {
	// CSSTypeNames returns the list of CSS type names for this node.
	CSSTypeNames() []string
	// CSSClasses returns the set of CSS class names on this node.
	CSSClasses() map[string]bool
	// NodeID returns the DOM id of the node (without the # prefix).
	NodeID() string
	// HasPseudoClasses reports whether all pseudo-classes in the set apply to this node.
	HasPseudoClasses(set map[string]bool) bool
	// HasClass reports whether the node has the given CSS class.
	HasClass(name string) bool
}

// Selector represents a single CSS selector component.
type Selector struct {
	Name         string
	Combinator   CombinatorType
	Type         SelectorType
	PseudoClasses map[string]bool
	Specificity  Specificity3
	Advance      int // 1 unless the next selector has CombinatorSame
}

// AddPseudoClass adds a pseudo-class and increments the class specificity.
func (s *Selector) AddPseudoClass(name string) {
	if s.PseudoClasses == nil {
		s.PseudoClasses = make(map[string]bool)
	}
	s.PseudoClasses[name] = true
	s.Specificity[1]++
}

// CSS returns the CSS representation of the selector.
func (s Selector) CSS() string {
	var pseudos []string
	for name := range s.PseudoClasses {
		pseudos = append(pseudos, ":"+name)
	}
	// Sort for determinism
	sortStrings(pseudos)
	suffix := strings.Join(pseudos, "")
	switch s.Type {
	case SelectorUniversal:
		return "*"
	case SelectorType_Type:
		return s.Name + suffix
	case SelectorClass:
		return "." + s.Name + suffix
	case SelectorID:
		return "#" + s.Name + suffix
	default:
		return s.Name + suffix
	}
}

// Check reports whether the selector matches the given node.
func (s Selector) Check(node SelectorNode) bool {
	var nameMatch bool
	switch s.Type {
	case SelectorUniversal, SelectorNested:
		nameMatch = !node.HasClass("-textual-system")
	case SelectorType_Type:
		for _, typeName := range node.CSSTypeNames() {
			if typeName == s.Name {
				nameMatch = true
				break
			}
		}
	case SelectorClass:
		nameMatch = node.HasClass(s.Name)
	case SelectorID:
		nameMatch = node.NodeID() == s.Name
	}
	if !nameMatch {
		return false
	}
	if len(s.PseudoClasses) > 0 {
		return node.HasPseudoClasses(s.PseudoClasses)
	}
	return true
}

// Declaration is a single (unparsed) CSS declaration.
type Declaration struct {
	Token  Token
	Name   string
	Tokens []Token
}

// SelectorSet is a set of selectors with a combined specificity.
type SelectorSet struct {
	Selectors   []Selector
	Specificity Specificity3
}

// TotalSpecificity computes and stores the combined specificity of all selectors.
func (ss *SelectorSet) TotalSpecificity() {
	var id, class, typ int
	for _, sel := range ss.Selectors {
		id += sel.Specificity[0]
		class += sel.Specificity[1]
		typ += sel.Specificity[2]
	}
	ss.Specificity = Specificity3{id, class, typ}
	// Update Advance fields
	for i := 0; i < len(ss.Selectors)-1; i++ {
		if ss.Selectors[i+1].Combinator != CombinatorSame {
			ss.Selectors[i].Advance = 1
		} else {
			ss.Selectors[i].Advance = 0
		}
	}
	if len(ss.Selectors) > 0 {
		ss.Selectors[len(ss.Selectors)-1].Advance = 1
	}
}

// CSS returns the CSS representation of the selector set.
func (ss SelectorSet) CSS() string {
	return selectorsToCSSString(ss.Selectors)
}

// IsSimple reports whether all selectors are simple (no pseudo-classes, only
// type or ID selectors).
func (ss SelectorSet) IsSimple() bool {
	for _, sel := range ss.Selectors {
		if len(sel.PseudoClasses) > 0 {
			return false
		}
		if sel.Type != SelectorType_Type && sel.Type != SelectorID {
			return false
		}
	}
	return true
}

// SelectorSetsFromSelectors builds a slice of SelectorSets from a list of selector lists.
func SelectorSetsFromSelectors(selectors [][]Selector) []SelectorSet {
	result := make([]SelectorSet, 0, len(selectors))
	for _, selList := range selectors {
		var id, class, typ int
		for _, sel := range selList {
			id += sel.Specificity[0]
			class += sel.Specificity[1]
			typ += sel.Specificity[2]
		}
		ss := SelectorSet{Selectors: selList, Specificity: Specificity3{id, class, typ}}
		ss.TotalSpecificity()
		result = append(result, ss)
	}
	return result
}

// RuleError is a CSS error found during rule-set parsing.
type RuleError struct {
	Token   Token
	Message string
}

// RuleSet is a parsed CSS rule set.
type RuleSet struct {
	SelectorSets   []SelectorSet
	Styles         *Styles
	Errors         []RuleError
	IsDefaultRules bool
	TieBreaker     int
	SelectorNames  map[string]bool
	PseudoClasses  map[string]bool
}

// CSS returns the CSS representation of the rule set.
func (rs *RuleSet) CSS() string {
	var selectorParts []string
	for _, ss := range rs.SelectorSets {
		selectorParts = append(selectorParts, ss.CSS())
	}
	selectors := strings.Join(selectorParts, ", ")
	var declLines []string
	if rs.Styles != nil {
		for _, line := range rs.Styles.CSSLines() {
			declLines = append(declLines, "    "+line)
		}
	}
	return selectors + " {\n" + strings.Join(declLines, "\n") + "\n}"
}

// PostParse builds the SelectorNames and PseudoClasses sets after parsing.
func (rs *RuleSet) PostParse() {
	rs.SelectorNames = make(map[string]bool)
	rs.PseudoClasses = make(map[string]bool)
	for _, ss := range rs.SelectorSets {
		for _, sel := range ss.Selectors {
			for pc := range sel.PseudoClasses {
				rs.PseudoClasses[pc] = true
			}
		}
		if len(ss.Selectors) == 0 {
			continue
		}
		last := ss.Selectors[len(ss.Selectors)-1]
		switch last.Type {
		case SelectorUniversal:
			rs.SelectorNames["*"] = true
		case SelectorType_Type:
			rs.SelectorNames[last.Name] = true
		case SelectorClass:
			rs.SelectorNames["."+last.Name] = true
		case SelectorID:
			rs.SelectorNames["#"+last.Name] = true
		}
	}
}

// selectorsToCSSString converts a list of selectors to a CSS string.
func selectorsToCSSString(selectors []Selector) string {
	var parts []string
	for i, sel := range selectors {
		if i == 0 {
			parts = append(parts, sel.CSS())
		} else {
			switch sel.Combinator {
			case CombinatorDescendant:
				parts = append(parts, " "+sel.CSS())
			case CombinatorChild:
				parts = append(parts, " > "+sel.CSS())
			default:
				parts = append(parts, sel.CSS())
			}
		}
	}
	return strings.TrimSpace(strings.Join(parts, ""))
}

// sortStrings sorts a string slice in place.
func sortStrings(ss []string) {
	for i := 1; i < len(ss); i++ {
		for j := i; j > 0 && ss[j] < ss[j-1]; j-- {
			ss[j], ss[j-1] = ss[j-1], ss[j]
		}
	}
}
