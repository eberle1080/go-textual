package css

// ValidVisibility is the set of valid visibility values.
var ValidVisibility = map[string]bool{"visible": true, "hidden": true}

// ValidDisplay is the set of valid display values.
var ValidDisplay = map[string]bool{"block": true, "none": true}

// ValidBorder is the set of valid border type names.
var ValidBorder = map[string]bool{
	"ascii": true, "blank": true, "dashed": true, "double": true,
	"heavy": true, "hidden": true, "hkey": true, "inner": true,
	"none": true, "outer": true, "panel": true, "round": true,
	"solid": true, "tall": true, "tab": true, "thick": true,
	"block": true, "vkey": true, "wide": true,
}

// ValidEdge is the set of valid edge values.
var ValidEdge = map[string]bool{"top": true, "right": true, "bottom": true, "left": true, "none": true}

// ValidLayout is the set of valid layout values.
var ValidLayout = map[string]bool{"vertical": true, "horizontal": true, "grid": true, "stream": true}

// ValidBoxSizing is the set of valid box-sizing values.
var ValidBoxSizing = map[string]bool{"border-box": true, "content-box": true}

// ValidOverflow is the set of valid overflow values.
var ValidOverflow = map[string]bool{"scroll": true, "hidden": true, "auto": true}

// ValidAlignHorizontal is the set of valid horizontal alignment values.
var ValidAlignHorizontal = map[string]bool{"left": true, "center": true, "right": true}

// ValidAlignVertical is the set of valid vertical alignment values.
var ValidAlignVertical = map[string]bool{"top": true, "middle": true, "bottom": true}

// ValidPosition is the set of valid position values.
var ValidPosition = map[string]bool{"relative": true, "absolute": true}

// ValidTextAlign is the set of valid text-align values.
var ValidTextAlign = map[string]bool{
	"start": true, "end": true, "left": true, "right": true,
	"center": true, "justify": true,
}

// ValidScrollbarGutter is the set of valid scrollbar-gutter values.
var ValidScrollbarGutter = map[string]bool{"auto": true, "stable": true}

// ValidStyleFlags is the set of valid text style flags.
var ValidStyleFlags = map[string]bool{
	"b": true, "blink": true, "bold": true, "dim": true,
	"i": true, "italic": true, "none": true, "not": true,
	"o": true, "overline": true, "reverse": true, "strike": true,
	"u": true, "underline": true, "uu": true,
}

// ValidPseudoClasses is the set of valid CSS pseudo-class names.
var ValidPseudoClasses = map[string]bool{
	"ansi": true, "blur": true, "can-focus": true, "dark": true,
	"disabled": true, "enabled": true, "focus-within": true, "focus": true,
	"hover": true, "inline": true, "light": true, "nocolor": true,
	"first-of-type": true, "last-of-type": true, "first-child": true,
	"last-child": true, "odd": true, "even": true, "empty": true,
}

// ValidOverlay is the set of valid overlay values.
var ValidOverlay = map[string]bool{"none": true, "screen": true}

// ValidConstrain is the set of valid constrain values.
var ValidConstrain = map[string]bool{"inflect": true, "inside": true, "none": true}

// ValidKeyline is the set of valid keyline values.
var ValidKeyline = map[string]bool{"none": true, "thin": true, "heavy": true, "double": true}

// ValidHatch is the set of valid hatch values.
var ValidHatch = map[string]bool{"left": true, "right": true, "cross": true, "vertical": true, "horizontal": true}

// ValidTextWrap is the set of valid text-wrap values.
var ValidTextWrap = map[string]bool{"wrap": true, "nowrap": true}

// ValidTextOverflow is the set of valid text-overflow values.
var ValidTextOverflow = map[string]bool{"clip": true, "fold": true, "ellipsis": true}

// ValidExpand is the set of valid expand values.
var ValidExpand = map[string]bool{"greedy": true, "optimal": true}

// ValidScrollbarVisibility is the set of valid scrollbar-visibility values.
var ValidScrollbarVisibility = map[string]bool{"visible": true, "hidden": true}

// ValidPointer is the set of valid pointer shape values.
var ValidPointer = map[string]bool{
	"alias": true, "cell": true, "copy": true, "crosshair": true,
	"default": true, "e-resize": true, "ew-resize": true, "grab": true,
	"grabbing": true, "help": true, "move": true, "n-resize": true,
	"ne-resize": true, "nesw-resize": true, "no-drop": true, "not-allowed": true,
	"ns-resize": true, "nw-resize": true, "nwse-resize": true, "pointer": true,
	"progress": true, "s-resize": true, "se-resize": true, "sw-resize": true,
	"text": true, "vertical-text": true, "w-resize": true, "wait": true,
	"zoom-in": true, "zoom-out": true,
}

// Hatches maps hatch names to their Unicode characters.
var Hatches = map[string]string{
	"left":       "╲",
	"right":      "╱",
	"cross":      "╳",
	"horizontal": "─",
	"vertical":   "│",
}
