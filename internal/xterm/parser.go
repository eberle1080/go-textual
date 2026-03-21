// Package xterm provides a state-machine based VT/xterm escape sequence
// parser. It replaces the coroutine-based generator pattern used in the
// Python Textual library with explicit states and a Feed/Tick API.
package xterm

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/eberle1080/go-textual/internal/ansi"
	"github.com/eberle1080/go-textual/keys"
	"github.com/eberle1080/go-textual/msg"
)

// parserState is the current state of the escape sequence state machine.
type parserState int

const (
	stateNormal parserState = iota
	stateEscape
	statePaste
)

// maxSequenceLength is the maximum number of bytes buffered before an
// unrecognised escape sequence is reissued as individual key events.
const maxSequenceLength = 32

// escapDelay is the maximum time the parser waits after receiving an isolated
// ESC before emitting it as the escape key.
const escapeDelay = 100 * time.Millisecond

// bracketed paste control sequences.
const (
	bracketedPasteStart = "\x1b[200~"
	bracketedPasteEnd   = "\x1b[201~"
	focusIn             = "\x1b[I"
	focusOut            = "\x1b[O"
)

// Pre-compiled regular expressions for sequence matching.
var (
	// In-band window resize: ESC [ 48 ; <rows> ; <cols> ; <px-h> ; <px-w> t
	reInBandResize = regexp.MustCompile(`^\x1b\[48;(\d+);(\d+);(\d+);(\d+)t$`)

	// Cursor position report: ESC [ row ; col R
	reCPR = regexp.MustCompile(`^\x1b\[(\d+);(\d+)R$`)

	// Extended key protocol: ESC [ [num[;mod]] final
	// final ∈ {u ~ A B C D E F H P Q R S}
	reExtKey = regexp.MustCompile(`^\x1b\[(?:(\d+)(?:;(\d+))?)?([u~ABCDEFHPQRS])$`)

	// SGR mouse: ESC [ < btn ; x ; y {M|m}
	reSGRMouse = regexp.MustCompile(`^\x1b\[<(\d+);(-?\d+);(-?\d+)([mM])$`)

	// Terminal mode response (DECRQM): ESC [ ? <mode> ; <setting> $ y
	reDecrqm = regexp.MustCompile(`^\x1b\[\?(\d+);(\d)\$y$`)
)

// Parser processes raw terminal input and converts it into a sequence of
// messages (key events, mouse events, paste events, etc.).
//
// Parser is not safe for concurrent use; call Feed and Tick from a single
// goroutine.
type Parser struct {
	lastX, lastY      float64
	mousePixels       bool
	terminalSize      *[2]int // (width, height) in cells
	terminalPixelSize *[2]int // (width, height) in pixels

	debug       bool
	debugLog    *os.File
	state       parserState
	seqBuf      string    // accumulated escape sequence (includes leading ESC)
	pasteBuf    []string  // bracketed paste accumulator
	lastEscTime time.Time // when we entered stateEscape
	escTimeout  time.Duration
}

// New constructs a Parser. When debug is true diagnostic output is written to
// os.Stderr.
func New(debug bool) *Parser {
	return &Parser{
		debug:      debug,
		escTimeout: escapeDelay,
		state:      stateNormal,
	}
}

// SetTerminalSize provides the current terminal cell size to the parser,
// enabling pixel-to-cell mouse coordinate conversion.
func (p *Parser) SetTerminalSize(w, h int) {
	p.terminalSize = &[2]int{w, h}
}

// SetTerminalPixelSize provides the current terminal pixel size to the parser.
func (p *Parser) SetTerminalPixelSize(w, h int) {
	p.terminalPixelSize = &[2]int{w, h}
}

// SetMousePixels enables or disables pixel-mode mouse coordinate conversion.
func (p *Parser) SetMousePixels(enabled bool) {
	p.mousePixels = enabled
}

// Feed processes the raw bytes read from the terminal and returns any messages
// that could be fully decoded.
func (p *Parser) Feed(data string) []msg.Msg {
	var msgs []msg.Msg
	for _, ch := range data {
		char := string(ch)
		msgs = append(msgs, p.processChar(char)...)
	}
	return msgs
}

// Tick is called periodically to flush any pending escape sequence that has timed out.
func (p *Parser) Tick() []msg.Msg {
	if p.state != stateEscape {
		return nil
	}
	if time.Since(p.lastEscTime) < p.escTimeout {
		return nil
	}
	seq := p.seqBuf
	p.seqBuf = ""
	p.state = stateNormal
	if seq == "\x1b" {
		return []msg.Msg{msg.NewKey(keys.Escape, nil)}
	}
	return p.reissueSequenceAsKeys(seq, true)
}

// processChar advances the state machine by one character.
func (p *Parser) processChar(ch string) []msg.Msg {
	switch p.state {
	case stateNormal:
		return p.handleNormal(ch)
	case stateEscape:
		return p.handleEscape(ch)
	case statePaste:
		return p.handlePaste(ch)
	}
	return nil
}

func (p *Parser) handleNormal(ch string) []msg.Msg {
	if ch == "\x1b" {
		p.state = stateEscape
		p.seqBuf = ch
		p.lastEscTime = time.Now()
		return nil
	}
	evts := p.sequenceToKeyMsgs(ch, false)
	msgs := make([]msg.Msg, len(evts))
	for i, e := range evts {
		msgs[i] = e
	}
	return msgs
}

func (p *Parser) handleEscape(ch string) []msg.Msg {
	if ch == "\x1b" {
		var msgs []msg.Msg
		if p.seqBuf != "" && p.seqBuf != "\x1b" {
			msgs = p.reissueSequenceAsKeys(p.seqBuf, false)
		} else if p.seqBuf == "\x1b" {
			msgs = append(msgs, msg.NewKey(keys.Escape, nil))
		}
		p.seqBuf = "\x1b"
		p.lastEscTime = time.Now()
		return msgs
	}

	p.seqBuf += ch

	switch p.seqBuf {
	case bracketedPasteStart:
		p.state = statePaste
		p.pasteBuf = p.pasteBuf[:0]
		p.seqBuf = ""
		return nil
	case focusIn:
		p.state = stateNormal
		p.seqBuf = ""
		return []msg.Msg{msg.AppFocusMsg{}}
	case focusOut:
		p.state = stateNormal
		p.seqBuf = ""
		return []msg.Msg{msg.AppBlurMsg{}}
	}

	if result, ok := ansi.Sequences[p.seqBuf]; ok {
		p.state = stateNormal
		seq := p.seqBuf
		p.seqBuf = ""
		if ansi.IsIgnored(result) {
			return nil
		}
		return p.sequenceResultToMsgs(result, seq)
	}

	if msgs := p.tryRegexMatchers(); msgs != nil {
		return msgs
	}

	if len(ch) == 1 && ch[0] >= 0x40 && ch[0] <= 0x7e && len(p.seqBuf) > 2 {
		msgs := p.reissueSequenceAsKeys(p.seqBuf, true)
		p.seqBuf = ""
		p.state = stateNormal
		return msgs
	}

	if len(p.seqBuf) > maxSequenceLength {
		msgs := p.reissueSequenceAsKeys(p.seqBuf, true)
		p.seqBuf = ""
		p.state = stateNormal
		return msgs
	}

	return nil
}

func (p *Parser) handlePaste(ch string) []msg.Msg {
	p.pasteBuf = append(p.pasteBuf, ch)
	combined := strings.Join(p.pasteBuf, "")
	if strings.HasSuffix(combined, bracketedPasteEnd) {
		text := strings.TrimSuffix(combined, bracketedPasteEnd)
		p.pasteBuf = p.pasteBuf[:0]
		p.state = stateNormal
		return []msg.Msg{msg.PasteMsg{Text: text}}
	}
	return nil
}

func (p *Parser) tryRegexMatchers() []msg.Msg {
	seq := p.seqBuf

	if m := reInBandResize.FindStringSubmatch(seq); m != nil {
		rows, _ := strconv.Atoi(m[1])
		cols, _ := strconv.Atoi(m[2])
		p.state = stateNormal
		p.seqBuf = ""
		resize := msg.FromDimensions([2]int{cols, rows}, nil)
		ibr := &InBandWindowResize{Supported: true, Enabled: true}
		return []msg.Msg{resize, ibr}
	}

	if m := reCPR.FindStringSubmatch(seq); m != nil {
		row, _ := strconv.Atoi(m[1])
		col, _ := strconv.Atoi(m[2])
		p.state = stateNormal
		p.seqBuf = ""
		return []msg.Msg{&CursorPosition{X: col - 1, Y: row - 1}}
	}

	if m := reSGRMouse.FindStringSubmatch(seq); m != nil {
		p.state = stateNormal
		p.seqBuf = ""
		if m := p.parseSGRMouse(seq); m != nil {
			return []msg.Msg{m}
		}
		return nil
	}

	if m := reDecrqm.FindStringSubmatch(seq); m != nil {
		mode, _ := strconv.Atoi(m[1])
		setting, _ := strconv.Atoi(m[2])
		p.state = stateNormal
		p.seqBuf = ""
		switch mode {
		case 2026:
			if setting >= 1 {
				return []msg.Msg{&TerminalSupportsSynchronizedOutput{}}
			}
			return nil
		case 2048:
			ibr := &InBandWindowResize{}
			ibr.FromSettingParameter(setting)
			return []msg.Msg{ibr}
		}
		return nil
	}

	if m := reExtKey.FindStringSubmatch(seq); m != nil {
		numStr, modStr, finalChar := m[1], m[2], m[3]
		p.state = stateNormal
		p.seqBuf = ""
		return p.parseExtKey(numStr, modStr, finalChar)
	}

	return nil
}

func (p *Parser) parseExtKey(numStr, modStr, finalChar string) []msg.Msg {
	token := numStr + finalChar

	modifier := 1
	if modStr != "" {
		modifier, _ = strconv.Atoi(modStr)
	}
	modBits := modifier - 1
	shift := modBits&1 != 0
	alt := modBits&2 != 0
	ctrl := modBits&4 != 0

	keyName := ""
	if k, ok := ansi.FunctionalKeys[token]; ok {
		keyName = k
	} else if k, ok := ansi.FunctionalKeys[finalChar]; ok {
		keyName = k
	} else if finalChar == "u" && numStr != "" {
		// Kitty protocol: codepoint-based key (e.g. \x1b[99;5u = ctrl+c).
		// FunctionalKeys only covers special keys; derive printable keys from
		// the Unicode codepoint directly.
		if cp, err := strconv.Atoi(numStr); err == nil {
			r := rune(cp)
			if unicode.IsPrint(r) {
				keyName = keys.CharacterToKey(string(r))
			}
		}
	}
	if keyName == "" {
		return p.reissueSequenceAsKeys(p.seqBuf, true)
	}

	keyName = applyModifiers(keyName, shift, alt, ctrl)
	if keyName == keys.Ignore {
		return nil
	}
	return []msg.Msg{msg.NewKey(keyName, nil)}
}

func applyModifiers(keyName string, shift, alt, ctrl bool) string {
	if shift && len([]rune(keyName)) == 1 {
		upper := strings.ToUpper(keyName)
		if upper != keyName {
			keyName = upper
			shift = false
		}
	}
	if ctrl {
		keyName = "ctrl+" + keyName
	}
	if shift {
		keyName = "shift+" + keyName
	}
	if alt {
		keyName = "escape+" + keyName
	}
	if rep, ok := keys.KeyNameReplacements[keyName]; ok {
		keyName = rep
	}
	return keyName
}

func (p *Parser) parseSGRMouse(seq string) msg.Msg {
	m := reSGRMouse.FindStringSubmatch(seq)
	if m == nil {
		return nil
	}
	btnCode, _ := strconv.Atoi(m[1])
	x, _ := strconv.Atoi(m[2])
	y, _ := strconv.Atoi(m[3])
	state := m[4]

	x--
	y--

	fx, fy := float64(x), float64(y)
	if p.mousePixels && p.terminalSize != nil && p.terminalPixelSize != nil {
		cellW := float64(p.terminalPixelSize[0]) / float64(p.terminalSize[0])
		cellH := float64(p.terminalPixelSize[1]) / float64(p.terminalSize[1])
		if cellW > 0 {
			fx = fx / cellW
		}
		if cellH > 0 {
			fy = fy / cellH
		}
	}

	deltaX := int(fx - p.lastX)
	deltaY := int(fy - p.lastY)

	shift := btnCode&4 != 0
	alt := btnCode&8 != 0
	ctrl := btnCode&16 != 0
	button := msg.MouseButton(btnCode & 3)

	var result msg.Msg

	switch {
	case btnCode&64 != 0:
		switch btnCode & 3 {
		case 0:
			result = msg.NewMouseScrollUp(nil, fx, fy, fx, fy, deltaX, deltaY, button, shift, alt, ctrl)
		case 1:
			result = msg.NewMouseScrollDown(nil, fx, fy, fx, fy, deltaX, deltaY, button, shift, alt, ctrl)
		case 2:
			result = msg.NewMouseScrollLeft(nil, fx, fy, fx, fy, deltaX, deltaY, button, shift, alt, ctrl)
		case 3:
			result = msg.NewMouseScrollRight(nil, fx, fy, fx, fy, deltaX, deltaY, button, shift, alt, ctrl)
		default:
			result = msg.NewMouseScrollDown(nil, fx, fy, fx, fy, deltaX, deltaY, button, shift, alt, ctrl)
		}
	case btnCode&32 != 0 || btnCode&3 == 3:
		b := button
		if btnCode&3 == 3 {
			b = 0
		}
		result = msg.NewMouseMove(nil, fx, fy, fx, fy, deltaX, deltaY, b, shift, alt, ctrl)
	default:
		if state == "M" {
			result = msg.NewMouseDown(nil, fx, fy, fx, fy, deltaX, deltaY, button, shift, alt, ctrl)
		} else {
			result = msg.NewMouseUp(nil, fx, fy, fx, fy, deltaX, deltaY, button, shift, alt, ctrl)
		}
	}

	p.lastX = fx
	p.lastY = fy
	return result
}

func (p *Parser) sequenceToKeyMsgs(sequence string, alt bool) []msg.KeyMsg {
	if result, ok := ansi.Sequences[sequence]; ok {
		if ansi.IsIgnored(result) {
			return nil
		}
		var msgs []msg.KeyMsg
		for _, k := range result.Keys {
			if k != keys.Ignore {
				msgs = append(msgs, msg.NewKey(k, nil))
			}
		}
		if result.Character != "" {
			msgs = append(msgs, msg.NewKey(result.Character, &result.Character))
		}
		return msgs
	}

	if utf8.RuneCountInString(sequence) == 1 {
		keyName := keys.CharacterToKey(sequence)
		if rep, ok := keys.KeyNameReplacements[keyName]; ok {
			keyName = rep
		}
		if alt {
			keyName = "escape+" + keyName
		}
		if keyName == keys.Ignore {
			return nil
		}
		return []msg.KeyMsg{msg.NewKey(keyName, nil)}
	}

	return nil
}

func (p *Parser) sequenceResultToMsgs(result ansi.SequenceResult, _ string) []msg.Msg {
	if ansi.IsIgnored(result) {
		return nil
	}
	var msgs []msg.Msg
	for _, k := range result.Keys {
		if k != keys.Ignore {
			msgs = append(msgs, msg.NewKey(k, nil))
		}
	}
	if result.Character != "" {
		ch := result.Character
		msgs = append(msgs, msg.NewKey(ch, &ch))
	}
	return msgs
}

func (p *Parser) reissueSequenceAsKeys(sequence string, processAlt bool) []msg.Msg {
	if sequence == "" {
		return nil
	}
	var msgs []msg.Msg
	if processAlt && strings.HasPrefix(sequence, "\x1b") {
		msgs = append(msgs, msg.NewKey(keys.Escape, nil))
		sequence = sequence[1:]
	}
	for _, ch := range sequence {
		char := string(ch)
		evts := p.sequenceToKeyMsgs(char, false)
		for _, e := range evts {
			msgs = append(msgs, e)
		}
	}
	return msgs
}

func (p *Parser) logDebug(format string, args ...any) {
	if !p.debug {
		return
	}
	_, _ = fmt.Fprintf(os.Stderr, "[xterm] "+format+"\n", args...)
}
