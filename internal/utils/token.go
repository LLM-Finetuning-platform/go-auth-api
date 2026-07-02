package utils

import "github.com/golang-jwt/jwt/v5"


type TokenConfig struct{
	secret string
}


func (config *TokenConfig) Encode (username string, email string) (*jwt.ClaimStrings){
	return  &jwt.ClaimStrings{}
}