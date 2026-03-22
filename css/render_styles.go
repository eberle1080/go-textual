package css

import (
	"strings"

	"github.com/eberle1080/go-textual/color"
	"github.com/eberle1080/go-textual/geometry"
)

// RenderStyles presents a combined view of a base Styles and inline Styles.
// Inline rules take precedence over base rules.
type RenderStyles struct {
	base    *Styles
	inline  *Styles
	updates int
}

// NewRenderStyles creates a RenderStyles that merges base and inline.
func NewRenderStyles(base, inline *Styles) *RenderStyles {
	return &RenderStyles{base: base, inline: inline}
}

// Base returns the base Styles.
func (rs *RenderStyles) Base() *Styles { return rs.base }

// Inline returns the inline Styles.
func (rs *RenderStyles) Inline() *Styles { return rs.inline }

// HasRule reports whether either inline or base has the named rule.
func (rs *RenderStyles) HasRule(name string) bool {
	return rs.inline.HasRule(name) || rs.base.HasRule(name)
}

// GetRule returns the value from inline if set, otherwise from base.
func (rs *RenderStyles) GetRule(name string) (any, bool) {
	if v, ok := rs.inline.GetRule(name); ok {
		return v, true
	}
	return rs.base.GetRule(name)
}

// SetRule sets a rule on the inline Styles.
func (rs *RenderStyles) SetRule(name string, value any) bool {
	changed := rs.inline.SetRule(name, value)
	if changed {
		rs.updates++
	}
	return changed
}

// ClearRule removes a rule from the inline Styles. Returns true if it was set.
func (rs *RenderStyles) ClearRule(name string) bool {
	changed := rs.inline.ClearRule(name)
	if changed {
		rs.updates++
	}
	return changed
}

// GetRules returns the merged rules: base overridden by inline.
func (rs *RenderStyles) GetRules() RulesMap {
	result := rs.base.GetRules()
	for k, v := range rs.inline.Rules {
		result[k] = v
	}
	return result
}

// merged returns a Styles representing the merged view.
func (rs *RenderStyles) merged() *Styles {
	m := rs.base.Copy()
	m.Merge(rs.inline)
	return m
}

// CSSLines returns the serialized CSS declarations for the merged styles.
func (rs *RenderStyles) CSSLines() []string {
	return rs.merged().CSSLines()
}

// CSS returns the serialized CSS for the merged styles.
func (rs *RenderStyles) CSS() string {
	return strings.Join(rs.CSSLines(), "\n")
}

// Equal reports whether two RenderStyles are equivalent.
func (rs *RenderStyles) Equal(other *RenderStyles) bool {
	return rs.base.Equal(other.base) && rs.inline.Equal(other.inline)
}

// ── Delegated typed getters (inline first, then base) ──────────────────────

func (rs *RenderStyles) getStringEnum(name, def string) string {
	if v, ok := rs.inline.Rules[name]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	if v, ok := rs.base.Rules[name]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return def
}

func (rs *RenderStyles) Display() Display                     { return rs.merged().Display() }
func (rs *RenderStyles) Visibility() Visibility               { return rs.merged().Visibility() }
func (rs *RenderStyles) Layout() string                       { return rs.merged().Layout() }
func (rs *RenderStyles) Color() color.Color                   { return rs.merged().Color() }
func (rs *RenderStyles) Background() color.Color              { return rs.merged().Background() }
func (rs *RenderStyles) Opacity() float64                     { return rs.merged().Opacity() }
func (rs *RenderStyles) TextOpacity() float64                 { return rs.merged().TextOpacity() }
func (rs *RenderStyles) Padding() geometry.Spacing            { return rs.merged().Padding() }
func (rs *RenderStyles) Margin() geometry.Spacing             { return rs.merged().Margin() }
func (rs *RenderStyles) Offset() ScalarOffset                 { return rs.merged().Offset() }
func (rs *RenderStyles) Position() Position                   { return rs.merged().Position() }
func (rs *RenderStyles) Width() *Scalar                       { return rs.merged().Width() }
func (rs *RenderStyles) Height() *Scalar                      { return rs.merged().Height() }
func (rs *RenderStyles) BorderTop() EdgeStyle                 { return rs.merged().BorderTop() }
func (rs *RenderStyles) BorderRight() EdgeStyle               { return rs.merged().BorderRight() }
func (rs *RenderStyles) BorderBottom() EdgeStyle              { return rs.merged().BorderBottom() }
func (rs *RenderStyles) BorderLeft() EdgeStyle                { return rs.merged().BorderLeft() }
func (rs *RenderStyles) Gutter() geometry.Spacing             { return rs.merged().Gutter() }
func (rs *RenderStyles) AlignHorizontal() AlignHorizontal     { return rs.merged().AlignHorizontal() }
func (rs *RenderStyles) AlignVertical() AlignVertical         { return rs.merged().AlignVertical() }
func (rs *RenderStyles) Dock() string                         { return rs.merged().Dock() }
func (rs *RenderStyles) OverflowX() Overflow                  { return rs.merged().OverflowX() }
func (rs *RenderStyles) OverflowY() Overflow                  { return rs.merged().OverflowY() }
func (rs *RenderStyles) TextAlign() TextAlign                 { return rs.merged().TextAlign() }
func (rs *RenderStyles) BoxSizing() BoxSizing                 { return rs.merged().BoxSizing() }
func (rs *RenderStyles) Overlay() Overlay                     { return rs.merged().Overlay() }
func (rs *RenderStyles) TextWrap() TextWrap                   { return rs.merged().TextWrap() }
func (rs *RenderStyles) TextOverflow() TextOverflow           { return rs.merged().TextOverflow() }
func (rs *RenderStyles) GetTransition(key string) *Transition { return rs.merged().GetTransition(key) }
func (rs *RenderStyles) TextStyle() string                    { return rs.merged().TextStyle() }
func (rs *RenderStyles) LinkStyle() string                    { return rs.merged().LinkStyle() }
func (rs *RenderStyles) LinkStyleHover() string               { return rs.merged().LinkStyleHover() }
func (rs *RenderStyles) BorderTitleStyle() string             { return rs.merged().BorderTitleStyle() }
func (rs *RenderStyles) BorderSubtitleStyle() string          { return rs.merged().BorderSubtitleStyle() }
