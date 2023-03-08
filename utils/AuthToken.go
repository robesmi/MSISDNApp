package utils

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/robesmi/MSISDNApp/config"
	"github.com/robesmi/MSISDNApp/model/errs"
)

// CreateAccessToken creates a JWT token with the custom claim "role" that will
// be used to check whether the bearer has the permissions to use certain routes
func CreateAccessToken(role string) (string, *errs.AppError){

	claims := make(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	claims["iat"] = time.Now().Unix()
	claims["nbf"] = time.Now().Unix()
	claims["role"] = role

	config, _ := config.LoadConfig()
	decodedPrivateKey, err := base64.StdEncoding.DecodeString(config.AccessTokenPrivateKey)
	if err != nil{
		return "", errs.TokenError(err.Error())
	}
	key, appErr := jwt.ParseRSAPrivateKeyFromPEM(decodedPrivateKey)
	if appErr != nil{
		return "", errs.TokenError(err.Error())
	}
	token, err:= jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(key)
	if err != nil{
		return "", errs.TokenError(err.Error())
	}

	return token, nil
}

// CreateRefreshToken creates a JWT token with the custom claim "id" that will
// be used to check whether the token has been revoked or not
func CreateRefreshToken(userid string) (string, *errs.AppError) {

	claims := make(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	claims["iat"] = time.Now().Unix()
	claims["nbf"] = time.Now().Unix()
	claims["id"] = userid

	config, _ := config.LoadConfig()
	decodedPrivateKey, err := base64.StdEncoding.DecodeString(config.RefreshTokenPrivateKey)
	if err != nil{
		return "", errs.TokenError(err.Error())
	}
	key, appErr := jwt.ParseRSAPrivateKeyFromPEM(decodedPrivateKey)
	if appErr != nil{
		return "", errs.TokenError(err.Error())
	}
	token, err:= jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(key)
	if err != nil{
		return "", errs.TokenError(err.Error())
	}

	return token, nil
}

func ValidateToken(token string) *errs.AppError{
	config, _ := config.LoadConfig()
	decodedPublicKey, err := base64.StdEncoding.DecodeString(config.AccessTokenPublicKey)
	if err != nil {
		return errs.UnexpectedError(err.Error())
	}

	key,err :=  jwt.ParseRSAPublicKeyFromPEM(decodedPublicKey)
	if err != nil {
		return errs.UnexpectedError(err.Error())
	}

	parsedToken, err := jwt.Parse(token, func(t *jwt.Token)(interface{}, error){

		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok{
			return nil, fmt.Errorf("unexpected method: %s", t.Header["alg"])
		}
		return key, nil
	})

	if err != nil{
		return errs.UnexpectedError(err.Error())
	}

	if parsedToken.Valid{
		return nil
	}else if ve, ok := err.(*jwt.ValidationError); ok{
		if ve.Errors&jwt.ValidationErrorMalformed != 0{
			return errs.MalformedToken()
		}else if ve.Errors&jwt.ValidationErrorExpired!= 0{
			return errs.ExpiredToken()
		}else {
			return errs.TokenError("Unexpected error")
		}
	}else{
		return errs.TokenError("idk man")
	}
}

func RefreshAccessToken(token string) (string, *errs.AppError){
	panic("panic!!!")
}