package css

import "fmt"

// Transition holds the parameters for a CSS transition animation.
type Transition struct {
	Duration float64
	Easing   string
	Delay    float64
}

// NewTransition creates a Transition with the given parameters.
// Default values: duration=1.0, easing="linear", delay=0.0.
func NewTransition(duration float64, easing string, delay float64) Transition {
	return Transition{Duration: duration, Easing: easing, Delay: delay}
}

// String returns the CSS shorthand for the transition, e.g. "0.5s ease-in".
func (t Transition) String() string {
	if t.Delay != 0 {
		return fmt.Sprintf("%.1fs %s %.1f", t.Duration, t.Easing, t.Delay)
	}
	if t.Easing != "linear" {
		return fmt.Sprintf("%.1fs %s", t.Duration, t.Easing)
	}
	return fmt.Sprintf("%.1fs", t.Duration)
}
