package msg

import "github.com/eberle1080/go-textual/geometry"

// MountMsg is sent when a widget is mounted into the tree.
type MountMsg struct{ BaseMsg }

// UnmountMsg is sent when a widget is removed from the tree.
type UnmountMsg struct{ BaseMsg }

// ResizeMsg is sent when the terminal or a widget is resized.
type ResizeMsg struct {
	BaseMsg
	// Size is the new size in cells.
	Size geometry.Size
	// VirtualSize is the virtual (scrollable) size.
	VirtualSize geometry.Size
	// PixelSize is the pixel dimensions, if known.
	PixelSize *geometry.Size
}

// FromDimensions constructs a ResizeMsg from raw cell and optional pixel dimensions.
func FromDimensions(cells [2]int, pixels *[2]int) ResizeMsg {
	size := geometry.Size{Width: cells[0], Height: cells[1]}
	var pixelSize *geometry.Size
	if pixels != nil {
		ps := geometry.Size{Width: pixels[0], Height: pixels[1]}
		pixelSize = &ps
	}
	return ResizeMsg{Size: size, VirtualSize: size, PixelSize: pixelSize}
}

// FocusMsg is sent when a widget gains keyboard focus.
type FocusMsg struct {
	BaseMsg
	// FromAppFocus is true when focus is gained because the terminal window
	// itself regained focus.
	FromAppFocus bool
}

// BlurMsg is sent when a widget loses keyboard focus.
type BlurMsg struct{ BaseMsg }

// AppFocusMsg is sent when the terminal window gains focus.
type AppFocusMsg struct{ BaseMsg }

// AppBlurMsg is sent when the terminal window loses focus.
type AppBlurMsg struct{ BaseMsg }

// PasteMsg is sent when text is pasted into the terminal.
type PasteMsg struct {
	BaseMsg
	Text string
}
