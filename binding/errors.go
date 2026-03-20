package binding

// BindingError is the base error for binding operations.
type BindingError struct {
	Msg string
}

func (e *BindingError) Error() string { return e.Msg }

// NoBinding is returned when no binding is found for a given key.
type NoBinding struct {
	BindingError
}

// InvalidBinding is returned when a binding key format is invalid.
type InvalidBinding struct {
	BindingError
}
