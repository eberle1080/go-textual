package css

import "fmt"

// TokenError is raised when the CSS cannot be tokenized (syntax error).
type TokenError struct {
	ReadFrom CSSLocation
	Code     string
	Start    [2]int
	End      [2]int
	Msg      string
}

func (e *TokenError) Error() string { return e.Msg }

// NewTokenError creates a TokenError.
// start and end are 1-indexed [line, col] pairs; if end is nil Start is used.
func NewTokenError(readFrom CSSLocation, code string, start [2]int, msg string, end *[2]int) *TokenError {
	e := end
	if e == nil {
		e = &start
	}
	return &TokenError{ReadFrom: readFrom, Code: code, Start: start, End: *e, Msg: msg}
}

// UnexpectedEndError indicates that the text being tokenized ended prematurely.
type UnexpectedEndError struct {
	TokenError
}

// DeclarationError is raised when a CSS declaration is invalid.
type DeclarationError struct {
	Name    string
	Token   Token
	Message string
}

func (e *DeclarationError) Error() string {
	return fmt.Sprintf("declaration error for %q: %s", e.Name, e.Message)
}

// StyleTypeError is raised when a style property receives the wrong type.
type StyleTypeError struct {
	Msg string
}

func (e *StyleTypeError) Error() string { return e.Msg }

// StyleValueError is raised when a style property receives an invalid value.
type StyleValueError struct {
	Msg string
}

func (e *StyleValueError) Error() string { return e.Msg }

// StylesheetError is raised for general stylesheet errors.
type StylesheetError struct {
	Msg string
}

func (e *StylesheetError) Error() string { return e.Msg }

// UnresolvedVariableError is raised when a CSS variable reference cannot be resolved.
type UnresolvedVariableError struct {
	TokenError
}
