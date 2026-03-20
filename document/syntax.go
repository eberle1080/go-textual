package document

import "unicode"

// SyntaxHighlight represents a single syntax-highlighted span within a document.
type SyntaxHighlight struct {
	// Start is the start location of the highlight.
	Start Location
	// End is the end location of the highlight.
	End Location
	// Name is the tree-sitter node type name (e.g. "keyword", "string").
	Name string
}

// SyntaxAwareDocument embeds a [Document] and adds syntax highlighting for
// supported languages. Currently supports: "go".
type SyntaxAwareDocument struct {
	*Document
	language string
}

// NewSyntaxAwareDocument creates a SyntaxAwareDocument for the given text
// and language identifier (e.g. "python", "go").
func NewSyntaxAwareDocument(text, language string) *SyntaxAwareDocument {
	return &SyntaxAwareDocument{
		Document: NewDocument(text),
		language: language,
	}
}

// Language returns the language identifier for this document.
func (s *SyntaxAwareDocument) Language() string { return s.language }

// QuerySyntaxTree returns syntax highlights for the given range.
// For "go", it performs a lexical scan for keyword tokens.
// For unsupported languages it returns nil.
func (s *SyntaxAwareDocument) QuerySyntaxTree(start, end Location) []SyntaxHighlight {
	switch s.language {
	case "go":
		return s.queryGo(start, end)
	default:
		return nil
	}
}

// goKeywords is the set of Go reserved keywords.
var goKeywords = map[string]bool{
	"break": true, "case": true, "chan": true, "const": true, "continue": true,
	"default": true, "defer": true, "else": true, "fallthrough": true, "for": true,
	"func": true, "go": true, "goto": true, "if": true, "import": true,
	"interface": true, "map": true, "package": true, "range": true, "return": true,
	"select": true, "struct": true, "switch": true, "type": true, "var": true,
}

// queryGo returns keyword highlights for rows start.Row..end.Row.
func (s *SyntaxAwareDocument) queryGo(start, end Location) []SyntaxHighlight {
	var out []SyntaxHighlight
	for row := start.Row; row <= end.Row && row < s.LineCount(); row++ {
		line := s.GetLine(row)
		runes := []rune(line)
		colStart := 0
		colEnd := len(runes)
		if row == start.Row && start.Col > colStart {
			colStart = start.Col
		}
		if row == end.Row && end.Col < colEnd {
			colEnd = end.Col
		}
		out = append(out, scanGoKeywords(runes, row, colStart, colEnd)...)
	}
	return out
}

// scanGoKeywords walks runes[colStart:colEnd] and returns a SyntaxHighlight for
// each Go keyword token found.
func scanGoKeywords(runes []rune, row, colStart, colEnd int) []SyntaxHighlight {
	var out []SyntaxHighlight
	i := colStart
	for i < colEnd {
		if isIdentStart(runes[i]) {
			j := i + 1
			for j < colEnd && isIdentPart(runes[j]) {
				j++
			}
			word := string(runes[i:j])
			if goKeywords[word] {
				out = append(out, SyntaxHighlight{
					Start: Location{Row: row, Col: i},
					End:   Location{Row: row, Col: j},
					Name:  "keyword",
				})
			}
			i = j
		} else {
			i++
		}
	}
	return out
}

func isIdentStart(r rune) bool { return r == '_' || unicode.IsLetter(r) }
func isIdentPart(r rune) bool  { return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r) }
