//go:build windows

package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/eberle1080/go-textual/msg"
)

func (a *App) runSignalHandler(ctx context.Context) {
	sigCh := make(chan os.Signal, 2)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
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
			}
		}
	}
}
