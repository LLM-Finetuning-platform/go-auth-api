package utils

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/redis/go-redis/v9"
)

func OTPGeneration() (string, error) {
	max := big.NewInt(900000)
	n, err := rand.Int(rand.Reader,max)
	if err != nil{
		return "", err
	}
	return fmt.Sprintf("%06d",n), nil
}

func OTPVerification (otp string, client *redis.Client, email string) (bool, error) {
	ctx := context.Background()
	ans,err := client.Get(ctx,otp).Result()
	if err != nil {
		return false, fmt.Errorf("OTP Verification failed %w", err)
	}
	if ans == email {
		return true, nil
	}
	return false, nil
}