package app

import (
	"github.com/eberle1080/go-textual/driver"
	"github.com/eberle1080/go-textual/msg"
)

// headlessDriver is a no-op driver for testing. It captures all written output
// and delivers the initial size as a ResizeMsg.
type headlessDriver struct {
	sink   driver.EventSink
	width  int
	height int
	output []string
}

func newHeadlessDriver(sink driver.EventSink, width, height int) *headlessDriver {
	return &headlessDriver{sink: sink, width: width, height: height}
}

func (d *headlessDriver) Write(data string)   { d.output = append(d.output, data) }
func (d *headlessDriver) Flush()              {}
func (d *headlessDriver) Close()              {}
func (d *headlessDriver) IsHeadless() bool    { return true }
func (d *headlessDriver) IsInline() bool      { return false }
func (d *headlessDriver) CanSuspend() bool    { return false }
func (d *headlessDriver) SuspendApplicationMode() {}
func (d *headlessDriver) ResumeApplicationMode()  {}
func (d *headlessDriver) OpenURL(_ string)         {}
func (d *headlessDriver) SetCursorOrigin(_, _ int) {}
func (d *headlessDriver) ClearCursorOrigin()       {}
func (d *headlessDriver) DisableInput()            {}

func (d *headlessDriver) StartApplicationMode() {
	d.sink.Send(msg.FromDimensions([2]int{d.width, d.height}, nil))
}

func (d *headlessDriver) StopApplicationMode() {}

// Output returns all data written to the driver.
func (d *headlessDriver) Output() []string { return d.output }
