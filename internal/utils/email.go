package utils

import "github.com/go-playground/validator/v10"

type Email struct {
	Id string `json:"email" validate:"required,email"`
}
func ValidateEmailFormat(email string) error{
	validate := validator.New()

	req := Email{Id: email}
	return validate.Struct(req)
}

