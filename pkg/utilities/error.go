// Package utilities Utility functions
package utilities

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

var (
	// ErrTooManyAffectedRows tells that affected rows is above actual effected row
	ErrTooManyAffectedRows = errors.New("too many affected rows")
)

// ErrValue validator value error interfaces
type ErrValue interface{}

// ValidationErrorMsg validator human friendly error messages
type ValidationErrorMsg struct {
	Field   string   `json:"field"`
	Message string   `json:"message"`
	Value   ErrValue `json:"value"`
}

// ParseValidationError validation error parser
func ParseValidationError(err validator.ValidationErrors) []ValidationErrorMsg {
	var msgs []ValidationErrorMsg
	for _, fe := range err {
		var errMsg string
		switch fe.Tag() {
		case "required":
			errMsg = fmt.Sprintf("%s is required", fe.Field())
		case "email":
			errMsg = "invalid email format"
		case "validname":
			errMsg = "invalid name (digit is not allowed)"
		case "gte":
			errMsg = fmt.Sprintf("value must be greater or equal than %s", fe.Param())
		case "lte":
			errMsg = fmt.Sprintf("value must be less or equal than %s", fe.Param())
		case "datetime":
			errMsg = fmt.Sprintf("datetime format must be %s", fe.Param())
		case "number":
			errMsg = "only accept positive number"
		default:
			errMsg = "invalid value"
		}

		fmt.Println(fe.Error())
		msgs = append(msgs, ValidationErrorMsg{
			Field:   fe.Field(),
			Message: errMsg,
			Value:   fe.Value(),
		})
	}

	return msgs
}
