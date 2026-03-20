// Package msg defines the message types used by the event loop.
// All communication between goroutines and the event loop flows through
// values implementing the Msg interface.
package msg

import (
	"context"
	"time"
)

// Msg is the sealed interface for all messages.
// Only types embedding BaseMsg satisfy this interface.
type Msg interface{ isMsg() }

// BaseMsg is embedded by all message types to satisfy Msg.
type BaseMsg struct{}

func (BaseMsg) isMsg() {}

// Cmd runs off the event loop goroutine; its return value (if non-nil) is sent
// back through the events channel.
type Cmd func(ctx context.Context) Msg

// Batch fans out all cmds concurrently. Each result is sent back independently.
func Batch(cmds ...Cmd) Cmd {
	if len(cmds) == 0 {
		return None()
	}
	return func(ctx context.Context) Msg {
		return batchMsg{cmds: cmds}
	}
}

// Sequence runs cmds one after another. Each cmd starts only after the previous
// one completes. The last result is returned.
func Sequence(cmds ...Cmd) Cmd {
	if len(cmds) == 0 {
		return None()
	}
	return func(ctx context.Context) Msg {
		return sequenceMsg{cmds: cmds}
	}
}

// Tick returns a Cmd that sends a TickMsg after duration d.
func Tick(d time.Duration) Cmd {
	return func(ctx context.Context) Msg {
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(d):
			return TickMsg{Duration: d}
		}
	}
}

// None returns a no-op Cmd.
func None() Cmd {
	return func(ctx context.Context) Msg { return nil }
}

// batchMsg is an internal message carrying multiple cmds to fan-out.
type batchMsg struct {
	BaseMsg
	cmds []Cmd
}

// Cmds returns the commands in a batchMsg. Used by the event loop.
func (b batchMsg) Cmds() []Cmd { return b.cmds }

// IsBatch reports whether m is a batch message and returns its commands.
func IsBatch(m Msg) ([]Cmd, bool) {
	b, ok := m.(batchMsg)
	if !ok {
		return nil, false
	}
	return b.cmds, true
}

// sequenceMsg is an internal message carrying sequential cmds.
type sequenceMsg struct {
	BaseMsg
	cmds []Cmd
}

// IsSequence reports whether m is a sequence message and returns its commands.
func IsSequence(m Msg) ([]Cmd, bool) {
	s, ok := m.(sequenceMsg)
	if !ok {
		return nil, false
	}
	return s.cmds, true
}
