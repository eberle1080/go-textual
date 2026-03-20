package keys

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// Key constants — Escape and control keys.
const (
	Escape      = "escape"
	ShiftEscape = "shift+escape"
	Return      = "return"

	ControlAt = "ctrl+@"

	ControlA = "ctrl+a"
	ControlB = "ctrl+b"
	ControlC = "ctrl+c"
	ControlD = "ctrl+d"
	ControlE = "ctrl+e"
	ControlF = "ctrl+f"
	ControlG = "ctrl+g"
	ControlH = "ctrl+h"
	ControlI = "ctrl+i" // Tab
	ControlJ = "ctrl+j" // Newline
	ControlK = "ctrl+k"
	ControlL = "ctrl+l"
	ControlM = "ctrl+m" // Carriage return
	ControlN = "ctrl+n"
	ControlO = "ctrl+o"
	ControlP = "ctrl+p"
	ControlQ = "ctrl+q"
	ControlR = "ctrl+r"
	ControlS = "ctrl+s"
	ControlT = "ctrl+t"
	ControlU = "ctrl+u"
	ControlV = "ctrl+v"
	ControlW = "ctrl+w"
	ControlX = "ctrl+x"
	ControlY = "ctrl+y"
	ControlZ = "ctrl+z"

	Control1 = "ctrl+1"
	Control2 = "ctrl+2"
	Control3 = "ctrl+3"
	Control4 = "ctrl+4"
	Control5 = "ctrl+5"
	Control6 = "ctrl+6"
	Control7 = "ctrl+7"
	Control8 = "ctrl+8"
	Control9 = "ctrl+9"
	Control0 = "ctrl+0"

	ControlShift1 = "ctrl+shift+1"
	ControlShift2 = "ctrl+shift+2"
	ControlShift3 = "ctrl+shift+3"
	ControlShift4 = "ctrl+shift+4"
	ControlShift5 = "ctrl+shift+5"
	ControlShift6 = "ctrl+shift+6"
	ControlShift7 = "ctrl+shift+7"
	ControlShift8 = "ctrl+shift+8"
	ControlShift9 = "ctrl+shift+9"
	ControlShift0 = "ctrl+shift+0"

	ControlBackslash    = "ctrl+backslash"
	ControlSquareClose  = "ctrl+right_square_bracket"
	ControlCircumflex   = "ctrl+circumflex_accent"
	ControlUnderscore   = "ctrl+underscore"
)

// Arrow and navigation keys.
const (
	Left     = "left"
	Right    = "right"
	Up       = "up"
	Down     = "down"
	Home     = "home"
	End      = "end"
	Insert   = "insert"
	Delete   = "delete"
	PageUp   = "pageup"
	PageDown = "pagedown"
)

// Ctrl+navigation keys.
const (
	ControlLeft     = "ctrl+left"
	ControlRight    = "ctrl+right"
	ControlUp       = "ctrl+up"
	ControlDown     = "ctrl+down"
	ControlHome     = "ctrl+home"
	ControlEnd      = "ctrl+end"
	ControlInsert   = "ctrl+insert"
	ControlDelete   = "ctrl+delete"
	ControlPageUp   = "ctrl+pageup"
	ControlPageDown = "ctrl+pagedown"
)

// Shift+navigation keys.
const (
	ShiftLeft     = "shift+left"
	ShiftRight    = "shift+right"
	ShiftUp       = "shift+up"
	ShiftDown     = "shift+down"
	ShiftHome     = "shift+home"
	ShiftEnd      = "shift+end"
	ShiftInsert   = "shift+insert"
	ShiftDelete   = "shift+delete"
	ShiftPageUp   = "shift+pageup"
	ShiftPageDown = "shift+pagedown"
)

// Ctrl+Shift+navigation keys.
const (
	ControlShiftLeft     = "ctrl+shift+left"
	ControlShiftRight    = "ctrl+shift+right"
	ControlShiftUp       = "ctrl+shift+up"
	ControlShiftDown     = "ctrl+shift+down"
	ControlShiftHome     = "ctrl+shift+home"
	ControlShiftEnd      = "ctrl+shift+end"
	ControlShiftInsert   = "ctrl+shift+insert"
	ControlShiftDelete   = "ctrl+shift+delete"
	ControlShiftPageUp   = "ctrl+shift+pageup"
	ControlShiftPageDown = "ctrl+shift+pagedown"

	BackTab = "shift+tab"
)

// Function keys.
const (
	F1  = "f1"
	F2  = "f2"
	F3  = "f3"
	F4  = "f4"
	F5  = "f5"
	F6  = "f6"
	F7  = "f7"
	F8  = "f8"
	F9  = "f9"
	F10 = "f10"
	F11 = "f11"
	F12 = "f12"
	F13 = "f13"
	F14 = "f14"
	F15 = "f15"
	F16 = "f16"
	F17 = "f17"
	F18 = "f18"
	F19 = "f19"
	F20 = "f20"
	F21 = "f21"
	F22 = "f22"
	F23 = "f23"
	F24 = "f24"
)

// Ctrl+function keys.
const (
	ControlF1  = "ctrl+f1"
	ControlF2  = "ctrl+f2"
	ControlF3  = "ctrl+f3"
	ControlF4  = "ctrl+f4"
	ControlF5  = "ctrl+f5"
	ControlF6  = "ctrl+f6"
	ControlF7  = "ctrl+f7"
	ControlF8  = "ctrl+f8"
	ControlF9  = "ctrl+f9"
	ControlF10 = "ctrl+f10"
	ControlF11 = "ctrl+f11"
	ControlF12 = "ctrl+f12"
	ControlF13 = "ctrl+f13"
	ControlF14 = "ctrl+f14"
	ControlF15 = "ctrl+f15"
	ControlF16 = "ctrl+f16"
	ControlF17 = "ctrl+f17"
	ControlF18 = "ctrl+f18"
	ControlF19 = "ctrl+f19"
	ControlF20 = "ctrl+f20"
	ControlF21 = "ctrl+f21"
	ControlF22 = "ctrl+f22"
	ControlF23 = "ctrl+f23"
	ControlF24 = "ctrl+f24"
)

// Special keys.
const (
	Any       = "<any>"
	ScrollUp  = "<scroll-up>"
	ScrollDown = "<scroll-down>"
	Ignore    = "<ignore>"

	// Aliases for common keys.
	ControlSpace = "ctrl-at"
	Tab          = "tab"
	Space        = "space"
	Enter        = "enter"
	Backspace    = "backspace"

	// ShiftControl aliases (renamed to ControlShift).
	ShiftControlLeft  = ControlShiftLeft
	ShiftControlRight = ControlShiftRight
	ShiftControlHome  = ControlShiftHome
	ShiftControlEnd   = ControlShiftEnd
)

// KeyNameReplacements replaces obscure Unicode names with common terms.
var KeyNameReplacements = map[string]string{
	"solidus":         "slash",
	"reverse_solidus": "backslash",
	"commercial_at":   "at",
	"hyphen_minus":    "minus",
	"plus_sign":       "plus",
	"low_line":        "underscore",
}

// ReplacedKeys is the reverse of KeyNameReplacements.
var ReplacedKeys = map[string]string{
	"slash":      "solidus",
	"backslash":  "reverse_solidus",
	"at":         "commercial_at",
	"minus":      "hyphen_minus",
	"plus":       "plus_sign",
	"underscore": "low_line",
}

// KeyToUnicodeName maps friendly key names back to their Unicode character names.
var KeyToUnicodeName = map[string]string{
	"exclamation_mark":     "EXCLAMATION MARK",
	"quotation_mark":       "QUOTATION MARK",
	"number_sign":          "NUMBER SIGN",
	"dollar_sign":          "DOLLAR SIGN",
	"percent_sign":         "PERCENT SIGN",
	"left_parenthesis":     "LEFT PARENTHESIS",
	"right_parenthesis":    "RIGHT PARENTHESIS",
	"plus_sign":            "PLUS SIGN",
	"hyphen_minus":         "HYPHEN-MINUS",
	"full_stop":            "FULL STOP",
	"less_than_sign":       "LESS-THAN SIGN",
	"equals_sign":          "EQUALS SIGN",
	"greater_than_sign":    "GREATER-THAN SIGN",
	"question_mark":        "QUESTION MARK",
	"commercial_at":        "COMMERCIAL AT",
	"left_square_bracket":  "LEFT SQUARE BRACKET",
	"reverse_solidus":      "REVERSE SOLIDUS",
	"right_square_bracket": "RIGHT SQUARE BRACKET",
	"circumflex_accent":    "CIRCUMFLEX ACCENT",
	"low_line":             "LOW LINE",
	"grave_accent":         "GRAVE ACCENT",
	"left_curly_bracket":   "LEFT CURLY BRACKET",
	"vertical_line":        "VERTICAL LINE",
	"right_curly_bracket":  "RIGHT CURLY BRACKET",
	"pound_sign":           "POUND SIGN",
	"tilde":                "TILDE",
}

// KeyAliases maps keys to their alternative names.
var KeyAliases = map[string][]string{
	"tab":    {"ctrl+i"},
	"enter":  {"ctrl+m"},
	"escape": {"ctrl+left_square_brace"},
	"ctrl+at": {"ctrl+space"},
	"ctrl+j": {"newline"},
}

// KeyDisplayAliases maps key names to their display representations.
var KeyDisplayAliases = map[string]string{
	"up":        "↑",
	"down":      "↓",
	"left":      "←",
	"right":     "→",
	"backspace": "⌫",
	"escape":    "esc",
	"enter":     "⏎",
	"minus":     "-",
	"space":     "space",
	"pagedown":  "pgdn",
	"pageup":    "pgup",
	"delete":    "del",
}

// unicodeNameToChar maps Unicode character names to their characters.
// This replaces Python's unicodedata.lookup() for the names we care about.
var unicodeNameToChar = map[string]rune{
	"EXCLAMATION MARK":     '!',
	"QUOTATION MARK":       '"',
	"NUMBER SIGN":          '#',
	"DOLLAR SIGN":          '$',
	"PERCENT SIGN":         '%',
	"LEFT PARENTHESIS":     '(',
	"RIGHT PARENTHESIS":    ')',
	"PLUS SIGN":            '+',
	"HYPHEN-MINUS":         '-',
	"FULL STOP":            '.',
	"LESS-THAN SIGN":       '<',
	"EQUALS SIGN":          '=',
	"GREATER-THAN SIGN":    '>',
	"QUESTION MARK":        '?',
	"COMMERCIAL AT":        '@',
	"LEFT SQUARE BRACKET":  '[',
	"REVERSE SOLIDUS":      '\\',
	"RIGHT SQUARE BRACKET": ']',
	"CIRCUMFLEX ACCENT":    '^',
	"LOW LINE":             '_',
	"GRAVE ACCENT":         '`',
	"LEFT CURLY BRACKET":   '{',
	"VERTICAL LINE":        '|',
	"RIGHT CURLY BRACKET":  '}',
	"SPACE":                ' ',
	"SOLIDUS":              '/',
	"TILDE":                '~',
	"APOSTROPHE":           '\'',
	"AMPERSAND":            '&',
	"ASTERISK":             '*',
	"COMMA":                ',',
	"COLON":                ':',
	"SEMICOLON":            ';',
	"POUND SIGN":           '£',
}

// GetKeyAliases returns all aliases for the given key, including the key itself.
func GetKeyAliases(key string) []string {
	return append([]string{key}, KeyAliases[key]...)
}

// FormatKey returns the display representation of a key.
// For keys with display aliases this returns the alias; otherwise attempts
// Unicode lookup, falling back to the key name itself.
func FormatKey(key string) string {
	if alias, ok := KeyDisplayAliases[key]; ok {
		return alias
	}
	originalKey, ok := ReplacedKeys[key]
	if !ok {
		originalKey = key
	}
	unicodeName, ok := KeyToUnicodeName[originalKey]
	if !ok {
		unicodeName = originalKey
	}
	if r, ok := unicodeNameToChar[strings.ToUpper(unicodeName)]; ok {
		if unicode.IsPrint(r) {
			return string(r)
		}
	}
	return unicodeName
}

// KeyToCharacter returns the character associated with a key identifier.
// Returns ("", false) if no character could be determined.
func KeyToCharacter(key string) (string, bool) {
	// Strip modifier prefix
	_, sep, _ := strings.Cut(key, "+")
	if sep != "" {
		// Modifier present (not just shift) — no printable character
		return "", false
	}
	if utf8.RuneCountInString(key) == 1 {
		return key, true
	}
	// Normalize friendly replacement names (e.g. "slash" → "solidus") before
	// consulting the Unicode name tables, so CharacterToKey and KeyToCharacter
	// are reliable inverses for every identifier KeyNameReplacements covers.
	lookupKey := key
	if original, ok := ReplacedKeys[key]; ok {
		lookupKey = original
	}
	// Try KeyToUnicodeName lookup
	if unicodeName, ok := KeyToUnicodeName[lookupKey]; ok {
		if r, ok := unicodeNameToChar[strings.ToUpper(unicodeName)]; ok {
			return string(r), true
		}
	}
	// Try replacing underscores with spaces and uppercasing
	candidate := strings.ToUpper(strings.ReplaceAll(lookupKey, "_", " "))
	if r, ok := unicodeNameToChar[candidate]; ok {
		return string(r), true
	}
	return "", false
}

// CharacterToKey converts a single character to its key identifier.
func CharacterToKey(char string) string {
	if utf8.RuneCountInString(char) != 1 {
		return char
	}
	c, _ := utf8.DecodeRuneInString(char)
	if unicode.IsLetter(c) || unicode.IsDigit(c) {
		return char
	}
	// Look up Unicode name
	for name, r := range unicodeNameToChar {
		if r == c {
			key := strings.ToLower(strings.ReplaceAll(strings.ReplaceAll(name, " ", "_"), "-", "_"))
			if replacement, ok := KeyNameReplacements[key]; ok {
				return replacement
			}
			return key
		}
	}
	if char == "\t" {
		return "tab"
	}
	return char
}
