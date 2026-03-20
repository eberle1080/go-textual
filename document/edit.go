package document

// Edit represents a single reversible editing operation on a [DocumentBase].
// Call [Edit.Do] to apply the edit, [Edit.Undo] to reverse it.
type Edit struct {
	// Text is the replacement text.
	Text string
	// From is the start of the replaced range.
	From Location
	// To is the end of the replaced range.
	To Location
	// MaintainSelectionOffset, when true, attempts to preserve relative cursor
	// offset within the edited region.
	MaintainSelectionOffset bool

	originalSelection *Selection
	updatedSelection  *Selection
	editResult        *EditResult
}

// NewEdit creates an Edit that replaces the range [from, to) with text.
func NewEdit(text string, from, to Location) *Edit {
	return &Edit{Text: text, From: from, To: to}
}

// NewInsert creates an Edit that inserts text at loc (zero-width replacement).
func NewInsert(text string, loc Location) *Edit {
	return &Edit{Text: text, From: loc, To: loc}
}

// NewDelete creates an Edit that deletes the range [from, to).
func NewDelete(from, to Location) *Edit {
	return &Edit{Text: "", From: from, To: to}
}

// Top returns the earlier of From and To.
func (e *Edit) Top() Location {
	if e.To.Less(e.From) {
		return e.To
	}
	return e.From
}

// Bottom returns the later of From and To.
func (e *Edit) Bottom() Location {
	if e.To.Less(e.From) {
		return e.From
	}
	return e.To
}

// Do applies the edit to doc. currentSelection is the selection before the
// edit; recordSelection controls whether it is stored for later undo.
// Returns the EditResult.
func (e *Edit) Do(doc DocumentBase, currentSelection Selection, recordSelection bool) EditResult {
	if recordSelection {
		sel := currentSelection
		e.originalSelection = &sel
	}

	result := doc.ReplaceRange(e.From, e.To, e.Text)
	e.editResult = &result

	// Compute updated selection after the edit.
	if recordSelection {
		updatedSel := Selection{
			Start: result.EndLocation,
			End:   result.EndLocation,
		}
		e.updatedSelection = &updatedSel
	}

	return result
}

// Undo reverses the edit using the stored [EditResult]. Panics if Do has not
// been called.
func (e *Edit) Undo(doc DocumentBase) EditResult {
	if e.editResult == nil {
		panic("document: Edit.Undo called before Do")
	}
	return doc.ReplaceRange(e.From, e.editResult.EndLocation, e.editResult.ReplacedText)
}

// OriginalSelection returns the selection stored when Do was called, or nil.
func (e *Edit) OriginalSelection() *Selection { return e.originalSelection }

// UpdatedSelection returns the selection after Do was called, or nil.
func (e *Edit) UpdatedSelection() *Selection { return e.updatedSelection }
