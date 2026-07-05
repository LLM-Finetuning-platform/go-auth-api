package login

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/LLM-Finetuning-platform/go-auth-api/internal/services/email"
	"github.com/LLM-Finetuning-platform/go-auth-api/internal/utils"
	"github.com/redis/go-redis/v9"
)



func Login(req *email.Request ,ctx context.Context,cache *redis.Client) ( error) {
	err := req.Verify()
	if err != nil {
		return  err
	}

	_, err = req.Client.EmailService(req.Params)
	if err != nil {
		return  err
	}
	cache.Set(ctx, req.OTP, req.Email, 10*time.Minute)

	return  nil
}

func LoginVerify(email string, otp string, cache *redis.Client, db *sql.DB, ctx context.Context) (bool, error) {
	//lets verify first
	verify, err := utils.Verify_otp(otp,email,cache)
	if err != nil {
		return false, err
	}
	if !verify{
		return  false, errors.New("Invalid OTP")
	}
	verify, err = utils.CheckExisting(email,ctx,db)	
	if err != nil{
		return false, err
	}
	return  verify, nil
}