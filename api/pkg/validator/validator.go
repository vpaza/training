package validator

import (
	govalidator "github.com/go-playground/validator/v10"
)

type CustomValidator struct {
	validator *govalidator.Validate
}

func Get() *CustomValidator {
	return &CustomValidator{
		validator: govalidator.New(),
	}
}

func (v *CustomValidator) Validate(i interface{}) error {
	if err := v.validator.Struct(i); err != nil {
		return err
	}

	return nil
}
