//go:build !windows

package driver

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/eberle1080/go-textual/msg"
)

// installSIGWINCH registers a SIGWINCH handler that sends ResizeSignalMsg to
// the sink. Returns a cleanup function that must be called to deregister.
func installSIGWINCH(sink EventSink) (stop func()) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGWINCH)
	done := make(chan struct{})
	go func() {
		defer close(done)
		for range ch {
			sink.Send(msg.ResizeSignalMsg{})
		}
	}()
	return func() {
		signal.Stop(ch)
		close(ch)
		<-done
	}
}
