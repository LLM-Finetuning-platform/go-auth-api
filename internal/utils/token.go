package utils

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)


type TokenConfig struct{
	privateKey *ecdsa.PrivateKey
	publicKey *ecdsa.PublicKey
}


func NewTokenConfig(privateKeyPEM []byte) (*TokenConfig, error){
	block, _ := pem.Decode(privateKeyPEM)
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block")
	}
	privKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil{
		return nil, err
	}
	return &TokenConfig{
		privateKey: privKey,
		publicKey: &privKey.PublicKey,
	}, nil
}

func (config *TokenConfig) Encode (username string,hrs int, email string) (string, error){
	token := jwt.NewWithClaims(jwt.SigningMethodES256,
		jwt.MapClaims{
			"username":username,
			"email": email,
			"exp":time.Now().Add(time.Hour*time.Duration(hrs)).Unix(),		
	})
	tokenStr, err := token.SignedString(config.privateKey)
	if err != nil {
		return "" , err
	}
	return tokenStr, nil

}

func (config *TokenConfig) Decode (tokenStr string) (map[string]interface{}, error){
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims,func (token *jwt.Token) (interface{},error){
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok{
			return nil,fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return config.publicKey, nil
	})
	if err != nil{
		return nil, err
	}	
	if !token.Valid{
		return nil, fmt.Errorf("invalid token")
	}
	return claims,nil
}