package utils

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func Verify_otp(otp string, email string, cache *redis.Client) (bool, error){
	ctx := context.Background()
	target, err := cache.Get(ctx, email).Result()
	if err != nil {
		return false,err
	}
	if target == otp {
		return true, nil
	}
	return false, nil
}
