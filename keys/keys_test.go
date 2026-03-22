package keys

import "testing"

func TestKeyConstants(t *testing.T) {
	if Escape != "escape" {
		t.Error("Escape constant wrong")
	}
	if Tab != "tab" {
		t.Error("Tab constant wrong")
	}
	if Enter != "enter" {
		t.Error("Enter constant wrong")
	}
	if Up != "up" {
		t.Error("Up constant wrong")
	}
	if F1 != "f1" {
		t.Error("F1 constant wrong")
	}
	if ControlA != "ctrl+a" {
		t.Error("ControlA constant wrong")
	}
	if BackTab != "shift+tab" {
		t.Error("BackTab constant wrong")
	}
}

func TestGetKeyAliases(t *testing.T) {
	aliases := GetKeyAliases("tab")
	if len(aliases) < 2 {
		t.Errorf("expected at least 2 aliases for tab, got %v", aliases)
	}
	found := false
	for _, a := range aliases {
		if a == "ctrl+i" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected ctrl+i in tab aliases, got %v", aliases)
	}

	aliases = GetKeyAliases("enter")
	found = false
	for _, a := range aliases {
		if a == "ctrl+m" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected ctrl+m in enter aliases, got %v", aliases)
	}

	// Key with no aliases returns just itself
	aliases = GetKeyAliases("f1")
	if len(aliases) != 1 || aliases[0] != "f1" {
		t.Errorf("f1 should have only itself as alias, got %v", aliases)
	}
}

func TestFormatKey(t *testing.T) {
	tests := []struct {
		key, want string
	}{
		{"up", "↑"},
		{"down", "↓"},
		{"left", "←"},
		{"right", "→"},
		{"backspace", "⌫"},
		{"escape", "esc"},
		{"enter", "⏎"},
		{"space", "space"},
		{"pagedown", "pgdn"},
		{"pageup", "pgup"},
		{"delete", "del"},
	}
	for _, tt := range tests {
		got := FormatKey(tt.key)
		if got != tt.want {
			t.Errorf("FormatKey(%q) = %q, want %q", tt.key, got, tt.want)
		}
	}
}

func TestKeyToCharacter(t *testing.T) {
	tests := []struct {
		key  string
		char string
		ok   bool
	}{
		{"a", "a", true},
		{"z", "z", true},
		{"1", "1", true},
		{"ctrl+a", "", false},  // modifier — no char
		{"shift+a", "", false}, // modifier — no char
		{"exclamation_mark", "!", true},
		{"question_mark", "?", true},
		// Parity cases from textual/tests/test_keys.py
		{"space", " ", true},
		{"tilde", "~", true},
		{"pound_sign", "£", true},
	}
	for _, tt := range tests {
		got, ok := KeyToCharacter(tt.key)
		if ok != tt.ok {
			t.Errorf("KeyToCharacter(%q) ok = %v, want %v", tt.key, ok, tt.ok)
		}
		if ok && got != tt.char {
			t.Errorf("KeyToCharacter(%q) = %q, want %q", tt.key, got, tt.char)
		}
	}
}

func TestCharacterToKey(t *testing.T) {
	tests := []struct {
		char string
		key  string
	}{
		{"a", "a"},
		{"1", "1"},
		{"~", "tilde"},
		{"£", "pound_sign"},
		{"!", "exclamation_mark"},
		{" ", "space"},
		{"/", "slash"},
		{"\\", "backslash"},
		{"@", "at"},
		{"+", "plus"},
		{"-", "minus"},
		{"_", "underscore"},
	}
	for _, tt := range tests {
		got := CharacterToKey(tt.char)
		if got != tt.key {
			t.Errorf("CharacterToKey(%q) = %q, want %q", tt.char, got, tt.key)
		}
	}
}

func TestKeyToCharacterReplacementNames(t *testing.T) {
	// Friendly replacement key names must reverse back to the correct character.
	tests := []struct {
		key  string
		char string
	}{
		{"slash", "/"},
		{"backslash", "\\"},
		{"at", "@"},
		{"plus", "+"},
		{"minus", "-"},
		{"underscore", "_"},
	}
	for _, tt := range tests {
		got, ok := KeyToCharacter(tt.key)
		if !ok {
			t.Errorf("KeyToCharacter(%q) returned ok=false", tt.key)
			continue
		}
		if got != tt.char {
			t.Errorf("KeyToCharacter(%q) = %q, want %q", tt.key, got, tt.char)
		}
	}
}

func TestReplacementNameRoundTrip(t *testing.T) {
	// CharacterToKey and KeyToCharacter must be inverses for replacement names.
	chars := []string{"/", "\\", "@", "+", "-", "_"}
	for _, ch := range chars {
		key := CharacterToKey(ch)
		got, ok := KeyToCharacter(key)
		if !ok {
			t.Errorf("KeyToCharacter(CharacterToKey(%q)=%q) returned ok=false", ch, key)
			continue
		}
		if got != ch {
			t.Errorf("round-trip %q -> %q -> %q", ch, key, got)
		}
	}
}

func TestKeyToCharacterMultibyte(t *testing.T) {
	// Single multibyte Unicode runes used directly as key identifiers must be
	// recognized as one-character keys, not rejected by a byte-length check.
	tests := []struct {
		key  string
		char string
		ok   bool
	}{
		{"£", "£", true}, // U+00A3 POUND SIGN — 2 bytes in UTF-8
		{"€", "€", true}, // U+20AC EURO SIGN — 3 bytes in UTF-8
	}
	for _, tt := range tests {
		got, ok := KeyToCharacter(tt.key)
		if ok != tt.ok {
			t.Errorf("KeyToCharacter(%q) ok = %v, want %v", tt.key, ok, tt.ok)
		}
		if ok && got != tt.char {
			t.Errorf("KeyToCharacter(%q) = %q, want %q", tt.key, got, tt.char)
		}
	}
}

func TestCharacterToKeyRoundTrip(t *testing.T) {
	// Verify that CharacterToKey and KeyToCharacter are inverses for key chars.
	chars := []string{"~", "£", "!", " ", "a", "z", "0"}
	for _, ch := range chars {
		key := CharacterToKey(ch)
		got, ok := KeyToCharacter(key)
		if !ok {
			t.Errorf("KeyToCharacter(CharacterToKey(%q)=%q) returned ok=false", ch, key)
			continue
		}
		if got != ch {
			t.Errorf("round-trip %q -> %q -> %q", ch, key, got)
		}
	}
}

func TestKeyToCharacterMultibyteRoundTrip(t *testing.T) {
	// Round-trip: raw multibyte rune as key → KeyToCharacter → same rune.
	// This covers the path where the key IS the character (no named lookup needed).
	raws := []string{"£", "€"}
	for _, ch := range raws {
		got, ok := KeyToCharacter(ch)
		if !ok {
			t.Errorf("KeyToCharacter(%q) returned ok=false", ch)
			continue
		}
		if got != ch {
			t.Errorf("KeyToCharacter(%q) = %q, want %q", ch, got, ch)
		}
	}
}
