package validator

import "regexp"

// Declare a regular expression for sanity checking the format of email addresses (we'll
// use this later in the book). If you're interested, this regular expression pattern is
// taken from https://html.spec.whatwg.org/#valid-e-mail-address. Note: if you're
// reading this in PDF or EPUB format and cannot see the full pattern, please see the
// note further down the page.
var (
	//nolint:lll
	EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

type Validator struct {
	Error map[string]string
}

func New() *Validator {
	return &Validator{
		Error: make(map[string]string),
	}
}

// Valid returns true if the error map doesn't contain any errors.
func (v *Validator) Valid() bool {
	return len(v.Error) == 0
}

// AddError Adds an error message to the map so long as no entry already exists for the given key.
func (v *Validator) AddError(key, message string) {
	if _, ok := v.Error[key]; !ok {
		v.Error[key] = message
	}
}

// Check adds an error message to the map only if a validation check is not 'ok'.
func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

// Generic function which returns true if a specific value is in a list.
func PermittedValue[T comparable](value T, permittedValues ...T) bool {
	for i := range permittedValues {
		if value == permittedValues[i] {
			return true
		}
	}

	return false
}

// Unique returns true if all values in a slice are unique.
func Unique[T comparable](values []T) bool {
	// Put the values in a map and compare the lengths of the slice and map
	// to know if there duplicates.
	uniqueValues := make(map[T]bool)

	for _, v := range values {
		uniqueValues[v] = true
	}

	return len(values) == len(uniqueValues)
}
