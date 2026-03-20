package strip

import (
	"testing"

	rich "github.com/eberle1080/go-rich"
)

func seg(text string) rich.Segment {
	return rich.Segment{Text: text, Style: rich.Style{}}
}

func boldSeg(text string) rich.Segment {
	return rich.Segment{Text: text, Style: rich.NewStyle().Bold()}
}

func TestStrip_CellLength(t *testing.T) {
	s := New(rich.Segments{seg("hello"), seg(" world")})
	if s.CellLength() != 11 {
		t.Fatalf("expected 11, got %d", s.CellLength())
	}
}

func TestStrip_Blank(t *testing.T) {
	s := Blank(5, rich.Style{})
	if s.CellLength() != 5 {
		t.Fatalf("expected 5, got %d", s.CellLength())
	}
	if s.Text() != "     " {
		t.Fatalf("expected 5 spaces, got %q", s.Text())
	}
}

func TestStrip_Crop(t *testing.T) {
	s := New(rich.Segments{seg("hello world")})

	cropped := s.Crop(6, 11)
	if cropped.Text() != "world" {
		t.Fatalf("expected 'world', got %q", cropped.Text())
	}
	if cropped.CellLength() != 5 {
		t.Fatalf("expected CellLength=5, got %d", cropped.CellLength())
	}
}

func TestStrip_Crop_EmptyResult(t *testing.T) {
	s := New(rich.Segments{seg("hello")})
	result := s.Crop(3, 3)
	if result.CellLength() != 0 {
		t.Fatalf("expected empty crop, got len=%d", result.CellLength())
	}
}

func TestStrip_Crop_MultiSegment(t *testing.T) {
	s := New(rich.Segments{seg("abc"), boldSeg("def"), seg("ghi")})
	// Crop across segment boundary: cells 2..7 = "cdefg"
	cropped := s.Crop(2, 7)
	if cropped.Text() != "cdefg" {
		t.Fatalf("expected 'cdefg', got %q", cropped.Text())
	}
}

func TestStrip_AdjustCellLength_Extend(t *testing.T) {
	s := New(rich.Segments{seg("hi")})
	extended := s.AdjustCellLength(5, rich.Style{})
	if extended.CellLength() != 5 {
		t.Fatalf("expected 5, got %d", extended.CellLength())
	}
	if extended.Text() != "hi   " {
		t.Fatalf("expected 'hi   ', got %q", extended.Text())
	}
}

func TestStrip_AdjustCellLength_Truncate(t *testing.T) {
	s := New(rich.Segments{seg("hello world")})
	truncated := s.AdjustCellLength(5, rich.Style{})
	if truncated.CellLength() != 5 {
		t.Fatalf("expected 5, got %d", truncated.CellLength())
	}
}

func TestStrip_Join(t *testing.T) {
	a := New(rich.Segments{seg("foo")})
	b := New(rich.Segments{seg("bar")})
	j := Join([]Strip{a, b})
	if j.Text() != "foobar" {
		t.Fatalf("expected 'foobar', got %q", j.Text())
	}
	if j.CellLength() != 6 {
		t.Fatalf("expected CellLength=6, got %d", j.CellLength())
	}
}

func TestStrip_Divide(t *testing.T) {
	s := New(rich.Segments{seg("hello world")})
	parts := s.Divide([]int{5, 6})
	if len(parts) != 3 {
		t.Fatalf("expected 3 parts, got %d", len(parts))
	}
	if parts[0].Text() != "hello" {
		t.Fatalf("parts[0]: expected 'hello', got %q", parts[0].Text())
	}
	if parts[1].Text() != " " {
		t.Fatalf("parts[1]: expected ' ', got %q", parts[1].Text())
	}
	if parts[2].Text() != "world" {
		t.Fatalf("parts[2]: expected 'world', got %q", parts[2].Text())
	}
}

func TestStrip_Simplify(t *testing.T) {
	style := rich.NewStyle().Bold()
	s := New(rich.Segments{
		{Text: "foo", Style: style},
		{Text: "bar", Style: style},
		{Text: "baz", Style: rich.Style{}},
	})
	simplified := s.Simplify()
	if simplified.Len() != 2 {
		t.Fatalf("expected 2 segments after simplify, got %d", simplified.Len())
	}
	if simplified.Segments()[0].Text != "foobar" {
		t.Fatalf("expected 'foobar', got %q", simplified.Segments()[0].Text)
	}
}

func TestLinePad(t *testing.T) {
	segs := rich.Segments{seg("hi")}
	padded := LinePad(segs, 2, 3, rich.Style{})
	if GetLineLength(padded) != 7 {
		t.Fatalf("expected length 7, got %d", GetLineLength(padded))
	}
}

func TestGetLineLength(t *testing.T) {
	segs := rich.Segments{seg("hello"), seg(" world")}
	if GetLineLength(segs) != 11 {
		t.Fatalf("expected 11, got %d", GetLineLength(segs))
	}
}

func TestStrip_Render(t *testing.T) {
	s := New(rich.Segments{seg("hello")})
	rendered := s.Render(rich.ColorModeNone)
	if rendered != "hello" {
		t.Fatalf("expected 'hello', got %q", rendered)
	}
}

func TestStrip_Equal(t *testing.T) {
	a := New(rich.Segments{seg("foo")})
	b := New(rich.Segments{seg("foo")})
	c := New(rich.Segments{seg("bar")})
	if !a.Equal(b) {
		t.Fatal("expected a == b")
	}
	if a.Equal(c) {
		t.Fatal("expected a != c")
	}
}
