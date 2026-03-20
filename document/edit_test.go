package document_test

import (
	"testing"

	"github.com/eberle1080/go-textual/document"
)

func TestEditDoUndo(t *testing.T) {
	d := document.NewDocument("hello world")
	sel := document.Cursor(document.Location{0, 0})

	edit := document.NewEdit("Go", document.Location{0, 6}, document.Location{0, 11})
	result := edit.Do(d, sel, true)

	if d.Text() != "hello Go" {
		t.Errorf("after Do: text = %q, want %q", d.Text(), "hello Go")
	}
	if !result.EndLocation.Equal((document.Location{0, 8})) {
		t.Errorf("EndLocation = %v, want {0,8}", result.EndLocation)
	}

	// Undo should restore original text.
	edit.Undo(d)
	if d.Text() != "hello world" {
		t.Errorf("after Undo: text = %q, want %q", d.Text(), "hello world")
	}
}

func TestNewInsert(t *testing.T) {
	d := document.NewDocument("world")
	edit := document.NewInsert("hello ", document.Location{0, 0})
	edit.Do(d, document.Cursor(document.Location{0, 0}), false)
	if d.Text() != "hello world" {
		t.Errorf("after insert: %q", d.Text())
	}
}

func TestNewDelete(t *testing.T) {
	d := document.NewDocument("hello world")
	edit := document.NewDelete(document.Location{0, 5}, document.Location{0, 11})
	edit.Do(d, document.Cursor(document.Location{0, 5}), false)
	if d.Text() != "hello" {
		t.Errorf("after delete: %q", d.Text())
	}
}

func TestEditTopBottom(t *testing.T) {
	edit := document.NewEdit("x", document.Location{3, 5}, document.Location{1, 2})
	top := edit.Top()
	bottom := edit.Bottom()
	if top.Row != 1 || top.Col != 2 {
		t.Errorf("Top = %v, want {1,2}", top)
	}
	if bottom.Row != 3 || bottom.Col != 5 {
		t.Errorf("Bottom = %v, want {3,5}", bottom)
	}
}

func TestEditRecordsSelection(t *testing.T) {
	d := document.NewDocument("hello")
	sel := document.Cursor(document.Location{0, 2})
	edit := document.NewInsert("X", document.Location{0, 2})
	edit.Do(d, sel, true)

	if edit.OriginalSelection() == nil {
		t.Error("original selection should be recorded")
	}
	if edit.UpdatedSelection() == nil {
		t.Error("updated selection should be recorded")
	}
}
