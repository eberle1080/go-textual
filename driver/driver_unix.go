//go:build !windows

package driver

import (
	"os"
	"time"
	"unicode/utf8"

	"golang.org/x/sys/unix"
	"golang.org/x/term"

	"github.com/eberle1080/go-textual/internal/xterm"
	"github.com/eberle1080/go-textual/msg"
)

// UnixDriver is the terminal driver for Linux and macOS.
type UnixDriver struct {
	base         *BaseDriver
	ttyFile      *os.File
	fileno       int
	attrsBefore  *term.State
	exitCh       chan struct{}
	writer       *Writer
	stopSIGWINCH func()
	inBandResize bool
	mousePixels  bool
}

// NewUnixDriver constructs a UnixDriver targeting the given sink.
func NewUnixDriver(sink EventSink, opts ...DriverOption) *UnixDriver {
	d := &UnixDriver{
		base:   NewBaseDriver(sink, opts...),
		exitCh: make(chan struct{}),
	}
	if tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0); err == nil {
		d.ttyFile = tty
		d.fileno = int(tty.Fd())
	} else {
		d.fileno = int(os.Stdin.Fd())
	}
	return d
}

func (d *UnixDriver) inputFile() *os.File {
	if d.ttyFile != nil {
		return d.ttyFile
	}
	return os.Stdin
}

func (d *UnixDriver) Write(data string) {
	if d.writer != nil {
		d.writer.Write(data)
	}
}

func (d *UnixDriver) Flush() {}

func (d *UnixDriver) Close() {
	if d.writer != nil {
		d.writer.Stop()
		d.writer = nil
	}
	if d.ttyFile != nil {
		_ = d.ttyFile.Close()
		d.ttyFile = nil
	}
}

func (d *UnixDriver) IsHeadless() bool { return false }
func (d *UnixDriver) IsInline() bool   { return false }
func (d *UnixDriver) CanSuspend() bool { return true }

func (d *UnixDriver) OpenURL(url string) {
	d.Write("\x1b]8;;" + url + "\x1b\\\x1b]8;;\x1b\\")
}

func (d *UnixDriver) SetCursorOrigin(x, y int) { d.base.SetCursorOrigin(x, y) }
func (d *UnixDriver) ClearCursorOrigin()        { d.base.ClearCursorOrigin() }

func (d *UnixDriver) getTerminalSize() (int, int) {
	w, h, err := term.GetSize(d.fileno)
	if err != nil {
		w, h = 80, 24
	}
	return w, h
}

func (d *UnixDriver) StartApplicationMode() {
	d.writer = NewWriter(os.Stdout)
	d.writer.Start()

	d.stopSIGWINCH = installSIGWINCH(d.base.sink)

	d.Write(altScreenEnter)

	if d.base.mouse {
		d.enableMouseSupport()
	}

	state, err := term.MakeRaw(d.fileno)
	if err == nil {
		d.attrsBefore = state
	}

	d.Write(hideCursor)
	d.Write(enableFocus)
	d.Write(kittyEnable)
	d.Write(enableBracketedPaste)
	d.Write(disableLineWrap)

	w, h := d.getTerminalSize()
	d.base.Send(msg.FromDimensions([2]int{w, h}, nil))

	go d.runInputThread()
}

func (d *UnixDriver) StopApplicationMode() {
	d.Write(disableBracketedPaste)
	d.Write(enableLineWrap)
	d.DisableInput()

	if d.attrsBefore != nil {
		_ = term.Restore(d.fileno, d.attrsBefore)
		d.attrsBefore = nil
	}

	if d.inBandResize {
		d.Write(disableInBandResize)
		d.inBandResize = false
	}
	d.Write(kittyDisable)
	d.Write(altScreenExit)
	d.Write(showCursor)
	d.Write(disableFocus)
}

func (d *UnixDriver) SuspendApplicationMode() {
	d.StopApplicationMode()
}

func (d *UnixDriver) ResumeApplicationMode() {
	d.exitCh = make(chan struct{})
	d.StartApplicationMode()
}

func (d *UnixDriver) DisableInput() {
	if d.stopSIGWINCH != nil {
		d.stopSIGWINCH()
		d.stopSIGWINCH = nil
	}
	if d.base.mouse {
		d.disableMouseSupport()
	}
	select {
	case <-d.exitCh:
	default:
		close(d.exitCh)
	}
}

func (d *UnixDriver) enableMouseSupport() {
	d.Write(enableMouse)
}

func (d *UnixDriver) disableMouseSupport() {
	d.Write(disableMouse)
	if d.mousePixels {
		d.Write(disableMousePixels)
	}
}

func (d *UnixDriver) handleParserMsg(m msg.Msg) {
	switch v := m.(type) {
	case *xterm.TerminalSupportsSynchronizedOutput:
		d.Write(enableSynchronizedOutput)
		d.base.Send(v)
		return
	case *xterm.InBandWindowResize:
		if v.Supported && !v.Enabled {
			d.Write(enableInBandResize)
			d.inBandResize = true
			if d.base.mouse && !d.mousePixels {
				d.Write(enableMousePixels)
				d.mousePixels = true
			}
		}
		d.base.Send(v)
		return
	}
	d.base.ProcessMsg(m)
}

func (d *UnixDriver) runInputThread() {
	parser := xterm.New(d.base.debug)
	w, h := d.getTerminalSize()
	parser.SetTerminalSize(w, h)

	input := d.inputFile()
	buf := make([]byte, 4096)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	pollFds := []unix.PollFd{{Fd: int32(input.Fd()), Events: unix.POLLIN}}

	for {
		select {
		case <-d.exitCh:
			n, _ := input.Read(buf)
			if n > 0 {
				str := bytesToUTF8(buf[:n])
				for _, m := range parser.Feed(str) {
					d.handleParserMsg(m)
				}
			}
			return
		case <-ticker.C:
			for _, m := range parser.Tick() {
				d.handleParserMsg(m)
			}
		default:
		}

		n, err := unix.Poll(pollFds, 100)
		if err != nil || n == 0 {
			continue
		}

		n, _ = input.Read(buf)
		if n > 0 {
			str := bytesToUTF8(buf[:n])
			for _, m := range parser.Feed(str) {
				d.handleParserMsg(m)
			}
		}
	}
}

func bytesToUTF8(b []byte) string {
	if utf8.Valid(b) {
		return string(b)
	}
	runes := make([]rune, 0, len(b))
	for len(b) > 0 {
		r, size := utf8.DecodeRune(b)
		runes = append(runes, r)
		b = b[size:]
	}
	return string(runes)
}
