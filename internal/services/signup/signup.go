package signup

import (
	"context"
	"fmt"
	"time"

	"github.com/eswarashish/go-auth-api/internal/services/email"
	"github.com/eswarashish/go-auth-api/internal/utils"
	"github.com/redis/go-redis/v9"
)

type SignUpRequest struct{
	Email utils.Email
	Client *email.ResendClient
	Params *email.EmailParams
	OTP string
}

func (req *SignUpRequest) verify() error{
	err:=utils.ValidateEmailFormat(req.Email.Id)
	if err != nil{
		return 	fmt.Errorf("Improper Email Format %s", err)
	}
	return nil
}

func SignUpEmail(req *SignUpRequest, cache *redis.Client, ttl time.Duration) (bool, error){
	err := req.verify()
	if err != nil{
		return false,fmt.Errorf("Email Verification failed %w",err)
	}
	
	_,err = req.Client.EmailService(req.Params)
	if err != nil {
		return  false, fmt.Errorf("Sign up failed %w",err)
	}
	ctx := context.Background()
	err = cache.Set(ctx,req.OTP,req.Email, ttl).Err()
	if err != nil {
		return  false, fmt.Errorf("Sign up failed %w",err)
	}
	
	return true, nil
}

func SignUpVerification (email string,otp string, cache *redis.Client) (bool, error) {
	err := utils.ValidateEmailFormat(email)
	if err != nil {
		return false, err	
	}
	return utils.OTPVerification(otp,cache,email)
}