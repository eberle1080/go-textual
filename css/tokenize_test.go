package css

import (
	"testing"
)

func tokenNames(tokens []Token) []string {
	names := make([]string, len(tokens))
	for i, t := range tokens {
		names[i] = t.Name
	}
	return names
}

func TestTokenizeSimpleRule(t *testing.T) {
	css := `Label { color: red; }`
	tokens, err := Tokenize(css, CSSLocation{})
	if err != nil {
		t.Fatalf("Tokenize error: %v", err)
	}
	// Should produce selector_start, declaration_set_start, declaration_name, token, declaration_end, declaration_set_end
	names := tokenNames(tokens)
	if len(names) == 0 {
		t.Fatal("no tokens produced")
	}
	// Check that we have a selector_start and declaration_name at minimum
	hasSelectorStart := false
	hasDeclName := false
	for _, n := range names {
		if n == "selector_start" {
			hasSelectorStart = true
		}
		if n == "declaration_name" {
			hasDeclName = true
		}
	}
	if !hasSelectorStart {
		t.Errorf("expected selector_start token, got %v", names)
	}
	if !hasDeclName {
		t.Errorf("expected declaration_name token, got %v", names)
	}
}

func TestTokenizeClassSelector(t *testing.T) {
	css := `.myclass { display: block; }`
	tokens, err := Tokenize(css, CSSLocation{})
	if err != nil {
		t.Fatalf("Tokenize error: %v", err)
	}
	names := tokenNames(tokens)
	hasSelectorStartClass := false
	for _, n := range names {
		if n == "selector_start_class" {
			hasSelectorStartClass = true
		}
	}
	if !hasSelectorStartClass {
		t.Errorf("expected selector_start_class token, got %v", names)
	}
}

func TestTokenizeIDSelector(t *testing.T) {
	css := `#main { color: blue; }`
	tokens, err := Tokenize(css, CSSLocation{})
	if err != nil {
		t.Fatalf("Tokenize error: %v", err)
	}
	names := tokenNames(tokens)
	hasSelectorStartID := false
	for _, n := range names {
		if n == "selector_start_id" {
			hasSelectorStartID = true
		}
	}
	if !hasSelectorStartID {
		t.Errorf("expected selector_start_id token, got %v", names)
	}
}

func TestTokenizeVariableDeclaration(t *testing.T) {
	css := "$primary: red;\nLabel { color: $primary; }"
	tokens, err := Tokenize(css, CSSLocation{})
	if err != nil {
		t.Fatalf("Tokenize error: %v", err)
	}
	names := tokenNames(tokens)
	hasVarName := false
	for _, n := range names {
		if n == "variable_name" {
			hasVarName = true
		}
	}
	if !hasVarName {
		t.Errorf("expected variable_name token, got %v", names)
	}
}

func TestTokenizeDeclarationsBasic(t *testing.T) {
	css := "color: red; display: block;"
	tokens, err := TokenizeDeclarations(css, CSSLocation{})
	if err != nil {
		t.Fatalf("TokenizeDeclarations error: %v", err)
	}
	names := tokenNames(tokens)
	hasDeclName := false
	for _, n := range names {
		if n == "declaration_name" {
			hasDeclName = true
		}
	}
	if !hasDeclName {
		t.Errorf("expected declaration_name token, got %v", names)
	}
}

func TestTokenizeValue(t *testing.T) {
	tokens, err := TokenizeValue("red", CSSLocation{})
	if err != nil {
		t.Fatalf("TokenizeValue error: %v", err)
	}
	if len(tokens) == 0 {
		t.Fatal("no tokens from TokenizeValue")
	}
}

// TestTokenizeInvalidLeadingChar verifies that a malformed declaration whose
// first character is not a valid token start is rejected with an error rather
// than silently skipped.  Prior to the anchored-match fix, FindStringSubmatch
// could advance past the bad character and produce a mis-located token.
func TestTokenizeInvalidLeadingChar(t *testing.T) {
	// "@" is not a valid leading character for a declaration name.
	_, err := TokenizeDeclarations("@color: red;", CSSLocation{})
	if err == nil {
		t.Error("expected error for declaration starting with '@', got nil")
	}
}

// TestTokenizeErrorLocation verifies that the error position reported when a
// bad character is encountered points to the actual cursor column rather than
// an offset that would only be correct after silently skipping characters.
func TestTokenizeErrorLocation(t *testing.T) {
	// "  !" — two spaces followed by "!" which is not a valid declaration name.
	// A correct anchored match should fail immediately at column 3 (1-based),
	// not at some later position.
	_, err := TokenizeDeclarations("  !bad: value;", CSSLocation{})
	if err == nil {
		t.Error("expected error for invalid leading '!', got nil")
	}
}

func TestTokenizeValues(t *testing.T) {
	vars := map[string]string{
		"primary": "red",
		"accent":  "#ff0000",
	}
	vt, err := TokenizeValues(vars)
	if err != nil {
		t.Fatalf("TokenizeValues error: %v", err)
	}
	if _, ok := vt["primary"]; !ok {
		t.Errorf("TokenizeValues missing 'primary' key")
	}
	if _, ok := vt["accent"]; !ok {
		t.Errorf("TokenizeValues missing 'accent' key")
	}
}
