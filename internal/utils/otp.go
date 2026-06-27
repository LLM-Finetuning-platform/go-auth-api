package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

func OTPGeneration() (string, error) {
	max := big.NewInt(900000)
	n, err := rand.Int(rand.Reader,max)
	if err != nil{
		return "", err
	}
	return fmt.Sprintf("%06d",n), nil
}