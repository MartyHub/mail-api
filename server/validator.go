package server

import "github.com/invopop/validation"

type validator struct{}

func (v validator) Validate(i interface{}) error {
	if a, ok := i.(validation.Validatable); ok {
		return a.Validate()
	}

	return nil
}
