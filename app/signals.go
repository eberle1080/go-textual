//go:build !windows

package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/eberle1080/go-textual/msg"
)

// runSignalHandler is a goroutine that handles OS signals and converts them
// to messages. It runs for the lifetime of the app.
func (a *App) runSignalHandler(ctx context.Context) {
	sigCh := make(chan os.Signal, 4)
	signal.Notify(sigCh,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGTSTP,
		syscall.SIGCONT,
	)
	defer signal.Stop(sigCh)

	for {
		select {
		case <-ctx.Done():
			return
		case sig, ok := <-sigCh:
			if !ok {
				return
			}
			switch sig {
			case syscall.SIGINT, syscall.SIGTERM:
				a.sendMsg(msg.QuitMsg{Signal: sig})
			case syscall.SIGTSTP:
				a.sendMsg(msg.SuspendMsg{})
			case syscall.SIGCONT:
				a.sendMsg(msg.ResumeMsg{})
			}
		}
	}
}
