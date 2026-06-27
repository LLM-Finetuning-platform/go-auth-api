package signup

import (
	"context"
	"fmt"
	"time"

	"github.com/LLM-Finetuning-platform/go-auth-api/internal/services/email"
	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
)

type SignUpRequest struct {
	Email  string `json:"email" validate:"required,email"`
	Client *email.ResendClient
	Params *email.EmailParams
	OTP    string
}

func validateEmailFormat(email string) error {
	validate := validator.New()
	req := SignUpRequest{Email: email}
	return validate.Struct(req)
}

func (req *SignUpRequest) verify() error {
	err := validateEmailFormat(req.Email)
	if err != nil {
		return fmt.Errorf("Improper Email Format %s", err)
	}

	return nil
}

func SignUp(req *SignUpRequest, cache *redis.Client) ( error) {
	err := req.verify()
	if err != nil {
		return  err
	}

	_, err = req.Client.EmailService(req.Params)
	if err != nil {
		return  err
	}
	ctx := context.Background()
	cache.Set(ctx, req.OTP, req.Email, 10*time.Minute)
	

	return  nil
}
