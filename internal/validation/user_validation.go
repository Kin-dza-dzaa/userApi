package validation

import (
	"regexp"
	"github.com/go-playground/validator/v10"
)

func InitValidators() (*validator.Validate, error) {
	validation := validator.New()
	if err := validation.RegisterValidation("passval", passValidation); err != nil {
		return nil, err
	}
	if err := validation.RegisterValidation("userval", userValidation); err != nil {
		return nil, err
	}
	return validation, nil
}

func passValidation(fl validator.FieldLevel) bool {
	reg, _ := regexp.Compile(`^[A-Za-z\d]{8,}$`)
	return reg.Match([]byte(fl.Field().String()))
}

func userValidation(fl validator.FieldLevel) bool {
	reg, _ := regexp.Compile(`^[\w\d-]{6,}$`)
	return reg.Match([]byte(fl.Field().String()))
}