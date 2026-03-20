package css

import "fmt"

// ReferencedBy records where a variable reference was used.
type ReferencedBy struct {
	Name     string
	Location [2]int // 0-indexed [line, col]
	Length   int
	Code     string
}

// Token is a single CSS token produced by the tokenizer.
// Location is 0-indexed [line, col].
type Token struct {
	Name         string
	Value        string
	ReadFrom     CSSLocation
	Code         string
	Location     [2]int // 0-indexed
	ReferencedBy *ReferencedBy
}

// Start returns the 1-indexed [line, col] start position of the token.
func (t Token) Start() [2]int {
	return [2]int{t.Location[0] + 1, t.Location[1] + 1}
}

// End returns the 1-indexed [line, col] end position of the token
// (exclusive — one past the last character).
func (t Token) End() [2]int {
	return [2]int{t.Location[0] + 1, t.Location[1] + len(t.Value) + 1}
}

// WithReference returns a copy of the token with the ReferencedBy field set.
func (t Token) WithReference(by *ReferencedBy) Token {
	t.ReferencedBy = by
	return t
}

// String returns the token's value.
func (t Token) String() string { return t.Value }

// GoString returns a human-readable debug representation.
func (t Token) GoString() string {
	return fmt.Sprintf("Token{Name:%q Value:%q Location:%v}", t.Name, t.Value, t.Location)
}
