package xterm

import "github.com/eberle1080/go-textual/msg"

// TerminalSupportsSynchronizedOutput is posted by the xterm parser when the
// terminal reports support for synchronized output (DECRQM response).
type TerminalSupportsSynchronizedOutput struct{ msg.BaseMsg }

// InBandWindowResize is posted when the terminal reports (or enables/disables)
// the in-band window resize extension.
type InBandWindowResize struct {
	msg.BaseMsg
	// Supported indicates whether the terminal supports the feature.
	Supported bool
	// Enabled indicates whether the feature is currently enabled.
	Enabled bool
}

// FromSettingParameter interprets the terminal's DECRQM setting parameter:
//
//	0 = not supported, 1 = not enabled, 2 = enabled
func (m *InBandWindowResize) FromSettingParameter(param int) {
	switch param {
	case 0:
		m.Supported = false
		m.Enabled = false
	case 1:
		m.Supported = true
		m.Enabled = false
	case 2:
		m.Supported = true
		m.Enabled = true
	}
}

// CursorPosition is posted when the terminal responds to a cursor position
// report request (CPR).
type CursorPosition struct {
	msg.BaseMsg
	// X is the zero-based column position.
	X int
	// Y is the zero-based row position.
	Y int
}
