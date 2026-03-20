package widgets_test

import (
	"context"
	"strings"
	"testing"

	rich "github.com/eberle1080/go-rich"

	"github.com/eberle1080/go-textual/geometry"
	"github.com/eberle1080/go-textual/keys"
	"github.com/eberle1080/go-textual/msg"
	"github.com/eberle1080/go-textual/widgets"
)

func TestButton_Render(t *testing.T) {
	b := widgets.NewButton("Click me")
	region := geometry.Region{Width: 20, Height: 1}
	strips := b.Render(region)

	if len(strips) != 1 {
		t.Fatalf("expected 1 strip, got %d", len(strips))
	}
	rendered := strips[0].Render(rich.ColorModeNone)
	if !strings.Contains(rendered, "Click me") {
		t.Errorf("rendered strip %q does not contain label text", rendered)
	}
}

func TestButton_KeyPress(t *testing.T) {
	b := widgets.NewButton("Click me")
	cmd := b.Update(context.Background(), msg.NewKey(keys.Enter, nil))
	if cmd == nil {
		t.Fatal("expected non-nil Cmd from Enter key press")
	}
	result := cmd(context.Background())
	pressed, ok := result.(widgets.ButtonPressedMsg)
	if !ok {
		t.Fatalf("expected ButtonPressedMsg, got %T", result)
	}
	if pressed.Button != b {
		t.Error("ButtonPressedMsg.Button should be the pressed button")
	}
}
