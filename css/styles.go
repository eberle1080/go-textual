package css

import (
	"fmt"
	"strings"

	"github.com/eberle1080/go-textual/color"
	"github.com/eberle1080/go-textual/geometry"
)

// RulesMap stores CSS property values keyed by rule name (underscore format).
type RulesMap map[string]any

// RuleNames is the ordered list of all supported CSS rule names.
var RuleNames = []string{
	"display", "visibility", "layout",
	"auto_color", "color", "background", "text_style",
	"background_tint",
	"opacity", "text_opacity",
	"padding", "margin", "offset", "position",
	"border_top", "border_right", "border_bottom", "border_left",
	"border_title_align", "border_subtitle_align",
	"outline_top", "outline_right", "outline_bottom", "outline_left",
	"keyline",
	"box_sizing", "width", "height", "min_width", "min_height", "max_width", "max_height",
	"dock", "split",
	"overflow_x", "overflow_y",
	"layers", "layer",
	"transitions",
	"tint",
	"scrollbar_color", "scrollbar_color_hover", "scrollbar_color_active",
	"scrollbar_corner_color",
	"scrollbar_background", "scrollbar_background_hover", "scrollbar_background_active",
	"scrollbar_gutter", "scrollbar_size_vertical", "scrollbar_size_horizontal",
	"scrollbar_visibility",
	"align_horizontal", "align_vertical",
	"content_align_horizontal", "content_align_vertical",
	"grid_size_rows", "grid_size_columns",
	"grid_gutter_horizontal", "grid_gutter_vertical",
	"grid_rows", "grid_columns",
	"row_span", "column_span",
	"text_align",
	"link_color", "auto_link_color", "link_background", "link_style",
	"link_color_hover", "auto_link_color_hover", "link_background_hover", "link_style_hover",
	"auto_border_title_color", "border_title_color", "border_title_background", "border_title_style",
	"auto_border_subtitle_color", "border_subtitle_color", "border_subtitle_background", "border_subtitle_style",
	"hatch",
	"overlay", "constrain_x", "constrain_y",
	"text_wrap", "text_overflow", "expand",
	"line_pad",
	"pointer",
}

// RuleNamesSet is the set of all supported CSS rule names.
var RuleNamesSet = func() map[string]bool {
	m := make(map[string]bool, len(RuleNames))
	for _, n := range RuleNames {
		m[n] = true
	}
	return m
}()

// Animatable is the set of CSS rule names that can be animated.
var Animatable = map[string]bool{
	"offset": true, "padding": true, "margin": true,
	"width": true, "height": true,
	"min_width": true, "min_height": true, "max_width": true, "max_height": true,
	"auto_color": true, "color": true, "background": true, "background_tint": true,
	"opacity": true, "position": true, "text_opacity": true, "tint": true,
	"scrollbar_color": true, "scrollbar_color_hover": true, "scrollbar_color_active": true,
	"scrollbar_background": true, "scrollbar_background_hover": true,
	"scrollbar_background_active": true, "scrollbar_visibility": true,
	"link_color": true, "link_background": true,
	"link_color_hover": true, "link_background_hover": true,
	"text_wrap": true, "text_overflow": true, "line_pad": true,
}

// ExtractedRule holds a single rule with its resolved specificity.
type ExtractedRule struct {
	Name        string
	Specificity Specificity6
	Value       any
}

// Styles holds the CSS properties for a single node.
type Styles struct {
	Rules     RulesMap
	Important map[string]bool
}

// NewStyles creates an empty Styles object.
func NewStyles() *Styles {
	return &Styles{
		Rules:     make(RulesMap),
		Important: make(map[string]bool),
	}
}

// HasRule reports whether the named rule is set.
func (s *Styles) HasRule(name string) bool {
	_, ok := s.Rules[name]
	return ok
}

// GetRule returns the value of the named rule and whether it was set.
func (s *Styles) GetRule(name string) (any, bool) {
	v, ok := s.Rules[name]
	return v, ok
}

// SetRule sets the value of the named rule. Returns true if the value changed.
func (s *Styles) SetRule(name string, value any) bool {
	current := s.Rules[name]
	s.Rules[name] = value
	return fmt.Sprintf("%v", current) != fmt.Sprintf("%v", value)
}

// ClearRule removes the named rule. Returns true if a rule was removed.
func (s *Styles) ClearRule(name string) bool {
	_, ok := s.Rules[name]
	if ok {
		delete(s.Rules, name)
	}
	return ok
}

// GetRules returns a copy of the rules map.
func (s *Styles) GetRules() RulesMap {
	copy := make(RulesMap, len(s.Rules))
	for k, v := range s.Rules {
		copy[k] = v
	}
	return copy
}

// Reset clears all rules.
func (s *Styles) Reset() {
	s.Rules = make(RulesMap)
}

// Copy returns a deep copy of the Styles.
func (s *Styles) Copy() *Styles {
	ns := &Styles{
		Rules:     s.GetRules(),
		Important: make(map[string]bool, len(s.Important)),
	}
	for k, v := range s.Important {
		ns.Important[k] = v
	}
	return ns
}

// Merge merges the rules from another Styles into this one.
func (s *Styles) Merge(other *Styles) {
	for k, v := range other.Rules {
		s.Rules[k] = v
	}
}

// MergeRules merges the given rules map into this Styles.
func (s *Styles) MergeRules(rules RulesMap) {
	for k, v := range rules {
		s.Rules[k] = v
	}
}

// ExtractRules returns all set rules as ExtractedRule values with resolved specificity.
// isDefault=true means these are widget-level defaults (lower priority).
func (s *Styles) ExtractRules(specificity Specificity3, isDefault bool, tieBreaker int) []ExtractedRule {
	defaultFlag := 1
	if isDefault {
		defaultFlag = 0
	}
	result := make([]ExtractedRule, 0, len(s.Rules))
	for name, value := range s.Rules {
		importantFlag := 0
		if s.Important[name] {
			importantFlag = 1
		}
		spec := Specificity6{
			defaultFlag, importantFlag,
			specificity[0], specificity[1], specificity[2],
			tieBreaker,
		}
		result = append(result, ExtractedRule{Name: name, Specificity: spec, Value: value})
	}
	return result
}

// IsAnimatable reports whether the named rule can be animated.
func IsAnimatable(rule string) bool {
	return Animatable[rule]
}

// GetTransition returns the Transition for the named key, or nil if not set.
func (s *Styles) GetTransition(key string) *Transition {
	if !IsAnimatable(key) {
		return nil
	}
	v, ok := s.Rules["transitions"]
	if !ok {
		return nil
	}
	transitions, ok := v.(map[string]Transition)
	if !ok {
		return nil
	}
	t, ok := transitions[key]
	if !ok {
		return nil
	}
	return &t
}

// Equal reports whether two Styles objects have the same rules.
func (s *Styles) Equal(other *Styles) bool {
	if len(s.Rules) != len(other.Rules) {
		return false
	}
	for k, v := range s.Rules {
		ov, ok := other.Rules[k]
		if !ok {
			return false
		}
		if fmt.Sprintf("%v", v) != fmt.Sprintf("%v", ov) {
			return false
		}
	}
	return true
}

// Gutter returns the total spacing (padding + border) around the widget.
func (s *Styles) Gutter() geometry.Spacing {
	pad := s.Padding()
	border := s.BorderSpacing()
	return geometry.Spacing{
		Top:    pad.Top + border.Top,
		Right:  pad.Right + border.Right,
		Bottom: pad.Bottom + border.Bottom,
		Left:   pad.Left + border.Left,
	}
}

// BorderSpacing returns the spacing taken up by the border.
func (s *Styles) BorderSpacing() geometry.Spacing {
	sp := geometry.Spacing{}
	if bt, ok := s.Rules["border_top"]; ok {
		if es, ok := bt.(EdgeStyle); ok && es.Type != "" && es.Type != "none" && es.Type != "hidden" {
			sp.Top = 1
		}
	}
	if br, ok := s.Rules["border_right"]; ok {
		if es, ok := br.(EdgeStyle); ok && es.Type != "" && es.Type != "none" && es.Type != "hidden" {
			sp.Right = 1
		}
	}
	if bb, ok := s.Rules["border_bottom"]; ok {
		if es, ok := bb.(EdgeStyle); ok && es.Type != "" && es.Type != "none" && es.Type != "hidden" {
			sp.Bottom = 1
		}
	}
	if bl, ok := s.Rules["border_left"]; ok {
		if es, ok := bl.(EdgeStyle); ok && es.Type != "" && es.Type != "none" && es.Type != "hidden" {
			sp.Left = 1
		}
	}
	return sp
}

// AutoDimensions reports whether width or height is set to "auto".
func (s *Styles) AutoDimensions() bool {
	if w := s.Width(); w != nil && w.IsAuto() {
		return true
	}
	if h := s.Height(); h != nil && h.IsAuto() {
		return true
	}
	return false
}

// IsRelativeWidth reports whether the width uses a relative unit (fraction or percent).
func (s *Styles) IsRelativeWidth() bool {
	w := s.Width()
	return w != nil && (w.Unit == UnitFraction || w.Unit == UnitPercent)
}

// IsRelativeHeight reports whether the height uses a relative unit.
func (s *Styles) IsRelativeHeight() bool {
	h := s.Height()
	return h != nil && (h.Unit == UnitFraction || h.Unit == UnitPercent)
}

// AlignWidth returns the X offset to align a child of the given width within parentWidth.
func (s *Styles) AlignWidth(width, parentWidth int) int {
	switch s.AlignHorizontal() {
	case "center":
		return (parentWidth - width) / 2
	case "right":
		return parentWidth - width
	default:
		return 0
	}
}

// AlignHeight returns the Y offset to align a child of the given height within parentHeight.
func (s *Styles) AlignHeight(height, parentHeight int) int {
	switch s.AlignVertical() {
	case "middle":
		return (parentHeight - height) / 2
	case "bottom":
		return parentHeight - height
	default:
		return 0
	}
}

// AlignSize returns the offset needed to align a child of the given size within parent.
func (s *Styles) AlignSize(child, parent [2]int) geometry.Offset {
	return geometry.Offset{
		X: s.AlignWidth(child[0], parent[0]),
		Y: s.AlignHeight(child[1], parent[1]),
	}
}

// CSSLines serializes the set rules to a list of CSS declaration strings.
func (s *Styles) CSSLines() []string {
	var lines []string
	appendDecl := func(name, value string) {
		if s.Important[name] {
			lines = append(lines, name+": "+value+" !important;")
		} else {
			lines = append(lines, name+": "+value+";")
		}
	}

	r := s.Rules

	if _, ok := r["display"]; ok {
		appendDecl("display", string(s.Display()))
	}
	if _, ok := r["visibility"]; ok {
		appendDecl("visibility", string(s.Visibility()))
	}
	if _, ok := r["padding"]; ok {
		p := s.Padding()
		appendDecl("padding", spacingCSS(p))
	}
	if _, ok := r["margin"]; ok {
		m := s.Margin()
		appendDecl("margin", spacingCSS(m))
	}
	// Border edges
	for _, side := range []string{"top", "right", "bottom", "left"} {
		key := "border_" + side
		cssKey := "border-" + side
		if v, ok := r[key]; ok {
			if es, ok := v.(EdgeStyle); ok {
				appendDecl(cssKey, string(es.Type)+" "+es.Color.Hex())
			}
		}
	}
	// Outline edges
	for _, side := range []string{"top", "right", "bottom", "left"} {
		key := "outline_" + side
		cssKey := "outline-" + side
		if v, ok := r[key]; ok {
			if es, ok := v.(EdgeStyle); ok {
				appendDecl(cssKey, string(es.Type)+" "+es.Color.Hex())
			}
		}
	}
	if _, ok := r["offset"]; ok {
		off := s.Offset()
		appendDecl("offset", off.X.String()+" "+off.Y.String())
	}
	if _, ok := r["position"]; ok {
		appendDecl("position", string(s.Position()))
	}
	if _, ok := r["dock"]; ok {
		appendDecl("dock", s.Dock())
	}
	if _, ok := r["split"]; ok {
		appendDecl("split", s.Split())
	}
	if _, ok := r["layers"]; ok {
		appendDecl("layers", strings.Join(s.Layers(), " "))
	}
	if _, ok := r["layer"]; ok {
		appendDecl("layer", s.Layer())
	}
	if _, ok := r["color"]; ok {
		appendDecl("color", s.Color().Hex())
	}
	if _, ok := r["background"]; ok {
		appendDecl("background", s.Background().Hex())
	}
	if _, ok := r["background_tint"]; ok {
		appendDecl("background-tint", s.BackgroundTint().Hex())
	}
	if _, ok := r["tint"]; ok {
		appendDecl("tint", s.Tint().CSS())
	}
	if _, ok := r["overflow_x"]; ok {
		appendDecl("overflow-x", string(s.OverflowX()))
	}
	if _, ok := r["overflow_y"]; ok {
		appendDecl("overflow-y", string(s.OverflowY()))
	}
	if _, ok := r["scrollbar_color"]; ok {
		appendDecl("scrollbar-color", s.ScrollbarColor().CSS())
	}
	if _, ok := r["scrollbar_color_hover"]; ok {
		appendDecl("scrollbar-color-hover", s.ScrollbarColorHover().CSS())
	}
	if _, ok := r["scrollbar_color_active"]; ok {
		appendDecl("scrollbar-color-active", s.ScrollbarColorActive().CSS())
	}
	if _, ok := r["scrollbar_corner_color"]; ok {
		appendDecl("scrollbar-corner-color", s.ScrollbarCornerColor().CSS())
	}
	if _, ok := r["scrollbar_background"]; ok {
		appendDecl("scrollbar-background", s.ScrollbarBackground().CSS())
	}
	if _, ok := r["scrollbar_background_hover"]; ok {
		appendDecl("scrollbar-background-hover", s.ScrollbarBackgroundHover().CSS())
	}
	if _, ok := r["scrollbar_background_active"]; ok {
		appendDecl("scrollbar-background-active", s.ScrollbarBackgroundActive().CSS())
	}
	if _, ok := r["scrollbar_gutter"]; ok {
		appendDecl("scrollbar-gutter", string(s.ScrollbarGutter()))
	}
	if _, ok := r["scrollbar_size_horizontal"]; ok {
		appendDecl("scrollbar-size-horizontal", fmt.Sprintf("%d", s.ScrollbarSizeHorizontal()))
	}
	if _, ok := r["scrollbar_size_vertical"]; ok {
		appendDecl("scrollbar-size-vertical", fmt.Sprintf("%d", s.ScrollbarSizeVertical()))
	}
	if _, ok := r["scrollbar_visibility"]; ok {
		appendDecl("scrollbar-visibility", string(s.ScrollbarVisibility()))
	}
	if _, ok := r["box_sizing"]; ok {
		appendDecl("box-sizing", string(s.BoxSizing()))
	}
	if w := s.Width(); w != nil {
		if _, ok := r["width"]; ok {
			appendDecl("width", w.String())
		}
	}
	if h := s.Height(); h != nil {
		if _, ok := r["height"]; ok {
			appendDecl("height", h.String())
		}
	}
	if mw := s.MinWidth(); mw != nil {
		if _, ok := r["min_width"]; ok {
			appendDecl("min-width", mw.String())
		}
	}
	if mh := s.MinHeight(); mh != nil {
		if _, ok := r["min_height"]; ok {
			appendDecl("min-height", mh.String())
		}
	}
	if mxw := s.MaxWidth(); mxw != nil {
		if _, ok := r["max_width"]; ok {
			appendDecl("max-width", mxw.String())
		}
	}
	if mxh := s.MaxHeight(); mxh != nil {
		if _, ok := r["max_height"]; ok {
			appendDecl("max-height", mxh.String())
		}
	}
	if _, ok := r["transitions"]; ok {
		var parts []string
		for k, t := range s.Transitions() {
			parts = append(parts, k+" "+t.String())
		}
		appendDecl("transition", strings.Join(parts, ", "))
	}
	hasAH := s.HasRule("align_horizontal")
	hasAV := s.HasRule("align_vertical")
	if hasAH && hasAV {
		appendDecl("align", string(s.AlignHorizontal())+" "+string(s.AlignVertical()))
	} else if hasAH {
		appendDecl("align-horizontal", string(s.AlignHorizontal()))
	} else if hasAV {
		appendDecl("align-vertical", string(s.AlignVertical()))
	}
	hasCAH := s.HasRule("content_align_horizontal")
	hasCAV := s.HasRule("content_align_vertical")
	if hasCAH && hasCAV {
		appendDecl("content-align", string(s.ContentAlignHorizontal())+" "+string(s.ContentAlignVertical()))
	} else if hasCAH {
		appendDecl("content-align-horizontal", string(s.ContentAlignHorizontal()))
	} else if hasCAV {
		appendDecl("content-align-vertical", string(s.ContentAlignVertical()))
	}
	if _, ok := r["text_align"]; ok {
		appendDecl("text-align", string(s.TextAlign()))
	}
	if _, ok := r["border_title_align"]; ok {
		appendDecl("border-title-align", string(s.BorderTitleAlign()))
	}
	if _, ok := r["border_subtitle_align"]; ok {
		appendDecl("border-subtitle-align", string(s.BorderSubtitleAlign()))
	}
	if _, ok := r["opacity"]; ok {
		appendDecl("opacity", fmt.Sprintf("%g", s.Opacity()))
	}
	if _, ok := r["text_opacity"]; ok {
		appendDecl("text-opacity", fmt.Sprintf("%g", s.TextOpacity()))
	}
	if _, ok := r["grid_columns"]; ok {
		var parts []string
		for _, sc := range s.GridColumns() {
			parts = append(parts, sc.String())
		}
		appendDecl("grid-columns", strings.Join(parts, " "))
	}
	if _, ok := r["grid_rows"]; ok {
		var parts []string
		for _, sc := range s.GridRows() {
			parts = append(parts, sc.String())
		}
		appendDecl("grid-rows", strings.Join(parts, " "))
	}
	if _, ok := r["grid_size_columns"]; ok {
		appendDecl("grid-size-columns", fmt.Sprintf("%d", s.GridSizeColumns()))
	}
	if _, ok := r["grid_size_rows"]; ok {
		appendDecl("grid-size-rows", fmt.Sprintf("%d", s.GridSizeRows()))
	}
	if _, ok := r["grid_gutter_horizontal"]; ok {
		appendDecl("grid-gutter-horizontal", fmt.Sprintf("%d", s.GridGutterHorizontal()))
	}
	if _, ok := r["grid_gutter_vertical"]; ok {
		appendDecl("grid-gutter-vertical", fmt.Sprintf("%d", s.GridGutterVertical()))
	}
	if _, ok := r["row_span"]; ok {
		appendDecl("row-span", fmt.Sprintf("%d", s.RowSpan()))
	}
	if _, ok := r["column_span"]; ok {
		appendDecl("column-span", fmt.Sprintf("%d", s.ColumnSpan()))
	}
	if _, ok := r["link_color"]; ok {
		appendDecl("link-color", s.LinkColor().CSS())
	}
	if _, ok := r["link_background"]; ok {
		appendDecl("link-background", s.LinkBackground().CSS())
	}
	if _, ok := r["link_color_hover"]; ok {
		appendDecl("link-color-hover", s.LinkColorHover().CSS())
	}
	if _, ok := r["link_background_hover"]; ok {
		appendDecl("link-background-hover", s.LinkBackgroundHover().CSS())
	}
	if _, ok := r["border_title_color"]; ok {
		appendDecl("title-color", s.BorderTitleColor().CSS())
	}
	if _, ok := r["border_title_background"]; ok {
		appendDecl("title-background", s.BorderTitleBackground().CSS())
	}
	if _, ok := r["border_subtitle_color"]; ok {
		appendDecl("subtitle-color", s.BorderSubtitleColor().CSS())
	}
	if _, ok := r["border_subtitle_background"]; ok {
		appendDecl("subtitle-background", s.BorderSubtitleBackground().CSS())
	}
	if _, ok := r["text_style"]; ok {
		if v := s.TextStyle(); v != "" {
			appendDecl("text-style", v)
		}
	}
	if _, ok := r["link_style"]; ok {
		if v := s.LinkStyle(); v != "" {
			appendDecl("link-style", v)
		}
	}
	if _, ok := r["link_style_hover"]; ok {
		if v := s.LinkStyleHover(); v != "" {
			appendDecl("link-style-hover", v)
		}
	}
	if _, ok := r["border_title_style"]; ok {
		if v := s.BorderTitleStyle(); v != "" {
			appendDecl("border-title-style", v)
		}
	}
	if _, ok := r["border_subtitle_style"]; ok {
		if v := s.BorderSubtitleStyle(); v != "" {
			appendDecl("border-subtitle-style", v)
		}
	}
	if _, ok := r["overlay"]; ok {
		appendDecl("overlay", string(s.Overlay()))
	}
	hasCX := s.HasRule("constrain_x")
	hasCY := s.HasRule("constrain_y")
	if hasCX && hasCY {
		cx, cy := s.ConstrainX(), s.ConstrainY()
		if cx == cy {
			appendDecl("constrain", string(cx))
		} else {
			appendDecl("constrain", string(cx)+" "+string(cy))
		}
	} else if hasCX {
		appendDecl("constrain-x", string(s.ConstrainX()))
	} else if hasCY {
		appendDecl("constrain-y", string(s.ConstrainY()))
	}
	if _, ok := r["keyline"]; ok {
		kl := s.Keyline()
		if kl.Type != "none" {
			appendDecl("keyline", string(kl.Type)+", "+kl.Color.CSS())
		}
	}
	if _, ok := r["hatch"]; ok {
		h := s.Hatch()
		appendDecl("hatch", `"`+string(h.Type)+`" `+h.Color.CSS())
	}
	if _, ok := r["text_wrap"]; ok {
		appendDecl("text-wrap", string(s.TextWrap()))
	}
	if _, ok := r["text_overflow"]; ok {
		appendDecl("text-overflow", string(s.TextOverflow()))
	}
	if _, ok := r["expand"]; ok {
		appendDecl("expand", string(s.Expand()))
	}
	if _, ok := r["line_pad"]; ok {
		appendDecl("line-pad", fmt.Sprintf("%d", s.LinePad()))
	}
	sortStrings(lines)
	return lines
}

// CSS returns the serialized CSS declaration string.
func (s *Styles) CSS() string {
	return strings.Join(s.CSSLines(), "\n")
}

// spacingCSS converts a geometry.Spacing to its CSS value string.
func spacingCSS(sp geometry.Spacing) string {
	if sp.Top == sp.Right && sp.Right == sp.Bottom && sp.Bottom == sp.Left {
		return fmt.Sprintf("%d", sp.Top)
	}
	if sp.Top == sp.Bottom && sp.Left == sp.Right {
		return fmt.Sprintf("%d %d", sp.Top, sp.Right)
	}
	if sp.Left == sp.Right {
		return fmt.Sprintf("%d %d %d", sp.Top, sp.Right, sp.Bottom)
	}
	return fmt.Sprintf("%d %d %d %d", sp.Top, sp.Right, sp.Bottom, sp.Left)
}

// ── Typed getters ──────────────────────────────────────────────────────────

func (s *Styles) Display() Display {
	if v, ok := s.Rules["display"]; ok {
		if d, ok := v.(Display); ok {
			return d
		}
		if str, ok := v.(string); ok {
			return Display(str)
		}
	}
	return "block"
}

func (s *Styles) Visibility() Visibility {
	if v, ok := s.Rules["visibility"]; ok {
		if vis, ok := v.(Visibility); ok {
			return vis
		}
		if str, ok := v.(string); ok {
			return Visibility(str)
		}
	}
	return "visible"
}

func (s *Styles) Layout() string {
	if v, ok := s.Rules["layout"]; ok {
		if str, ok := v.(string); ok {
			return str
		}
	}
	return ""
}

func (s *Styles) AutoColor() bool {
	if v, ok := s.Rules["auto_color"]; ok {
		if b, ok := v.(bool); ok {
			return b
		}
	}
	return false
}

func (s *Styles) Color() color.Color {
	if v, ok := s.Rules["color"]; ok {
		if c, ok := v.(color.Color); ok {
			return c
		}
	}
	return color.New(255, 255, 255)
}

func (s *Styles) Background() color.Color {
	if v, ok := s.Rules["background"]; ok {
		if c, ok := v.(color.Color); ok {
			return c
		}
	}
	return color.Transparent
}

func (s *Styles) BackgroundTint() color.Color {
	if v, ok := s.Rules["background_tint"]; ok {
		if c, ok := v.(color.Color); ok {
			return c
		}
	}
	return color.Transparent
}

func (s *Styles) Opacity() float64 {
	if v, ok := s.Rules["opacity"]; ok {
		if f, ok := v.(float64); ok {
			return f
		}
	}
	return 1.0
}

func (s *Styles) TextOpacity() float64 {
	if v, ok := s.Rules["text_opacity"]; ok {
		if f, ok := v.(float64); ok {
			return f
		}
	}
	return 1.0
}

func (s *Styles) Padding() geometry.Spacing {
	if v, ok := s.Rules["padding"]; ok {
		if sp, ok := v.(geometry.Spacing); ok {
			return sp
		}
	}
	return geometry.Spacing{}
}

func (s *Styles) Margin() geometry.Spacing {
	if v, ok := s.Rules["margin"]; ok {
		if sp, ok := v.(geometry.Spacing); ok {
			return sp
		}
	}
	return geometry.Spacing{}
}

func (s *Styles) Offset() ScalarOffset {
	if v, ok := s.Rules["offset"]; ok {
		if so, ok := v.(ScalarOffset); ok {
			return so
		}
	}
	return NullScalarOffset
}

func (s *Styles) Position() Position {
	if v, ok := s.Rules["position"]; ok {
		if p, ok := v.(Position); ok {
			return p
		}
		if str, ok := v.(string); ok {
			return Position(str)
		}
	}
	return "relative"
}

func (s *Styles) borderEdge(name string) EdgeStyle {
	if v, ok := s.Rules[name]; ok {
		if es, ok := v.(EdgeStyle); ok {
			return es
		}
	}
	return EdgeStyle{Type: "none", Color: color.New(0, 255, 0)}
}

func (s *Styles) BorderTop() EdgeStyle     { return s.borderEdge("border_top") }
func (s *Styles) BorderRight() EdgeStyle   { return s.borderEdge("border_right") }
func (s *Styles) BorderBottom() EdgeStyle  { return s.borderEdge("border_bottom") }
func (s *Styles) BorderLeft() EdgeStyle    { return s.borderEdge("border_left") }
func (s *Styles) OutlineTop() EdgeStyle    { return s.borderEdge("outline_top") }
func (s *Styles) OutlineRight() EdgeStyle  { return s.borderEdge("outline_right") }
func (s *Styles) OutlineBottom() EdgeStyle { return s.borderEdge("outline_bottom") }
func (s *Styles) OutlineLeft() EdgeStyle   { return s.borderEdge("outline_left") }

func (s *Styles) Keyline() EdgeStyle {
	if v, ok := s.Rules["keyline"]; ok {
		if es, ok := v.(EdgeStyle); ok {
			return es
		}
	}
	return EdgeStyle{Type: "none"}
}

func (s *Styles) BorderTitleAlign() AlignHorizontal {
	if v, ok := s.Rules["border_title_align"]; ok {
		if a, ok := v.(AlignHorizontal); ok {
			return a
		}
		if str, ok := v.(string); ok {
			return AlignHorizontal(str)
		}
	}
	return "left"
}

func (s *Styles) BorderSubtitleAlign() AlignHorizontal {
	if v, ok := s.Rules["border_subtitle_align"]; ok {
		if a, ok := v.(AlignHorizontal); ok {
			return a
		}
		if str, ok := v.(string); ok {
			return AlignHorizontal(str)
		}
	}
	return "right"
}

func (s *Styles) BoxSizing() BoxSizing {
	if v, ok := s.Rules["box_sizing"]; ok {
		if b, ok := v.(BoxSizing); ok {
			return b
		}
		if str, ok := v.(string); ok {
			return BoxSizing(str)
		}
	}
	return "border-box"
}

func (s *Styles) scalar(name string) *Scalar {
	if v, ok := s.Rules[name]; ok {
		if sc, ok := v.(Scalar); ok {
			return &sc
		}
	}
	return nil
}

func (s *Styles) Width() *Scalar     { return s.scalar("width") }
func (s *Styles) Height() *Scalar    { return s.scalar("height") }
func (s *Styles) MinWidth() *Scalar  { return s.scalar("min_width") }
func (s *Styles) MinHeight() *Scalar { return s.scalar("min_height") }
func (s *Styles) MaxWidth() *Scalar  { return s.scalar("max_width") }
func (s *Styles) MaxHeight() *Scalar { return s.scalar("max_height") }

func (s *Styles) Dock() string {
	if v, ok := s.Rules["dock"]; ok {
		if str, ok := v.(string); ok {
			return str
		}
	}
	return "none"
}

func (s *Styles) Split() string {
	if v, ok := s.Rules["split"]; ok {
		if str, ok := v.(string); ok {
			return str
		}
	}
	return "none"
}

func (s *Styles) OverflowX() Overflow {
	if v, ok := s.Rules["overflow_x"]; ok {
		if o, ok := v.(Overflow); ok {
			return o
		}
		if str, ok := v.(string); ok {
			return Overflow(str)
		}
	}
	return "hidden"
}

func (s *Styles) OverflowY() Overflow {
	if v, ok := s.Rules["overflow_y"]; ok {
		if o, ok := v.(Overflow); ok {
			return o
		}
		if str, ok := v.(string); ok {
			return Overflow(str)
		}
	}
	return "hidden"
}

func (s *Styles) Layers() []string {
	if v, ok := s.Rules["layers"]; ok {
		if layers, ok := v.([]string); ok {
			return layers
		}
	}
	return nil
}

func (s *Styles) Layer() string {
	if v, ok := s.Rules["layer"]; ok {
		if str, ok := v.(string); ok {
			return str
		}
	}
	return ""
}

func (s *Styles) Transitions() map[string]Transition {
	if v, ok := s.Rules["transitions"]; ok {
		if t, ok := v.(map[string]Transition); ok {
			return t
		}
	}
	return nil
}

func (s *Styles) Tint() color.Color {
	if v, ok := s.Rules["tint"]; ok {
		if c, ok := v.(color.Color); ok {
			return c
		}
	}
	return color.Transparent
}

func (s *Styles) colorRule(name string, def color.Color) color.Color {
	if v, ok := s.Rules[name]; ok {
		if c, ok := v.(color.Color); ok {
			return c
		}
	}
	return def
}

func (s *Styles) ScrollbarColor() color.Color {
	return s.colorRule("scrollbar_color", color.New(188, 0, 188)) // ansi_bright_magenta approx
}

func (s *Styles) ScrollbarColorHover() color.Color {
	return s.colorRule("scrollbar_color_hover", color.New(128, 128, 0)) // ansi_yellow approx
}

func (s *Styles) ScrollbarColorActive() color.Color {
	return s.colorRule("scrollbar_color_active", color.New(255, 255, 0)) // ansi_bright_yellow approx
}

func (s *Styles) ScrollbarCornerColor() color.Color {
	return s.colorRule("scrollbar_corner_color", color.New(102, 102, 102))
}

func (s *Styles) ScrollbarBackground() color.Color {
	return s.colorRule("scrollbar_background", color.New(85, 85, 85))
}

func (s *Styles) ScrollbarBackgroundHover() color.Color {
	return s.colorRule("scrollbar_background_hover", color.New(68, 68, 68))
}

func (s *Styles) ScrollbarBackgroundActive() color.Color {
	return s.colorRule("scrollbar_background_active", color.Black)
}

func (s *Styles) ScrollbarGutter() ScrollbarGutter {
	if v, ok := s.Rules["scrollbar_gutter"]; ok {
		if sg, ok := v.(ScrollbarGutter); ok {
			return sg
		}
		if str, ok := v.(string); ok {
			return ScrollbarGutter(str)
		}
	}
	return "auto"
}

func (s *Styles) intRule(name string, def int) int {
	if v, ok := s.Rules[name]; ok {
		if i, ok := v.(int); ok {
			return i
		}
	}
	return def
}

func (s *Styles) ScrollbarSizeVertical() int   { return s.intRule("scrollbar_size_vertical", 2) }
func (s *Styles) ScrollbarSizeHorizontal() int { return s.intRule("scrollbar_size_horizontal", 1) }

func (s *Styles) ScrollbarVisibility() ScrollbarVisibility {
	if v, ok := s.Rules["scrollbar_visibility"]; ok {
		if sv, ok := v.(ScrollbarVisibility); ok {
			return sv
		}
		if str, ok := v.(string); ok {
			return ScrollbarVisibility(str)
		}
	}
	return "visible"
}

func (s *Styles) AlignHorizontal() AlignHorizontal {
	if v, ok := s.Rules["align_horizontal"]; ok {
		if a, ok := v.(AlignHorizontal); ok {
			return a
		}
		if str, ok := v.(string); ok {
			return AlignHorizontal(str)
		}
	}
	return "left"
}

func (s *Styles) AlignVertical() AlignVertical {
	if v, ok := s.Rules["align_vertical"]; ok {
		if a, ok := v.(AlignVertical); ok {
			return a
		}
		if str, ok := v.(string); ok {
			return AlignVertical(str)
		}
	}
	return "top"
}

func (s *Styles) ContentAlignHorizontal() AlignHorizontal {
	if v, ok := s.Rules["content_align_horizontal"]; ok {
		if a, ok := v.(AlignHorizontal); ok {
			return a
		}
		if str, ok := v.(string); ok {
			return AlignHorizontal(str)
		}
	}
	return "left"
}

func (s *Styles) ContentAlignVertical() AlignVertical {
	if v, ok := s.Rules["content_align_vertical"]; ok {
		if a, ok := v.(AlignVertical); ok {
			return a
		}
		if str, ok := v.(string); ok {
			return AlignVertical(str)
		}
	}
	return "top"
}

func (s *Styles) GridSizeRows() int         { return s.intRule("grid_size_rows", 0) }
func (s *Styles) GridSizeColumns() int      { return s.intRule("grid_size_columns", 1) }
func (s *Styles) GridGutterHorizontal() int { return s.intRule("grid_gutter_horizontal", 0) }
func (s *Styles) GridGutterVertical() int   { return s.intRule("grid_gutter_vertical", 0) }
func (s *Styles) RowSpan() int              { return s.intRule("row_span", 1) }
func (s *Styles) ColumnSpan() int           { return s.intRule("column_span", 1) }

func (s *Styles) scalarSlice(name string) []Scalar {
	if v, ok := s.Rules[name]; ok {
		if sl, ok := v.([]Scalar); ok {
			return sl
		}
	}
	return nil
}

func (s *Styles) GridRows() []Scalar    { return s.scalarSlice("grid_rows") }
func (s *Styles) GridColumns() []Scalar { return s.scalarSlice("grid_columns") }

func (s *Styles) TextAlign() TextAlign {
	if v, ok := s.Rules["text_align"]; ok {
		if t, ok := v.(TextAlign); ok {
			return t
		}
		if str, ok := v.(string); ok {
			return TextAlign(str)
		}
	}
	return "start"
}

func (s *Styles) LinkColor() color.Color { return s.colorRule("link_color", color.Transparent) }
func (s *Styles) AutoLinkColor() bool    { return s.boolRule("auto_link_color", false) }
func (s *Styles) LinkBackground() color.Color {
	return s.colorRule("link_background", color.Transparent)
}
func (s *Styles) LinkColorHover() color.Color {
	return s.colorRule("link_color_hover", color.Transparent)
}
func (s *Styles) AutoLinkColorHover() bool { return s.boolRule("auto_link_color_hover", false) }
func (s *Styles) LinkBackgroundHover() color.Color {
	return s.colorRule("link_background_hover", color.Transparent)
}

func (s *Styles) AutoBorderTitleColor() bool { return s.boolRule("auto_border_title_color", false) }
func (s *Styles) BorderTitleColor() color.Color {
	return s.colorRule("border_title_color", color.NewWithAlpha(255, 255, 255, 0))
}
func (s *Styles) BorderTitleBackground() color.Color {
	return s.colorRule("border_title_background", color.Transparent)
}
func (s *Styles) AutoBorderSubtitleColor() bool {
	return s.boolRule("auto_border_subtitle_color", false)
}
func (s *Styles) BorderSubtitleColor() color.Color {
	return s.colorRule("border_subtitle_color", color.NewWithAlpha(255, 255, 255, 0))
}
func (s *Styles) BorderSubtitleBackground() color.Color {
	return s.colorRule("border_subtitle_background", color.Transparent)
}

func (s *Styles) boolRule(name string, def bool) bool {
	if v, ok := s.Rules[name]; ok {
		if b, ok := v.(bool); ok {
			return b
		}
	}
	return def
}

func (s *Styles) Hatch() EdgeStyle {
	if v, ok := s.Rules["hatch"]; ok {
		if es, ok := v.(EdgeStyle); ok {
			return es
		}
	}
	return EdgeStyle{Type: "none"}
}

func (s *Styles) Overlay() Overlay {
	if v, ok := s.Rules["overlay"]; ok {
		if o, ok := v.(Overlay); ok {
			return o
		}
		if str, ok := v.(string); ok {
			return Overlay(str)
		}
	}
	return "none"
}

func (s *Styles) ConstrainX() Constrain {
	if v, ok := s.Rules["constrain_x"]; ok {
		if c, ok := v.(Constrain); ok {
			return c
		}
		if str, ok := v.(string); ok {
			return Constrain(str)
		}
	}
	return "none"
}

func (s *Styles) ConstrainY() Constrain {
	if v, ok := s.Rules["constrain_y"]; ok {
		if c, ok := v.(Constrain); ok {
			return c
		}
		if str, ok := v.(string); ok {
			return Constrain(str)
		}
	}
	return "none"
}

func (s *Styles) TextWrap() TextWrap {
	if v, ok := s.Rules["text_wrap"]; ok {
		if t, ok := v.(TextWrap); ok {
			return t
		}
		if str, ok := v.(string); ok {
			return TextWrap(str)
		}
	}
	return "wrap"
}

func (s *Styles) TextOverflow() TextOverflow {
	if v, ok := s.Rules["text_overflow"]; ok {
		if t, ok := v.(TextOverflow); ok {
			return t
		}
		if str, ok := v.(string); ok {
			return TextOverflow(str)
		}
	}
	return "fold"
}

func (s *Styles) Expand() Expand {
	if v, ok := s.Rules["expand"]; ok {
		if e, ok := v.(Expand); ok {
			return e
		}
		if str, ok := v.(string); ok {
			return Expand(str)
		}
	}
	return "greedy"
}

func (s *Styles) LinePad() int { return s.intRule("line_pad", 0) }

func (s *Styles) TextStyle() string {
	if v, ok := s.Rules["text_style"]; ok {
		if str, ok := v.(string); ok {
			return str
		}
	}
	return ""
}

func (s *Styles) LinkStyle() string {
	if v, ok := s.Rules["link_style"]; ok {
		if str, ok := v.(string); ok {
			return str
		}
	}
	return ""
}

func (s *Styles) LinkStyleHover() string {
	if v, ok := s.Rules["link_style_hover"]; ok {
		if str, ok := v.(string); ok {
			return str
		}
	}
	return ""
}

func (s *Styles) BorderTitleStyle() string {
	if v, ok := s.Rules["border_title_style"]; ok {
		if str, ok := v.(string); ok {
			return str
		}
	}
	return ""
}

func (s *Styles) BorderSubtitleStyle() string {
	if v, ok := s.Rules["border_subtitle_style"]; ok {
		if str, ok := v.(string); ok {
			return str
		}
	}
	return ""
}

func (s *Styles) Pointer() PointerShape {
	if v, ok := s.Rules["pointer"]; ok {
		if p, ok := v.(PointerShape); ok {
			return p
		}
		if str, ok := v.(string); ok {
			return PointerShape(str)
		}
	}
	return "default"
}

// ── Typed setters ──────────────────────────────────────────────────────────

func (s *Styles) SetDisplay(v Display)                     { s.Rules["display"] = v }
func (s *Styles) SetVisibility(v Visibility)               { s.Rules["visibility"] = v }
func (s *Styles) SetLayout(v string)                       { s.Rules["layout"] = v }
func (s *Styles) SetAutoColor(v bool)                      { s.Rules["auto_color"] = v }
func (s *Styles) SetColor(v color.Color)                   { s.Rules["color"] = v }
func (s *Styles) SetBackground(v color.Color)              { s.Rules["background"] = v }
func (s *Styles) SetBackgroundTint(v color.Color)          { s.Rules["background_tint"] = v }
func (s *Styles) SetOpacity(v float64)                     { s.Rules["opacity"] = v }
func (s *Styles) SetTextOpacity(v float64)                 { s.Rules["text_opacity"] = v }
func (s *Styles) SetPadding(v geometry.Spacing)            { s.Rules["padding"] = v }
func (s *Styles) SetMargin(v geometry.Spacing)             { s.Rules["margin"] = v }
func (s *Styles) SetOffset(v ScalarOffset)                 { s.Rules["offset"] = v }
func (s *Styles) SetPosition(v Position)                   { s.Rules["position"] = v }
func (s *Styles) SetBorderTop(v EdgeStyle)                 { s.Rules["border_top"] = v }
func (s *Styles) SetBorderRight(v EdgeStyle)               { s.Rules["border_right"] = v }
func (s *Styles) SetBorderBottom(v EdgeStyle)              { s.Rules["border_bottom"] = v }
func (s *Styles) SetBorderLeft(v EdgeStyle)                { s.Rules["border_left"] = v }
func (s *Styles) SetOutlineTop(v EdgeStyle)                { s.Rules["outline_top"] = v }
func (s *Styles) SetOutlineRight(v EdgeStyle)              { s.Rules["outline_right"] = v }
func (s *Styles) SetOutlineBottom(v EdgeStyle)             { s.Rules["outline_bottom"] = v }
func (s *Styles) SetOutlineLeft(v EdgeStyle)               { s.Rules["outline_left"] = v }
func (s *Styles) SetKeyline(v EdgeStyle)                   { s.Rules["keyline"] = v }
func (s *Styles) SetBorderTitleAlign(v AlignHorizontal)    { s.Rules["border_title_align"] = v }
func (s *Styles) SetBorderSubtitleAlign(v AlignHorizontal) { s.Rules["border_subtitle_align"] = v }
func (s *Styles) SetBoxSizing(v BoxSizing)                 { s.Rules["box_sizing"] = v }
func (s *Styles) SetWidth(v Scalar)                        { s.Rules["width"] = v }
func (s *Styles) SetHeight(v Scalar)                       { s.Rules["height"] = v }
func (s *Styles) SetMinWidth(v Scalar)                     { s.Rules["min_width"] = v }
func (s *Styles) SetMinHeight(v Scalar)                    { s.Rules["min_height"] = v }
func (s *Styles) SetMaxWidth(v Scalar)                     { s.Rules["max_width"] = v }
func (s *Styles) SetMaxHeight(v Scalar)                    { s.Rules["max_height"] = v }
func (s *Styles) SetDock(v string)                         { s.Rules["dock"] = v }
func (s *Styles) SetSplit(v string)                        { s.Rules["split"] = v }
func (s *Styles) SetOverflowX(v Overflow)                  { s.Rules["overflow_x"] = v }
func (s *Styles) SetOverflowY(v Overflow)                  { s.Rules["overflow_y"] = v }
func (s *Styles) SetLayers(v []string)                     { s.Rules["layers"] = v }
func (s *Styles) SetLayer(v string)                        { s.Rules["layer"] = v }
func (s *Styles) SetTransitions(v map[string]Transition)   { s.Rules["transitions"] = v }
func (s *Styles) SetTint(v color.Color)                    { s.Rules["tint"] = v }
func (s *Styles) SetScrollbarColor(v color.Color)          { s.Rules["scrollbar_color"] = v }
func (s *Styles) SetScrollbarColorHover(v color.Color)     { s.Rules["scrollbar_color_hover"] = v }
func (s *Styles) SetScrollbarColorActive(v color.Color)    { s.Rules["scrollbar_color_active"] = v }
func (s *Styles) SetScrollbarCornerColor(v color.Color)    { s.Rules["scrollbar_corner_color"] = v }
func (s *Styles) SetScrollbarBackground(v color.Color)     { s.Rules["scrollbar_background"] = v }
func (s *Styles) SetScrollbarBackgroundHover(v color.Color) {
	s.Rules["scrollbar_background_hover"] = v
}
func (s *Styles) SetScrollbarBackgroundActive(v color.Color) {
	s.Rules["scrollbar_background_active"] = v
}
func (s *Styles) SetScrollbarGutter(v ScrollbarGutter)         { s.Rules["scrollbar_gutter"] = v }
func (s *Styles) SetScrollbarSizeVertical(v int)               { s.Rules["scrollbar_size_vertical"] = v }
func (s *Styles) SetScrollbarSizeHorizontal(v int)             { s.Rules["scrollbar_size_horizontal"] = v }
func (s *Styles) SetScrollbarVisibility(v ScrollbarVisibility) { s.Rules["scrollbar_visibility"] = v }
func (s *Styles) SetAlignHorizontal(v AlignHorizontal)         { s.Rules["align_horizontal"] = v }
func (s *Styles) SetAlignVertical(v AlignVertical)             { s.Rules["align_vertical"] = v }
func (s *Styles) SetContentAlignHorizontal(v AlignHorizontal) {
	s.Rules["content_align_horizontal"] = v
}
func (s *Styles) SetContentAlignVertical(v AlignVertical) { s.Rules["content_align_vertical"] = v }
func (s *Styles) SetGridSizeRows(v int)                   { s.Rules["grid_size_rows"] = v }
func (s *Styles) SetGridSizeColumns(v int)                { s.Rules["grid_size_columns"] = v }
func (s *Styles) SetGridGutterHorizontal(v int)           { s.Rules["grid_gutter_horizontal"] = v }
func (s *Styles) SetGridGutterVertical(v int)             { s.Rules["grid_gutter_vertical"] = v }
func (s *Styles) SetGridRows(v []Scalar)                  { s.Rules["grid_rows"] = v }
func (s *Styles) SetGridColumns(v []Scalar)               { s.Rules["grid_columns"] = v }
func (s *Styles) SetRowSpan(v int)                        { s.Rules["row_span"] = v }
func (s *Styles) SetColumnSpan(v int)                     { s.Rules["column_span"] = v }
func (s *Styles) SetTextAlign(v TextAlign)                { s.Rules["text_align"] = v }
func (s *Styles) SetLinkColor(v color.Color)              { s.Rules["link_color"] = v }
func (s *Styles) SetAutoLinkColor(v bool)                 { s.Rules["auto_link_color"] = v }
func (s *Styles) SetLinkBackground(v color.Color)         { s.Rules["link_background"] = v }
func (s *Styles) SetLinkColorHover(v color.Color)         { s.Rules["link_color_hover"] = v }
func (s *Styles) SetAutoLinkColorHover(v bool)            { s.Rules["auto_link_color_hover"] = v }
func (s *Styles) SetLinkBackgroundHover(v color.Color)    { s.Rules["link_background_hover"] = v }
func (s *Styles) SetAutoBorderTitleColor(v bool)          { s.Rules["auto_border_title_color"] = v }
func (s *Styles) SetBorderTitleColor(v color.Color)       { s.Rules["border_title_color"] = v }
func (s *Styles) SetBorderTitleBackground(v color.Color)  { s.Rules["border_title_background"] = v }
func (s *Styles) SetAutoBorderSubtitleColor(v bool)       { s.Rules["auto_border_subtitle_color"] = v }
func (s *Styles) SetBorderSubtitleColor(v color.Color)    { s.Rules["border_subtitle_color"] = v }
func (s *Styles) SetBorderSubtitleBackground(v color.Color) {
	s.Rules["border_subtitle_background"] = v
}
func (s *Styles) SetHatch(v EdgeStyle)            { s.Rules["hatch"] = v }
func (s *Styles) SetOverlay(v Overlay)            { s.Rules["overlay"] = v }
func (s *Styles) SetConstrainX(v Constrain)       { s.Rules["constrain_x"] = v }
func (s *Styles) SetConstrainY(v Constrain)       { s.Rules["constrain_y"] = v }
func (s *Styles) SetTextWrap(v TextWrap)          { s.Rules["text_wrap"] = v }
func (s *Styles) SetTextOverflow(v TextOverflow)  { s.Rules["text_overflow"] = v }
func (s *Styles) SetExpand(v Expand)              { s.Rules["expand"] = v }
func (s *Styles) SetLinePad(v int)                { s.Rules["line_pad"] = v }
func (s *Styles) SetPointer(v PointerShape)       { s.Rules["pointer"] = v }
func (s *Styles) SetTextStyle(v string)           { s.Rules["text_style"] = v }
func (s *Styles) SetLinkStyle(v string)           { s.Rules["link_style"] = v }
func (s *Styles) SetLinkStyleHover(v string)      { s.Rules["link_style_hover"] = v }
func (s *Styles) SetBorderTitleStyle(v string)    { s.Rules["border_title_style"] = v }
func (s *Styles) SetBorderSubtitleStyle(v string) { s.Rules["border_subtitle_style"] = v }

// builtInStyleDefault returns the package-level default value for a CSS rule
// name, mirroring the return values of the typed getter methods on Styles.
// The second return value is false for rules that have no meaningful built-in
// default (e.g. width/height scalar pointers that are nil by design).
func builtInStyleDefault(name string) (any, bool) {
	borderNone := EdgeStyle{Type: "none", Color: color.New(0, 255, 0)}
	switch name {
	case "display":
		return Display("block"), true
	case "visibility":
		return Visibility("visible"), true
	case "layout":
		return "", true
	case "auto_color":
		return false, true
	case "color":
		return color.New(255, 255, 255), true
	case "background":
		return color.Transparent, true
	case "background_tint":
		return color.Transparent, true
	case "opacity":
		return 1.0, true
	case "text_opacity":
		return 1.0, true
	case "padding":
		return geometry.Spacing{}, true
	case "margin":
		return geometry.Spacing{}, true
	case "offset":
		return NullScalarOffset, true
	case "position":
		return Position("relative"), true
	case "border_top", "border_right", "border_bottom", "border_left":
		return borderNone, true
	case "outline_top", "outline_right", "outline_bottom", "outline_left":
		return borderNone, true
	case "keyline":
		return EdgeStyle{Type: "none"}, true
	case "border_title_align":
		return AlignHorizontal("left"), true
	case "border_subtitle_align":
		return AlignHorizontal("right"), true
	case "box_sizing":
		return BoxSizing("border-box"), true
	// width, height, min_width, min_height, max_width, max_height: nil *Scalar — no built-in default.
	case "dock":
		return "none", true
	case "split":
		return "none", true
	case "overflow_x":
		return Overflow("hidden"), true
	case "overflow_y":
		return Overflow("hidden"), true
	// layers, layer, transitions, grid_rows, grid_columns: nil slices/maps — no built-in default.
	case "tint":
		return color.Transparent, true
	case "scrollbar_color":
		return color.New(188, 0, 188), true
	case "scrollbar_color_hover":
		return color.New(128, 128, 0), true
	case "scrollbar_color_active":
		return color.New(255, 255, 0), true
	case "scrollbar_corner_color":
		return color.New(102, 102, 102), true
	case "scrollbar_background":
		return color.New(85, 85, 85), true
	case "scrollbar_background_hover":
		return color.New(68, 68, 68), true
	case "scrollbar_background_active":
		return color.Black, true
	case "scrollbar_gutter":
		return ScrollbarGutter("auto"), true
	case "scrollbar_size_vertical":
		return 2, true
	case "scrollbar_size_horizontal":
		return 1, true
	case "scrollbar_visibility":
		return ScrollbarVisibility("visible"), true
	case "align_horizontal":
		return AlignHorizontal("left"), true
	case "align_vertical":
		return AlignVertical("top"), true
	case "content_align_horizontal":
		return AlignHorizontal("left"), true
	case "content_align_vertical":
		return AlignVertical("top"), true
	case "grid_size_rows":
		return 0, true
	case "grid_size_columns":
		return 1, true
	case "grid_gutter_horizontal":
		return 0, true
	case "grid_gutter_vertical":
		return 0, true
	case "row_span":
		return 1, true
	case "column_span":
		return 1, true
	case "text_align":
		return TextAlign("start"), true
	case "link_color":
		return color.Transparent, true
	case "auto_link_color":
		return false, true
	case "link_background":
		return color.Transparent, true
	case "link_style":
		return "", true
	case "link_color_hover":
		return color.Transparent, true
	case "auto_link_color_hover":
		return false, true
	case "link_background_hover":
		return color.Transparent, true
	case "link_style_hover":
		return "", true
	case "auto_border_title_color":
		return false, true
	case "border_title_color":
		return color.NewWithAlpha(255, 255, 255, 0), true
	case "border_title_background":
		return color.Transparent, true
	case "border_title_style":
		return "", true
	case "auto_border_subtitle_color":
		return false, true
	case "border_subtitle_color":
		return color.NewWithAlpha(255, 255, 255, 0), true
	case "border_subtitle_background":
		return color.Transparent, true
	case "border_subtitle_style":
		return "", true
	case "hatch":
		return EdgeStyle{Type: "none"}, true
	case "overlay":
		return Overlay("none"), true
	case "constrain_x":
		return Constrain("none"), true
	case "constrain_y":
		return Constrain("none"), true
	case "text_wrap":
		return TextWrap("wrap"), true
	case "text_overflow":
		return TextOverflow("fold"), true
	case "expand":
		return Expand("greedy"), true
	case "line_pad":
		return 0, true
	case "pointer":
		return PointerShape("default"), true
	case "text_style":
		return "", true
	}
	return nil, false
}
