// Package validation provides a composable framework for validating string input.
//
// Validators implement the [Validator] interface and return a [ValidationResult]
// describing any failures. Multiple validators can be combined via [Merge].
//
// # Basic usage
//
//	minLen := 3
//	v := &validation.LengthValidator{Minimum: &minLen}
//	result := v.Validate("hi")
//	if !result.IsValid() {
//	    for _, f := range result.Failures {
//	        fmt.Println(*f.Description)
//	    }
//	}
//
// # Concrete validators
//
//   - [RegexValidator] — full-match regular expression
//   - [NumberValidator] — floating-point with optional min/max
//   - [IntegerValidator] — integer with optional min/max
//   - [LengthValidator] — rune-length with optional min/max
//   - [FunctionValidator] — custom predicate function
//   - [URLValidator] — requires scheme and host
package validation
