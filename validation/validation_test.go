package validation_test

import (
	"testing"

	"github.com/eberle1080/go-textual/validation"
)

func ptr[T any](v T) *T { return &v }

func TestValidationResultIsValid(t *testing.T) {
	r := validation.ValidationResult{}
	if !r.IsValid() {
		t.Error("empty result should be valid")
	}
	r.Failures = []validation.Failure{{Description: ptr("oops")}}
	if r.IsValid() {
		t.Error("result with failures should not be valid")
	}
}

func TestMerge(t *testing.T) {
	v := &validation.LengthValidator{Minimum: ptr(5)}
	a := v.Validate("hi")
	b := v.Validate("hello")
	merged := validation.Merge(a, b)
	if len(merged.Failures) != 1 {
		t.Errorf("expected 1 failure, got %d", len(merged.Failures))
	}
}

func TestRegexValidator(t *testing.T) {
	v, err := validation.NewRegexValidator(`^\d+$`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tests := []struct {
		input string
		valid bool
	}{
		{"123", true},
		{"abc", false},
		{"12a", false},
		{"", false},
	}
	for _, tt := range tests {
		r := v.Validate(tt.input)
		if r.IsValid() != tt.valid {
			t.Errorf("input=%q: expected valid=%v, got %v", tt.input, tt.valid, r.IsValid())
		}
	}
}

func TestNumberValidator(t *testing.T) {
	min := 0.0
	max := 100.0
	v := &validation.NumberValidator{Minimum: &min, Maximum: &max}
	tests := []struct {
		input string
		valid bool
	}{
		{"50", true},
		{"0", true},
		{"100", true},
		{"-1", false},
		{"101", false},
		{"abc", false},
	}
	for _, tt := range tests {
		r := v.Validate(tt.input)
		if r.IsValid() != tt.valid {
			t.Errorf("input=%q: expected valid=%v, got %v", tt.input, tt.valid, r.IsValid())
		}
	}
}

func TestIntegerValidator(t *testing.T) {
	min := 1
	max := 10
	v := &validation.IntegerValidator{Minimum: &min, Maximum: &max}
	tests := []struct {
		input string
		valid bool
	}{
		{"5", true},
		{"1", true},
		{"10", true},
		{"0", false},
		{"11", false},
		{"3.14", false},
	}
	for _, tt := range tests {
		r := v.Validate(tt.input)
		if r.IsValid() != tt.valid {
			t.Errorf("input=%q: expected valid=%v, got %v", tt.input, tt.valid, r.IsValid())
		}
	}
}

func TestLengthValidator(t *testing.T) {
	min := 2
	max := 5
	v := &validation.LengthValidator{Minimum: &min, Maximum: &max}
	tests := []struct {
		input string
		valid bool
	}{
		{"hi", true},
		{"hello", true},
		{"a", false},
		{"toolong", false},
	}
	for _, tt := range tests {
		r := v.Validate(tt.input)
		if r.IsValid() != tt.valid {
			t.Errorf("input=%q: expected valid=%v, got %v", tt.input, tt.valid, r.IsValid())
		}
	}
}

func TestFunctionValidator(t *testing.T) {
	v := &validation.FunctionValidator{
		Func: func(s string) bool { return s == "secret" },
	}
	if !v.Validate("secret").IsValid() {
		t.Error("expected valid for 'secret'")
	}
	if v.Validate("wrong").IsValid() {
		t.Error("expected invalid for 'wrong'")
	}
}

func TestURLValidator(t *testing.T) {
	v := &validation.URLValidator{}
	tests := []struct {
		input string
		valid bool
	}{
		{"https://example.com", true},
		{"http://foo.bar/path", true},
		{"example.com", false},
		{"not a url", false},
		{"", false},
	}
	for _, tt := range tests {
		r := v.Validate(tt.input)
		if r.IsValid() != tt.valid {
			t.Errorf("input=%q: expected valid=%v, got %v", tt.input, tt.valid, r.IsValid())
		}
	}
}

func TestFailureDescription(t *testing.T) {
	desc := "custom error"
	v := &validation.LengthValidator{
		BaseValidator: validation.BaseValidator{FailureDescription: &desc},
		Minimum:       ptr(10),
	}
	r := v.Validate("hi")
	if r.IsValid() {
		t.Fatal("expected failure")
	}
	if r.Failures[0].Description == nil || *r.Failures[0].Description != desc {
		t.Errorf("expected custom description %q, got %v", desc, r.Failures[0].Description)
	}
}
