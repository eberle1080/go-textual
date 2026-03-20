package validation

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

// Validator is the interface implemented by all validators.
type Validator interface {
	// Validate returns a ValidationResult for the given input string.
	Validate(value string) ValidationResult
}

// Failure describes a single validation failure.
type Failure struct {
	// Validator is the validator that produced this failure.
	Validator Validator
	// Value is the input that failed validation (may be nil).
	Value *string
	// Description is a human-readable explanation (may be nil).
	Description *string
}

// ValidationResult holds the outcome of one or more validators.
type ValidationResult struct {
	// Failures lists every constraint that was violated.
	Failures []Failure
}

// IsValid reports whether there are no failures.
func (r ValidationResult) IsValid() bool { return len(r.Failures) == 0 }

// BaseValidator is an embeddable helper for building concrete validators.
// Set FailureDescription to override the default error message.
type BaseValidator struct {
	// FailureDescription, when non-nil, replaces the auto-generated message.
	FailureDescription *string
}

// Success returns a successful (empty) ValidationResult.
func (b *BaseValidator) Success() ValidationResult { return ValidationResult{} }

// Fail constructs a failure result with a single Failure entry plus any
// additional failures supplied in extra.
func (b *BaseValidator) Fail(validator Validator, description, value *string, extra ...Failure) ValidationResult {
	f := Failure{Validator: validator, Value: value, Description: description}
	return ValidationResult{Failures: append([]Failure{f}, extra...)}
}

// Merge combines the failures of every result into a single ValidationResult.
func Merge(results ...ValidationResult) ValidationResult {
	var failures []Failure
	for _, r := range results {
		failures = append(failures, r.Failures...)
	}
	return ValidationResult{Failures: failures}
}

// ─── Concrete validators ──────────────────────────────────────────────────────

// RegexValidator checks that a string fully matches a compiled regular expression.
type RegexValidator struct {
	BaseValidator
	// Pattern is the compiled regular expression to match against.
	Pattern *regexp.Regexp
}

// NewRegexValidator compiles pattern and returns a RegexValidator.
// An optional failureDescription overrides the auto-generated error message.
func NewRegexValidator(pattern string, failureDescription ...*string) (*RegexValidator, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	v := &RegexValidator{Pattern: re}
	if len(failureDescription) > 0 {
		v.FailureDescription = failureDescription[0]
	}
	return v, nil
}

// Validate returns a failure if value does not fully match the pattern.
// A full match requires the pattern to cover the entire string from index 0
// to len(value).
func (v *RegexValidator) Validate(value string) ValidationResult {
	loc := v.Pattern.FindStringIndex(value)
	if loc != nil && loc[0] == 0 && loc[1] == len(value) {
		return v.Success()
	}
	desc := fmt.Sprintf("%q does not match pattern /%s/", value, v.Pattern.String())
	if v.FailureDescription != nil {
		desc = *v.FailureDescription
	}
	return v.Fail(v, &desc, &value)
}

// NumberValidator checks that a string represents a valid floating-point number
// within optional bounds.
type NumberValidator struct {
	BaseValidator
	// Minimum, when set, is the inclusive lower bound.
	Minimum *float64
	// Maximum, when set, is the inclusive upper bound.
	Maximum *float64
}

// Validate returns a failure if value is not a number or is out of bounds.
func (v *NumberValidator) Validate(value string) ValidationResult {
	n, err := strconv.ParseFloat(strings.TrimSpace(value), 64)
	if err != nil {
		desc := fmt.Sprintf("%q is not a valid number", value)
		if v.FailureDescription != nil {
			desc = *v.FailureDescription
		}
		return v.Fail(v, &desc, &value)
	}
	if v.Minimum != nil && n < *v.Minimum {
		desc := fmt.Sprintf("%g is less than minimum %g", n, *v.Minimum)
		return v.Fail(v, &desc, &value)
	}
	if v.Maximum != nil && n > *v.Maximum {
		desc := fmt.Sprintf("%g is greater than maximum %g", n, *v.Maximum)
		return v.Fail(v, &desc, &value)
	}
	return v.Success()
}

// IntegerValidator checks that a string represents a valid integer within
// optional bounds.
type IntegerValidator struct {
	BaseValidator
	// Minimum, when set, is the inclusive lower bound.
	Minimum *int
	// Maximum, when set, is the inclusive upper bound.
	Maximum *int
}

// Validate returns a failure if value is not an integer or is out of bounds.
func (v *IntegerValidator) Validate(value string) ValidationResult {
	n, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil {
		desc := fmt.Sprintf("%q is not a valid integer", value)
		if v.FailureDescription != nil {
			desc = *v.FailureDescription
		}
		return v.Fail(v, &desc, &value)
	}
	if v.Minimum != nil && n < *v.Minimum {
		desc := fmt.Sprintf("%d is less than minimum %d", n, *v.Minimum)
		return v.Fail(v, &desc, &value)
	}
	if v.Maximum != nil && n > *v.Maximum {
		desc := fmt.Sprintf("%d is greater than maximum %d", n, *v.Maximum)
		return v.Fail(v, &desc, &value)
	}
	return v.Success()
}

// LengthValidator checks that the rune-count of a string lies within optional
// bounds.
type LengthValidator struct {
	BaseValidator
	// Minimum, when set, is the inclusive minimum rune count.
	Minimum *int
	// Maximum, when set, is the inclusive maximum rune count.
	Maximum *int
}

// Validate returns a failure if the rune-length of value is out of bounds.
func (v *LengthValidator) Validate(value string) ValidationResult {
	length := len([]rune(value))
	if v.Minimum != nil && length < *v.Minimum {
		desc := fmt.Sprintf("length %d is less than minimum %d", length, *v.Minimum)
		if v.FailureDescription != nil {
			desc = *v.FailureDescription
		}
		return v.Fail(v, &desc, &value)
	}
	if v.Maximum != nil && length > *v.Maximum {
		desc := fmt.Sprintf("length %d is greater than maximum %d", length, *v.Maximum)
		if v.FailureDescription != nil {
			desc = *v.FailureDescription
		}
		return v.Fail(v, &desc, &value)
	}
	return v.Success()
}

// FunctionValidator wraps a custom predicate function.
type FunctionValidator struct {
	BaseValidator
	// Func is the predicate; Validate succeeds when it returns true.
	Func func(string) bool
}

// Validate calls Func and returns a failure if it returns false.
func (v *FunctionValidator) Validate(value string) ValidationResult {
	if v.Func != nil && v.Func(value) {
		return v.Success()
	}
	desc := fmt.Sprintf("%q did not pass custom validation", value)
	if v.FailureDescription != nil {
		desc = *v.FailureDescription
	}
	return v.Fail(v, &desc, &value)
}

// URLValidator checks that a string is a parseable URL with both a scheme and
// a host component.
type URLValidator struct {
	BaseValidator
}

// Validate returns a failure if value is not a valid URL.
func (v *URLValidator) Validate(value string) ValidationResult {
	u, err := url.Parse(value)
	if err != nil || u.Scheme == "" || u.Host == "" {
		desc := fmt.Sprintf("%q is not a valid URL", value)
		if v.FailureDescription != nil {
			desc = *v.FailureDescription
		}
		return v.Fail(v, &desc, &value)
	}
	return v.Success()
}
