package css

import "testing"

// buildStyles parses a simple CSS declaration string and returns the Styles.
func buildStyles(t *testing.T, decl string) *Styles {
	t.Helper()
	styles, err := ParseDeclarations(decl, CSSLocation{})
	if err != nil {
		t.Fatalf("ParseDeclarations(%q) error: %v", decl, err)
	}
	return styles
}

func TestBuilderDisplay(t *testing.T) {
	s := buildStyles(t, "display: block;")
	v, ok := s.GetRule("display")
	if !ok {
		t.Fatal("display rule not set")
	}
	if v.(Display) != "block" {
		t.Errorf("display = %v, want block", v)
	}
}

func TestBuilderVisibility(t *testing.T) {
	s := buildStyles(t, "visibility: hidden;")
	v, ok := s.GetRule("visibility")
	if !ok {
		t.Fatal("visibility rule not set")
	}
	if v.(Visibility) != "hidden" {
		t.Errorf("visibility = %v, want hidden", v)
	}
}

func TestBuilderColor(t *testing.T) {
	s := buildStyles(t, "color: red;")
	if !s.HasRule("color") {
		t.Fatal("color rule not set")
	}
}

func TestBuilderBackground(t *testing.T) {
	s := buildStyles(t, "background: blue;")
	if !s.HasRule("background") {
		t.Fatal("background rule not set")
	}
}

func TestBuilderWidth(t *testing.T) {
	s := buildStyles(t, "width: 100;")
	if !s.HasRule("width") {
		t.Fatal("width rule not set")
	}
}

func TestBuilderHeight(t *testing.T) {
	s := buildStyles(t, "height: 50;")
	if !s.HasRule("height") {
		t.Fatal("height rule not set")
	}
}

func TestBuilderWidthPercent(t *testing.T) {
	s := buildStyles(t, "width: 50%;")
	if !s.HasRule("width") {
		t.Fatal("width rule not set")
	}
	v := s.Rules["width"].(Scalar)
	if v.Unit != UnitPercent {
		t.Errorf("width unit = %v, want UnitPercent", v.Unit)
	}
}

func TestBuilderPadding(t *testing.T) {
	s := buildStyles(t, "padding: 2;")
	if !s.HasRule("padding") {
		t.Fatal("padding rule not set")
	}
}

func TestBuilderPaddingFourValues(t *testing.T) {
	s := buildStyles(t, "padding: 1 2 3 4;")
	if !s.HasRule("padding") {
		t.Fatal("padding rule not set")
	}
}

func TestBuilderMargin(t *testing.T) {
	s := buildStyles(t, "margin: 1;")
	if !s.HasRule("margin") {
		t.Fatal("margin rule not set")
	}
}

func TestBuilderDock(t *testing.T) {
	s := buildStyles(t, "dock: left;")
	if !s.HasRule("dock") {
		t.Fatal("dock rule not set")
	}
}

func TestBuilderOverflow(t *testing.T) {
	s := buildStyles(t, "overflow-x: scroll;")
	if !s.HasRule("overflow_x") {
		t.Fatal("overflow_x rule not set")
	}
}

func TestBuilderTextAlign(t *testing.T) {
	s := buildStyles(t, "text-align: center;")
	if !s.HasRule("text_align") {
		t.Fatal("text_align rule not set")
	}
}

func TestBuilderOpacity(t *testing.T) {
	s := buildStyles(t, "opacity: 0.5;")
	if !s.HasRule("opacity") {
		t.Fatal("opacity rule not set")
	}
}

func TestBuilderTransition(t *testing.T) {
	s := buildStyles(t, "transition: color 0.5s linear;")
	if !s.HasRule("transitions") {
		t.Fatal("transitions rule not set")
	}
}

func TestBuilderUnknownProperty(t *testing.T) {
	_, err := ParseDeclarations("unknown-property: value;", CSSLocation{})
	if err == nil {
		t.Error("expected error for unknown CSS property")
	}
}

func TestBuilderInitialKeyword(t *testing.T) {
	s := buildStyles(t, "color: initial;")
	if !s.HasRule("color") {
		t.Fatal("expected color rule to be set to nil via 'initial'")
	}
	v, _ := s.GetRule("color")
	if v != nil {
		t.Errorf("color should be nil after 'initial', got %v", v)
	}
}

func TestBuilderBorderTop(t *testing.T) {
	s := buildStyles(t, "border-top: solid red;")
	if !s.HasRule("border_top") {
		t.Fatal("border_top rule not set")
	}
}

func TestBuilderLayer(t *testing.T) {
	s := buildStyles(t, "layer: base;")
	if !s.HasRule("layer") {
		t.Fatal("layer rule not set")
	}
}

func TestBuilderLayers(t *testing.T) {
	s := buildStyles(t, "layers: base above;")
	if !s.HasRule("layers") {
		t.Fatal("layers rule not set")
	}
}

func TestBuilderScrollbarGutter(t *testing.T) {
	s := buildStyles(t, "scrollbar-gutter: stable;")
	if !s.HasRule("scrollbar_gutter") {
		t.Fatal("scrollbar_gutter rule not set")
	}
}

func TestBuilderAlignHorizontal(t *testing.T) {
	s := buildStyles(t, "align-horizontal: center;")
	if !s.HasRule("align_horizontal") {
		t.Fatal("align_horizontal rule not set")
	}
}

func TestBuilderAlignVertical(t *testing.T) {
	s := buildStyles(t, "align-vertical: middle;")
	if !s.HasRule("align_vertical") {
		t.Fatal("align_vertical rule not set")
	}
}

func TestBuilderTextWrap(t *testing.T) {
	s := buildStyles(t, "text-wrap: wrap;")
	if !s.HasRule("text_wrap") {
		t.Fatal("text_wrap rule not set")
	}
}

func TestBuilderExpand(t *testing.T) {
	s := buildStyles(t, "expand: greedy;")
	if !s.HasRule("expand") {
		t.Fatal("expand rule not set")
	}
}

func TestBuilderTextStyle(t *testing.T) {
	s := buildStyles(t, "text-style: bold;")
	if !s.HasRule("text_style") {
		t.Fatal("text_style rule not set")
	}
	if s.TextStyle() != "bold" {
		t.Errorf("TextStyle() = %q, want %q", s.TextStyle(), "bold")
	}
}

func TestBuilderLinkStyle(t *testing.T) {
	s := buildStyles(t, "link-style: underline;")
	if !s.HasRule("link_style") {
		t.Fatal("link_style rule not set")
	}
	if s.LinkStyle() != "underline" {
		t.Errorf("LinkStyle() = %q, want %q", s.LinkStyle(), "underline")
	}
}

func TestBuilderLinkStyleHover(t *testing.T) {
	s := buildStyles(t, "link-style-hover: bold underline;")
	if !s.HasRule("link_style_hover") {
		t.Fatal("link_style_hover rule not set")
	}
	if s.LinkStyleHover() != "bold underline" {
		t.Errorf("LinkStyleHover() = %q, want %q", s.LinkStyleHover(), "bold underline")
	}
}

func TestBuilderBorderTitleStyle(t *testing.T) {
	s := buildStyles(t, "border-title-style: bold;")
	if !s.HasRule("border_title_style") {
		t.Fatal("border_title_style rule not set")
	}
	if s.BorderTitleStyle() != "bold" {
		t.Errorf("BorderTitleStyle() = %q, want %q", s.BorderTitleStyle(), "bold")
	}
}

func TestBuilderBorderSubtitleStyle(t *testing.T) {
	s := buildStyles(t, "border-subtitle-style: italic;")
	if !s.HasRule("border_subtitle_style") {
		t.Fatal("border_subtitle_style rule not set")
	}
	if s.BorderSubtitleStyle() != "italic" {
		t.Errorf("BorderSubtitleStyle() = %q, want %q", s.BorderSubtitleStyle(), "italic")
	}
}
