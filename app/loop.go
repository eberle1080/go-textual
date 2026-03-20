package app

import (
	"context"
	"runtime/debug"

	"github.com/eberle1080/go-textual/msg"
	"github.com/eberle1080/go-textual/widget"
)

// runLoop is the main event loop goroutine. It owns all widget state.
func (a *App) runLoop() error {
	// First defer: restore terminal on panic.
	defer func() {
		if r := recover(); r != nil {
			a.drv.StopApplicationMode()
			panic(r)
		}
	}()

	for {
		// Drain pending messages and collect resulting cmds.
		cmds := a.drainAndDispatch()
		for _, cmd := range cmds {
			go a.runCmd(a.ctx, cmd)
		}

		if a.quitting {
			return nil
		}

		a.renderIfDirty()

		// Block until next message arrives.
		select {
		case <-a.ctx.Done():
			return context.Cause(a.ctx)
		case m := <-a.events:
			a.sendMsg(m)
		}
	}
}

// drainAndDispatch reads all currently queued messages and dispatches them.
// Returns the aggregate list of Cmds produced.
func (a *App) drainAndDispatch() []msg.Cmd {
	var cmds []msg.Cmd
	for {
		select {
		case m := <-a.events:
			if c := a.dispatch(m); c != nil {
				cmds = append(cmds, c)
			}
		default:
			return cmds
		}
	}
}

// dispatch handles a single message.
func (a *App) dispatch(m msg.Msg) msg.Cmd {
	switch v := m.(type) {
	case msg.QuitMsg:
		a.quitting = true
		return nil

	case msg.ResizeSignalMsg:
		return nil

	case msg.ResizeMsg:
		a.size = v.Size
		if a.screen != nil {
			a.screen.MarkDirty()
		}
		return nil

	case msg.SuspendMsg:
		a.drv.SuspendApplicationMode()
		return nil

	case msg.ResumeMsg:
		a.drv.ResumeApplicationMode()
		if a.screen != nil {
			a.screen.MarkDirty()
		}
		return nil

	case msg.PanicMsg:
		a.logger.Error("panic in cmd",
			"recovered", v.Recovered,
			"stack", string(v.Stack))
		// Fall through to let screen handle it (e.g. quit on unrecoverable errors).

	case pushScreenMsg:
		a.screen = v.screen
		return a.mountScreen(a.ctx, v.screen)

	case msg.KeyMsg:
		if a.screen != nil {
			switch v.Key {
			case "tab":
				a.focused = widget.NextFocus(a.screen, a.focused)
				return nil
			case "shift+tab":
				a.focused = widget.PrevFocus(a.screen, a.focused)
				return nil
			}
		}

	case msg.MouseDownMsg, msg.MouseUpMsg, msg.MouseMoveMsg,
		msg.MouseScrollUpMsg, msg.MouseScrollDownMsg,
		msg.MouseScrollLeftMsg, msg.MouseScrollRightMsg:
		return a.dispatchMouse(m)
	}

	// Route to focused widget first, then screen.
	if a.focused != nil {
		if cmd := a.focused.Update(a.ctx, m); cmd != nil {
			return cmd
		}
	}
	if a.screen != nil {
		if cmd := a.screen.Update(a.ctx, m); cmd != nil {
			return cmd
		}
	}
	return nil
}

// dispatchMouse performs hit-testing and routes a mouse message to the widget
// under the cursor, translating coordinates to widget-local space.
func (a *App) dispatchMouse(m msg.Msg) msg.Cmd {
	if a.screen == nil {
		return nil
	}
	sx, sy, ok := mouseScreenXY(m)
	if !ok {
		return nil
	}
	target := widgetAt(a.screen, sx, sy)

	switch m.(type) {
	case msg.MouseDownMsg:
		// Move focus to clicked widget if it's focusable.
		if target != nil && target.CanFocus() {
			a.focused = target
		}
		a.mouseDownTarget = target

	case msg.MouseUpMsg:
		// Clear the down-target after processing below.
		defer func() { a.mouseDownTarget = nil }()
	}

	if target == nil {
		// No widget hit — let the screen handle it.
		if cmd := a.screen.Update(a.ctx, m); cmd != nil {
			return cmd
		}
		return nil
	}

	translated := translateMouseMsg(m, target)
	if cmd := target.Update(a.ctx, translated); cmd != nil {
		return cmd
	}
	// Also give the screen a chance to react (e.g. for global scroll handling).
	if cmd := a.screen.Update(a.ctx, translated); cmd != nil {
		return cmd
	}
	return nil
}

// runCmd executes a single Cmd with panic recovery.
func (a *App) runCmd(ctx context.Context, cmd msg.Cmd) {
	if cmd == nil {
		return
	}
	defer func() {
		if r := recover(); r != nil {
			a.sendMsg(msg.PanicMsg{Recovered: r, Stack: debug.Stack()})
		}
	}()

	result := cmd(ctx)
	if result == nil {
		return
	}

	// Handle batch/sequence transparently.
	if batchCmds, ok := msg.IsBatch(result); ok {
		for _, c := range batchCmds {
			go a.runCmd(ctx, c)
		}
		return
	}
	if seqCmds, ok := msg.IsSequence(result); ok {
		go a.runSequence(ctx, seqCmds)
		return
	}

	a.sendMsg(result)
}

// runSequence runs commands sequentially.
func (a *App) runSequence(ctx context.Context, cmds []msg.Cmd) {
	for _, cmd := range cmds {
		if cmd == nil {
			continue
		}
		result := func() (r msg.Msg) {
			defer func() {
				if rec := recover(); rec != nil {
					r = msg.PanicMsg{Recovered: rec, Stack: debug.Stack()}
				}
			}()
			return cmd(ctx)
		}()
		if result != nil {
			a.sendMsg(result)
		}
	}
}

// sendMsg delivers a message to the events channel.
func (a *App) sendMsg(m msg.Msg) {
	select {
	case a.events <- m:
	default:
		a.logger.Debug("event channel full, dropping message")
	}
}
