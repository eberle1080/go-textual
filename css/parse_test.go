package css

import (
	"strings"
	"testing"
)

func TestParseSelectors(t *testing.T) {
	tests := []struct {
		input     string
		wantCount int
		wantErr   bool
	}{
		{"Label", 1, false},
		{".myclass", 1, false},
		{"#main", 1, false},
		{"Label .child", 1, false},
		{"Label, .other", 2, false},
		{"Label > .child", 1, false},
		{"", 0, false},
	}
	for _, tt := range tests {
		ss, err := ParseSelectors(tt.input)
		if tt.wantErr {
			if err == nil {
				t.Errorf("ParseSelectors(%q) expected error", tt.input)
			}
			continue
		}
		if err != nil {
			t.Errorf("ParseSelectors(%q) unexpected error: %v", tt.input, err)
			continue
		}
		if len(ss) != tt.wantCount {
			t.Errorf("ParseSelectors(%q) = %d selector sets, want %d", tt.input, len(ss), tt.wantCount)
		}
	}
}

func TestParseDeclarations(t *testing.T) {
	tests := []struct {
		css      string
		wantRule string
	}{
		{"display: block;", "display"},
		{"color: red;", "color"},
		{"width: 100;", "width"},
	}
	for _, tt := range tests {
		styles, err := ParseDeclarations(tt.css, CSSLocation{})
		if err != nil {
			t.Errorf("ParseDeclarations(%q) error: %v", tt.css, err)
			continue
		}
		if styles == nil {
			t.Errorf("ParseDeclarations(%q) returned nil styles", tt.css)
			continue
		}
		if !styles.HasRule(tt.wantRule) {
			t.Errorf("ParseDeclarations(%q) missing rule %q", tt.css, tt.wantRule)
		}
	}
}

func TestParseFullCSS(t *testing.T) {
	css := `Label {
    color: red;
    display: block;
}`
	ruleSets, err := Parse("", css, CSSLocation{}, nil, nil, false, 0)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if len(ruleSets) == 0 {
		t.Fatal("Parse returned no rule sets")
	}
	rs := ruleSets[0]
	if len(rs.SelectorSets) == 0 {
		t.Fatal("RuleSet has no selector sets")
	}
	if rs.Styles == nil {
		t.Fatal("RuleSet has nil Styles")
	}
}

func TestParseMultipleRules(t *testing.T) {
	css := `
Label { color: red; }
.myclass { display: block; }
`
	ruleSets, err := Parse("", css, CSSLocation{}, nil, nil, false, 0)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if len(ruleSets) != 2 {
		t.Errorf("Parse returned %d rule sets, want 2", len(ruleSets))
	}
}

func TestParseWithVariables(t *testing.T) {
	css := `Label { color: $primary; }`
	vars := map[string]string{"primary": "red"}
	ruleSets, err := Parse("", css, CSSLocation{}, vars, nil, false, 0)
	if err != nil {
		t.Fatalf("Parse with variables error: %v", err)
	}
	if len(ruleSets) == 0 {
		t.Fatal("Parse returned no rule sets")
	}
}

func TestParseWithScope(t *testing.T) {
	css := `Label { color: red; }`
	ruleSets, err := Parse("Screen", css, CSSLocation{}, nil, nil, false, 0)
	if err != nil {
		t.Fatalf("Parse with scope error: %v", err)
	}
	if len(ruleSets) == 0 {
		t.Fatal("Parse returned no rule sets")
	}
	// The first selector should be the scope
	rs := ruleSets[0]
	if len(rs.SelectorSets) == 0 || len(rs.SelectorSets[0].Selectors) < 2 {
		t.Errorf("expected scoped selectors, got %v", rs.SelectorSets)
	}
	first := rs.SelectorSets[0].Selectors[0]
	if first.Name != "Screen" {
		t.Errorf("scope selector name = %q, want %q", first.Name, "Screen")
	}
}

func TestSubstituteReferences(t *testing.T) {
	vars := map[string][]Token{
		"primary": {{Name: "token", Value: "red"}},
	}
	tokens := []Token{
		{Name: "variable_ref", Value: "$primary"},
	}
	result, err := SubstituteReferences(tokens, vars)
	if err != nil {
		t.Fatalf("SubstituteReferences error: %v", err)
	}
	if len(result) != 1 || result[0].Value != "red" {
		t.Errorf("SubstituteReferences = %v, want [{red}]", result)
	}
}

func TestSubstituteReferencesUnresolved(t *testing.T) {
	vars := map[string][]Token{}
	tokens := []Token{
		{Name: "variable_ref", Value: "$missing"},
	}
	_, err := SubstituteReferences(tokens, vars)
	if err == nil {
		t.Error("SubstituteReferences expected error for missing variable")
	}
	if !strings.Contains(err.Error(), "missing") {
		t.Errorf("error message should mention 'missing', got: %v", err)
	}
}

func TestRuleSetCSS(t *testing.T) {
	css := `Label { color: red; }`
	ruleSets, err := Parse("", css, CSSLocation{}, nil, nil, false, 0)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if len(ruleSets) == 0 {
		t.Fatal("no rule sets")
	}
	cssStr := ruleSets[0].CSS()
	if !strings.Contains(cssStr, "Label") {
		t.Errorf("RuleSet.CSS() missing 'Label', got: %q", cssStr)
	}
}

func TestCombineSelectorsNested(t *testing.T) {
	parent := []Selector{
		{Name: "Label", Type: SelectorType_Type, Combinator: CombinatorDescendant, Specificity: Specificity3{0, 0, 1}, Advance: 1},
	}
	child := []Selector{
		{Name: "&", Type: SelectorNested, Combinator: CombinatorDescendant, Specificity: Specificity3{0, 0, 0}},
		{Name: "focused", Type: SelectorClass, Combinator: CombinatorSame, Specificity: Specificity3{0, 1, 0}, Advance: 1},
	}
	result := combineSelectors(parent, child)
	if len(result) != 2 {
		t.Errorf("combineSelectors result len = %d, want 2", len(result))
	}
}
