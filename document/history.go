package document

import "time"

// EditHistory manages an undo/redo stack with batching and checkpoints.
//
// Edits are grouped into batches. A new batch is started when:
//   - A newline is inserted
//   - The timer has expired since the last edit
//   - The maximum character count per batch is exceeded
//   - [EditHistory.Checkpoint] is called explicitly
//   - A replacement follows an insertion (or vice versa)
type EditHistory struct {
	maxCheckpoints     int
	checkpointTimer    time.Duration
	checkpointMaxChars int

	undoStack          [][]Edit
	redoStack          [][]Edit
	lastEditTime       time.Time
	charCount          int
	forceEndBatch      bool
	previouslyReplaced bool
	currentBatch       []Edit
}

// NewEditHistory creates an EditHistory with sensible defaults.
//
//   - maxCheckpoints: 50
//   - checkpointTimer: 15 seconds
//   - checkpointMaxChars: 500
func NewEditHistory(opts ...EditHistoryOption) *EditHistory {
	h := &EditHistory{
		maxCheckpoints:     50,
		checkpointTimer:    15 * time.Second,
		checkpointMaxChars: 500,
	}
	for _, o := range opts {
		o(h)
	}
	return h
}

// EditHistoryOption is a functional option for [NewEditHistory].
type EditHistoryOption func(*EditHistory)

// WithMaxCheckpoints sets the maximum undo depth.
func WithMaxCheckpoints(n int) EditHistoryOption {
	return func(h *EditHistory) { h.maxCheckpoints = n }
}

// WithCheckpointTimer sets the inactivity timer before a new batch starts.
func WithCheckpointTimer(d time.Duration) EditHistoryOption {
	return func(h *EditHistory) { h.checkpointTimer = d }
}

// WithCheckpointMaxChars sets the character-count ceiling per batch.
func WithCheckpointMaxChars(n int) EditHistoryOption {
	return func(h *EditHistory) { h.checkpointMaxChars = n }
}

// Record adds edit to the history, possibly starting a new batch.
func (h *EditHistory) Record(edit *Edit) {
	if edit == nil {
		return
	}

	now := time.Now()
	isReplacement := !edit.From.Equal(edit.To)

	// Decide whether to start a new batch.
	newBatch := h.forceEndBatch ||
		len(h.currentBatch) == 0 ||
		(h.checkpointTimer > 0 && !h.lastEditTime.IsZero() && now.Sub(h.lastEditTime) > h.checkpointTimer) ||
		(h.checkpointMaxChars > 0 && h.charCount >= h.checkpointMaxChars) ||
		(isReplacement != h.previouslyReplaced) ||
		edit.Text == "\n" || edit.Text == "\r\n" || edit.Text == "\r"

	if newBatch && len(h.currentBatch) > 0 {
		h.commitBatch()
	}

	h.currentBatch = append(h.currentBatch, *edit)
	h.lastEditTime = now
	h.charCount += len(edit.Text)
	h.forceEndBatch = false
	h.previouslyReplaced = isReplacement

	// Clear redo stack on any new edit.
	h.redoStack = nil
}

// Checkpoint forces the next Record call to start a new batch.
func (h *EditHistory) Checkpoint() {
	if len(h.currentBatch) > 0 {
		h.commitBatch()
	}
	h.forceEndBatch = true
}

// PopUndo pops the top undo batch (newest committed group of edits).
// Returns nil if the undo stack is empty.
func (h *EditHistory) PopUndo() []Edit {
	// Flush current batch first.
	if len(h.currentBatch) > 0 {
		h.commitBatch()
	}
	if len(h.undoStack) == 0 {
		return nil
	}
	batch := h.undoStack[len(h.undoStack)-1]
	h.undoStack = h.undoStack[:len(h.undoStack)-1]
	return batch
}

// PopRedo pops the top redo batch.
// Returns nil if the redo stack is empty.
func (h *EditHistory) PopRedo() []Edit {
	if len(h.redoStack) == 0 {
		return nil
	}
	batch := h.redoStack[len(h.redoStack)-1]
	h.redoStack = h.redoStack[:len(h.redoStack)-1]
	return batch
}

// PushRedo pushes a batch onto the redo stack (called after an undo).
func (h *EditHistory) PushRedo(batch []Edit) {
	h.redoStack = append(h.redoStack, batch)
}

// Clear empties both undo and redo stacks.
func (h *EditHistory) Clear() {
	h.undoStack = nil
	h.redoStack = nil
	h.currentBatch = nil
	h.charCount = 0
	h.forceEndBatch = false
	h.previouslyReplaced = false
}

// commitBatch pushes the current batch onto the undo stack and resets state.
func (h *EditHistory) commitBatch() {
	if len(h.currentBatch) == 0 {
		return
	}
	batch := make([]Edit, len(h.currentBatch))
	copy(batch, h.currentBatch)
	h.undoStack = append(h.undoStack, batch)
	// Trim to max checkpoints.
	if h.maxCheckpoints > 0 && len(h.undoStack) > h.maxCheckpoints {
		h.undoStack = h.undoStack[len(h.undoStack)-h.maxCheckpoints:]
	}
	h.currentBatch = h.currentBatch[:0]
	h.charCount = 0
}
