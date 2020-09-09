package accounts

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator"
)

type Validation struct {
	validate *validator.Validate
}

func NewValidate() *Validation {
	validate := validator.New()

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return &Validation{validate}
}

func (v *Validation) Validate(i interface{}) map[string]string {
	err := v.validate.Struct(i)

	if err == nil {
		return nil
	}

	errs := err.(validator.ValidationErrors)

	formattedErrs := make(map[string]string)

	for _, err := range errs {
		ve := err.(validator.FieldError)
		formattedErrs[ve.Field()] = ve.Tag()
	}

	return formattedErrs
}
