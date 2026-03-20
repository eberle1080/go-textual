// Package ansi provides ANSI escape sequence constants, functional key
// mappings, and helper functions for terminal interactions such as
// synchronized output and pointer shape changes.
//
// The [Sequences] map covers the ~250 VT/xterm sequences that Textual
// recognises. [FunctionalKeys] covers extended sequences from the Kitty
// keyboard protocol and legacy VT sequences with numeric parameters.
//
// # Sequence Lookup
//
// Translate a raw byte sequence received from the terminal into a key name:
//
//	if name, ok := ansi.Sequences["\x1b[A"]; ok {
//	    fmt.Println(name) // "up"
//	}
//
// # Functional Keys
//
// Look up extended sequences such as Kitty keyboard protocol entries and
// numeric-parameter VT sequences:
//
//	if key, ok := ansi.FunctionalKeys["\x1b[1;5A"]; ok {
//	    fmt.Println(key) // "ctrl+up"
//	}
//
// # Synchronized Output
//
// Wrap terminal writes in a synchronized-update block to eliminate tearing:
//
//	os.Stdout.WriteString(ansi.SynchronizedOutputStart)
//	os.Stdout.WriteString(renderedFrame)
//	os.Stdout.WriteString(ansi.SynchronizedOutputEnd)
package ansi
