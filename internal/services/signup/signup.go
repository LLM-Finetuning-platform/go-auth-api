package signup

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/LLM-Finetuning-platform/go-auth-api/internal/services/email"
	"github.com/LLM-Finetuning-platform/go-auth-api/internal/utils"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
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

func SignUp(req *SignUpRequest ,ctx context.Context,cache *redis.Client) ( error) {
	err := req.verify()
	if err != nil {
		return  err
	}

	_, err = req.Client.EmailService(req.Params)
	if err != nil {
		return  err
	}
	cache.Set(ctx, req.OTP, req.Email, 10*time.Minute)

	return  nil
}//verify otp standalone function save to db standalone function  and then signupverify fucntion 

func SignUpVerify(email string, otp string,username string, cache *redis.Client, db *sql.DB) (bool, error) {
	//lets verify first
	verify, err := utils.Verify_otp(otp,email,cache)
	if err != nil {
		return false, err
	}
	if !verify{
		return  false, errors.New("Invalid OTP")
	}
	userID := uuid.New().String()

	query := `
		INSERT INTO users (id, username, email) 
		VALUES ($1, $2, $3);
	`
	_, err = db.Exec(query,userID,username,email)
	if err != nil{
		return false, err
	}
	return  true, nil
}