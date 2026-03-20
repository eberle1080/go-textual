package app

import "fmt"

// PanicError wraps a value recovered from a panic inside a Cmd goroutine.
type PanicError struct {
	Recovered any
	Stack     []byte
}

func (e *PanicError) Error() string {
	return fmt.Sprintf("panic in cmd: %v", e.Recovered)
}

// DriverError wraps an error from the terminal driver.
type DriverError struct {
	Op  string
	Err error
}

func (e *DriverError) Error() string {
	return fmt.Sprintf("driver %s: %v", e.Op, e.Err)
}

func (e *DriverError) Unwrap() error { return e.Err }

// CSSError wraps a CSS parsing or application error.
type CSSError struct {
	Source string
	Err    error
}

func (e *CSSError) Error() string {
	return fmt.Sprintf("css %s: %v", e.Source, e.Err)
}

func (e *CSSError) Unwrap() error { return e.Err }
