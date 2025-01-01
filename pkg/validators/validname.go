// Package validators Custom validator provider
package validators

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

// ValidName Common name validator including non standard latin character
func ValidName(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	// This regex allowed all common name including non standard
	// latin character
	re := regexp.MustCompile(`^[\p{L} \.'\-]+$`)
	return re.Match([]byte(value))
}
