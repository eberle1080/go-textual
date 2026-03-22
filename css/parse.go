package css

import (
	"strings"
)

// SelectorMapEntry maps a token name to a selector type and specificity.
type SelectorMapEntry struct {
	Type        SelectorType
	Specificity Specificity3
}

// selectorTokenMap maps tokenizer token names to selector types and specificities.
var selectorTokenMap = map[string]SelectorMapEntry{
	"selector":                 {SelectorType_Type, Specificity3{0, 0, 1}},
	"selector_start":           {SelectorType_Type, Specificity3{0, 0, 1}},
	"selector_class":           {SelectorClass, Specificity3{0, 1, 0}},
	"selector_start_class":     {SelectorClass, Specificity3{0, 1, 0}},
	"selector_id":              {SelectorID, Specificity3{1, 0, 0}},
	"selector_start_id":        {SelectorID, Specificity3{1, 0, 0}},
	"selector_universal":       {SelectorUniversal, Specificity3{0, 0, 0}},
	"selector_start_universal": {SelectorUniversal, Specificity3{0, 0, 0}},
	"nested":                   {SelectorNested, Specificity3{0, 0, 0}},
}

// TokenIterator is a simple iterator over a slice of tokens.
type TokenIterator struct {
	tokens []Token
	pos    int
}

// NewTokenIterator creates a new TokenIterator.
func NewTokenIterator(tokens []Token) *TokenIterator {
	return &TokenIterator{tokens: tokens}
}

// Next returns the next token and true, or a zero Token and false at the end.
func (it *TokenIterator) Next() (Token, bool) {
	if it.pos >= len(it.tokens) {
		return Token{}, false
	}
	t := it.tokens[it.pos]
	it.pos++
	return t, true
}

// Peek returns the next token without consuming it.
func (it *TokenIterator) Peek() (Token, bool) {
	if it.pos >= len(it.tokens) {
		return Token{}, false
	}
	return it.tokens[it.pos], true
}

// ParseSelectors parses a CSS selector string into SelectorSets.
func ParseSelectors(cssSelectors string) ([]SelectorSet, error) {
	if strings.TrimSpace(cssSelectors) == "" {
		return nil, nil
	}
	tokens, err := Tokenize(cssSelectors, CSSLocation{})
	if err != nil {
		return nil, err
	}

	var combinator *CombinatorType
	desc := CombinatorDescendant
	same := CombinatorSame
	combinator = &desc

	var selectors []Selector
	var ruleSelectorLists [][]Selector

	for _, tok := range tokens {
		switch tok.Name {
		case "pseudo_class":
			if len(selectors) > 0 {
				selectors[len(selectors)-1].AddPseudoClass(strings.TrimLeft(tok.Value, ":"))
			}
		case "whitespace":
			if combinator == nil || *combinator == CombinatorSame {
				combinator = &desc
			}
		case "new_selector":
			ruleSelectorLists = append(ruleSelectorLists, append([]Selector{}, selectors...))
			selectors = selectors[:0]
			combinator = nil
		case "declaration_set_start":
			goto done
		case "combinator_child":
			child := CombinatorChild
			combinator = &child
		default:
			entry, ok := selectorTokenMap[tok.Name]
			if !ok {
				continue
			}
			comb := desc
			if combinator != nil {
				comb = *combinator
			}
			sel := Selector{
				Name:        strings.TrimLeft(tok.Value, ".#"),
				Combinator:  comb,
				Type:        entry.Type,
				Specificity: entry.Specificity,
				Advance:     1,
			}
			selectors = append(selectors, sel)
			combinator = &same
		}
	}
done:
	if len(selectors) > 0 {
		ruleSelectorLists = append(ruleSelectorLists, selectors)
	}
	return SelectorSetsFromSelectors(ruleSelectorLists), nil
}

// ParseDeclarations parses inline CSS declarations and returns a Styles object.
func ParseDeclarations(css string, readFrom CSSLocation) (*Styles, error) {
	tokens, err := TokenizeDeclarations(css, readFrom)
	if err != nil {
		return nil, err
	}

	builder := NewStylesBuilder()
	var decl *Declaration

	for _, tok := range tokens {
		switch tok.Name {
		case "whitespace", "declaration_end", "eof":
			continue
		case "declaration_name":
			if decl != nil {
				if err := builder.AddDeclaration(*decl); err != nil {
					return nil, err
				}
			}
			d := Declaration{Token: tok, Name: strings.TrimRight(tok.Value, ":")}
			decl = &d
		case "declaration_set_end":
			goto parseDone
		default:
			if decl != nil {
				decl.Tokens = append(decl.Tokens, tok)
			}
		}
	}
parseDone:
	if decl != nil {
		if err := builder.AddDeclaration(*decl); err != nil {
			return nil, err
		}
	}
	return builder.Styles, nil
}

// SubstituteReferences replaces variable_ref tokens with their values.
func SubstituteReferences(tokens []Token, variables map[string][]Token) ([]Token, error) {
	vars := make(map[string][]Token)
	for k, v := range variables {
		vars[k] = v
	}

	var result []Token
	i := 0
	for i < len(tokens) {
		tok := tokens[i]
		i++

		if tok.Name == "variable_name" {
			// Variable definition: trim "$" and ":"
			varName := strings.TrimRight(strings.TrimLeft(tok.Value, "$"), ":")
			varTokens := vars[varName]
			result = append(result, tok)

			// Skip whitespace after variable_name
			for i < len(tokens) && tokens[i].Name == "whitespace" {
				result = append(result, tokens[i])
				i++
			}

			// Collect variable value tokens until variable_value_end
			for i < len(tokens) {
				t := tokens[i]
				i++
				if t.Name == "whitespace" {
					varTokens = append(varTokens, t)
					result = append(result, t)
				} else if t.Name == "variable_value_end" {
					result = append(result, t)
					break
				} else if t.Name == "variable_ref" {
					refName := strings.TrimLeft(t.Value, "$")
					if refToks, ok := vars[refName]; ok {
						varTokens = append(varTokens, refToks...)
						loc := t.Location
						length := len(t.Value)
						for _, rt := range refToks {
							result = append(result, rt.WithReference(&ReferencedBy{
								Name:     refName,
								Location: loc,
								Length:   length,
								Code:     t.Code,
							}))
						}
					} else {
						return nil, &UnresolvedVariableError{
							TokenError: TokenError{
								ReadFrom: t.ReadFrom,
								Code:     t.Code,
								Start:    t.Start(),
								End:      t.End(),
								Msg:      "reference to undefined variable '$" + refName + "'",
							},
						}
					}
				} else {
					varTokens = append(varTokens, t)
					result = append(result, t)
				}
			}
			vars[varName] = varTokens

		} else if tok.Name == "variable_ref" {
			varName := strings.TrimLeft(tok.Value, "$")
			if refToks, ok := vars[varName]; ok {
				loc := tok.Location
				length := len(tok.Value)
				for _, rt := range refToks {
					result = append(result, rt.WithReference(&ReferencedBy{
						Name:     varName,
						Location: loc,
						Length:   length,
						Code:     tok.Code,
					}))
				}
			} else {
				return nil, &UnresolvedVariableError{
					TokenError: TokenError{
						ReadFrom: tok.ReadFrom,
						Code:     tok.Code,
						Start:    tok.Start(),
						End:      tok.End(),
						Msg:      "reference to undefined variable '$" + varName + "'",
					},
				}
			}
		} else {
			result = append(result, tok)
		}
	}
	return result, nil
}

// Parse tokenizes and parses a full TCSS stylesheet, returning a list of RuleSets.
func Parse(
	scope string,
	css string,
	readFrom CSSLocation,
	variables map[string]string,
	variableTokens map[string][]Token,
	isDefault bool,
	tieBreaker int,
) ([]RuleSet, error) {
	// Tokenize variable values
	refTokens := make(map[string][]Token)
	if variables != nil {
		vt, err := TokenizeValues(variables)
		if err != nil {
			return nil, err
		}
		for k, v := range vt {
			refTokens[k] = v
		}
	}
	if variableTokens != nil {
		for k, v := range variableTokens {
			refTokens[k] = v
		}
	}

	// Tokenize the CSS
	rawTokens, err := Tokenize(css, readFrom)
	if err != nil {
		return nil, err
	}

	// Substitute variables
	tokens, err := SubstituteReferences(rawTokens, refTokens)
	if err != nil {
		return nil, err
	}

	// Parse rule sets
	it := NewTokenIterator(tokens)
	var ruleSets []RuleSet
	for {
		tok, ok := it.Next()
		if !ok {
			break
		}
		if strings.HasPrefix(tok.Name, "selector_start") {
			rules, err := parseRuleSet(scope, it, tok, isDefault, tieBreaker)
			if err != nil {
				return nil, err
			}
			ruleSets = append(ruleSets, rules...)
		}
	}
	return ruleSets, nil
}

// parseRuleSet parses a single rule set (possibly with nested rules).
func parseRuleSet(
	scope string,
	tokens *TokenIterator,
	firstToken Token,
	isDefault bool,
	tieBreaker int,
) ([]RuleSet, error) {
	desc := CombinatorDescendant
	same := CombinatorSame

	var combinator *CombinatorType = &desc
	var selectors []Selector
	var ruleSelectorLists [][]Selector
	builder := NewStylesBuilder()

	tok := firstToken
	// Process selector tokens until declaration_set_start
selectorLoop:
	for {
		switch tok.Name {
		case "pseudo_class":
			if len(selectors) > 0 {
				selectors[len(selectors)-1].AddPseudoClass(strings.TrimLeft(tok.Value, ":"))
			}
		case "whitespace":
			if combinator == nil || *combinator == CombinatorSame {
				combinator = &desc
			}
		case "new_selector":
			ruleSelectorLists = append(ruleSelectorLists, append([]Selector{}, selectors...))
			selectors = selectors[:0]
			combinator = nil
		case "declaration_set_start":
			break selectorLoop
		case "combinator_child":
			child := CombinatorChild
			combinator = &child
		default:
			entry, ok := selectorTokenMap[tok.Name]
			if ok {
				comb := desc
				if combinator != nil {
					comb = *combinator
				}
				sel := Selector{
					Name:        strings.TrimLeft(tok.Value, ".#"),
					Combinator:  comb,
					Type:        entry.Type,
					Specificity: entry.Specificity,
					Advance:     1,
				}
				selectors = append(selectors, sel)
				combinator = &same
			}
		}
		var ok bool
		tok, ok = tokens.Next()
		if !ok {
			break
		}
	}

	// Apply scope prefix if needed
	if len(selectors) > 0 {
		if scope != "" && (len(selectors) == 0 || selectors[0].Name != scope) {
			entry, ok := selectorTokenMap[scope]
			scopeType := SelectorType_Type
			scopeSpec := Specificity3{0, 0, 1}
			if ok {
				scopeType = entry.Type
				scopeSpec = entry.Specificity
			}
			scopeSel := Selector{
				Name:        scope,
				Combinator:  CombinatorDescendant,
				Type:        scopeType,
				Specificity: scopeSpec,
				Advance:     1,
			}
			selectors = append([]Selector{scopeSel}, selectors...)
		}
		ruleSelectorLists = append(ruleSelectorLists, selectors)
	}

	// Parse declarations
	var decl Declaration
	var ruleErrors []RuleError
	var nestedRules []RuleSet

	for {
		tok, ok := tokens.Next()
		if !ok {
			break
		}

		switch tok.Name {
		case "whitespace", "declaration_end":
			continue
		case "selector_start_id", "selector_start_class", "selector_start_universal", "selector_start", "nested":
			// Nested rule
			recursiveRules, err := parseRuleSet("", tokens, tok, isDefault, tieBreaker)
			if err != nil {
				return nil, err
			}
			// Combine selectors
			for _, ruleSelector := range ruleSelectorLists {
				for _, ruleSet := range recursiveRules {
					var combined []SelectorSet
					for _, rss := range ruleSet.SelectorSets {
						mergedSels := combineSelectors(ruleSelector, rss.Selectors)
						ss := SelectorSet{Selectors: mergedSels}
						ss.TotalSpecificity()
						combined = append(combined, ss)
					}
					nestedRules = append(nestedRules, RuleSet{
						SelectorSets:   combined,
						Styles:         ruleSet.Styles,
						Errors:         ruleSet.Errors,
						IsDefaultRules: ruleSet.IsDefaultRules,
						TieBreaker:     ruleSet.TieBreaker + tieBreaker,
					})
				}
			}
			continue
		case "declaration_name":
			if err := builder.AddDeclaration(decl); err != nil && decl.Name != "" {
				ruleErrors = append(ruleErrors, RuleError{Token: decl.Token, Message: err.Error()})
			}
			decl = Declaration{Token: tok, Name: strings.TrimRight(tok.Value, ":")}
		case "declaration_set_end":
			goto ruleDone
		default:
			decl.Tokens = append(decl.Tokens, tok)
		}
	}
ruleDone:
	if err := builder.AddDeclaration(decl); err != nil && decl.Name != "" {
		ruleErrors = append(ruleErrors, RuleError{Token: decl.Token, Message: err.Error()})
	}

	rs := RuleSet{
		SelectorSets:   SelectorSetsFromSelectors(ruleSelectorLists),
		Styles:         builder.Styles,
		Errors:         ruleErrors,
		IsDefaultRules: isDefault,
		TieBreaker:     tieBreaker,
	}
	rs.PostParse()

	result := []RuleSet{rs}
	for i := range nestedRules {
		nestedRules[i].PostParse()
		result = append(result, nestedRules[i])
	}
	return result, nil
}

// combineSelectors merges two selector lists, handling nested selectors (&).
func combineSelectors(parent, child []Selector) []Selector {
	if len(child) > 0 && child[0].Type == SelectorNested {
		// Merge: replace the nested placeholder with the parent's last selector +
		// any pseudo-classes from the child's first selector.
		if len(parent) == 0 {
			return child[1:]
		}
		last := parent[len(parent)-1]
		nested := child[0]
		// Merge pseudo-classes
		merged := Selector{
			Name:        last.Name,
			Combinator:  last.Combinator,
			Type:        last.Type,
			Specificity: AddSpecificity(last.Specificity, nested.Specificity),
			Advance:     last.Advance,
		}
		if merged.PseudoClasses == nil {
			merged.PseudoClasses = make(map[string]bool)
		}
		for pc := range last.PseudoClasses {
			merged.PseudoClasses[pc] = true
		}
		for pc := range nested.PseudoClasses {
			merged.PseudoClasses[pc] = true
		}
		result := append(parent[:len(parent)-1:len(parent)-1], merged)
		result = append(result, child[1:]...)
		return result
	}
	return append(append([]Selector{}, parent...), child...)
}
