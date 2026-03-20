package layout

import "fmt"

// MissingLayoutError is returned by GetLayout when no layout with the given
// name is registered.
type MissingLayoutError struct{ Name string }

func (e *MissingLayoutError) Error() string {
	return fmt.Sprintf("unknown layout %q", e.Name)
}

// GetLayout returns the Layout implementation for the given CSS layout name.
// The supported names are "vertical", "horizontal", "grid", and "stream".
func GetLayout(name string) (Layout, error) {
	switch name {
	case "vertical", "":
		return VerticalLayout{}, nil
	case "horizontal":
		return HorizontalLayout{}, nil
	case "grid":
		return &GridLayout{}, nil
	case "stream":
		return StreamLayout{}, nil
	default:
		return nil, &MissingLayoutError{Name: name}
	}
}
