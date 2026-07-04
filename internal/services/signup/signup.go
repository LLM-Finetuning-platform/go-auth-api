package signup

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/LLM-Finetuning-platform/go-auth-api/internal/services/email"
	"github.com/LLM-Finetuning-platform/go-auth-api/internal/utils"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

func SignUp(req *email.Request ,ctx context.Context,cache *redis.Client) ( error) {
	err := req.Verify()
	if err != nil {
		return  err
	}

	_, err = req.Client.EmailService(req.Params)
	if err != nil {
		return  err
	}
	cache.Set(ctx, req.Email, req.OTP, 10*time.Minute)

	return  nil
}//verify otp standalone function save to db standalone function  and then signupverify fucntion 

func SignUpVerify(ctx context.Context,email string, otp string,username string, cache *redis.Client, db *sql.DB) (bool, error) {
	//lets verify first
	//Cache checked check in db as well
	verify, err := utils.Verify_otp(otp,email,cache)
	if err != nil {
		return false, err
	}
	if !verify{
		return  false, errors.New("Invalid OTP")
	}
	user_id:= uuid.New().String()

	query := `
		INSERT INTO users (user_id, username, email) 
		VALUES ($1, $2, $3);
	`
	_, err = db.ExecContext(ctx,query,user_id,username,email)
	if err != nil{
		return false, err
	}
	return  true, nil
}