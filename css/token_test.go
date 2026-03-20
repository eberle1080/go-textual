package css

import "testing"

func TestTokenStart(t *testing.T) {
	tok := Token{Location: [2]int{0, 0}, Value: "foo"}
	s := tok.Start()
	if s[0] != 1 || s[1] != 1 {
		t.Errorf("Start() = %v, want [1 1]", s)
	}

	tok2 := Token{Location: [2]int{2, 4}, Value: "bar"}
	s2 := tok2.Start()
	if s2[0] != 3 || s2[1] != 5 {
		t.Errorf("Start() = %v, want [3 5]", s2)
	}
}

func TestTokenEnd(t *testing.T) {
	tok := Token{Location: [2]int{0, 0}, Value: "foo"}
	e := tok.End()
	// End is 1-indexed; col = location[1]+len(value)+1 = 0+3+1 = 4
	if e[0] != 1 || e[1] != 4 {
		t.Errorf("End() = %v, want [1 4]", e)
	}

	tok2 := Token{Location: [2]int{1, 2}, Value: "ab"}
	e2 := tok2.End()
	if e2[0] != 2 || e2[1] != 5 {
		t.Errorf("End() = %v, want [2 5]", e2)
	}
}

func TestTokenWithReference(t *testing.T) {
	tok := Token{Name: "scalar", Value: "10px"}
	ref := &ReferencedBy{Name: "myvar", Location: [2]int{0, 0}, Length: 6, Code: "$myvar"}
	tok2 := tok.WithReference(ref)
	if tok2.ReferencedBy != ref {
		t.Errorf("WithReference did not set ReferencedBy")
	}
	if tok.ReferencedBy != nil {
		t.Errorf("WithReference mutated original token")
	}
}

func TestTokenString(t *testing.T) {
	tok := Token{Value: "hello"}
	if tok.String() != "hello" {
		t.Errorf("String() = %q, want %q", tok.String(), "hello")
	}
}
