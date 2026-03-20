package css

import (
	"testing"

	"github.com/eberle1080/go-textual/geometry"
)

func linearEasing(t float64) float64 { return t }

// TestScalarAnimationInterpolation verifies that Tick computes a blended value
// between Start and Destination based on elapsed time.
func TestScalarAnimationInterpolation(t *testing.T) {
	var applied []geometry.Offset
	anim := NewScalarAnimation(
		0.0, 1.0,
		"offset",
		nil,
		geometry.Offset{X: 0, Y: 0},
		geometry.Offset{X: 100, Y: 50},
		linearEasing,
		nil,
		"",
	)
	anim.ApplyFunc = func(v geometry.Offset) { applied = append(applied, v) }

	// Tick at 50% elapsed → should apply ~midpoint.
	done := anim.Tick(0.5, "")
	if done {
		t.Fatal("animation should not be complete at 50% duration")
	}
	if len(applied) == 0 {
		t.Fatal("ApplyFunc was not called")
	}
	mid := applied[len(applied)-1]
	if mid.X != 50 || mid.Y != 25 {
		t.Errorf("midpoint = %v, want {50 25}", mid)
	}
}

// TestScalarAnimationCompletion verifies that Tick returns true, fires
// OnComplete, and applies the destination value when time >= duration.
func TestScalarAnimationCompletion(t *testing.T) {
	completed := false
	var lastApplied geometry.Offset
	dst := geometry.Offset{X: 10, Y: 20}

	anim := NewScalarAnimation(0.0, 1.0, "offset", nil,
		geometry.Offset{}, dst, linearEasing,
		func() { completed = true },
		"",
	)
	anim.ApplyFunc = func(v geometry.Offset) { lastApplied = v }

	done := anim.Tick(1.0, "") // exactly at duration
	if !done {
		t.Fatal("expected animation to be complete")
	}
	if !completed {
		t.Error("OnComplete was not called")
	}
	if lastApplied != dst {
		t.Errorf("ApplyFunc received %v, want %v", lastApplied, dst)
	}
	if anim.FinalValue != dst {
		t.Errorf("FinalValue = %v, want %v", anim.FinalValue, dst)
	}
}

// TestScalarAnimationLevelSuppression verifies that when Level is set and does
// not match appLevel, ApplyFunc is not called for that tick.
func TestScalarAnimationLevelSuppression(t *testing.T) {
	applied := false
	anim := NewScalarAnimation(0.0, 1.0, "offset", nil,
		geometry.Offset{}, geometry.Offset{X: 100, Y: 0},
		linearEasing, nil, "app",
	)
	anim.ApplyFunc = func(v geometry.Offset) { applied = true }

	// Tick with a different appLevel — apply should be suppressed.
	done := anim.Tick(0.5, "screen")
	if done {
		t.Fatal("animation should not be done at 50%")
	}
	if applied {
		t.Error("ApplyFunc should not be called when Level != appLevel")
	}

	// Tick with matching level — apply should fire.
	anim.Tick(0.5, "app")
	if !applied {
		t.Error("ApplyFunc should be called when Level == appLevel")
	}
}

// TestScalarAnimationStopComplete verifies that Stop(true) applies the
// destination value, sets FinalValue, and fires OnComplete.
func TestScalarAnimationStopComplete(t *testing.T) {
	completed := false
	var lastApplied geometry.Offset
	dst := geometry.Offset{X: 7, Y: 3}

	anim := NewScalarAnimation(0.0, 5.0, "offset", nil,
		geometry.Offset{}, dst, linearEasing,
		func() { completed = true },
		"",
	)
	anim.ApplyFunc = func(v geometry.Offset) { lastApplied = v }

	anim.Stop(true)
	if !completed {
		t.Error("OnComplete not called by Stop(true)")
	}
	if lastApplied != dst {
		t.Errorf("ApplyFunc received %v, want %v", lastApplied, dst)
	}
	if anim.FinalValue != dst {
		t.Errorf("FinalValue = %v, want %v", anim.FinalValue, dst)
	}
}

// TestScalarAnimationZeroDuration verifies that a zero-duration animation
// completes immediately on the first Tick.
func TestScalarAnimationZeroDuration(t *testing.T) {
	completed := false
	dst := geometry.Offset{X: 5, Y: 5}
	anim := NewScalarAnimation(0.0, 0.0, "offset", nil,
		geometry.Offset{}, dst, nil,
		func() { completed = true },
		"",
	)
	done := anim.Tick(0.0, "")
	if !done {
		t.Fatal("zero-duration animation should complete immediately")
	}
	if !completed {
		t.Error("OnComplete not called for zero-duration animation")
	}
}
