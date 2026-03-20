package document_test

import (
	"testing"

	"github.com/eberle1080/go-textual/document"
)

func TestNewDocumentEmpty(t *testing.T) {
	d := document.NewDocument("")
	if d.LineCount() != 1 {
		t.Errorf("empty document should have 1 line, got %d", d.LineCount())
	}
	if d.Text() != "" {
		t.Errorf("empty document text should be '', got %q", d.Text())
	}
}

func TestNewDocumentDetectsLF(t *testing.T) {
	d := document.NewDocument("hello\nworld")
	if d.Newline() != document.NewlineLF {
		t.Errorf("expected LF newline, got %q", d.Newline())
	}
	if d.LineCount() != 2 {
		t.Errorf("expected 2 lines, got %d", d.LineCount())
	}
}

func TestNewDocumentDetectsCRLF(t *testing.T) {
	d := document.NewDocument("hello\r\nworld")
	if d.Newline() != document.NewlineCRLF {
		t.Errorf("expected CRLF newline, got %q", d.Newline())
	}
}

func TestDocumentGetLine(t *testing.T) {
	d := document.NewDocument("line0\nline1\nline2")
	tests := []struct {
		row  int
		want string
	}{
		{0, "line0"},
		{1, "line1"},
		{2, "line2"},
		{-1, ""},
		{99, ""},
	}
	for _, tt := range tests {
		got := d.GetLine(tt.row)
		if got != tt.want {
			t.Errorf("GetLine(%d) = %q, want %q", tt.row, got, tt.want)
		}
	}
}

func TestDocumentGetTextRange(t *testing.T) {
	d := document.NewDocument("hello\nworld")
	tests := []struct {
		start, end document.Location
		want       string
	}{
		{document.Location{0, 0}, document.Location{0, 5}, "hello"},
		{document.Location{0, 0}, document.Location{1, 5}, "hello\nworld"},
		{document.Location{0, 1}, document.Location{0, 4}, "ell"},
		{document.Location{0, 0}, document.Location{0, 0}, ""},
	}
	for _, tt := range tests {
		got := d.GetTextRange(tt.start, tt.end)
		if got != tt.want {
			t.Errorf("GetTextRange(%v, %v) = %q, want %q", tt.start, tt.end, got, tt.want)
		}
	}
}

func TestDocumentReplaceRange(t *testing.T) {
	tests := []struct {
		name     string
		initial  string
		start    document.Location
		end      document.Location
		text     string
		wantText string
		wantEnd  document.Location
	}{
		{
			name:     "insert at start",
			initial:  "world",
			start:    document.Location{0, 0},
			end:      document.Location{0, 0},
			text:     "hello ",
			wantText: "hello world",
			wantEnd:  document.Location{0, 6},
		},
		{
			name:     "replace word",
			initial:  "hello world",
			start:    document.Location{0, 6},
			end:      document.Location{0, 11},
			text:     "Go",
			wantText: "hello Go",
			wantEnd:  document.Location{0, 8},
		},
		{
			name:     "delete characters",
			initial:  "hello",
			start:    document.Location{0, 0},
			end:      document.Location{0, 5},
			text:     "",
			wantText: "",
			wantEnd:  document.Location{0, 0},
		},
		{
			name:     "insert newline",
			initial:  "hello world",
			start:    document.Location{0, 5},
			end:      document.Location{0, 6},
			text:     "\n",
			wantText: "hello\nworld",
			wantEnd:  document.Location{1, 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := document.NewDocument(tt.initial)
			result := d.ReplaceRange(tt.start, tt.end, tt.text)
			if d.Text() != tt.wantText {
				t.Errorf("Text() = %q, want %q", d.Text(), tt.wantText)
			}
			if !result.EndLocation.Equal(tt.wantEnd) {
				t.Errorf("EndLocation = %v, want %v", result.EndLocation, tt.wantEnd)
			}
		})
	}
}

func TestDocumentEndLocation(t *testing.T) {
	d := document.NewDocument("hello\nworld")
	end := d.End()
	if end.Row != 1 || end.Col != 5 {
		t.Errorf("End() = %v, want {1, 5}", end)
	}
}

func TestLocationLess(t *testing.T) {
	a := document.Location{Row: 0, Col: 0}
	b := document.Location{Row: 0, Col: 1}
	c := document.Location{Row: 1, Col: 0}

	if !a.Less(b) {
		t.Error("a should be less than b")
	}
	if !b.Less(c) {
		t.Error("b should be less than c")
	}
	if a.Less(a) {
		t.Error("a should not be less than itself")
	}
}

func TestSelectionIsEmpty(t *testing.T) {
	loc := document.Location{Row: 1, Col: 3}
	sel := document.Cursor(loc)
	if !sel.IsEmpty() {
		t.Error("cursor selection should be empty")
	}
	sel.End = document.Location{Row: 1, Col: 5}
	if sel.IsEmpty() {
		t.Error("non-degenerate selection should not be empty")
	}
}
