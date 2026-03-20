package css

import (
	"regexp"
	"strings"
)

// Expect describes the set of tokens the tokenizer may encounter at a given point.
type Expect struct {
	Description     string
	Names           []string
	Regexes         []string
	compiled        *regexp.Regexp
	ExpectEOF       bool
	ExpectSemicolon bool
	ExtractText     bool
}

// NewExpect creates an Expect that recognises the named token regexes.
// tokens is an ordered slice of [name, pattern] pairs.
func NewExpect(description string, tokens [][2]string) *Expect {
	e := &Expect{
		Description:     "Expected " + description,
		ExpectSemicolon: true,
	}
	var parts []string
	for _, kv := range tokens {
		e.Names = append(e.Names, kv[0])
		e.Regexes = append(e.Regexes, kv[1])
		parts = append(parts, "(?P<"+kv[0]+">"+kv[1]+")")
	}
	pattern := "(" + strings.Join(parts, "|") + ")"
	e.compiled = regexp.MustCompile(pattern)
	return e
}

// WithEOF returns a copy of the Expect with ExpectEOF set.
func (e *Expect) WithEOF(eof bool) *Expect {
	e.ExpectEOF = eof
	return e
}

// WithSemicolon returns a copy of the Expect with ExpectSemicolon set.
func (e *Expect) WithSemicolon(semi bool) *Expect {
	e.ExpectSemicolon = semi
	return e
}

// WithExtractText returns a copy of the Expect with ExtractText set.
func (e *Expect) WithExtractText(extract bool) *Expect {
	e.ExtractText = extract
	return e
}

// Match attempts to match the pattern at position pos in line.
// Returns the match indices from regexp.FindStringIndex, or nil if no match.
func (e *Expect) Match(line string, pos int) []int {
	if pos > len(line) {
		return nil
	}
	loc := e.compiled.FindStringIndex(line[pos:])
	if loc == nil {
		return nil
	}
	// Anchored: must start at pos
	if loc[0] != 0 {
		return nil
	}
	return []int{loc[0] + pos, loc[1] + pos}
}

// Search finds the first occurrence of the pattern at or after pos in line.
func (e *Expect) Search(line string, pos int) []int {
	if pos > len(line) {
		return nil
	}
	loc := e.compiled.FindStringIndex(line[pos:])
	if loc == nil {
		return nil
	}
	return []int{loc[0] + pos, loc[1] + pos}
}

// matchGroups returns the name and matched value for the first matching group
// in an anchored match at position pos in line.  The match must start at pos
// (offset 0 within the sub-slice) so that invalid leading characters cannot
// be silently skipped.
func (e *Expect) matchGroups(line string, pos int) (string, string, bool) {
	if pos > len(line) {
		return "", "", false
	}
	sub := line[pos:]
	idx := e.compiled.FindStringSubmatchIndex(sub)
	if idx == nil || idx[0] != 0 {
		// No match or match does not start at the current cursor position.
		return "", "", false
	}
	names := e.compiled.SubexpNames()
	// idx is pairs [start,end] for each sub-expression; skip indices 0 (full
	// match) and 1 (outer wrapper group), then scan named groups.
	for i := 2; i < len(idx)/2; i++ {
		start, end := idx[i*2], idx[i*2+1]
		if start >= 0 && names[i] != "" {
			return names[i], sub[start:end], true
		}
	}
	return "", "", false
}

// Tokenizer tokenizes Textual CSS.
type Tokenizer struct {
	ReadFrom CSSLocation
	code     string
	lines    []string
	lineNo   int
	colNo    int
}

// NewTokenizer creates a tokenizer for the given CSS code.
func NewTokenizer(code string, readFrom CSSLocation) *Tokenizer {
	return &Tokenizer{
		ReadFrom: readFrom,
		code:     code,
		lines:    splitLines(code),
	}
}

// splitLines splits text into lines preserving line endings (like Python's splitlines(keepends=True)).
func splitLines(s string) []string {
	if s == "" {
		return nil
	}
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, s[start:i+1])
			start = i + 1
		} else if s[i] == '\r' {
			if i+1 < len(s) && s[i+1] == '\n' {
				lines = append(lines, s[start:i+2])
				i++
			} else {
				lines = append(lines, s[start:i+1])
			}
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}

// GetToken reads and returns the next token from the CSS stream.
func (t *Tokenizer) GetToken(expect *Expect) (Token, error) {
	lineNo := t.lineNo
	colNo := t.colNo

	if lineNo >= len(t.lines) {
		if expect.ExpectEOF {
			return Token{
				Name:     "eof",
				Value:    "",
				ReadFrom: t.ReadFrom,
				Code:     t.code,
				Location: [2]int{lineNo, colNo},
			}, nil
		}
		msg := "Unexpected end of file; did you forget a '}' ?"
		if !expect.ExpectSemicolon {
			msg = "Unexpected end of text"
		}
		return Token{}, &UnexpectedEndError{
			TokenError: TokenError{
				ReadFrom: t.ReadFrom,
				Code:     t.code,
				Start:    [2]int{lineNo + 1, colNo + 1},
				End:      [2]int{lineNo + 1, colNo + 1},
				Msg:      msg,
			},
		}
	}

	line := t.lines[lineNo]

	if expect.ExtractText {
		loc := expect.Search(line, colNo)
		var precedingText string
		if loc == nil {
			precedingText = line[colNo:]
			t.lineNo++
			t.colNo = 0
		} else {
			matchStart := loc[0]
			precedingText = line[colNo:matchStart]
			t.colNo = matchStart
		}
		if precedingText != "" {
			return Token{
				Name:     "text",
				Value:    precedingText,
				ReadFrom: t.ReadFrom,
				Code:     t.code,
				Location: [2]int{lineNo, colNo},
			}, nil
		}
	}

	// Non-extract mode: match at current position
	name, value, ok := expect.matchGroups(line, colNo)
	if !ok {
		errLine := line[colNo:]
		semi := strings.SplitN(errLine, ";", 2)[0]
		errMsg := expect.Description + " (found " + quote(semi) + ")."
		if expect.ExpectSemicolon && !strings.HasSuffix(strings.TrimSpace(errLine), ";") {
			errMsg += "; Did you forget a semicolon at the end of a line?"
		}
		return Token{}, NewTokenError(t.ReadFrom, t.code, [2]int{lineNo + 1, colNo + 1}, errMsg, nil)
	}

	token := Token{
		Name:     name,
		Value:    value,
		ReadFrom: t.ReadFrom,
		Code:     t.code,
		Location: [2]int{lineNo, colNo},
	}

	// Validate pseudo-classes
	if token.Name == "pseudo_class" {
		pc := strings.TrimLeft(token.Value, ":")
		if !ValidPseudoClasses[pc] {
			suggestion := getSuggestion(pc, ValidPseudoClasses)
			msg := "unknown pseudo-class " + quote(pc)
			if suggestion != "" {
				msg += "; did you mean " + quote(suggestion) + "?"
			}
			return Token{}, NewTokenError(t.ReadFrom, t.code, [2]int{lineNo + 1, colNo + 1}, msg, nil)
		}
	}

	// Advance cursor
	colNo += len(value)
	if colNo >= len(line) {
		t.lineNo = lineNo + 1
		t.colNo = 0
	} else {
		t.lineNo = lineNo
		t.colNo = colNo
	}
	return token, nil
}

// SkipTo advances the tokenizer until a token matching expect is found.
func (t *Tokenizer) SkipTo(expect *Expect) (Token, error) {
	lineNo := t.lineNo
	colNo := t.colNo

	for {
		if lineNo >= len(t.lines) {
			msg := "Unexpected end of file; did you forget a '}' ?"
			if !expect.ExpectSemicolon {
				msg = "Unexpected end of markup"
			}
			return Token{}, &UnexpectedEndError{
				TokenError: TokenError{
					ReadFrom: t.ReadFrom,
					Code:     t.code,
					Start:    [2]int{lineNo, colNo},
					End:      [2]int{lineNo, colNo},
					Msg:      msg,
				},
			}
		}
		line := t.lines[lineNo]
		loc := expect.Search(line, colNo)
		if loc == nil {
			lineNo++
			colNo = 0
		} else {
			t.lineNo = lineNo
			t.colNo = loc[0]
			return t.GetToken(expect)
		}
	}
}

// quote returns a Go-style double-quoted representation of a string.
func quote(s string) string {
	return `"` + strings.ReplaceAll(s, `"`, `\"`) + `"`
}

// getSuggestion returns the closest match for s from the set of valid strings.
func getSuggestion(s string, valid map[string]bool) string {
	best := ""
	bestScore := -1
	for candidate := range valid {
		score := editDistanceScore(s, candidate)
		if score > bestScore {
			bestScore = score
			best = candidate
		}
	}
	if bestScore <= 0 {
		return ""
	}
	return best
}

// editDistanceScore returns a score indicating how similar two strings are.
// Higher is better. Returns 0 if there is no reasonable match.
func editDistanceScore(a, b string) int {
	// Simple prefix/substring matching heuristic
	if a == b {
		return 100
	}
	if strings.HasPrefix(b, a) || strings.HasPrefix(a, b) {
		return 50
	}
	if strings.Contains(b, a) || strings.Contains(a, b) {
		return 25
	}
	return 0
}
