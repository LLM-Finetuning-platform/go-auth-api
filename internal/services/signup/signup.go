package signup

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type SignUpRequest struct{
	Email string `json:"email" validate:"required,email"`
}

func validateEmailFormat(email string) error{
	validate := validator.New()
	req := SignUpRequest{Email: email}
	return validate.Struct(req)
}




func (req *SignUpRequest) verify() error{
	err:=validateEmailFormat(req.Email)
	if err != nil{
		return 	fmt.Errorf("Improper Email Format %s", err)
	}
	

	return nil
}

