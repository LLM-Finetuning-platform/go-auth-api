package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)


type TokenConfig struct{
	secret []byte
}


func (config *TokenConfig) Encode (username string,hrs int, email string) (string, error){
	token := jwt.NewWithClaims(jwt.SigningMethodES256,
		jwt.MapClaims{
			"username":username,
			"email": email,
			"exp":time.Now().Add(time.Hour*time.Duration(hrs)).Unix(),		
	})
	tokenStr, err := token.SignedString(config.secret)
	if err != nil {
		return "" , err
	}
	return tokenStr, nil

}