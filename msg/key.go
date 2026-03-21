package msg

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/eberle1080/go-textual/keys"
)

// KeyMsg represents a keyboard event.
type KeyMsg struct {
	BaseMsg
	// Key is the canonical key name, e.g. "ctrl+a", "f1", "a".
	Key string
	// Character is the printable character produced by this key, if any.
	Character *string
	// Aliases is the list of key aliases from the keys package.
	Aliases []string
}

// NewKey constructs a KeyMsg. When character is nil and key is a single Unicode
// code point, Character is set to that code point. For named keys like "space"
// or "full_stop", KeyToCharacter is used as a fallback.
func NewKey(key string, character *string) KeyMsg {
	if character == nil {
		if runes := []rune(key); len(runes) == 1 {
			ch := string(runes[0])
			character = &ch
		} else if ch, ok := keys.KeyToCharacter(key); ok {
			character = &ch
		}
	}
	return KeyMsg{
		Key:       key,
		Character: character,
		Aliases:   keys.GetKeyAliases(key),
	}
}

// IsPrintable reports whether this key produces a printable character.
func (k KeyMsg) IsPrintable() bool {
	if k.Character == nil {
		return false
	}
	r, size := utf8.DecodeRuneInString(*k.Character)
	if size == 0 || r == utf8.RuneError {
		return false
	}
	return unicode.IsPrint(r)
}

// Name returns the key name normalised to a Go identifier.
func (k KeyMsg) Name() string {
	return keyToIdentifier(k.Key)
}

// keyToIdentifier normalises a key name to a Go identifier following Textual
// convention.
func keyToIdentifier(keyName string) string {
	runes := []rune(keyName)
	if len(runes) == 1 && unicode.IsUpper(runes[0]) {
		return "upper_" + strings.ToLower(keyName)
	}
	s := strings.ReplaceAll(keyName, "+", "_")
	s = strings.ReplaceAll(s, "-", "_")
	return strings.ToLower(s)
}
