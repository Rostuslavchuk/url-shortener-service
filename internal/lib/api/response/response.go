package response

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type ValidatorErrors struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type Response struct {
	Status           string            `json:"status"`
	Error            string            `json:"error,omitempty"`
	ValidationErrors []ValidatorErrors `json:"errors,omitempty"`
}

const (
	StatusOk    = "OK"
	StatusError = "Error"
)

func OK() *Response {
	return &Response{
		Status: StatusOk,
	}
}

func Error(msg string) *Response {
	return &Response{
		Status: StatusError,
		Error:  msg,
	}
}

func Validate(errs validator.ValidationErrors) *Response {
	var errors []ValidatorErrors
	for _, err := range errs {
		var msg string
		switch err.ActualTag() {
		case "required":
			msg = fmt.Sprintf("filled %s is required filed", err.Field())
		case "url":
			msg = fmt.Sprintf("filled %s is not valid URL", err.Field())
		default:
			msg = fmt.Sprintf("filled %s is not valid", err.Field())
		}

		errors = append(errors, ValidatorErrors{
			Field:   err.Field(),
			Message: msg,
		})
	}

	return &Response{
		Status:           StatusError,
		Error:            "Validation Error",
		ValidationErrors: errors,
	}
}
