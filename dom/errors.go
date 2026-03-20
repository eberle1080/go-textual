package dom

import "fmt"

// NoMatchesError is returned by QueryOne when no nodes match the selector.
type NoMatchesError struct{ Selector string }

func (e *NoMatchesError) Error() string {
	return fmt.Sprintf("no matches for selector %q", e.Selector)
}

// TooManyMatchesError is returned by QueryOne when more than one node matches.
type TooManyMatchesError struct {
	Selector string
	Count    int
}

func (e *TooManyMatchesError) Error() string {
	return fmt.Sprintf("selector %q matched %d nodes; expected exactly 1", e.Selector, e.Count)
}

// InvalidQueryFormatError is returned when a CSS selector string cannot be parsed.
type InvalidQueryFormatError struct {
	Selector string
	Cause    error
}

func (e *InvalidQueryFormatError) Error() string {
	return fmt.Sprintf("invalid query %q: %v", e.Selector, e.Cause)
}

func (e *InvalidQueryFormatError) Unwrap() error { return e.Cause }
