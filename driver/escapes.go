package driver

// Shared ANSI/VT escape sequences used by both Unix and Windows drivers.
const (
	altScreenEnter  = "\x1b[?1049h"
	altScreenExit   = "\x1b[?1049l"
	hideCursor      = "\x1b[?25l"
	showCursor      = "\x1b[?25h"
	enableFocus     = "\x1b[?1004h"
	disableFocus    = "\x1b[?1004l"
	kittyEnable     = "\x1b[>1u"
	kittyDisable    = "\x1b[<u"
	disableLineWrap = "\x1b[?7l"
	enableLineWrap  = "\x1b[?7h"

	enableBracketedPaste  = "\x1b[?2004h"
	disableBracketedPaste = "\x1b[?2004l"

	// Mouse: VT200 + any-event + highlight + SGR extension.
	enableMouse  = "\x1b[?1000h\x1b[?1003h\x1b[?1004h\x1b[?1006h"
	disableMouse = "\x1b[?1006l\x1b[?1004l\x1b[?1003l\x1b[?1000l"

	// Pixel mouse coordinates (urxvt extension).
	enableMousePixels  = "\x1b[?1016h"
	disableMousePixels = "\x1b[?1016l"

	// Synchronized output (request mode + enable/disable).
	querySynchronizedOutput  = "\x1b[?2026$p"
	enableSynchronizedOutput = "\x1b[?2026h"

	// In-band window resize (request mode).
	queryInBandResize   = "\x1b[?2048$p"
	enableInBandResize  = "\x1b[?2048h"
	disableInBandResize = "\x1b[?2048l"
)
