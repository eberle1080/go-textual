package css

// Pattern constants used to build tokenizer Expect objects.
const (
	PatternPercent         = `-?\d+\.?\d*%`
	PatternDecimal         = `-?\d+\.?\d*`
	PatternComma           = `\s*,\s*`
	PatternOpenBrace       = `\(\s*`
	PatternCloseBrace      = `\s*\)`
	PatternHexColor        = `\#[0-9a-fA-F]{8}|\#[0-9a-fA-F]{6}|\#[0-9a-fA-F]{4}|\#[0-9a-fA-F]{3}`
	PatternRGBColor        = `rgb\(\s*` + PatternDecimal + `\s*,\s*` + PatternDecimal + `\s*,\s*` + PatternDecimal + `\s*\)|rgba\(\s*` + PatternDecimal + `\s*,\s*` + PatternDecimal + `\s*,\s*` + PatternDecimal + `\s*,\s*` + PatternDecimal + `\s*\)`
	PatternHSLColor        = `hsl\(\s*` + PatternDecimal + `\s*,\s*` + PatternPercent + `\s*,\s*` + PatternPercent + `\s*\)|hsla\(\s*` + PatternDecimal + `\s*,\s*` + PatternPercent + `\s*,\s*` + PatternPercent + `\s*,\s*` + PatternDecimal + `\s*\)`
	PatternScalar          = PatternDecimal + `(?:fr|%|w|h|vw|vh)`
	PatternDuration        = `\d+\.?\d*(?:ms|s)`
	PatternNumber          = `\-?\d+\.?\d*`
	PatternColor           = PatternHexColor + `|` + PatternRGBColor + `|` + PatternHSLColor
	PatternKeyValue        = `[a-zA-Z_-][a-zA-Z0-9_-]*=[0-9a-zA-Z_\-\/]+`
	PatternToken           = `[a-zA-Z_][a-zA-Z0-9_-]*`
	PatternString          = `".*?"`
	PatternVariableRef     = `\$[a-zA-Z0-9_\-]+`
	PatternIdentifier      = `[a-zA-Z_\-][a-zA-Z0-9_\-]*`
	PatternSelectorTypeName = `[A-Z_][a-zA-Z0-9_]*`
	PatternDeclarationName = `[a-z][a-zA-Z0-9_\-]*`
)

// declarationValues is the ordered set of token patterns valid inside a declaration.
var declarationValues = [][2]string{
	{"scalar", PatternScalar},
	{"duration", PatternDuration},
	{"number", PatternNumber},
	{"color", PatternColor},
	{"key_value", PatternKeyValue},
	{"token", PatternToken},
	{"string", PatternString},
	{"variable_ref", PatternVariableRef},
}

// Package-level Expect variables for the tokenizer state machine.
var (
	ExpectRootScope = buildRootScope()
	ExpectRootNested = buildRootNested()
	ExpectVariableNameContinue = buildVariableNameContinue()
	ExpectCommentEnd = NewExpect("comment end", [][2]string{
		{"comment_end", `\*/`},
	})
	ExpectSelectorContinue = buildSelectorContinue()
	ExpectDeclaration      = buildDeclaration()
	ExpectDeclarationSolo  = buildDeclarationSolo()
	ExpectDeclarationContent     = buildDeclarationContent()
	ExpectDeclarationContentSolo = buildDeclarationContentSolo()
)

func buildRootScope() *Expect {
	e := NewExpect("selector or end of file", [][2]string{
		{"whitespace", `\s+`},
		{"comment_start", `\/\*`},
		{"comment_line", `\# .*$`},
		{"selector_start_id", `\#` + PatternIdentifier},
		{"selector_start_class", `\.` + PatternIdentifier},
		{"selector_start_universal", `\*`},
		{"selector_start", PatternSelectorTypeName},
		{"variable_name", `\$[a-zA-Z0-9_\-]+:`},
		{"declaration_set_end", `\}`},
	})
	return e.WithEOF(true)
}

func buildRootNested() *Expect {
	return NewExpect("selector or end of file", [][2]string{
		{"whitespace", `\s+`},
		{"comment_start", `\/\*`},
		{"comment_line", `\# .*$`},
		{"declaration_name", PatternDeclarationName + `\:`},
		{"selector_start_id", `\#` + PatternIdentifier},
		{"selector_start_class", `\.` + PatternIdentifier},
		{"selector_start_universal", `\*`},
		{"selector_start", PatternSelectorTypeName},
		{"variable_name", `\$[a-zA-Z0-9_\-]+:`},
		{"declaration_set_end", `\}`},
		{"nested", `\&`},
	})
}

func buildVariableNameContinue() *Expect {
	tokens := [][2]string{
		{"variable_value_end", `\n|;`},
		{"whitespace", `\s+`},
		{"comment_start", `\/\*`},
		{"comment_line", `\# .*$`},
	}
	tokens = append(tokens, declarationValues...)
	e := NewExpect("variable value", tokens)
	return e.WithEOF(true)
}

func buildSelectorContinue() *Expect {
	e := NewExpect("selector or {", [][2]string{
		{"whitespace", `\s+`},
		{"comment_start", `\/\*`},
		{"comment_line", `\# .*$`},
		{"pseudo_class", `\:[a-zA-Z_-]+`},
		{"selector_id", `\#` + PatternIdentifier},
		{"selector_class", `\.` + PatternIdentifier},
		{"selector_universal", `\*`},
		{"selector", PatternSelectorTypeName},
		{"combinator_child", `>`},
		{"new_selector", `,`},
		{"declaration_set_start", `\{`},
		{"declaration_set_end", `\}`},
		{"nested", `\&`},
	})
	return e.WithEOF(true)
}

func buildDeclaration() *Expect {
	return NewExpect("rule or selector", [][2]string{
		{"nested", `\&`},
		{"whitespace", `\s+`},
		{"comment_start", `\/\*`},
		{"comment_line", `\# .*$`},
		{"declaration_name", PatternDeclarationName + `\:`},
		{"declaration_set_end", `\}`},
		{"selector_start_id", `\#` + PatternIdentifier},
		{"selector_start_class", `\.` + PatternIdentifier},
		{"selector_start_universal", `\*`},
		{"selector_start", PatternSelectorTypeName},
	})
}

func buildDeclarationSolo() *Expect {
	e := NewExpect("rule declaration", [][2]string{
		{"whitespace", `\s+`},
		{"comment_start", `\/\*`},
		{"comment_line", `\# .*$`},
		{"declaration_name", PatternDeclarationName + `\:`},
		{"declaration_set_end", `\}`},
	})
	return e.WithEOF(true)
}

func buildDeclarationContent() *Expect {
	tokens := [][2]string{
		{"declaration_end", `;`},
		{"whitespace", `\s+`},
		{"comment_start", `\/\*`},
		{"comment_line", `\# .*$`},
	}
	tokens = append(tokens, declarationValues...)
	tokens = append(tokens,
		[2]string{"important", `\!important`},
		[2]string{"comma", `,`},
		[2]string{"declaration_set_end", `\}`},
	)
	return NewExpect("rule value or end of declaration", tokens)
}

func buildDeclarationContentSolo() *Expect {
	tokens := [][2]string{
		{"declaration_end", `;`},
		{"whitespace", `\s+`},
		{"comment_start", `\/\*`},
		{"comment_line", `\# .*$`},
	}
	tokens = append(tokens, declarationValues...)
	tokens = append(tokens,
		[2]string{"important", `\!important`},
		[2]string{"comma", `,`},
		[2]string{"declaration_set_end", `\}`},
	)
	e := NewExpect("rule value or end of declaration", tokens)
	return e.WithEOF(true)
}

// Tokenize tokenizes a full TCSS stylesheet string.
func Tokenize(code string, readFrom CSSLocation) ([]Token, error) {
	return tcssTokenize(code, readFrom)
}

// TokenizeDeclarations tokenizes a CSS declarations block (no selectors).
func TokenizeDeclarations(code string, readFrom CSSLocation) ([]Token, error) {
	return declarationTokenize(code, readFrom)
}

// TokenizeValue tokenizes a single declaration value string.
func TokenizeValue(code string, readFrom CSSLocation) ([]Token, error) {
	return valueTokenize(code, readFrom)
}

// TokenizeValues tokenizes a map of variable name -> value strings.
func TokenizeValues(values map[string]string) (map[string][]Token, error) {
	result := make(map[string][]Token, len(values))
	for name, val := range values {
		tokens, err := TokenizeValue(val, CSSLocation{Path: "__name__"})
		if err != nil {
			return nil, err
		}
		result[name] = tokens
	}
	return result, nil
}

// tcssTokenize implements the TCSSTokenizerState state machine.
func tcssTokenize(code string, readFrom CSSLocation) ([]Token, error) {
	tokenizer := NewTokenizer(code, readFrom)
	stateMap := map[string]*Expect{
		"variable_name":             ExpectVariableNameContinue,
		"variable_value_end":        ExpectRootScope,
		"selector_start":            ExpectSelectorContinue,
		"selector_start_id":         ExpectSelectorContinue,
		"selector_start_class":      ExpectSelectorContinue,
		"selector_start_universal":  ExpectSelectorContinue,
		"selector_id":               ExpectSelectorContinue,
		"selector_class":            ExpectSelectorContinue,
		"selector_universal":        ExpectSelectorContinue,
		"declaration_set_start":     ExpectDeclaration,
		"declaration_name":          ExpectDeclarationContent,
		"declaration_end":           ExpectDeclaration,
		"declaration_set_end":       ExpectRootNested,
		"nested":                    ExpectSelectorContinue,
	}

	expect := ExpectRootScope
	nestLevel := 0
	var tokens []Token

	for {
		token, err := tokenizer.GetToken(expect)
		if err != nil {
			return nil, err
		}
		name := token.Name

		if name == "eof" {
			break
		}
		switch name {
		case "comment_line":
			continue
		case "comment_start":
			if _, err := tokenizer.SkipTo(ExpectCommentEnd); err != nil {
				return nil, err
			}
			continue
		case "declaration_set_start":
			nestLevel++
		case "declaration_set_end":
			nestLevel--
			if nestLevel > 0 {
				expect = ExpectDeclaration
			} else {
				expect = ExpectRootScope
			}
			tokens = append(tokens, token)
			continue
		}

		if next, ok := stateMap[name]; ok {
			expect = next
		}
		tokens = append(tokens, token)
	}
	return tokens, nil
}

// declarationTokenize implements the DeclarationTokenizerState state machine.
func declarationTokenize(code string, readFrom CSSLocation) ([]Token, error) {
	tokenizer := NewTokenizer(code, readFrom)
	stateMap := map[string]*Expect{
		"declaration_name": ExpectDeclarationContent,
		"declaration_end":  ExpectDeclarationSolo,
	}
	expect := ExpectDeclarationSolo
	var tokens []Token

	for {
		token, err := tokenizer.GetToken(expect)
		if err != nil {
			return nil, err
		}
		name := token.Name
		if name == "eof" {
			break
		}
		if next, ok := stateMap[name]; ok {
			expect = next
		}
		tokens = append(tokens, token)
	}
	return tokens, nil
}

// valueTokenize implements the ValueTokenizerState state machine.
func valueTokenize(code string, readFrom CSSLocation) ([]Token, error) {
	tokenizer := NewTokenizer(code, readFrom)
	expect := ExpectDeclarationContentSolo
	var tokens []Token

	for {
		token, err := tokenizer.GetToken(expect)
		if err != nil {
			return nil, err
		}
		if token.Name == "eof" {
			break
		}
		tokens = append(tokens, token)
	}
	return tokens, nil
}
