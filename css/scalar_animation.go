package css

import "github.com/eberle1080/go-textual/geometry"

// EasingFunction is a function that maps a progress value [0,1] to an eased value.
type EasingFunction func(t float64) float64

// StylesResolver is an interface used by ScalarAnimation to resolve sizes without
// coupling to the widget layer.
type StylesResolver interface {
	// ContainerSize returns the current container size.
	ContainerSize() geometry.Size
	// ViewportSize returns the current viewport size.
	ViewportSize() geometry.Size
}

// ScalarAnimation holds the parameters for an in-progress scalar animation.
// Set ApplyFunc before calling Tick to receive intermediate and final values.
type ScalarAnimation struct {
	StartTime   float64
	Duration    float64
	Attribute   string
	FinalValue  interface{}
	Start       geometry.Offset
	Destination geometry.Offset
	Easing      EasingFunction
	OnComplete  func()
	Level       string
	// ApplyFunc is called by Tick with the interpolated offset at each step.
	// It is also called with the final destination value on completion.
	ApplyFunc func(value geometry.Offset)
}

// NewScalarAnimation creates a new ScalarAnimation.
func NewScalarAnimation(
	startTime, duration float64,
	attribute string,
	finalValue interface{},
	start, destination geometry.Offset,
	easing EasingFunction,
	onComplete func(),
	level string,
) ScalarAnimation {
	return ScalarAnimation{
		StartTime:   startTime,
		Duration:    duration,
		Attribute:   attribute,
		FinalValue:  finalValue,
		Start:       start,
		Destination: destination,
		Easing:      easing,
		OnComplete:  onComplete,
		Level:       level,
	}
}

// lerp linearly interpolates between a and b by t (clamped to [0,1]).
func lerp(a, b float64, t float64) float64 {
	if t <= 0 {
		return a
	}
	if t >= 1 {
		return b
	}
	return a + (b-a)*t
}

// Tick advances the animation to the given time.  It computes the eased
// progress, blends Start toward Destination, and calls ApplyFunc with the
// current value.  Returns true (and fires OnComplete) when the animation is
// complete.
//
// If Level is non-empty and does not match appLevel the animation is suppressed
// for this tick (ApplyFunc is not called) but the timer still advances.
func (a *ScalarAnimation) Tick(time float64, appLevel string) bool {
	// Zero-duration animation: jump to final value immediately.
	if a.Duration <= 0 {
		a.FinalValue = a.Destination
		if a.Level == "" || a.Level == appLevel {
			if a.ApplyFunc != nil {
				a.ApplyFunc(a.Destination)
			}
		}
		if a.OnComplete != nil {
			a.OnComplete()
		}
		return true
	}

	elapsed := time - a.StartTime
	if elapsed >= a.Duration {
		// Animation finished: apply the exact destination value.
		a.FinalValue = a.Destination
		if a.Level == "" || a.Level == appLevel {
			if a.ApplyFunc != nil {
				a.ApplyFunc(a.Destination)
			}
		}
		if a.OnComplete != nil {
			a.OnComplete()
		}
		return true
	}

	// Compute eased progress and apply intermediate value.
	progress := elapsed / a.Duration
	if a.Easing != nil {
		progress = a.Easing(progress)
	}
	if a.Level == "" || a.Level == appLevel {
		current := geometry.Offset{
			X: int(lerp(float64(a.Start.X), float64(a.Destination.X), progress)),
			Y: int(lerp(float64(a.Start.Y), float64(a.Destination.Y), progress)),
		}
		if a.ApplyFunc != nil {
			a.ApplyFunc(current)
		}
	}
	return false
}

// Stop halts the animation.  If complete is true the destination value is
// applied via ApplyFunc, FinalValue is set, and the completion callback fires.
func (a *ScalarAnimation) Stop(complete bool) {
	if complete {
		a.FinalValue = a.Destination
		if a.ApplyFunc != nil {
			a.ApplyFunc(a.Destination)
		}
		if a.OnComplete != nil {
			a.OnComplete()
		}
	}
}

// Equal reports whether two ScalarAnimations target the same attribute.
func (a ScalarAnimation) Equal(other ScalarAnimation) bool {
	return a.Attribute == other.Attribute
}
