package css

import "testing"

func TestTransitionString(t *testing.T) {
	tests := []struct {
		tr   Transition
		want string
	}{
		{Transition{Duration: 0.5, Easing: "linear", Delay: 0}, "0.5s"},
		{Transition{Duration: 1.0, Easing: "linear", Delay: 0}, "1.0s"},
		{Transition{Duration: 0.3, Easing: "ease-in", Delay: 0}, "0.3s ease-in"},
		{Transition{Duration: 0.5, Easing: "ease-out", Delay: 0}, "0.5s ease-out"},
		{Transition{Duration: 0.5, Easing: "linear", Delay: 0.2}, "0.5s linear 0.2"},
		{Transition{Duration: 1.0, Easing: "ease-in-out", Delay: 0.5}, "1.0s ease-in-out 0.5"},
	}
	for _, tt := range tests {
		got := tt.tr.String()
		if got != tt.want {
			t.Errorf("Transition%+v.String() = %q, want %q", tt.tr, got, tt.want)
		}
	}
}

func TestNewTransition(t *testing.T) {
	tr := NewTransition(0.5, "ease-in", 0.1)
	if tr.Duration != 0.5 || tr.Easing != "ease-in" || tr.Delay != 0.1 {
		t.Errorf("NewTransition fields = %+v", tr)
	}
}
