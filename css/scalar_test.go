package css

import (
	"math/big"
	"testing"

	"github.com/eberle1080/go-textual/geometry"
)

func TestParseScalar(t *testing.T) {
	tests := []struct {
		input       string
		percentUnit Unit
		wantValue   float64
		wantUnit    Unit
		wantErr     bool
	}{
		{"auto", UnitWidth, 1.0, UnitAuto, false},
		{"AUTO", UnitWidth, 1.0, UnitAuto, false},
		{"10", UnitWidth, 10, UnitCells, false},
		{"3.14", UnitWidth, 3.14, UnitCells, false},
		{"2fr", UnitWidth, 2, UnitFraction, false},
		{"50%", UnitWidth, 50, UnitPercent, false},
		{"50%", UnitHeight, 50, UnitPercent, false},
		{"25w", UnitWidth, 25, UnitWidth, false},
		{"10h", UnitHeight, 10, UnitHeight, false},
		{"30vw", UnitWidth, 30, UnitViewWidth, false},
		{"20vh", UnitHeight, 20, UnitViewHeight, false},
		{"-5", UnitWidth, -5, UnitCells, false},
		{"", UnitWidth, 0, 0, true},
		{"foo", UnitWidth, 0, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := ParseScalar(tt.input, tt.percentUnit)
			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseScalar(%q) expected error", tt.input)
				}
				return
			}
			if err != nil {
				t.Fatalf("ParseScalar(%q) unexpected error: %v", tt.input, err)
			}
			if got.Value != tt.wantValue {
				t.Errorf("ParseScalar(%q).Value = %v, want %v", tt.input, got.Value, tt.wantValue)
			}
			if got.Unit != tt.wantUnit {
				t.Errorf("ParseScalar(%q).Unit = %v, want %v", tt.input, got.Unit, tt.wantUnit)
			}
		})
	}
}

func TestScalarString(t *testing.T) {
	tests := []struct {
		s    Scalar
		want string
	}{
		{Scalar{Value: 1.0, Unit: UnitAuto}, "auto"},
		{Scalar{Value: 10, Unit: UnitCells}, "10"},
		{Scalar{Value: 3.5, Unit: UnitCells}, "3.5"},
		{Scalar{Value: 2, Unit: UnitFraction}, "2fr"},
		{Scalar{Value: 50, Unit: UnitPercent}, "50%"},
		{Scalar{Value: 25, Unit: UnitViewWidth}, "25vw"},
	}
	for _, tt := range tests {
		got := tt.s.String()
		if got != tt.want {
			t.Errorf("Scalar(%v).String() = %q, want %q", tt.s, got, tt.want)
		}
	}
}

func TestScalarResolve(t *testing.T) {
	size := geometry.Size{Width: 100, Height: 50}
	viewport := geometry.Size{Width: 200, Height: 100}

	tests := []struct {
		s    Scalar
		want float64
	}{
		{Scalar{Value: 10, Unit: UnitCells}, 10},
		{Scalar{Value: 2, Unit: UnitFraction, PercentUnit: UnitWidth}, 2},   // fractionUnit=1 → 2*1=2
		{Scalar{Value: 50, Unit: UnitPercent, PercentUnit: UnitWidth}, 50},  // 50% of 100 = 50
		{Scalar{Value: 50, Unit: UnitPercent, PercentUnit: UnitHeight}, 25}, // 50% of 50 = 25
		{Scalar{Value: 10, Unit: UnitWidth}, 10},                            // 10% of width=100
		{Scalar{Value: 10, Unit: UnitHeight}, 5},                            // 10% of height=50
		{Scalar{Value: 10, Unit: UnitViewWidth}, 20},                        // 10% of vw=200
		{Scalar{Value: 10, Unit: UnitViewHeight}, 10},                       // 10% of vh=100
	}
	for _, tt := range tests {
		r, err := tt.s.Resolve(size, viewport, nil)
		if err != nil {
			t.Errorf("Scalar(%v).Resolve() error: %v", tt.s, err)
			continue
		}
		f, _ := r.Float64()
		if !approxEqScalar(f, tt.want, 0.001) {
			t.Errorf("Scalar(%v).Resolve() = %v, want %v", tt.s, f, tt.want)
		}
	}
}

func TestScalarResolveFraction(t *testing.T) {
	size := geometry.Size{Width: 100, Height: 50}
	viewport := geometry.Size{Width: 200, Height: 100}
	// 1fr = 20 cells
	fr := new(big.Rat).SetFrac64(20, 1)

	s := Scalar{Value: 3, Unit: UnitFraction}
	r, err := s.Resolve(size, viewport, fr)
	if err != nil {
		t.Fatal(err)
	}
	f, _ := r.Float64()
	if f != 60 {
		t.Errorf("3fr with fractionUnit=20 → %v, want 60", f)
	}
}

func TestScalarOffsetResolve(t *testing.T) {
	so := ScalarOffsetFromOffset(5, 10)
	size := geometry.Size{Width: 100, Height: 50}
	viewport := geometry.Size{Width: 200, Height: 100}

	off, err := so.Resolve(size, viewport)
	if err != nil {
		t.Fatal(err)
	}
	if off.X != 5 || off.Y != 10 {
		t.Errorf("ScalarOffset.Resolve() = %v, want {5 10}", off)
	}
}

func TestPercentageStringToFloat(t *testing.T) {
	tests := []struct {
		input string
		want  float64
	}{
		{"50%", 0.5},
		{"100%", 1.0},
		{"0%", 0.0},
		{"0.5", 0.5},
		{"1.0", 1.0},
	}
	for _, tt := range tests {
		got, err := PercentageStringToFloat(tt.input)
		if err != nil {
			t.Errorf("PercentageStringToFloat(%q) error: %v", tt.input, err)
			continue
		}
		if !approxEqScalar(got, tt.want, 0.001) {
			t.Errorf("PercentageStringToFloat(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func approxEqScalar(a, b, tol float64) bool {
	d := a - b
	if d < 0 {
		d = -d
	}
	return d <= tol
}
