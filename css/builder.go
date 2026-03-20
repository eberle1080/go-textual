package css

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/eberle1080/go-textual/color"
	"github.com/eberle1080/go-textual/geometry"
)

// StylesBuilder converts parsed CSS declarations into a Styles object.
type StylesBuilder struct {
	Styles *Styles
}

// NewStylesBuilder creates a new StylesBuilder.
func NewStylesBuilder() *StylesBuilder {
	return &StylesBuilder{Styles: NewStyles()}
}

// AddDeclaration processes a single Declaration and updates the Styles.
func (b *StylesBuilder) AddDeclaration(decl Declaration) error {
	if decl.Name == "" {
		return nil
	}
	ruleName := strings.ReplaceAll(decl.Name, "-", "_")
	if len(decl.Tokens) == 0 {
		return &DeclarationError{
			Name:    ruleName,
			Token:   decl.Token,
			Message: fmt.Sprintf("Missing property value for '%s:'", decl.Name),
		}
	}

	tokens := decl.Tokens

	// Check for !important
	if len(tokens) > 0 && tokens[len(tokens)-1].Name == "important" {
		tokens = tokens[:len(tokens)-1]
		b.Styles.Important[ruleName] = true
	}

	// Handle "initial" keyword
	if len(tokens) > 0 && tokens[0].Name == "token" && tokens[0].Value == "initial" {
		b.Styles.Rules[ruleName] = nil
		return nil
	}

	processors := b.processorMap()
	proc, ok := processors[ruleName]
	if !ok {
		return &DeclarationError{
			Name:    decl.Name,
			Token:   decl.Token,
			Message: fmt.Sprintf("unknown CSS property %q", decl.Name),
		}
	}

	if err := proc(decl.Name, tokens); err != nil {
		return &DeclarationError{Name: decl.Name, Token: decl.Token, Message: err.Error()}
	}
	return nil
}

// processorMap returns the dispatch table mapping rule names to processor functions.
func (b *StylesBuilder) processorMap() map[string]func(name string, tokens []Token) error {
	return map[string]func(name string, tokens []Token) error{
		"display":                   b.processDisplay,
		"visibility":                b.processVisibility,
		"layout":                    b.processLayout,
		"color":                     b.processColor,
		"background":                b.processColor,
		"background_tint":           b.processColor,
		"tint":                      b.processColor,
		"opacity":                   b.processFractional,
		"text_opacity":              b.processFractional,
		"padding":                   b.processSpace,
		"margin":                    b.processSpace,
		"padding_top":               b.processSpacePartial,
		"padding_right":             b.processSpacePartial,
		"padding_bottom":            b.processSpacePartial,
		"padding_left":              b.processSpacePartial,
		"margin_top":                b.processSpacePartial,
		"margin_right":              b.processSpacePartial,
		"margin_bottom":             b.processSpacePartial,
		"margin_left":               b.processSpacePartial,
		"border":                    b.processBorder,
		"border_top":                b.processBorderTop,
		"border_right":              b.processBorderRight,
		"border_bottom":             b.processBorderBottom,
		"border_left":               b.processBorderLeft,
		"outline":                   b.processOutline,
		"outline_top":               b.processOutlineTop,
		"outline_right":             b.processOutlineRight,
		"outline_bottom":            b.processOutlineBottom,
		"outline_left":              b.processOutlineLeft,
		"keyline":                   b.processKeyline,
		"offset":                    b.processOffset,
		"offset_x":                  b.processOffsetX,
		"offset_y":                  b.processOffsetY,
		"position":                  b.processPosition,
		"box_sizing":                b.processBoxSizing,
		"width":                     b.processScalar,
		"height":                    b.processScalar,
		"min_width":                 b.processScalar,
		"min_height":                b.processScalar,
		"max_width":                 b.processScalar,
		"max_height":                b.processScalar,
		"overflow":                  b.processOverflow,
		"overflow_x":                b.processOverflowX,
		"overflow_y":                b.processOverflowY,
		"dock":                      b.processDock,
		"split":                     b.processSplit,
		"layer":                     b.processLayer,
		"layers":                    b.processLayers,
		"transition":                b.processTransition,
		"align":                     b.processAlign,
		"align_horizontal":          b.processAlignHorizontal,
		"align_vertical":            b.processAlignVertical,
		"content_align":             b.processAlign,
		"content_align_horizontal":  b.processAlignHorizontal,
		"content_align_vertical":    b.processAlignVertical,
		"border_title_align":        b.processAlignHorizontal,
		"border_subtitle_align":     b.processAlignHorizontal,
		"scrollbar_color":           b.processColor,
		"scrollbar_color_hover":     b.processColor,
		"scrollbar_color_active":    b.processColor,
		"scrollbar_corner_color":    b.processColor,
		"scrollbar_background":      b.processColor,
		"scrollbar_background_hover":  b.processColor,
		"scrollbar_background_active": b.processColor,
		"scrollbar_gutter":          b.processScrollbarGutter,
		"scrollbar_size":            b.processScrollbarSize,
		"scrollbar_size_vertical":   b.processScrollbarSizeVertical,
		"scrollbar_size_horizontal": b.processScrollbarSizeHorizontal,
		"scrollbar_visibility":      b.processScrollbarVisibility,
		"grid_rows":                 b.processGridRowsOrColumns,
		"grid_columns":              b.processGridRowsOrColumns,
		"grid_size":                 b.processGridSize,
		"grid_size_rows":            b.processInteger,
		"grid_size_columns":         b.processInteger,
		"grid_gutter":               b.processGridGutter,
		"grid_gutter_horizontal":    b.processInteger,
		"grid_gutter_vertical":      b.processInteger,
		"row_span":                  b.processInteger,
		"column_span":               b.processInteger,
		"text_align":                b.processTextAlign,
		"text_style":                b.processTextStyle,
		"link_color":                b.processColor,
		"link_background":           b.processColor,
		"link_style":                b.processTextStyle,
		"link_color_hover":          b.processColor,
		"link_background_hover":     b.processColor,
		"link_style_hover":          b.processTextStyle,
		"border_title_color":        b.processColor,
		"border_title_background":   b.processColor,
		"border_title_style":        b.processTextStyle,
		"border_subtitle_color":     b.processColor,
		"border_subtitle_background": b.processColor,
		"border_subtitle_style":     b.processTextStyle,
		"hatch":                     b.processHatch,
		"overlay":                   b.processOverlay,
		"constrain":                 b.processConstrain,
		"constrain_x":               b.processConstrainX,
		"constrain_y":               b.processConstrainY,
		"text_wrap":                 b.processTextWrap,
		"text_overflow":             b.processTextOverflow,
		"expand":                    b.processExpand,
		"line_pad":                  b.processInteger,
		"pointer":                   b.processPointer,
	}
}

// ── Helpers ────────────────────────────────────────────────────────────────

func (b *StylesBuilder) errorf(name string, tok Token, format string, args ...any) error {
	return fmt.Errorf(format, args...)
}

func (b *StylesBuilder) processEnum(name string, tokens []Token, valid map[string]bool) (string, error) {
	if len(tokens) != 1 {
		return "", fmt.Errorf("expected 1 token for %s, got %d", name, len(tokens))
	}
	tok := tokens[0]
	if tok.Name != "token" {
		return "", fmt.Errorf("expected a token value for %s; found %q", name, tok.Value)
	}
	v := strings.ToLower(tok.Value)
	if !valid[v] {
		return "", fmt.Errorf("invalid value %q for %s", v, name)
	}
	return v, nil
}

func (b *StylesBuilder) processEnumMultiple(name string, tokens []Token, valid map[string]bool, count int) ([]string, error) {
	if len(tokens) > count || len(tokens) == 0 {
		return nil, fmt.Errorf("expected 1 to %d tokens for %s", count, name)
	}
	var results []string
	for _, tok := range tokens {
		if tok.Name != "token" {
			return nil, fmt.Errorf("invalid token %q for %s", tok.Value, name)
		}
		results = append(results, tok.Value)
	}
	// Pad
	short := append([]string{}, results...)
	for len(results) < count {
		results = append(results, short...)
	}
	return results[:count], nil
}

func (b *StylesBuilder) parseColor(name string, tokens []Token) (color.Color, error) {
	var c *color.Color
	var alpha *float64

	for _, tok := range tokens {
		switch tok.Name {
		case "token":
			if !strings.Contains(name, "background") && tok.Value == "auto" {
				b.Styles.Rules["auto_"+strings.ReplaceAll(name, "-", "_")] = true
				continue
			}
			parsed, err := color.Parse(tok.Value)
			if err != nil {
				return color.Color{}, fmt.Errorf("invalid color %q for %s: %v", tok.Value, name, err)
			}
			c = &parsed
		case "color":
			parsed, err := color.Parse(tok.Value)
			if err != nil {
				return color.Color{}, fmt.Errorf("invalid color %q for %s: %v", tok.Value, name, err)
			}
			c = &parsed
		case "scalar":
			sc, err := ParseScalar(tok.Value, UnitPercent)
			if err != nil {
				return color.Color{}, err
			}
			if sc.Unit != UnitPercent {
				return color.Color{}, fmt.Errorf("alpha must be given as a percentage for %s", name)
			}
			a := sc.Value / 100.0
			alpha = &a
		default:
			return color.Color{}, fmt.Errorf("unexpected token %q for color %s", tok.Value, name)
		}
	}

	if c == nil && alpha == nil {
		return color.Color{}, fmt.Errorf("no color value for %s", name)
	}
	var result color.Color
	if c != nil {
		result = *c
	} else {
		result = color.New(255, 255, 255)
	}
	if alpha != nil {
		result = result.MultiplyAlpha(*alpha)
	}
	return result, nil
}

func (b *StylesBuilder) parseBorderEdge(name string, tokens []Token) (EdgeStyle, error) {
	edgeType := EdgeType("solid")
	edgeColor := color.New(0, 255, 0)
	var alpha *float64

	for _, tok := range tokens {
		switch tok.Name {
		case "token":
			if ValidBorder[tok.Value] {
				edgeType = EdgeType(tok.Value)
			} else {
				parsed, err := color.Parse(tok.Value)
				if err != nil {
					return EdgeStyle{}, fmt.Errorf("unknown border type %q for %s", tok.Value, name)
				}
				edgeColor = parsed
			}
		case "color":
			parsed, err := color.Parse(tok.Value)
			if err != nil {
				return EdgeStyle{}, fmt.Errorf("invalid color %q for %s: %v", tok.Value, name, err)
			}
			edgeColor = parsed
		case "scalar":
			sc, err := ParseScalar(tok.Value, UnitPercent)
			if err != nil {
				return EdgeStyle{}, err
			}
			if sc.Unit != UnitPercent {
				return EdgeStyle{}, fmt.Errorf("border alpha must be a percentage for %s", name)
			}
			a := sc.Value / 100.0
			alpha = &a
		default:
			return EdgeStyle{}, fmt.Errorf("unexpected token %q for border %s", tok.Value, name)
		}
	}

	if alpha != nil {
		edgeColor = edgeColor.MultiplyAlpha(*alpha)
	}
	return EdgeStyle{Type: edgeType, Color: edgeColor}, nil
}

func (b *StylesBuilder) distributeImportance(prefix string, suffixes []string) {
	if b.Styles.Important[prefix] {
		delete(b.Styles.Important, prefix)
		for _, suf := range suffixes {
			b.Styles.Important[prefix+"_"+suf] = true
		}
	}
}

func durationAsSeconds(s string) (float64, error) {
	s = strings.TrimSpace(s)
	if strings.HasSuffix(s, "ms") {
		v, err := strconv.ParseFloat(s[:len(s)-2], 64)
		if err != nil {
			return 0, err
		}
		return v / 1000.0, nil
	}
	if strings.HasSuffix(s, "s") {
		v, err := strconv.ParseFloat(s[:len(s)-1], 64)
		if err != nil {
			return 0, err
		}
		return v, nil
	}
	return strconv.ParseFloat(s, 64)
}

// ── Processor methods ──────────────────────────────────────────────────────

func (b *StylesBuilder) processDisplay(name string, tokens []Token) error {
	for _, tok := range tokens {
		if tok.Name != "token" {
			return fmt.Errorf("expected token for display, got %q", tok.Value)
		}
		v := strings.ToLower(tok.Value)
		if !ValidDisplay[v] {
			return fmt.Errorf("invalid display value %q", v)
		}
		b.Styles.Rules["display"] = Display(v)
	}
	return nil
}

func (b *StylesBuilder) processVisibility(name string, tokens []Token) error {
	for _, tok := range tokens {
		if tok.Name != "token" {
			return fmt.Errorf("expected token for visibility")
		}
		v := strings.ToLower(tok.Value)
		if !ValidVisibility[v] {
			return fmt.Errorf("invalid visibility value %q", v)
		}
		b.Styles.Rules["visibility"] = Visibility(v)
	}
	return nil
}

func (b *StylesBuilder) processLayout(name string, tokens []Token) error {
	if len(tokens) == 0 {
		return nil
	}
	if len(tokens) != 1 {
		return fmt.Errorf("expected 1 token for layout")
	}
	v := tokens[0].Value
	if !ValidLayout[v] {
		return fmt.Errorf("invalid layout %q", v)
	}
	b.Styles.Rules["layout"] = v
	return nil
}

func (b *StylesBuilder) processColor(name string, tokens []Token) error {
	ruleName := strings.ReplaceAll(name, "-", "_")
	b.Styles.Rules["auto_"+ruleName] = false
	c, err := b.parseColor(name, tokens)
	if err != nil {
		return err
	}
	b.Styles.Rules[ruleName] = c
	return nil
}

func (b *StylesBuilder) processFractional(name string, tokens []Token) error {
	if len(tokens) != 1 {
		return fmt.Errorf("expected 1 token for %s", name)
	}
	tok := tokens[0]
	ruleName := strings.ReplaceAll(name, "-", "_")
	if tok.Name == "scalar" && strings.HasSuffix(tok.Value, "%") {
		f, err := PercentageStringToFloat(tok.Value)
		if err != nil {
			return err
		}
		b.Styles.Rules[ruleName] = f
	} else if tok.Name == "number" {
		f, err := strconv.ParseFloat(tok.Value, 64)
		if err != nil {
			return err
		}
		if f < 0 {
			f = 0
		}
		if f > 1 {
			f = 1
		}
		b.Styles.Rules[ruleName] = f
	} else {
		return fmt.Errorf("expected a number or percentage for %s", name)
	}
	return nil
}

func (b *StylesBuilder) processSpace(name string, tokens []Token) error {
	var space []int
	for _, tok := range tokens {
		if tok.Name != "number" {
			return fmt.Errorf("expected number for spacing %s, got %q", name, tok.Value)
		}
		v, err := strconv.Atoi(tok.Value)
		if err != nil {
			return err
		}
		space = append(space, v)
	}
	if len(space) != 1 && len(space) != 2 && len(space) != 4 {
		return fmt.Errorf("expected 1, 2, or 4 values for spacing %s, got %d", name, len(space))
	}
	sp := unpackSpacing(space)
	b.Styles.Rules[strings.ReplaceAll(name, "-", "_")] = sp
	return nil
}

func (b *StylesBuilder) processSpacePartial(name string, tokens []Token) error {
	if len(tokens) != 1 {
		return fmt.Errorf("expected 1 value for %s", name)
	}
	tok := tokens[0]
	if tok.Name != "number" {
		return fmt.Errorf("expected number for %s", name)
	}
	v, err := strconv.Atoi(tok.Value)
	if err != nil {
		return err
	}
	edgeMap := map[string]int{"top": 0, "right": 1, "bottom": 2, "left": 3}
	ruleName := strings.ReplaceAll(name, "-", "_")
	parts := strings.SplitN(ruleName, "_", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid partial spacing name %s", name)
	}
	styleName, edge := parts[0], parts[1]
	idx, ok := edgeMap[edge]
	if !ok {
		return fmt.Errorf("unknown edge %q in %s", edge, name)
	}
	current := geometry.Spacing{}
	if cv, ok := b.Styles.Rules[styleName]; ok {
		if sp, ok := cv.(geometry.Spacing); ok {
			current = sp
		}
	}
	arr := [4]int{current.Top, current.Right, current.Bottom, current.Left}
	arr[idx] = v
	b.Styles.Rules[styleName] = geometry.Spacing{Top: arr[0], Right: arr[1], Bottom: arr[2], Left: arr[3]}
	return nil
}

func unpackSpacing(vals []int) geometry.Spacing {
	switch len(vals) {
	case 1:
		return geometry.Spacing{Top: vals[0], Right: vals[0], Bottom: vals[0], Left: vals[0]}
	case 2:
		return geometry.Spacing{Top: vals[0], Right: vals[1], Bottom: vals[0], Left: vals[1]}
	case 4:
		return geometry.Spacing{Top: vals[0], Right: vals[1], Bottom: vals[2], Left: vals[3]}
	default:
		return geometry.Spacing{}
	}
}

func (b *StylesBuilder) processBorder(name string, tokens []Token) error {
	edge, err := b.parseBorderEdge(name, tokens)
	if err != nil {
		return err
	}
	b.Styles.Rules["border_top"] = edge
	b.Styles.Rules["border_right"] = edge
	b.Styles.Rules["border_bottom"] = edge
	b.Styles.Rules["border_left"] = edge
	b.distributeImportance("border", []string{"top", "left", "bottom", "right"})
	return nil
}

func (b *StylesBuilder) processBorderTop(name string, tokens []Token) error {
	edge, err := b.parseBorderEdge(name, tokens)
	if err != nil {
		return err
	}
	b.Styles.Rules["border_top"] = edge
	return nil
}

func (b *StylesBuilder) processBorderRight(name string, tokens []Token) error {
	edge, err := b.parseBorderEdge(name, tokens)
	if err != nil {
		return err
	}
	b.Styles.Rules["border_right"] = edge
	return nil
}

func (b *StylesBuilder) processBorderBottom(name string, tokens []Token) error {
	edge, err := b.parseBorderEdge(name, tokens)
	if err != nil {
		return err
	}
	b.Styles.Rules["border_bottom"] = edge
	return nil
}

func (b *StylesBuilder) processBorderLeft(name string, tokens []Token) error {
	edge, err := b.parseBorderEdge(name, tokens)
	if err != nil {
		return err
	}
	b.Styles.Rules["border_left"] = edge
	return nil
}

func (b *StylesBuilder) processOutline(name string, tokens []Token) error {
	edge, err := b.parseBorderEdge(name, tokens)
	if err != nil {
		return err
	}
	b.Styles.Rules["outline_top"] = edge
	b.Styles.Rules["outline_right"] = edge
	b.Styles.Rules["outline_bottom"] = edge
	b.Styles.Rules["outline_left"] = edge
	b.distributeImportance("outline", []string{"top", "left", "bottom", "right"})
	return nil
}

func (b *StylesBuilder) processOutlineTop(name string, tokens []Token) error {
	edge, err := b.parseBorderEdge(name, tokens)
	if err != nil {
		return err
	}
	b.Styles.Rules["outline_top"] = edge
	return nil
}

func (b *StylesBuilder) processOutlineRight(name string, tokens []Token) error {
	edge, err := b.parseBorderEdge(name, tokens)
	if err != nil {
		return err
	}
	b.Styles.Rules["outline_right"] = edge
	return nil
}

func (b *StylesBuilder) processOutlineBottom(name string, tokens []Token) error {
	edge, err := b.parseBorderEdge(name, tokens)
	if err != nil {
		return err
	}
	b.Styles.Rules["outline_bottom"] = edge
	return nil
}

func (b *StylesBuilder) processOutlineLeft(name string, tokens []Token) error {
	edge, err := b.parseBorderEdge(name, tokens)
	if err != nil {
		return err
	}
	b.Styles.Rules["outline_left"] = edge
	return nil
}

func (b *StylesBuilder) processKeyline(name string, tokens []Token) error {
	if len(tokens) == 0 {
		return nil
	}
	keylineStyle := EdgeType("none")
	keylineColor := color.New(0, 128, 0)
	var alpha *float64

	for _, tok := range tokens {
		switch tok.Name {
		case "color":
			c, err := color.Parse(tok.Value)
			if err != nil {
				return err
			}
			keylineColor = c
		case "token":
			if ValidKeyline[tok.Value] {
				keylineStyle = EdgeType(tok.Value)
			} else {
				c, err := color.Parse(tok.Value)
				if err != nil {
					return fmt.Errorf("unknown keyline value %q", tok.Value)
				}
				keylineColor = c
			}
		case "scalar":
			sc, err := ParseScalar(tok.Value, UnitPercent)
			if err != nil {
				return err
			}
			a := sc.Value / 100.0
			alpha = &a
		}
	}

	if alpha != nil {
		keylineColor = keylineColor.MultiplyAlpha(*alpha)
	}
	b.Styles.Rules["keyline"] = EdgeStyle{Type: keylineStyle, Color: keylineColor}
	return nil
}

func (b *StylesBuilder) processOffset(name string, tokens []Token) error {
	if len(tokens) == 0 {
		return nil
	}
	if len(tokens) != 2 {
		return fmt.Errorf("offset requires exactly 2 values")
	}
	tok1, tok2 := tokens[0], tokens[1]
	if tok1.Name != "scalar" && tok1.Name != "number" {
		return fmt.Errorf("expected scalar or number for offset x")
	}
	if tok2.Name != "scalar" && tok2.Name != "number" {
		return fmt.Errorf("expected scalar or number for offset y")
	}
	x, err := ParseScalar(tok1.Value, UnitWidth)
	if err != nil {
		return err
	}
	y, err := ParseScalar(tok2.Value, UnitHeight)
	if err != nil {
		return err
	}
	b.Styles.Rules["offset"] = ScalarOffset{X: x, Y: y}
	return nil
}

func (b *StylesBuilder) processOffsetX(name string, tokens []Token) error {
	if len(tokens) != 1 {
		return fmt.Errorf("offset-x requires 1 value")
	}
	tok := tokens[0]
	if tok.Name != "scalar" && tok.Name != "number" {
		return fmt.Errorf("expected scalar or number for offset-x")
	}
	x, err := ParseScalar(tok.Value, UnitWidth)
	if err != nil {
		return err
	}
	var y Scalar
	if cv, ok := b.Styles.Rules["offset"]; ok {
		if so, ok := cv.(ScalarOffset); ok {
			y = so.Y
		}
	}
	b.Styles.Rules["offset"] = ScalarOffset{X: x, Y: y}
	return nil
}

func (b *StylesBuilder) processOffsetY(name string, tokens []Token) error {
	if len(tokens) != 1 {
		return fmt.Errorf("offset-y requires 1 value")
	}
	tok := tokens[0]
	if tok.Name != "scalar" && tok.Name != "number" {
		return fmt.Errorf("expected scalar or number for offset-y")
	}
	y, err := ParseScalar(tok.Value, UnitHeight)
	if err != nil {
		return err
	}
	var x Scalar
	if cv, ok := b.Styles.Rules["offset"]; ok {
		if so, ok := cv.(ScalarOffset); ok {
			x = so.X
		}
	}
	b.Styles.Rules["offset"] = ScalarOffset{X: x, Y: y}
	return nil
}

func (b *StylesBuilder) processPosition(name string, tokens []Token) error {
	if len(tokens) != 1 {
		return fmt.Errorf("expected 1 token for position")
	}
	v := tokens[0].Value
	if !ValidPosition[v] {
		return fmt.Errorf("invalid position %q", v)
	}
	b.Styles.Rules["position"] = Position(v)
	return nil
}

func (b *StylesBuilder) processBoxSizing(name string, tokens []Token) error {
	for _, tok := range tokens {
		if tok.Name != "token" {
			return fmt.Errorf("expected token for box-sizing")
		}
		v := strings.ToLower(tok.Value)
		if !ValidBoxSizing[v] {
			return fmt.Errorf("invalid box-sizing %q", v)
		}
		b.Styles.Rules["box_sizing"] = BoxSizing(v)
	}
	return nil
}

func (b *StylesBuilder) processScalar(name string, tokens []Token) error {
	if len(tokens) != 1 {
		return fmt.Errorf("expected 1 token for %s", name)
	}
	ruleName := strings.ReplaceAll(name, "-", "_")
	sc, err := ParseScalar(tokens[0].Value, UnitWidth)
	if err != nil {
		return err
	}
	b.Styles.Rules[ruleName] = sc
	return nil
}

func (b *StylesBuilder) processOverflow(name string, tokens []Token) error {
	results, err := b.processEnumMultiple(name, tokens, ValidOverflow, 2)
	if err != nil {
		return err
	}
	b.Styles.Rules["overflow_x"] = Overflow(results[0])
	b.Styles.Rules["overflow_y"] = Overflow(results[1])
	b.distributeImportance("overflow", []string{"x", "y"})
	return nil
}

func (b *StylesBuilder) processOverflowX(name string, tokens []Token) error {
	v, err := b.processEnum(name, tokens, ValidOverflow)
	if err != nil {
		return err
	}
	b.Styles.Rules["overflow_x"] = Overflow(v)
	return nil
}

func (b *StylesBuilder) processOverflowY(name string, tokens []Token) error {
	v, err := b.processEnum(name, tokens, ValidOverflow)
	if err != nil {
		return err
	}
	b.Styles.Rules["overflow_y"] = Overflow(v)
	return nil
}

func (b *StylesBuilder) processDock(name string, tokens []Token) error {
	if len(tokens) == 0 {
		return nil
	}
	if len(tokens) > 1 || !ValidEdge[tokens[0].Value] {
		return fmt.Errorf("invalid dock value")
	}
	b.Styles.Rules["dock"] = tokens[0].Value
	return nil
}

func (b *StylesBuilder) processSplit(name string, tokens []Token) error {
	if len(tokens) == 0 {
		return nil
	}
	if len(tokens) > 1 || !ValidEdge[tokens[0].Value] {
		return fmt.Errorf("invalid split value")
	}
	b.Styles.Rules["split"] = tokens[0].Value
	return nil
}

func (b *StylesBuilder) processLayer(name string, tokens []Token) error {
	if len(tokens) == 0 {
		return nil
	}
	b.Styles.Rules["layer"] = tokens[0].Value
	return nil
}

func (b *StylesBuilder) processLayers(name string, tokens []Token) error {
	var layers []string
	for _, tok := range tokens {
		layers = append(layers, tok.Value)
	}
	b.Styles.Rules["layers"] = layers
	return nil
}

func (b *StylesBuilder) processTransition(name string, tokens []Token) error {
	transitions := make(map[string]Transition)

	// Split tokens by comma
	var groups [][]Token
	var current []Token
	for _, tok := range tokens {
		if tok.Name == "comma" {
			if len(current) > 0 {
				groups = append(groups, current)
				current = nil
			}
		} else {
			current = append(current, tok)
		}
	}
	if len(current) > 0 {
		groups = append(groups, current)
	}

	for _, group := range groups {
		property := ""
		duration := 1.0
		easing := "linear"
		delay := 0.0

		i := 0
		if i < len(group) {
			if group[i].Name == "token" {
				property = group[i].Value
				i++
			} else {
				return fmt.Errorf("expected property name in transition")
			}
		}
		if i < len(group) {
			d, err := durationAsSeconds(group[i].Value)
			if err != nil {
				return fmt.Errorf("expected duration in transition: %v", err)
			}
			duration = d
			i++
		}
		if i < len(group) {
			easing = group[i].Value
			i++
		}
		if i < len(group) {
			d, err := durationAsSeconds(group[i].Value)
			if err == nil {
				delay = d
			}
		}

		transitions[property] = NewTransition(duration, easing, delay)
	}

	b.Styles.Rules["transitions"] = transitions
	return nil
}

func (b *StylesBuilder) processAlign(name string, tokens []Token) error {
	if len(tokens) != 2 {
		return fmt.Errorf("align requires 2 values")
	}
	tokH, tokV := tokens[0], tokens[1]
	if !ValidAlignHorizontal[tokH.Value] {
		return fmt.Errorf("invalid horizontal align %q", tokH.Value)
	}
	if !ValidAlignVertical[tokV.Value] {
		return fmt.Errorf("invalid vertical align %q", tokV.Value)
	}
	ruleName := strings.ReplaceAll(name, "-", "_")
	b.Styles.Rules[ruleName+"_horizontal"] = AlignHorizontal(tokH.Value)
	b.Styles.Rules[ruleName+"_vertical"] = AlignVertical(tokV.Value)
	b.distributeImportance(ruleName, []string{"horizontal", "vertical"})
	return nil
}

func (b *StylesBuilder) processAlignHorizontal(name string, tokens []Token) error {
	v, err := b.processEnum(name, tokens, ValidAlignHorizontal)
	if err != nil {
		return err
	}
	b.Styles.Rules[strings.ReplaceAll(name, "-", "_")] = AlignHorizontal(v)
	return nil
}

func (b *StylesBuilder) processAlignVertical(name string, tokens []Token) error {
	v, err := b.processEnum(name, tokens, ValidAlignVertical)
	if err != nil {
		return err
	}
	b.Styles.Rules[strings.ReplaceAll(name, "-", "_")] = AlignVertical(v)
	return nil
}

func (b *StylesBuilder) processScrollbarGutter(name string, tokens []Token) error {
	v, err := b.processEnum(name, tokens, ValidScrollbarGutter)
	if err != nil {
		return err
	}
	b.Styles.Rules["scrollbar_gutter"] = ScrollbarGutter(v)
	return nil
}

func (b *StylesBuilder) processScrollbarSize(name string, tokens []Token) error {
	if len(tokens) != 2 {
		return fmt.Errorf("scrollbar-size requires 2 values")
	}
	h, err := strconv.Atoi(tokens[0].Value)
	if err != nil {
		return fmt.Errorf("invalid scrollbar-size horizontal: %v", err)
	}
	v, err := strconv.Atoi(tokens[1].Value)
	if err != nil {
		return fmt.Errorf("invalid scrollbar-size vertical: %v", err)
	}
	b.Styles.Rules["scrollbar_size_horizontal"] = h
	b.Styles.Rules["scrollbar_size_vertical"] = v
	b.distributeImportance("scrollbar_size", []string{"horizontal", "vertical"})
	return nil
}

func (b *StylesBuilder) processScrollbarSizeVertical(name string, tokens []Token) error {
	if len(tokens) != 1 {
		return fmt.Errorf("expected 1 value for scrollbar-size-vertical")
	}
	v, err := strconv.Atoi(tokens[0].Value)
	if err != nil {
		return err
	}
	b.Styles.Rules["scrollbar_size_vertical"] = v
	return nil
}

func (b *StylesBuilder) processScrollbarSizeHorizontal(name string, tokens []Token) error {
	if len(tokens) != 1 {
		return fmt.Errorf("expected 1 value for scrollbar-size-horizontal")
	}
	v, err := strconv.Atoi(tokens[0].Value)
	if err != nil {
		return err
	}
	b.Styles.Rules["scrollbar_size_horizontal"] = v
	return nil
}

func (b *StylesBuilder) processScrollbarVisibility(name string, tokens []Token) error {
	v, err := b.processEnum(name, tokens, ValidScrollbarVisibility)
	if err != nil {
		return err
	}
	b.Styles.Rules["scrollbar_visibility"] = ScrollbarVisibility(v)
	return nil
}

func (b *StylesBuilder) processGridRowsOrColumns(name string, tokens []Token) error {
	ruleName := strings.ReplaceAll(name, "-", "_")
	percentUnit := UnitWidth
	if ruleName == "grid_rows" {
		percentUnit = UnitHeight
	}
	var scalars []Scalar
	for _, tok := range tokens {
		switch tok.Name {
		case "number":
			f, err := strconv.ParseFloat(tok.Value, 64)
			if err != nil {
				return err
			}
			scalars = append(scalars, FromNumber(f))
		case "scalar":
			sc, err := ParseScalar(tok.Value, percentUnit)
			if err != nil {
				return err
			}
			scalars = append(scalars, sc)
		case "token":
			if tok.Value == "auto" {
				sc, _ := ParseScalar("auto", percentUnit)
				scalars = append(scalars, sc)
			} else {
				return fmt.Errorf("unexpected token %q in %s", tok.Value, name)
			}
		default:
			return fmt.Errorf("unexpected token %q in %s", tok.Value, name)
		}
	}
	b.Styles.Rules[ruleName] = scalars
	return nil
}

func (b *StylesBuilder) processInteger(name string, tokens []Token) error {
	if len(tokens) != 1 {
		return fmt.Errorf("expected 1 value for %s", name)
	}
	tok := tokens[0]
	v, err := strconv.Atoi(tok.Value)
	if err != nil {
		return fmt.Errorf("expected integer for %s: %v", name, err)
	}
	b.Styles.Rules[strings.ReplaceAll(name, "-", "_")] = v
	return nil
}

func (b *StylesBuilder) processGridGutter(name string, tokens []Token) error {
	if len(tokens) == 1 {
		v, err := strconv.Atoi(tokens[0].Value)
		if err != nil {
			return err
		}
		if v < 0 {
			v = 0
		}
		b.Styles.Rules["grid_gutter_horizontal"] = v
		b.Styles.Rules["grid_gutter_vertical"] = v
	} else if len(tokens) == 2 {
		h, err := strconv.Atoi(tokens[0].Value)
		if err != nil {
			return err
		}
		vv, err := strconv.Atoi(tokens[1].Value)
		if err != nil {
			return err
		}
		b.Styles.Rules["grid_gutter_horizontal"] = h
		b.Styles.Rules["grid_gutter_vertical"] = vv
	} else {
		return fmt.Errorf("expected 1 or 2 values for grid-gutter")
	}
	return nil
}

func (b *StylesBuilder) processGridSize(name string, tokens []Token) error {
	if len(tokens) == 1 {
		v, err := strconv.Atoi(tokens[0].Value)
		if err != nil {
			return err
		}
		b.Styles.Rules["grid_size_columns"] = v
		b.Styles.Rules["grid_size_rows"] = 0
	} else if len(tokens) == 2 {
		cols, err := strconv.Atoi(tokens[0].Value)
		if err != nil {
			return err
		}
		rows, err := strconv.Atoi(tokens[1].Value)
		if err != nil {
			return err
		}
		b.Styles.Rules["grid_size_columns"] = cols
		b.Styles.Rules["grid_size_rows"] = rows
	} else {
		return fmt.Errorf("expected 1 or 2 values for grid-size")
	}
	return nil
}

func (b *StylesBuilder) processTextAlign(name string, tokens []Token) error {
	if len(tokens) == 0 {
		return nil
	}
	if len(tokens) > 1 || !ValidTextAlign[tokens[0].Value] {
		return fmt.Errorf("invalid text-align %q", tokens[0].Value)
	}
	b.Styles.Rules["text_align"] = TextAlign(tokens[0].Value)
	return nil
}

func (b *StylesBuilder) processTextStyle(name string, tokens []Token) error {
	for _, tok := range tokens {
		if !ValidStyleFlags[tok.Value] {
			return fmt.Errorf("invalid style flag %q", tok.Value)
		}
	}
	var parts []string
	for _, tok := range tokens {
		parts = append(parts, tok.Value)
	}
	b.Styles.Rules[strings.ReplaceAll(name, "-", "_")] = strings.Join(parts, " ")
	return nil
}

func (b *StylesBuilder) processTextWrap(name string, tokens []Token) error {
	for _, tok := range tokens {
		if tok.Name != "token" {
			return fmt.Errorf("expected token for text-wrap")
		}
		v := strings.ToLower(tok.Value)
		if !ValidTextWrap[v] {
			return fmt.Errorf("invalid text-wrap %q", v)
		}
		b.Styles.Rules["text_wrap"] = TextWrap(v)
	}
	return nil
}

func (b *StylesBuilder) processTextOverflow(name string, tokens []Token) error {
	for _, tok := range tokens {
		if tok.Name != "token" {
			return fmt.Errorf("expected token for text-overflow")
		}
		v := strings.ToLower(tok.Value)
		if !ValidTextOverflow[v] {
			return fmt.Errorf("invalid text-overflow %q", v)
		}
		b.Styles.Rules["text_overflow"] = TextOverflow(v)
	}
	return nil
}

func (b *StylesBuilder) processHatch(name string, tokens []Token) error {
	if len(tokens) == 0 {
		return nil
	}
	if len(tokens) == 1 && tokens[0].Value == "none" {
		b.Styles.Rules["hatch"] = EdgeStyle{Type: "none"}
		return nil
	}
	if len(tokens) < 2 || len(tokens) > 3 {
		return fmt.Errorf("hatch requires 2 or 3 values")
	}
	charToken := tokens[0]
	colorToken := tokens[1]

	var character string
	if charToken.Name == "token" {
		if !ValidHatch[charToken.Value] {
			return fmt.Errorf("invalid hatch type %q", charToken.Value)
		}
		character = Hatches[charToken.Value]
	} else if charToken.Name == "string" {
		character = charToken.Value[1 : len(charToken.Value)-1]
	} else {
		return fmt.Errorf("expected token or string for hatch character")
	}

	c, err := color.Parse(colorToken.Value)
	if err != nil {
		return fmt.Errorf("invalid hatch color: %v", err)
	}

	opacity := 1.0
	if len(tokens) == 3 {
		sc, err := ParseScalar(tokens[2].Value, UnitPercent)
		if err != nil {
			return err
		}
		opacity = sc.Value / 100.0
		if opacity < 0 {
			opacity = 0
		}
		if opacity > 1 {
			opacity = 1
		}
	}
	c = c.MultiplyAlpha(opacity)
	b.Styles.Rules["hatch"] = EdgeStyle{Type: EdgeType(character), Color: c}
	return nil
}

func (b *StylesBuilder) processOverlay(name string, tokens []Token) error {
	v, err := b.processEnum(name, tokens, ValidOverlay)
	if err != nil {
		return err
	}
	b.Styles.Rules["overlay"] = Overlay(v)
	return nil
}

func (b *StylesBuilder) processConstrain(name string, tokens []Token) error {
	if len(tokens) == 1 {
		v, err := b.processEnum(name, tokens, ValidConstrain)
		if err != nil {
			return err
		}
		b.Styles.Rules["constrain_x"] = Constrain(v)
		b.Styles.Rules["constrain_y"] = Constrain(v)
	} else if len(tokens) == 2 {
		results, err := b.processEnumMultiple(name, tokens, ValidConstrain, 2)
		if err != nil {
			return err
		}
		b.Styles.Rules["constrain_x"] = Constrain(results[0])
		b.Styles.Rules["constrain_y"] = Constrain(results[1])
	} else {
		return fmt.Errorf("expected 1 or 2 values for constrain")
	}
	return nil
}

func (b *StylesBuilder) processConstrainX(name string, tokens []Token) error {
	v, err := b.processEnum(name, tokens, ValidConstrain)
	if err != nil {
		return err
	}
	b.Styles.Rules["constrain_x"] = Constrain(v)
	return nil
}

func (b *StylesBuilder) processConstrainY(name string, tokens []Token) error {
	v, err := b.processEnum(name, tokens, ValidConstrain)
	if err != nil {
		return err
	}
	b.Styles.Rules["constrain_y"] = Constrain(v)
	return nil
}

func (b *StylesBuilder) processExpand(name string, tokens []Token) error {
	if len(tokens) != 1 {
		return fmt.Errorf("expected 1 value for expand")
	}
	if !ValidExpand[tokens[0].Value] {
		return fmt.Errorf("invalid expand value %q", tokens[0].Value)
	}
	b.Styles.Rules["expand"] = Expand(tokens[0].Value)
	return nil
}

func (b *StylesBuilder) processPointer(name string, tokens []Token) error {
	for _, tok := range tokens {
		if tok.Name != "token" {
			return fmt.Errorf("expected token for pointer")
		}
		v := strings.ToLower(tok.Value)
		if !ValidPointer[v] {
			return fmt.Errorf("invalid pointer value %q", v)
		}
		b.Styles.Rules["pointer"] = PointerShape(v)
	}
	return nil
}
