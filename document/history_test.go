package document_test

import (
	"testing"

	"github.com/eberle1080/go-textual/document"
)

func TestHistoryUndoRedo(t *testing.T) {
	h := document.NewEditHistory()
	d := document.NewDocument("hello")

	e1 := document.NewInsert(" world", document.Location{0, 5})
	e1.Do(d, document.Cursor(document.Location{0, 5}), true)
	h.Record(e1)

	h.Checkpoint()

	e2 := document.NewInsert("!", document.Location{0, 11})
	e2.Do(d, document.Cursor(document.Location{0, 11}), true)
	h.Record(e2)

	h.Checkpoint()

	// Undo e2 batch.
	batch := h.PopUndo()
	if len(batch) == 0 {
		t.Fatal("expected undo batch")
	}
	for i := len(batch) - 1; i >= 0; i-- {
		edit := batch[i]
		edit.Undo(d)
	}
	if d.Text() != "hello world" {
		t.Errorf("after undo: %q", d.Text())
	}

	// Undo e1 batch.
	batch = h.PopUndo()
	if len(batch) == 0 {
		t.Fatal("expected second undo batch")
	}
	for i := len(batch) - 1; i >= 0; i-- {
		edit := batch[i]
		edit.Undo(d)
	}
	if d.Text() != "hello" {
		t.Errorf("after second undo: %q", d.Text())
	}
}

func TestHistoryClear(t *testing.T) {
	h := document.NewEditHistory()
	d := document.NewDocument("hello")
	e := document.NewInsert(" world", document.Location{0, 5})
	e.Do(d, document.Cursor(document.Location{0, 5}), true)
	h.Record(e)
	h.Checkpoint()
	h.Clear()

	if h.PopUndo() != nil {
		t.Error("undo stack should be empty after clear")
	}
}

func TestHistoryMaxCheckpoints(t *testing.T) {
	h := document.NewEditHistory(document.WithMaxCheckpoints(2))
	d := document.NewDocument("")

	for i := 0; i < 5; i++ {
		e := document.NewInsert("x", document.Location{0, i})
		e.Do(d, document.Cursor(document.Location{0, i}), false)
		h.Record(e)
		h.Checkpoint()
	}

	count := 0
	for h.PopUndo() != nil {
		count++
	}
	if count > 2 {
		t.Errorf("expected max 2 checkpoints, got %d", count)
	}
}
