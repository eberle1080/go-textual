package msg

import (
	"os"
	"time"
)

// QuitMsg tells the event loop to exit cleanly.
type QuitMsg struct {
	BaseMsg
	// Signal is the OS signal that triggered the quit, if any.
	Signal os.Signal
}

// SuspendMsg tells the event loop to suspend (SIGTSTP).
type SuspendMsg struct{ BaseMsg }

// ResumeMsg tells the event loop to resume after suspension (SIGCONT).
type ResumeMsg struct{ BaseMsg }

// PanicMsg is sent when a Cmd panics. The event loop can log it and decide
// whether to exit or continue.
type PanicMsg struct {
	BaseMsg
	// Recovered is the value passed to panic().
	Recovered any
	// Stack is the formatted goroutine stack trace.
	Stack []byte
}

// TickMsg is sent by the Tick Cmd after the timer fires.
type TickMsg struct {
	BaseMsg
	// Duration is the duration that was passed to Tick.
	Duration time.Duration
}

// ResizeSignalMsg is an internal message posted by the SIGWINCH handler.
// The event loop re-queries the terminal size when it receives this.
type ResizeSignalMsg struct{ BaseMsg }
