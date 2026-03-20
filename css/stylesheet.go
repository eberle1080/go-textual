package css

import (
	"fmt"
	"os"
	"strings"
)

// CSSSource holds a CSS string and its metadata.
type CSSSource struct {
	Content    string
	IsDefaults bool
	TieBreaker int
	Scope      string
}

// StylesheetNode is the interface implemented by DOM nodes that can be styled.
type StylesheetNode interface {
	SelectorMatchNode
	// SelectorNames returns the set of selector names that may match this node.
	SelectorNames() map[string]bool
	// Styles returns the node's RenderStyles.
	NodeStyles() *RenderStyles
	// NotifyStyleUpdate is called after styles are updated.
	NotifyStyleUpdate()
	// Refresh causes the node to redraw.
	Refresh()
	// PseudoClassesCacheKey returns an opaque key for pseudo-class caching.
	PseudoClassesCacheKey() any
	// CSSTypeName returns the primary CSS type name.
	CSSTypeName() string
	// Parent returns the parent node, or nil if at the root.
	Parent() StylesheetNode
	// Children returns the direct children of this node in the DOM tree.
	Children() []StylesheetNode
}

// StylesheetParseError is raised when a stylesheet fails to parse.
type StylesheetParseError struct {
	Errors []string
}

func (e *StylesheetParseError) Error() string {
	return "CSS parsing failed: " + strings.Join(e.Errors, "; ")
}

// Stylesheet manages multiple CSS sources and applies them to DOM nodes.
type Stylesheet struct {
	rules          []RuleSet
	rulesMap       map[string][]RuleSet
	variables      map[string]string
	variableTokens map[string][]Token
	source         map[CSSLocation]CSSSource
	requireParse   bool
	invalidCSS     map[string]bool
}

// NewStylesheet creates a new Stylesheet with optional CSS variables.
func NewStylesheet(variables map[string]string) *Stylesheet {
	if variables == nil {
		variables = make(map[string]string)
	}
	return &Stylesheet{
		variables:    variables,
		source:       make(map[CSSLocation]CSSSource),
		invalidCSS:   make(map[string]bool),
	}
}

// getVariableTokens lazily tokenizes the CSS variables.
func (s *Stylesheet) getVariableTokens() (map[string][]Token, error) {
	if s.variableTokens != nil {
		return s.variableTokens, nil
	}
	vt, err := TokenizeValues(s.variables)
	if err != nil {
		return nil, err
	}
	s.variableTokens = vt
	return vt, nil
}

// Rules returns the parsed rule sets, triggering a parse if needed.
func (s *Stylesheet) Rules() ([]RuleSet, error) {
	if s.requireParse {
		if err := s.ParseAll(); err != nil {
			return nil, err
		}
		s.requireParse = false
	}
	return s.rules, nil
}

// RulesMap returns a map from selector name to matching rule sets.
func (s *Stylesheet) RulesMap() (map[string][]RuleSet, error) {
	if s.rulesMap != nil {
		return s.rulesMap, nil
	}
	rules, err := s.Rules()
	if err != nil {
		return nil, err
	}
	rm := make(map[string][]RuleSet)
	for _, rule := range rules {
		for name := range rule.SelectorNames {
			rm[name] = append(rm[name], rule)
		}
	}
	s.rulesMap = rm
	return rm, nil
}

// CSS returns the equivalent TCSS for all parsed rules.
func (s *Stylesheet) CSS() (string, error) {
	rules, err := s.Rules()
	if err != nil {
		return "", err
	}
	var parts []string
	for _, r := range rules {
		parts = append(parts, r.CSS())
	}
	return strings.Join(parts, "\n\n"), nil
}

// Copy creates a copy of the stylesheet (without parsed rule cache).
func (s *Stylesheet) Copy() *Stylesheet {
	vars := make(map[string]string, len(s.variables))
	for k, v := range s.variables {
		vars[k] = v
	}
	ns := NewStylesheet(vars)
	for k, v := range s.source {
		ns.source[k] = v
	}
	if len(ns.source) > 0 {
		ns.requireParse = true
	}
	return ns
}

// SetVariables replaces the CSS variable map and invalidates the parse cache.
func (s *Stylesheet) SetVariables(variables map[string]string) {
	s.variables = variables
	s.variableTokens = nil
	s.invalidCSS = make(map[string]bool)
	s.rulesMap = nil
}

// Read reads and adds a TCSS file to the stylesheet.
func (s *Stylesheet) Read(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return &StylesheetError{Msg: fmt.Sprintf("unable to read CSS file %q: %v", filename, err)}
	}
	absPath, err := absPath(filename)
	if err != nil {
		absPath = filename
	}
	loc := CSSLocation{Path: absPath}
	s.source[loc] = CSSSource{Content: string(data), IsDefaults: false, TieBreaker: 0}
	s.requireParse = true
	return nil
}

// ReadAll reads and adds multiple TCSS files to the stylesheet.
func (s *Stylesheet) ReadAll(paths []string) error {
	for _, p := range paths {
		if err := s.Read(p); err != nil {
			return err
		}
	}
	return nil
}

// HasSource reports whether the stylesheet already has a source at the given location.
func (s *Stylesheet) HasSource(path, classVar string) bool {
	_, ok := s.source[CSSLocation{Path: path, Variable: classVar}]
	return ok
}

// AddSource adds a CSS string to the stylesheet.
func (s *Stylesheet) AddSource(css string, readFrom *CSSLocation, isDefault bool, tieBreaker int, scope string) error {
	var loc CSSLocation
	if readFrom != nil {
		loc = *readFrom
	} else {
		loc = CSSLocation{Variable: fmt.Sprintf("%d", hashString(css))}
	}

	if existing, ok := s.source[loc]; ok && existing.Content == css {
		if existing.TieBreaker > tieBreaker {
			s.source[loc] = CSSSource{Content: existing.Content, IsDefaults: existing.IsDefaults, TieBreaker: tieBreaker, Scope: existing.Scope}
		}
		return nil
	}

	s.source[loc] = CSSSource{Content: css, IsDefaults: isDefault, TieBreaker: tieBreaker, Scope: scope}
	s.requireParse = true
	s.rulesMap = nil
	return nil
}

// ParseAll parses all CSS sources.
func (s *Stylesheet) ParseAll() error {
	vt, err := s.getVariableTokens()
	if err != nil {
		return err
	}

	var rules []RuleSet
	var parseErrors []string

	for loc, src := range s.source {
		if s.invalidCSS[src.Content] {
			continue
		}
		parsed, err := Parse(src.Scope, src.Content, loc, nil, vt, src.IsDefaults, src.TieBreaker)
		if err != nil {
			s.invalidCSS[src.Content] = true
			return &StylesheetParseError{Errors: []string{err.Error()}}
		}
		// Check for rule errors
		for _, rule := range parsed {
			if len(rule.Errors) > 0 {
				for _, re := range rule.Errors {
					parseErrors = append(parseErrors, re.Message)
				}
				s.invalidCSS[src.Content] = true
			}
		}
		rules = append(rules, parsed...)
	}

	if len(parseErrors) > 0 {
		return &StylesheetParseError{Errors: parseErrors}
	}

	s.rules = rules
	s.requireParse = false
	s.rulesMap = nil
	return nil
}

// Reparse re-parses all sources applying current variables.
func (s *Stylesheet) Reparse() error {
	// Parse in a fresh copy to avoid corrupting self on error
	fresh := s.Copy()
	if err := fresh.ParseAll(); err != nil {
		for css := range fresh.invalidCSS {
			s.invalidCSS[css] = true
		}
		return err
	}
	s.rules = fresh.rules
	s.rulesMap = nil
	s.source = fresh.source
	s.requireParse = false
	return nil
}

// checkRule returns the specificities of all matching selector sets for a given node path.
func checkRule(rule RuleSet, pathNodes []SelectorNode) []Specificity3 {
	var result []Specificity3
	for _, ss := range rule.SelectorSets {
		if CheckSelectors(ss.Selectors, pathNodes) {
			result = append(result, ss.Specificity)
		}
	}
	return result
}

// Apply applies the stylesheet rules to a node.
func (s *Stylesheet) Apply(node StylesheetNode, animate bool, cache map[any]RulesMap) error {
	rulesMap, err := s.RulesMap()
	if err != nil {
		return err
	}
	rules, err := s.Rules()
	if err != nil {
		return err
	}

	// Limit candidate rules by selector names, iterating in reverse for priority.
	selectorNames := node.SelectorNames()
	_ = rulesMap // used above for cache warmup only
	var candidates []RuleSet
	for i := len(rules) - 1; i >= 0; i-- {
		rule := rules[i]
		for name := range rule.SelectorNames {
			if selectorNames[name] {
				candidates = append(candidates, rule)
				break
			}
		}
	}

	cssPathNodes := node.CSSPathNodes()

	type ruleEntry struct {
		spec  Specificity6
		value any
	}
	ruleAttributes := make(map[string][]ruleEntry)

	for _, rule := range candidates {
		for _, baseSpec := range checkRule(rule, cssPathNodes) {
			for _, extracted := range rule.Styles.ExtractRules(baseSpec, rule.IsDefaultRules, rule.TieBreaker) {
				ruleAttributes[extracted.Name] = append(ruleAttributes[extracted.Name], ruleEntry{
					spec:  extracted.Specificity,
					value: extracted.Value,
				})
			}
		}
	}

	// For each rule, keep the most specific value, then apply two-stage
	// `initial` resolution mirroring Textual's cascade semantics.
	//
	// Specificity6[0] encodes CSS origin: 1 = user CSS, 0 = default CSS
	// (set by ExtractRules via its defaultFlag).
	//
	//   User `initial`    (best.spec[0]==1, value==nil):
	//     Only fall back to default-CSS entries (spec[0]==0) with a non-nil
	//     value.  Lower-specificity *user* declarations for the same property
	//     are ignored, so a user `initial` cannot be revived by another user
	//     rule with lower priority.  If no default-CSS fallback exists, the
	//     built-in property default is materialised.
	//
	//   Default-rule `initial`  (best.spec[0]==0, value==nil):
	//     Resolve directly to the built-in property default.
	//
	// In both cases the resolved value is stored in nodeRules before
	// ReplaceRules is called, making it visible to HasRule/GetRules/CSSLines.
	nodeRules := make(RulesMap)
	for name, entries := range ruleAttributes {
		best := entries[0]
		for _, e := range entries[1:] {
			if specGreater(e.spec, best.spec) {
				best = e
			}
		}
		if best.value != nil {
			nodeRules[name] = best.value
			continue
		}

		// best.value == nil — the winner used the `initial` keyword.
		isUserInitial := best.spec[0] == 1

		if isUserInitial {
			// Stage 1: look for the best non-nil entry that comes from default CSS.
			var defaultFallback *ruleEntry
			for i := range entries {
				e := &entries[i]
				if e.value != nil && e.spec[0] == 0 {
					if defaultFallback == nil || specGreater(e.spec, defaultFallback.spec) {
						defaultFallback = e
					}
				}
			}
			if defaultFallback != nil {
				nodeRules[name] = defaultFallback.value
				continue
			}
		}

		// Stage 2 (user initial with no default-CSS fallback, or default-rule
		// initial): materialise the package built-in default so downstream APIs
		// observe a concrete value rather than an absent rule.
		if defVal, ok := builtInStyleDefault(name); ok {
			nodeRules[name] = defVal
		}
		// Properties with no built-in default (e.g. width scalar pointer) are
		// intentionally left absent from nodeRules.
	}

	s.ReplaceRules(node, nodeRules, animate)
	return nil
}

// ReplaceRules replaces the node's base styles with the given rules map.
func (s *Stylesheet) ReplaceRules(node StylesheetNode, rules RulesMap, animate bool) {
	ns := node.NodeStyles()
	if ns == nil {
		return
	}
	ns.Base().Reset()
	for k, v := range rules {
		if v != nil {
			ns.Base().Rules[k] = v
		}
	}
	node.NotifyStyleUpdate()
}

// Update applies the stylesheet to a root node and all its descendants by
// performing a depth-first walk of the DOM tree via Children().
func (s *Stylesheet) Update(root StylesheetNode, animate bool) {
	s.walkAndApply(root, animate)
}

// walkAndApply applies styles to node and recursively to all its descendants.
func (s *Stylesheet) walkAndApply(node StylesheetNode, animate bool) {
	_ = s.Apply(node, animate, nil)
	for _, child := range node.Children() {
		s.walkAndApply(child, animate)
	}
}

// UpdateNodes applies the stylesheet to a specific set of nodes.
func (s *Stylesheet) UpdateNodes(nodes []StylesheetNode, animate bool) {
	for _, node := range nodes {
		_ = s.Apply(node, animate, nil)
	}
}

// specGreater reports whether a is strictly greater than b.
func specGreater(a, b Specificity6) bool {
	for i := range a {
		if a[i] > b[i] {
			return true
		}
		if a[i] < b[i] {
			return false
		}
	}
	return false
}

// absPath returns the absolute path for a filename.
func absPath(filename string) (string, error) {
	if len(filename) > 0 && filename[0] == '/' {
		return filename, nil
	}
	wd, err := os.Getwd()
	if err != nil {
		return filename, err
	}
	return wd + "/" + filename, nil
}

// hashString returns a simple hash of a string.
func hashString(s string) int {
	h := 0
	for _, c := range s {
		h = h*31 + int(c)
	}
	return h
}
