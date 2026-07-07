package utils

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func Verify_otp(otp string,ctx context.Context, email string, cache *redis.Client) (bool, error){
	target, err := cache.Get(ctx, email).Result()
	if err != nil {
		return false,err
	}
	if target == otp {
		return true, nil
	}
	return false, nil
}
