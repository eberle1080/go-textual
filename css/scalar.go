package css

import (
	"fmt"
	"math/big"
	"regexp"
	"strconv"
	"strings"

	"github.com/eberle1080/go-textual/geometry"
)

// Unit is the measurement unit of a Scalar.
type Unit int

const (
	UnitCells      Unit = iota + 1 // Explicit cell count (no suffix)
	UnitFraction                   // Fractional unit "fr"
	UnitPercent                    // Percentage "%"
	UnitWidth                      // Width-relative "w"
	UnitHeight                     // Height-relative "h"
	UnitViewWidth                  // Viewport-width "vw"
	UnitViewHeight                 // Viewport-height "vh"
	UnitAuto                       // Automatic "auto"
)

// UnitSymbol maps each Unit to its CSS suffix string.
var UnitSymbol = map[Unit]string{
	UnitCells:      "",
	UnitFraction:   "fr",
	UnitPercent:    "%",
	UnitWidth:      "w",
	UnitHeight:     "h",
	UnitViewWidth:  "vw",
	UnitViewHeight: "vh",
}

// SymbolUnit maps CSS suffix strings back to Unit constants.
var SymbolUnit = map[string]Unit{
	"":   UnitCells,
	"fr": UnitFraction,
	"%":  UnitPercent,
	"w":  UnitWidth,
	"h":  UnitHeight,
	"vw": UnitViewWidth,
	"vh": UnitViewHeight,
}

var matchScalar = regexp.MustCompile(`^(-?\d+\.?\d*)(fr|%|w|h|vw|vh)?$`)

// ScalarError is the base error for Scalar operations.
type ScalarError struct{ Msg string }

func (e *ScalarError) Error() string { return e.Msg }

// ScalarResolveError is raised when a scalar cannot be resolved.
type ScalarResolveError struct{ Msg string }

func (e *ScalarResolveError) Error() string { return e.Msg }

// ScalarParseError is raised when a scalar string cannot be parsed.
type ScalarParseError struct{ Msg string }

func (e *ScalarParseError) Error() string { return e.Msg }

// Scalar is a numeric value with a CSS unit.
type Scalar struct {
	Value       float64
	Unit        Unit
	PercentUnit Unit // The unit used when Unit == UnitPercent
}

// String returns the CSS representation of the scalar.
func (s Scalar) String() string {
	if s.Unit == UnitAuto {
		return "auto"
	}
	v := s.Value
	var vs string
	if v == float64(int(v)) {
		vs = strconv.Itoa(int(v))
	} else {
		vs = strconv.FormatFloat(v, 'f', -1, 64)
	}
	return vs + s.Symbol()
}

// IsCells reports whether the unit is explicit cells.
func (s Scalar) IsCells() bool { return s.Unit == UnitCells }

// IsPercent reports whether the unit is a percentage.
func (s Scalar) IsPercent() bool { return s.Unit == UnitPercent }

// IsFraction reports whether the unit is a fraction.
func (s Scalar) IsFraction() bool { return s.Unit == UnitFraction }

// IsAuto reports whether the unit is automatic.
func (s Scalar) IsAuto() bool { return s.Unit == UnitAuto }

// Cells returns the integer cell count if the unit is UnitCells, and ok=false otherwise.
func (s Scalar) Cells() (int, bool) {
	if s.Unit == UnitCells {
		return int(s.Value), true
	}
	return 0, false
}

// Fraction returns the integer fraction value if the unit is UnitFraction, and ok=false otherwise.
func (s Scalar) Fraction() (int, bool) {
	if s.Unit == UnitFraction {
		return int(s.Value), true
	}
	return 0, false
}

// Symbol returns the CSS suffix for the scalar's unit.
func (s Scalar) Symbol() string {
	sym, ok := UnitSymbol[s.Unit]
	if !ok {
		return ""
	}
	return sym
}

// Resolve converts the scalar to a rational cell count.
// size is the container size, viewport is the terminal size, fractionUnit is the
// size of 1fr (pass nil to default to 1).
func (s Scalar) Resolve(size, viewport geometry.Size, fractionUnit *big.Rat) (*big.Rat, error) {
	if fractionUnit == nil {
		fractionUnit = new(big.Rat).SetInt64(1)
	}
	unit := s.Unit
	if unit == UnitPercent {
		unit = s.PercentUnit
	}
	value := new(big.Rat).SetFloat64(s.Value)
	switch unit {
	case UnitCells:
		return value, nil
	case UnitFraction:
		return new(big.Rat).Mul(fractionUnit, value), nil
	case UnitWidth:
		// value * size.Width / 100
		return new(big.Rat).Mul(value, new(big.Rat).SetFrac(
			new(big.Int).SetInt64(int64(size.Width)),
			new(big.Int).SetInt64(100),
		)), nil
	case UnitHeight:
		return new(big.Rat).Mul(value, new(big.Rat).SetFrac(
			new(big.Int).SetInt64(int64(size.Height)),
			new(big.Int).SetInt64(100),
		)), nil
	case UnitViewWidth:
		return new(big.Rat).Mul(value, new(big.Rat).SetFrac(
			new(big.Int).SetInt64(int64(viewport.Width)),
			new(big.Int).SetInt64(100),
		)), nil
	case UnitViewHeight:
		return new(big.Rat).Mul(value, new(big.Rat).SetFrac(
			new(big.Int).SetInt64(int64(viewport.Height)),
			new(big.Int).SetInt64(100),
		)), nil
	default:
		return nil, &ScalarResolveError{Msg: fmt.Sprintf("expected dimensions; found %q", s.String())}
	}
}

// CopyWith returns a copy of the scalar with optional field overrides.
func (s Scalar) CopyWith(value *float64, unit *Unit, percentUnit *Unit) Scalar {
	result := s
	if value != nil {
		result.Value = *value
	}
	if unit != nil {
		result.Unit = *unit
	}
	if percentUnit != nil {
		result.PercentUnit = *percentUnit
	}
	return result
}

// ParseScalar parses a CSS scalar string into a Scalar.
// percentUnit specifies which unit to use when "%" is encountered.
func ParseScalar(token string, percentUnit Unit) (Scalar, error) {
	token = strings.TrimSpace(token)
	if strings.ToLower(token) == "auto" {
		return Scalar{Value: 1.0, Unit: UnitAuto, PercentUnit: UnitAuto}, nil
	}
	m := matchScalar.FindStringSubmatch(token)
	if m == nil {
		return Scalar{}, &ScalarParseError{Msg: fmt.Sprintf("%q is not a valid scalar", token)}
	}
	val, err := strconv.ParseFloat(m[1], 64)
	if err != nil {
		return Scalar{}, &ScalarParseError{Msg: fmt.Sprintf("invalid scalar value %q", m[1])}
	}
	sym := m[2]
	unit, ok := SymbolUnit[sym]
	if !ok {
		unit = UnitCells
	}
	if unit == UnitPercent {
		return Scalar{Value: val, Unit: UnitPercent, PercentUnit: percentUnit}, nil
	}
	return Scalar{Value: val, Unit: unit, PercentUnit: percentUnit}, nil
}

// FromNumber creates a Scalar with UnitCells from a numeric value.
func FromNumber(value float64) Scalar {
	return Scalar{Value: value, Unit: UnitCells, PercentUnit: UnitWidth}
}

// ScalarOffset is a 2D offset expressed as two Scalars.
type ScalarOffset struct {
	X Scalar
	Y Scalar
}

// NullScalarOffset is the zero scalar offset (0 cells, 0 cells).
var NullScalarOffset = ScalarOffset{
	X: FromNumber(0),
	Y: FromNumber(0),
}

// NullScalarOffsetValue returns a ScalarOffset of zero.
func NullScalarOffsetValue() ScalarOffset { return NullScalarOffset }

// ScalarOffsetFromOffset creates a ScalarOffset from integer cell coordinates.
func ScalarOffsetFromOffset(x, y int) ScalarOffset {
	return ScalarOffset{
		X: Scalar{Value: float64(x), Unit: UnitCells, PercentUnit: UnitWidth},
		Y: Scalar{Value: float64(y), Unit: UnitCells, PercentUnit: UnitHeight},
	}
}

// IsZero reports whether both components are zero.
func (so ScalarOffset) IsZero() bool {
	return so.X.Value == 0 && so.Y.Value == 0
}

// Resolve resolves the ScalarOffset to an integer Offset.
func (so ScalarOffset) Resolve(size, viewport geometry.Size) (geometry.Offset, error) {
	xr, err := so.X.Resolve(size, viewport, nil)
	if err != nil {
		return geometry.Offset{}, err
	}
	yr, err := so.Y.Resolve(size, viewport, nil)
	if err != nil {
		return geometry.Offset{}, err
	}
	xf, _ := xr.Float64()
	yf, _ := yr.Float64()
	return geometry.Offset{X: int(xf + 0.5), Y: int(yf + 0.5)}, nil
}

// PercentageStringToFloat converts a percentage string like "20%" to a float
// in the range [0, 1], or parses the string directly as a float.
func PercentageStringToFloat(s string) (float64, error) {
	s = strings.TrimSpace(s)
	if strings.HasSuffix(s, "%") {
		f, err := strconv.ParseFloat(s[:len(s)-1], 64)
		if err != nil {
			return 0, err
		}
		if f < 0 {
			f = 0
		}
		if f > 100 {
			f = 100
		}
		return f / 100.0, nil
	}
	return strconv.ParseFloat(s, 64)
}
