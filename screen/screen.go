// Package screen defines the Screen interface. A Screen is a root widget that
// composes the initial widget tree and hooks into mount/unmount lifecycle.
package screen

import (
	"context"

	"github.com/eberle1080/go-textual/msg"
	"github.com/eberle1080/go-textual/widget"
)

// Screen is a root widget that owns a widget tree.
type Screen interface {
	widget.Widget

	// Compose returns the initial child widgets. Called once at mount time.
	Compose() []widget.Widget

	// OnMount is called after the screen is mounted. It may return a Cmd.
	OnMount(ctx context.Context) msg.Cmd

	// OnUnmount is called before the screen is removed.
	OnUnmount(ctx context.Context)
}

// BaseScreen provides default implementations for Screen lifecycle methods.
// Embed *BaseScreen in your concrete screen types.
type BaseScreen struct {
	widget.BaseWidget
}

// Compose returns an empty child list by default.
func (s *BaseScreen) Compose() []widget.Widget { return nil }

// OnMount is a no-op by default.
func (s *BaseScreen) OnMount(_ context.Context) msg.Cmd { return nil }

// OnUnmount is a no-op by default.
func (s *BaseScreen) OnUnmount(_ context.Context) {}
