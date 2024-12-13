package utils

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type ErrorResponse struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Value string `json:"value,omitempty"`
}
type ErrorDetail struct {
	TypeError        string `json:"type_error"`
	ErrorDescription string `json:"error_description"`
}

func ValidateStruct(payload interface{}) []*ErrorResponse {
	var errors []*ErrorResponse
	validate := validator.New()
	err := validate.Struct(payload)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.Field = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}
func ShowErrors(errors []*ErrorResponse) ErrorDetail {
	var stringError string
	count := 0
	for _, err := range errors {
		count += 1
		stringError += "Field: " + err.Field + " [tag-error: " + err.Tag + "]"
		if count < len(errors) {
			stringError += ", "
		}
	}
	return ErrorDetail{
		TypeError:        "Error fields",
		ErrorDescription: stringError,
	}
}
func IsDateValid(input string) bool {
	layout := "2006-01-02"
	_, err := time.Parse(layout, input)
	return err == nil
}
