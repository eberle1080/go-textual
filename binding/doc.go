// Package binding provides key binding management for the Textual TUI framework.
//
// # Overview
//
// Key bindings associate keyboard input with named actions. This package is a
// Go port of Python Textual's binding system, adapted to use Go structs and
// maps instead of Python dataclasses and dicts.
//
// # Key types
//
//   - [Binding] is an immutable struct describing a single key-to-action mapping.
//   - [BindingsMap] manages a map from key strings to slices of [Binding] values.
//   - [ActiveBinding] pairs a [Binding] with the DOM node where it is defined and
//     whether it is currently enabled.
//   - [KeymapApplyResult] is returned by [BindingsMap.ApplyKeymap] and reports
//     any bindings displaced (clashed) by the keymap substitution.
//
// # Binding lifecycle
//
// Bindings are declared as [Binding] values, optionally using [MakeBindings] to
// expand compound keys (comma-separated) and normalise single-character keys via
// the keys package. A [BindingsMap] is constructed from the expanded bindings.
// At application startup, [BindingsMap.ApplyKeymap] substitutes keys according to
// user-configured keymap overrides. At runtime, [BindingsMap.GetBindingsForKey]
// looks up the active bindings for an incoming key event.
package binding
