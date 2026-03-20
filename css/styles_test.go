package css

import "testing"

func TestStylesSetAndGet(t *testing.T) {
	s := NewStyles()
	if s.HasRule("color") {
		t.Error("expected no 'color' rule initially")
	}
	s.SetRule("color", "red")
	if !s.HasRule("color") {
		t.Error("expected 'color' rule after SetRule")
	}
	v, ok := s.GetRule("color")
	if !ok || v != "red" {
		t.Errorf("GetRule('color') = %v, %v; want 'red', true", v, ok)
	}
}

func TestStylesClearRule(t *testing.T) {
	s := NewStyles()
	s.SetRule("display", "block")
	removed := s.ClearRule("display")
	if !removed {
		t.Error("ClearRule should return true when rule existed")
	}
	if s.HasRule("display") {
		t.Error("rule should be gone after ClearRule")
	}
	removed2 := s.ClearRule("display")
	if removed2 {
		t.Error("ClearRule should return false when rule doesn't exist")
	}
}

func TestStylesReset(t *testing.T) {
	s := NewStyles()
	s.SetRule("color", "red")
	s.SetRule("display", "block")
	s.Reset()
	if s.HasRule("color") || s.HasRule("display") {
		t.Error("Reset should clear all rules")
	}
}

func TestStylesCopy(t *testing.T) {
	s := NewStyles()
	s.SetRule("color", "red")
	c := s.Copy()
	if !c.HasRule("color") {
		t.Error("Copy should include all rules")
	}
	// Mutating original should not affect copy
	s.SetRule("color", "blue")
	v, _ := c.GetRule("color")
	if v != "red" {
		t.Errorf("Copy should be independent, got color=%v", v)
	}
}

func TestStylesMerge(t *testing.T) {
	s1 := NewStyles()
	s1.SetRule("color", "red")
	s2 := NewStyles()
	s2.SetRule("display", "block")
	s1.Merge(s2)
	if !s1.HasRule("display") {
		t.Error("Merge should add rules from other")
	}
	if !s1.HasRule("color") {
		t.Error("Merge should preserve existing rules")
	}
}

func TestStylesMergeRules(t *testing.T) {
	s := NewStyles()
	rules := RulesMap{"color": "green", "width": 100}
	s.MergeRules(rules)
	if !s.HasRule("color") || !s.HasRule("width") {
		t.Error("MergeRules should add all rules")
	}
}

func TestStylesGetRules(t *testing.T) {
	s := NewStyles()
	s.SetRule("color", "red")
	rules := s.GetRules()
	if rules["color"] != "red" {
		t.Errorf("GetRules missing 'color'")
	}
	// Modifying copy should not affect original
	rules["color"] = "blue"
	v, _ := s.GetRule("color")
	if v != "red" {
		t.Error("GetRules should return a copy")
	}
}

func TestStylesExtractRules(t *testing.T) {
	s := NewStyles()
	s.SetRule("color", "red")
	s.SetRule("display", "block")
	spec := Specificity3{0, 1, 0}
	extracted := s.ExtractRules(spec, false, 0)
	if len(extracted) != 2 {
		t.Errorf("ExtractRules returned %d rules, want 2", len(extracted))
	}
	for _, e := range extracted {
		if e.Name != "color" && e.Name != "display" {
			t.Errorf("unexpected rule name %q", e.Name)
		}
	}
}

func TestStylesExtractRulesDefault(t *testing.T) {
	s := NewStyles()
	s.SetRule("color", "blue")
	normalExtracted := s.ExtractRules(Specificity3{0, 0, 1}, false, 0)
	defaultExtracted := s.ExtractRules(Specificity3{0, 0, 1}, true, 0)
	// Default rules have lower specificity (defaultFlag=0 vs 1)
	if len(normalExtracted) == 0 || len(defaultExtracted) == 0 {
		t.Fatal("expected extracted rules")
	}
	normal := normalExtracted[0].Specificity
	def := defaultExtracted[0].Specificity
	// Normal (defaultFlag=1) should be > default (defaultFlag=0) in some dimension
	allEqual := true
	for i := range normal {
		if normal[i] != def[i] {
			allEqual = false
			break
		}
	}
	if allEqual {
		t.Errorf("normal specificity %v should differ from default specificity %v", normal, def)
	}
}

func TestStylesTextStyleAccessors(t *testing.T) {
	s := NewStyles()
	if s.TextStyle() != "" {
		t.Errorf("TextStyle() default = %q, want empty", s.TextStyle())
	}
	s.SetTextStyle("bold")
	if s.TextStyle() != "bold" {
		t.Errorf("TextStyle() = %q, want bold", s.TextStyle())
	}
}

func TestStylesLinkStyleAccessors(t *testing.T) {
	s := NewStyles()
	s.SetLinkStyle("underline")
	if s.LinkStyle() != "underline" {
		t.Errorf("LinkStyle() = %q, want underline", s.LinkStyle())
	}
}

func TestStylesLinkStyleHoverAccessors(t *testing.T) {
	s := NewStyles()
	s.SetLinkStyleHover("bold underline")
	if s.LinkStyleHover() != "bold underline" {
		t.Errorf("LinkStyleHover() = %q, want 'bold underline'", s.LinkStyleHover())
	}
}

func TestStylesBorderTitleStyleAccessors(t *testing.T) {
	s := NewStyles()
	s.SetBorderTitleStyle("italic")
	if s.BorderTitleStyle() != "italic" {
		t.Errorf("BorderTitleStyle() = %q, want italic", s.BorderTitleStyle())
	}
}

func TestStylesBorderSubtitleStyleAccessors(t *testing.T) {
	s := NewStyles()
	s.SetBorderSubtitleStyle("bold")
	if s.BorderSubtitleStyle() != "bold" {
		t.Errorf("BorderSubtitleStyle() = %q, want bold", s.BorderSubtitleStyle())
	}
}

func TestStylesTextStyleCSSLines(t *testing.T) {
	s := NewStyles()
	s.SetTextStyle("bold")
	s.SetLinkStyle("underline")
	lines := s.CSSLines()
	found := make(map[string]bool)
	for _, l := range lines {
		found[l] = true
	}
	if !found["text-style: bold;"] {
		t.Errorf("CSSLines missing text-style; got %v", lines)
	}
	if !found["link-style: underline;"] {
		t.Errorf("CSSLines missing link-style; got %v", lines)
	}
}

func TestStylesCSSLines(t *testing.T) {
	css := `Label { color: red; }`
	ruleSets, err := Parse("", css, CSSLocation{}, nil, nil, false, 0)
	if err != nil || len(ruleSets) == 0 {
		t.Fatalf("Parse error or no rule sets: %v", err)
	}
	lines := ruleSets[0].Styles.CSSLines()
	if len(lines) == 0 {
		t.Error("CSSLines should return at least one line")
	}
}
