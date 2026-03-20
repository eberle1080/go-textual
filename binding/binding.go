package binding

import (
	"strings"

	"github.com/eberle1080/go-textual/keys"
)

// Group represents a named group of related bindings displayed together.
type Group struct {
	Description string
	Compact     bool
}

// Binding describes a single key-to-action mapping. It is immutable by convention.
type Binding struct {
	Key         string  // Key or comma-separated keys
	Action      string  // Action identifier
	Description string  // Display description
	Show        bool    // Show in footer
	KeyDisplay  *string // Custom footer display (nil = auto)
	Priority    bool    // Priority binding
	Tooltip     string  // Footer tooltip
	ID          *string // Globally unique binding ID
	System      bool    // System binding (hidden from key panel)
	Group       *Group  // Binding group
}

// ParseKey splits the binding's Key on "+" and returns the modifiers and final key.
// For example, "ctrl+shift+a" returns (["ctrl", "shift"], "a").
func (b Binding) ParseKey() (modifiers []string, key string) {
	parts := strings.Split(b.Key, "+")
	if len(parts) == 1 {
		return nil, parts[0]
	}
	return parts[:len(parts)-1], parts[len(parts)-1]
}

// WithKey returns a copy of this Binding with the key and optional display updated.
func (b Binding) WithKey(key string, keyDisplay *string) Binding {
	c := b
	c.Key = key
	c.KeyDisplay = keyDisplay
	return c
}

// FromTuple2 constructs a Binding from a (key, action) pair with no description.
func FromTuple2(key, action string) Binding {
	return Binding{Key: key, Action: action, Show: true}
}

// FromTuple3 constructs a Binding from a (key, action, description) triple.
func FromTuple3(key, action, description string) Binding {
	return Binding{Key: key, Action: action, Description: description, Show: description != ""}
}

// MakeBindings expands a slice of binding inputs into individual Binding values.
// Each element of inputs must be one of:
//   - Binding — used directly after key expansion
//   - [2]string — treated as (key, action)
//   - [3]string — treated as (key, action, description)
//
// Compound keys (comma-separated) are split into separate Binding values.
// Single-character keys are normalised via keys.CharacterToKey.
func MakeBindings(inputs []any) ([]Binding, error) {
	var result []Binding
	for _, input := range inputs {
		var b Binding
		switch v := input.(type) {
		case Binding:
			b = v
		case [2]string:
			b = FromTuple2(v[0], v[1])
		case [3]string:
			b = FromTuple3(v[0], v[1], v[2])
		default:
			continue
		}
		expanded, err := expandBinding(b)
		if err != nil {
			return nil, err
		}
		result = append(result, expanded...)
	}
	return result, nil
}

// expandBinding splits a binding's Key on "," and returns one Binding per key.
func expandBinding(b Binding) ([]Binding, error) {
	rawKeys := strings.Split(b.Key, ",")
	var out []Binding
	for _, k := range rawKeys {
		k = strings.TrimSpace(k)
		if k == "" {
			return nil, &InvalidBinding{BindingError{Msg: "empty key in binding: " + b.Key}}
		}
		// Normalise single-character keys.
		runes := []rune(k)
		if len(runes) == 1 {
			k = keys.CharacterToKey(k)
		}
		out = append(out, b.WithKey(k, b.KeyDisplay))
	}
	return out, nil
}
