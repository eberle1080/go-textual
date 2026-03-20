//go:build windows

package driver

import (
	"os"
	"unicode"
	"unsafe"

	"golang.org/x/sys/windows"

	"github.com/eberle1080/go-textual/keys"
	"github.com/eberle1080/go-textual/msg"
)

// Windows console input record event types.
const (
	winKeyEvent              uint16 = 0x0001
	winMouseEvent            uint16 = 0x0002
	winWindowBufferSizeEvent uint16 = 0x0004
	winFocusEvent            uint16 = 0x0010
)

// Mouse event flags.
const (
	winMouseMoved    = uint32(0x0001)
	winMouseWheeled  = uint32(0x0004)
	winMouseHWheeled = uint32(0x0008)
)

// Mouse button state bits.
const (
	winFromLeft1stButton = uint32(0x0001)
	winRightmostButton   = uint32(0x0002)
	winFromLeft2ndButton = uint32(0x0004)
)

// Control key state bits.
const (
	winRightAltPressed  = uint32(0x0001)
	winLeftAltPressed   = uint32(0x0002)
	winRightCtrlPressed = uint32(0x0004)
	winLeftCtrlPressed  = uint32(0x0008)
	winShiftPressed     = uint32(0x0010)
)

type winCoord struct{ X, Y int16 }

type winInputRecord struct {
	EventType uint16
	_         [2]byte
	Event     [16]byte
}

type winMouseEventRecord struct {
	MousePosition   winCoord
	ButtonState     uint32
	ControlKeyState uint32
	EventFlags      uint32
}

type winWindowBufferSizeRecord struct{ Size winCoord }

type winFocusEventRecord struct{ SetFocus int32 }

type winKeyEventRecord struct {
	KeyDown         int32
	RepeatCount     uint16
	VirtualKeyCode  uint16
	VirtualScanCode uint16
	UnicodeChar     uint16
	ControlKeyState uint32
}

var vkKeyMap = map[uint16]string{
	0x08: keys.Backspace,
	0x09: keys.Tab,
	0x0D: keys.Enter,
	0x1B: keys.Escape,
	0x21: keys.PageUp,
	0x22: keys.PageDown,
	0x23: keys.End,
	0x24: keys.Home,
	0x25: keys.Left,
	0x26: keys.Up,
	0x27: keys.Right,
	0x28: keys.Down,
	0x2D: keys.Insert,
	0x2E: keys.Delete,
	0x70: keys.F1,
	0x71: keys.F2,
	0x72: keys.F3,
	0x73: keys.F4,
	0x74: keys.F5,
	0x75: keys.F6,
	0x76: keys.F7,
	0x77: keys.F8,
	0x78: keys.F9,
	0x79: keys.F10,
	0x7A: keys.F11,
	0x7B: keys.F12,
}

var (
	modKernel32          = windows.NewLazySystemDLL("kernel32.dll")
	procReadConsoleInput = modKernel32.NewProc("ReadConsoleInputW")
)

const (
	winWaitObject0 = uint32(0x00000000)
	winWaitTimeout = uint32(0x00000102)
)

// WindowsDriver is the terminal driver for Windows.
type WindowsDriver struct {
	base             *BaseDriver
	exitCh           chan struct{}
	writer           *Writer
	restoreConsole   func()
	restoreInput     func()
	consoleMouseBtns uint32
}

// NewWindowsDriver constructs a WindowsDriver targeting the given sink.
func NewWindowsDriver(sink EventSink, opts ...DriverOption) *WindowsDriver {
	return &WindowsDriver{
		base:   NewBaseDriver(sink, opts...),
		exitCh: make(chan struct{}),
	}
}

func (d *WindowsDriver) Write(data string) {
	if d.writer != nil {
		d.writer.Write(data)
	}
}

func (d *WindowsDriver) Flush() {}

func (d *WindowsDriver) Close() {
	if d.writer != nil {
		d.writer.Stop()
		d.writer = nil
	}
}

func (d *WindowsDriver) IsHeadless() bool { return false }
func (d *WindowsDriver) IsInline() bool   { return false }
func (d *WindowsDriver) CanSuspend() bool { return false }

func (d *WindowsDriver) SuspendApplicationMode() {}
func (d *WindowsDriver) ResumeApplicationMode()  {}
func (d *WindowsDriver) OpenURL(_ string)         {}

func (d *WindowsDriver) SetCursorOrigin(x, y int) { d.base.SetCursorOrigin(x, y) }
func (d *WindowsDriver) ClearCursorOrigin()        { d.base.ClearCursorOrigin() }

func (d *WindowsDriver) StartApplicationMode() {
	stdoutHandle := windows.Handle(os.Stdout.Fd())
	var originalOutMode uint32
	if err := windows.GetConsoleMode(stdoutHandle, &originalOutMode); err == nil {
		const vtProcessing = windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING
		_ = windows.SetConsoleMode(stdoutHandle, originalOutMode|vtProcessing)
		d.restoreConsole = func() {
			_ = windows.SetConsoleMode(stdoutHandle, originalOutMode)
		}
	}

	stdinHandle := windows.Handle(os.Stdin.Fd())
	var originalInMode uint32
	if err := windows.GetConsoleMode(stdinHandle, &originalInMode); err == nil {
		const (
			enableMouseInput    = 0x0010
			enableWindowInput   = 0x0008
			enableExtendedFlags = 0x0080
		)
		newMode := uint32(enableMouseInput | enableWindowInput | enableExtendedFlags)
		_ = windows.SetConsoleMode(stdinHandle, newMode)
		d.restoreInput = func() {
			_ = windows.SetConsoleMode(stdinHandle, originalInMode)
		}
	}

	d.writer = NewWriter(os.Stdout)
	d.writer.Start()

	d.Write(altScreenEnter)
	if d.base.mouse {
		d.Write(enableMouse)
	}
	d.Write(hideCursor)
	d.Write(enableFocus)
	d.Write(enableBracketedPaste)
	d.Write(disableLineWrap)

	if w, h, err := getWindowsTerminalSize(); err == nil {
		d.base.Send(msg.FromDimensions([2]int{w, h}, nil))
	}

	go d.runConsoleEventThread(stdinHandle)
}

func (d *WindowsDriver) StopApplicationMode() {
	d.Write(disableBracketedPaste)
	d.Write(enableLineWrap)
	d.DisableInput()
	d.Write(altScreenExit)
	d.Write(showCursor)
	d.Write(disableFocus)

	if d.restoreInput != nil {
		d.restoreInput()
		d.restoreInput = nil
	}
	if d.restoreConsole != nil {
		d.restoreConsole()
		d.restoreConsole = nil
	}
}

func (d *WindowsDriver) DisableInput() {
	if d.base.mouse {
		d.Write(disableMouse)
	}
	select {
	case <-d.exitCh:
	default:
		close(d.exitCh)
	}
}

func (d *WindowsDriver) runConsoleEventThread(handle windows.Handle) {
	const waitTimeoutMs = uint32(100)
	buf := make([]winInputRecord, 16)

	for {
		select {
		case <-d.exitCh:
			return
		default:
		}

		ret, _ := windows.WaitForSingleObject(handle, waitTimeoutMs)
		if ret == winWaitTimeout {
			continue
		}
		if ret != winWaitObject0 {
			return
		}

		select {
		case <-d.exitCh:
			return
		default:
		}

		var count uint32
		r1, _, _ := procReadConsoleInput.Call(
			uintptr(handle),
			uintptr(unsafe.Pointer(&buf[0])),
			uintptr(len(buf)),
			uintptr(unsafe.Pointer(&count)),
		)
		if r1 == 0 {
			return
		}

		for i := uint32(0); i < count; i++ {
			d.handleInputRecord(&buf[i])
		}
	}
}

func (d *WindowsDriver) handleInputRecord(rec *winInputRecord) {
	switch rec.EventType {
	case winKeyEvent:
		ker := (*winKeyEventRecord)(unsafe.Pointer(&rec.Event[0]))
		if m := d.translateKeyRecord(ker); m != nil {
			d.base.ProcessMsg(m)
		}
	case winMouseEvent:
		mer := (*winMouseEventRecord)(unsafe.Pointer(&rec.Event[0]))
		d.translateMouseRecord(mer)
	case winWindowBufferSizeEvent:
		wbsr := (*winWindowBufferSizeRecord)(unsafe.Pointer(&rec.Event[0]))
		w, h := int(wbsr.Size.X), int(wbsr.Size.Y)
		d.base.Send(msg.FromDimensions([2]int{w, h}, nil))
	case winFocusEvent:
		fer := (*winFocusEventRecord)(unsafe.Pointer(&rec.Event[0]))
		if fer.SetFocus != 0 {
			d.base.Send(msg.AppFocusMsg{})
		} else {
			d.base.Send(msg.AppBlurMsg{})
		}
	}
}

func (d *WindowsDriver) translateKeyRecord(ker *winKeyEventRecord) msg.Msg {
	if ker.KeyDown == 0 {
		return nil
	}

	ck := ker.ControlKeyState
	altGr := ck&(winRightAltPressed|winLeftCtrlPressed) == (winRightAltPressed | winLeftCtrlPressed)
	ctrl := !altGr && ck&(winRightCtrlPressed|winLeftCtrlPressed) != 0
	alt := !altGr && ck&(winRightAltPressed|winLeftAltPressed) != 0
	shift := ck&winShiftPressed != 0

	vk := ker.VirtualKeyCode

	if baseName, ok := vkKeyMap[vk]; ok {
		if baseName == keys.Tab && shift {
			k := msg.NewKey(keys.BackTab, nil)
			return k
		}
		name := baseName
		if ctrl {
			name = "ctrl+" + name
		}
		if shift {
			name = "shift+" + name
		}
		if alt {
			name = "escape+" + name
		}
		k := msg.NewKey(name, nil)
		return k
	}

	if ctrl && vk >= 0x41 && vk <= 0x5A {
		letter := string(rune(vk + 32))
		name := "ctrl+" + letter
		if shift {
			name = "ctrl+shift+" + letter
		}
		if alt {
			name = "escape+" + name
		}
		k := msg.NewKey(name, nil)
		return k
	}

	if ctrl && vk >= 0x30 && vk <= 0x39 {
		k := msg.NewKey("ctrl+"+string(rune(vk)), nil)
		return k
	}

	uc := ker.UnicodeChar
	if uc == 0 {
		return nil
	}
	r := rune(uc)
	if !unicode.IsPrint(r) {
		return nil
	}
	ch := string(r)
	keyName := keys.CharacterToKey(ch)
	if rep, ok := keys.KeyNameReplacements[keyName]; ok {
		keyName = rep
	}
	if alt {
		keyName = "escape+" + keyName
	}
	k := msg.NewKey(keyName, nil)
	return k
}

func (d *WindowsDriver) translateMouseRecord(mer *winMouseEventRecord) {
	x := float64(mer.MousePosition.X)
	y := float64(mer.MousePosition.Y)

	ck := mer.ControlKeyState
	shift := ck&winShiftPressed != 0
	altGr := ck&(winRightAltPressed|winLeftCtrlPressed) == (winRightAltPressed | winLeftCtrlPressed)
	alt := !altGr && ck&(winRightAltPressed|winLeftAltPressed) != 0
	ctrl := !altGr && ck&(winRightCtrlPressed|winLeftCtrlPressed) != 0

	btns := mer.ButtonState

	switch {
	case mer.EventFlags&winMouseWheeled != 0:
		if int16(btns>>16) > 0 {
			d.base.ProcessMsg(msg.NewMouseScrollUp(nil, x, y, x, y, 0, 0, 0, shift, alt, ctrl))
		} else {
			d.base.ProcessMsg(msg.NewMouseScrollDown(nil, x, y, x, y, 0, 0, 0, shift, alt, ctrl))
		}

	case mer.EventFlags&winMouseHWheeled != 0:
		if int16(btns>>16) > 0 {
			d.base.ProcessMsg(msg.NewMouseScrollRight(nil, x, y, x, y, 0, 0, 0, shift, alt, ctrl))
		} else {
			d.base.ProcessMsg(msg.NewMouseScrollLeft(nil, x, y, x, y, 0, 0, 0, shift, alt, ctrl))
		}

	case mer.EventFlags&winMouseMoved != 0:
		btn := msg.MouseButton(0)
		switch {
		case btns&winFromLeft1stButton != 0:
			btn = 0
		case btns&winFromLeft2ndButton != 0:
			btn = 1
		case btns&winRightmostButton != 0:
			btn = 2
		}
		d.consoleMouseBtns = btns & (winFromLeft1stButton | winFromLeft2ndButton | winRightmostButton)
		d.base.ProcessMsg(msg.NewMouseMove(nil, x, y, x, y, 0, 0, btn, shift, alt, ctrl))

	default:
		prev := d.consoleMouseBtns
		curr := btns & (winFromLeft1stButton | winFromLeft2ndButton | winRightmostButton)
		d.consoleMouseBtns = curr

		type btnInfo struct {
			bit uint32
			idx msg.MouseButton
		}
		btnDefs := [3]btnInfo{
			{winFromLeft1stButton, 0},
			{winFromLeft2ndButton, 1},
			{winRightmostButton, 2},
		}
		for _, bi := range btnDefs {
			wasDown := prev&bi.bit != 0
			isDown := curr&bi.bit != 0
			switch {
			case !wasDown && isDown:
				d.base.ProcessMsg(msg.NewMouseDown(nil, x, y, x, y, 0, 0, bi.idx, shift, alt, ctrl))
			case wasDown && !isDown:
				d.base.ProcessMsg(msg.NewMouseUp(nil, x, y, x, y, 0, 0, bi.idx, shift, alt, ctrl))
			}
		}
	}
}

func getWindowsTerminalSize() (int, int, error) {
	handle := windows.Handle(os.Stdout.Fd())
	var info windows.ConsoleScreenBufferInfo
	if err := windows.GetConsoleScreenBufferInfo(handle, &info); err != nil {
		return 0, 0, err
	}
	w := int(info.Window.Right - info.Window.Left + 1)
	h := int(info.Window.Bottom - info.Window.Top + 1)
	return w, h, nil
}
