//go:build !windows

package app

import "github.com/eberle1080/go-textual/driver"

func newPlatformDriver(sink driver.EventSink) driver.Driver {
	return driver.NewUnixDriver(sink)
}
