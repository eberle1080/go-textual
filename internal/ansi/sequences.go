package ansi

import (
	"fmt"

	"github.com/eberle1080/go-textual/keys"
)

// SequenceResult is the value stored in the [Sequences] map.
// Either Keys or Character will be populated; an ignored sequence has neither.
type SequenceResult struct {
	// Keys is the ordered list of key names to emit, e.g. [keys.Escape, keys.Up].
	Keys []string
	// Character is the printable character produced by this sequence (keypad
	// mappings like "*", "+", etc.).
	Character string
	// ignored marks sequences that should be silently consumed.
	ignored bool
}

// IgnoreSequence is a sentinel SequenceResult for sequences that should be
// consumed without producing any events.
var IgnoreSequence = SequenceResult{ignored: true}

// IsIgnored reports whether r represents a sequence that should be discarded.
func IsIgnored(r SequenceResult) bool { return r.ignored }

// keys1 returns a SequenceResult with a single key name.
func keys1(k string) SequenceResult { return SequenceResult{Keys: []string{k}} }

// keys2 returns a SequenceResult with two key names (multi-key sequences such
// as escape + letter emitted as a pair).
func keys2(k1, k2 string) SequenceResult { return SequenceResult{Keys: []string{k1, k2}} }

// char returns a SequenceResult for a printable character mapping.
func char(c string) SequenceResult { return SequenceResult{Character: c} }

// SyncStart is the escape sequence to begin a synchronized update block.
const SyncStart = "\x1b[?2026h"

// SyncEnd is the escape sequence to end a synchronized update block.
const SyncEnd = "\x1b[?2026l"

// SetPointerShape returns the OSC escape sequence that changes the mouse
// pointer shape to the named CSS cursor value (e.g. "default", "pointer").
func SetPointerShape(shape string) string {
	return fmt.Sprintf("\x1b]22;%s\x07", shape)
}

// Sequences maps raw escape-sequence strings to their SequenceResult.
// Populated by init().
var Sequences map[string]SequenceResult

// FunctionalKeys maps extended key protocol tokens (number + final character)
// to canonical key names.  Tokens are formed by concatenating the numeric
// parameter (if present) with the final byte, e.g. "27u", "1A", "3~".
// Populated by init().
var FunctionalKeys map[string]string

func init() {
	initSequences()
	initFunctionalKeys()
}

func initSequences() {
	Sequences = map[string]SequenceResult{
		// --- Single-byte control characters ---
		"\x00": keys1(keys.ControlAt),
		"\x01": keys1(keys.ControlA),
		"\x02": keys1(keys.ControlB),
		"\x03": keys1(keys.ControlC),
		"\x04": keys1(keys.ControlD),
		"\x05": keys1(keys.ControlE),
		"\x06": keys1(keys.ControlF),
		"\x07": keys1(keys.ControlG),
		"\x08": keys1(keys.ControlH),
		"\x09": keys1(keys.Tab),
		"\x0a": keys1(keys.ControlJ),
		"\x0b": keys1(keys.ControlK),
		"\x0c": keys1(keys.ControlL),
		"\x0d": keys1(keys.Enter),
		"\x0e": keys1(keys.ControlN),
		"\x0f": keys1(keys.ControlO),
		"\x10": keys1(keys.ControlP),
		"\x11": keys1(keys.ControlQ),
		"\x12": keys1(keys.ControlR),
		"\x13": keys1(keys.ControlS),
		"\x14": keys1(keys.ControlT),
		"\x15": keys1(keys.ControlU),
		"\x16": keys1(keys.ControlV),
		"\x17": keys1(keys.ControlW),
		"\x18": keys1(keys.ControlX),
		"\x19": keys1(keys.ControlY),
		"\x1a": keys1(keys.ControlZ),
		"\x1c": keys1(keys.ControlBackslash),
		"\x1d": keys1(keys.ControlSquareClose),
		"\x1e": keys1(keys.ControlCircumflex),
		"\x1f": keys1(keys.ControlUnderscore),
		"\x7f": keys1(keys.Backspace),

		// --- Alt + single character (ESC prefix) ---
		"\x1ba": keys2(keys.Escape, "a"),
		"\x1bb": keys2(keys.Escape, "b"),
		"\x1bc": keys2(keys.Escape, "c"),
		"\x1bd": keys2(keys.Escape, "d"),
		"\x1be": keys2(keys.Escape, "e"),
		"\x1bf": keys2(keys.Escape, "f"),
		"\x1bg": keys2(keys.Escape, "g"),
		"\x1bh": keys2(keys.Escape, "h"),
		"\x1bi": keys2(keys.Escape, "i"),
		"\x1bj": keys2(keys.Escape, "j"),
		"\x1bk": keys2(keys.Escape, "k"),
		"\x1bl": keys2(keys.Escape, "l"),
		"\x1bm": keys2(keys.Escape, "m"),
		"\x1bn": keys2(keys.Escape, "n"),
		"\x1bo": keys2(keys.Escape, "o"),
		"\x1bp": keys2(keys.Escape, "p"),
		"\x1bq": keys2(keys.Escape, "q"),
		"\x1br": keys2(keys.Escape, "r"),
		"\x1bs": keys2(keys.Escape, "s"),
		"\x1bt": keys2(keys.Escape, "t"),
		"\x1bu": keys2(keys.Escape, "u"),
		"\x1bv": keys2(keys.Escape, "v"),
		"\x1bw": keys2(keys.Escape, "w"),
		"\x1bx": keys2(keys.Escape, "x"),
		"\x1by": keys2(keys.Escape, "y"),
		"\x1bz": keys2(keys.Escape, "z"),

		// --- CSI arrow keys ---
		"\x1b[A": keys1(keys.Up),
		"\x1b[B": keys1(keys.Down),
		"\x1b[C": keys1(keys.Right),
		"\x1b[D": keys1(keys.Left),

		// --- Application mode arrow keys (SS3) ---
		"\x1bOA": keys1(keys.Up),
		"\x1bOB": keys1(keys.Down),
		"\x1bOC": keys1(keys.Right),
		"\x1bOD": keys1(keys.Left),

		// --- Navigation keys ---
		"\x1b[H": keys1(keys.Home),
		"\x1b[F": keys1(keys.End),
		"\x1b[1~": keys1(keys.Home),
		"\x1b[4~": keys1(keys.End),
		"\x1b[2~": keys1(keys.Insert),
		"\x1b[3~": keys1(keys.Delete),
		"\x1b[5~": keys1(keys.PageUp),
		"\x1b[6~": keys1(keys.PageDown),
		"\x1b[7~": keys1(keys.Home),
		"\x1b[8~": keys1(keys.End),

		// Application mode home/end
		"\x1bOH": keys1(keys.Home),
		"\x1bOF": keys1(keys.End),

		// --- Function keys (application mode) ---
		"\x1bOP": keys1(keys.F1),
		"\x1bOQ": keys1(keys.F2),
		"\x1bOR": keys1(keys.F3),
		"\x1bOS": keys1(keys.F4),

		// --- Function keys (VT220 style) ---
		"\x1b[11~": keys1(keys.F1),
		"\x1b[12~": keys1(keys.F2),
		"\x1b[13~": keys1(keys.F3),
		"\x1b[14~": keys1(keys.F4),
		"\x1b[15~": keys1(keys.F5),
		"\x1b[17~": keys1(keys.F6),
		"\x1b[18~": keys1(keys.F7),
		"\x1b[19~": keys1(keys.F8),
		"\x1b[20~": keys1(keys.F9),
		"\x1b[21~": keys1(keys.F10),
		"\x1b[23~": keys1(keys.F11),
		"\x1b[24~": keys1(keys.F12),
		"\x1b[25~": keys1(keys.F13),
		"\x1b[26~": keys1(keys.F14),
		"\x1b[28~": keys1(keys.F15),
		"\x1b[29~": keys1(keys.F16),
		"\x1b[31~": keys1(keys.F17),
		"\x1b[32~": keys1(keys.F18),
		"\x1b[33~": keys1(keys.F19),
		"\x1b[34~": keys1(keys.F20),

		// --- Shift + arrows ---
		"\x1b[1;2A": keys1(keys.ShiftUp),
		"\x1b[1;2B": keys1(keys.ShiftDown),
		"\x1b[1;2C": keys1(keys.ShiftRight),
		"\x1b[1;2D": keys1(keys.ShiftLeft),
		"\x1b[1;2H": keys1(keys.ShiftHome),
		"\x1b[1;2F": keys1(keys.ShiftEnd),

		// Older xterm shift codes
		"\x1b[a": keys1(keys.ShiftUp),
		"\x1b[b": keys1(keys.ShiftDown),
		"\x1b[c": keys1(keys.ShiftRight),
		"\x1b[d": keys1(keys.ShiftLeft),

		// --- Ctrl + arrows ---
		"\x1b[1;5A": keys1(keys.ControlUp),
		"\x1b[1;5B": keys1(keys.ControlDown),
		"\x1b[1;5C": keys1(keys.ControlRight),
		"\x1b[1;5D": keys1(keys.ControlLeft),
		"\x1b[1;5H": keys1(keys.ControlHome),
		"\x1b[1;5F": keys1(keys.ControlEnd),

		// --- Ctrl+Shift + arrows ---
		"\x1b[1;6A": keys1(keys.ControlShiftUp),
		"\x1b[1;6B": keys1(keys.ControlShiftDown),
		"\x1b[1;6C": keys1(keys.ControlShiftRight),
		"\x1b[1;6D": keys1(keys.ControlShiftLeft),
		"\x1b[1;6H": keys1(keys.ControlShiftHome),
		"\x1b[1;6F": keys1(keys.ControlShiftEnd),

		// --- Shift + F-keys ---
		"\x1b[1;2P": keys1(keys.F13),
		"\x1b[1;2Q": keys1(keys.F14),
		"\x1b[1;2R": keys1(keys.F15),
		"\x1b[1;2S": keys1(keys.F16),

		// --- Shift + navigation ---
		"\x1b[2;2~": keys1(keys.ShiftInsert),
		"\x1b[3;2~": keys1(keys.ShiftDelete),
		"\x1b[5;2~": keys1(keys.ShiftPageUp),
		"\x1b[6;2~": keys1(keys.ShiftPageDown),

		// --- Ctrl + navigation ---
		"\x1b[2;5~": keys1(keys.ControlInsert),
		"\x1b[3;5~": keys1(keys.ControlDelete),
		"\x1b[5;5~": keys1(keys.ControlPageUp),
		"\x1b[6;5~": keys1(keys.ControlPageDown),

		// --- Keypad characters (application mode) ---
		"\x1bOj": char("*"),
		"\x1bOk": char("+"),
		"\x1bOl": char(","),
		"\x1bOm": char("-"),
		"\x1bOn": char("."),
		"\x1bOo": char("/"),
		"\x1bOp": char("0"),
		"\x1bOq": char("1"),
		"\x1bOr": char("2"),
		"\x1bOs": char("3"),
		"\x1bOt": char("4"),
		"\x1bOu": char("5"),
		"\x1bOv": char("6"),
		"\x1bOw": char("7"),
		"\x1bOx": char("8"),
		"\x1bOy": char("9"),
		"\x1bOX": char("="),
		"\x1bOM": keys1(keys.Enter), // keypad enter

		// --- Back tab ---
		"\x1b[Z": keys1(keys.BackTab),

		// --- Shift+Escape ---
		"\x1b\x1b": keys1(keys.ShiftEscape),

		// --- Sequences to ignore (terminal responses, etc.) ---
		"\x1b[?1;0c": IgnoreSequence,
		"\x1b[?1;2c": IgnoreSequence,
		"\x1b[?6c":   IgnoreSequence,
		"\x1b[0c":    IgnoreSequence,
		"\x1b[200~":  IgnoreSequence, // bracketed paste start (handled by parser)
		"\x1b[201~":  IgnoreSequence, // bracketed paste end (handled by parser)
		"\x1b[I":     IgnoreSequence, // focus in (handled by parser)
		"\x1b[O":     IgnoreSequence, // focus out (handled by parser)
	}
}

func initFunctionalKeys() {
	FunctionalKeys = map[string]string{
		// --- VT/CSI final characters without number ---
		"A": keys.Up,
		"B": keys.Down,
		"C": keys.Right,
		"D": keys.Left,
		"E": keys.Up, // keypad center (some terminals)
		"F": keys.End,
		"H": keys.Home,
		"P": keys.F1,
		"Q": keys.F2,
		"R": keys.F3,
		"S": keys.F4,

		// --- CSI sequences with number 1 (same key, explicit parameter) ---
		"1A": keys.Up,
		"1B": keys.Down,
		"1C": keys.Right,
		"1D": keys.Left,
		"1F": keys.End,
		"1H": keys.Home,
		"1P": keys.F1,
		"1Q": keys.F2,
		"1R": keys.F3,
		"1S": keys.F4,

		// --- Tilde-terminated navigation ---
		"2~":  keys.Insert,
		"3~":  keys.Delete,
		"5~":  keys.PageUp,
		"6~":  keys.PageDown,
		"7~":  keys.Home,
		"8~":  keys.End,
		"11~": keys.F1,
		"12~": keys.F2,
		"13~": keys.F3,
		"14~": keys.F4,
		"15~": keys.F5,
		"17~": keys.F6,
		"18~": keys.F7,
		"19~": keys.F8,
		"20~": keys.F9,
		"21~": keys.F10,
		"23~": keys.F11,
		"24~": keys.F12,
		"25~": keys.F13,
		"26~": keys.F14,
		"28~": keys.F15,
		"29~": keys.F16,
		"31~": keys.F17,
		"32~": keys.F18,
		"33~": keys.F19,
		"34~": keys.F20,

		// --- Kitty protocol: unicode codepoint + 'u' ---
		"9u":  keys.Tab,
		"13u": keys.Enter,
		"27u": keys.Escape,
		"32u": keys.Space,

		// Backspace variants
		"8u":   keys.Backspace,
		"127u": keys.Backspace,

		// Function keys via Kitty protocol (private-use area)
		"57344u": keys.F1,
		"57345u": keys.F2,
		"57346u": keys.F3,
		"57347u": keys.F4,
		"57348u": keys.F5,
		"57349u": keys.F6,
		"57350u": keys.F7,
		"57351u": keys.F8,
		"57352u": keys.F9,
		"57353u": keys.F10,
		"57354u": keys.F11,
		"57355u": keys.F12,
		"57356u": keys.F13,
		"57357u": keys.F14,
		"57358u": keys.F15,
		"57359u": keys.F16,
		"57360u": keys.F17,
		"57361u": keys.F18,
		"57362u": keys.F19,
		"57363u": keys.F20,

		// Navigation via Kitty
		"57399u": keys.PageUp,
		"57400u": keys.PageDown,
		"57401u": keys.Home,
		"57402u": keys.End,
		"57403u": keys.Insert,
		"57404u": keys.Delete,

		// Arrow keys via Kitty (codepoints)
		"57419u": keys.Up,
		"57420u": keys.Down,
		"57421u": keys.Right,
		"57422u": keys.Left,
	}
}
