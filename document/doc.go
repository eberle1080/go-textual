// Package document provides the document model used by [widgets.TextArea].
//
// It implements text storage, editing, undo/redo history, and line wrapping.
// The core types are:
//
//   - [Document] — in-memory text storage with newline-style detection
//   - [Edit] — a single reversible editing operation
//   - [EditHistory] — undo/redo stack with batching and checkpoints
//   - [WrappedDocument] — line-wrapping layer over a Document
//   - [DocumentNavigator] — cursor movement in a wrapped document
//   - [SyntaxAwareDocument] — placeholder for tree-sitter integration
//
// # Basic Usage
//
// Create a document, insert text, and inspect the result:
//
//	doc := document.NewDocument("Hello, World!")
//	edit := document.NewInsert(" Go", document.Location{Row: 0, Col: 5})
//	_ = edit.Do(doc, document.Cursor(document.Location{}), true)
//	fmt.Println(doc.Text()) // "Hello, Go World!"
//
// # Undo/Redo
//
// Record edits in an [EditHistory] and replay them in reverse:
//
//	history := document.NewEditHistory()
//	history.Record(edit)
//
//	// Undo: replay the batch in reverse.
//	batch := history.PopUndo()
//	for i := len(batch) - 1; i >= 0; i-- {
//	    batch[i].Undo(doc)
//	}
//
// # Line Wrapping
//
// Wrap a document to a fixed column width for soft-wrap rendering:
//
//	wrapped := document.NewWrappedDocument(doc, 80)
//	nav := document.NewDocumentNavigator(wrapped)
//	below := nav.GetLocationBelow(document.Location{Row: 0, Col: 0})
package document
