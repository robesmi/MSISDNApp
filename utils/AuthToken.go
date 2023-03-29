package utils

import (
	"encoding/base64"
	"fmt"
	"log"
	"time"

	"github.com/MicahParks/keyfunc"
	"github.com/golang-jwt/jwt/v4"
	"github.com/robesmi/MSISDNApp/model/errs"
	"github.com/robesmi/MSISDNApp/vault"
)

var GoogleJwkUrl = "https://www.googleapis.com/oauth2/v3/certs"

// CreateAccessToken creates a JWT access token with the custom claim "role" that will
// be used to check whether the bearer has the permissions to use certain routes
func CreateAccessToken(role string, vault *vault.Vault) (string, error){

	claims := make(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	claims["iat"] = time.Now().Unix()
	claims["nbf"] = time.Now().Unix()
	claims["role"] = role

	data, fetchErr := vault.Fetch("appvars", "AccessTokenPrivateKey")
	if fetchErr != nil{
		return "", fetchErr
	}
	decodedPrivateKey, err := base64.StdEncoding.DecodeString(data["AccessTokenPrivateKey"])
	if err != nil{
		return "", errs.NewTokenError(err.Error())
	}
	key, appErr := jwt.ParseRSAPrivateKeyFromPEM(decodedPrivateKey)
	if appErr != nil{
		return "", errs.NewTokenError(err.Error())
	}
	token, err:= jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(key)
	if err != nil{
		return "", errs.NewTokenError(err.Error())
	}

	return token, nil
}

// CreateRefreshToken creates a JWT refresh token with the custom claim "id" that will
// be used to check whether the token has been revoked or not
func CreateRefreshToken(userid string, vault *vault.Vault) (string, error) {

	claims := make(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	claims["iat"] = time.Now().Unix()
	claims["nbf"] = time.Now().Unix()
	claims["id"] = userid

	data, fetchErr := vault.Fetch("appvars", "RefreshTokenPrivateKey")
	if fetchErr != nil{
		return "", fetchErr
	}
	decodedPrivateKey, err := base64.StdEncoding.DecodeString(data["RefreshTokenPrivateKey"])
	if err != nil{
		return "", errs.NewTokenError(err.Error())
	}
	key, appErr := jwt.ParseRSAPrivateKeyFromPEM(decodedPrivateKey)
	if appErr != nil{
		return "", errs.NewTokenError(err.Error())
	}
	token, err:= jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(key)
	if err != nil{
		return "", errs.NewTokenError(err.Error())
	}

	return token, nil
}

// ValidateRefreshToken takes a jwt access token as input and validates it. Returns a jwt.MapClaims of the user role or
// an error otherwise
func ValidateAccessToken(client *vault.Vault, token string) (jwt.MapClaims,error){

	publicKey, fetchErr := client.Fetch("appvars", "AccessTokenPublicKey")
	if fetchErr != nil{
		return nil, fetchErr
	}

	decodedPublicKey, err := base64.StdEncoding.DecodeString(publicKey["AccessTokenPublicKey"])
	if err != nil{
		return nil, err
	}

	key,err :=  jwt.ParseRSAPublicKeyFromPEM(decodedPublicKey)
	if err != nil {
		return nil, errs.NewUnexpectedError(err.Error())
	}

	parsedToken, err := jwt.Parse(token, func(t *jwt.Token)(interface{}, error){
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok{
			return nil, fmt.Errorf("unexpected method: %s", t.Header["alg"])
		}
		return key, nil
	})

	if ve, ok := err.(*jwt.ValidationError); ok{
		if ve.Errors&jwt.ValidationErrorMalformed != 0{
			return nil, errs.NewMalformedTokenError()
		}else if ve.Errors&jwt.ValidationErrorExpired!= 0{
			return nil, errs.NewExpiredTokenError()
		}else {
			return nil, errs.NewUnexpectedError(ve.Error())
		}
	}else if err != nil{
		return nil, errs.NewUnexpectedError(ve.Error())
	}else if claims, valid := parsedToken.Claims.(jwt.MapClaims); valid && parsedToken.Valid{
		return claims, nil
	}
	return nil, errs.NewExpiredTokenError()
}

// ValidateRefreshToken takes a jwt refresh token as input and validates it. Returns a jwt.MapClaims of the user uuid or
// an error
func ValidateRefreshToken(client *vault.Vault, token string) (jwt.MapClaims,error){

	publicKey, fetchErr := client.Fetch("appvars", "RefreshTokenPublicKey")
	if fetchErr != nil{
		return nil, fetchErr
	}
	decodedPublicKey, err := base64.StdEncoding.DecodeString(publicKey["RefreshTokenPublicKey"])
	if err != nil{
		return nil, err
	}

	key,err :=  jwt.ParseRSAPublicKeyFromPEM(decodedPublicKey)
	if err != nil {
		return nil, errs.NewUnexpectedError(err.Error())
	}

	parsedToken, err := jwt.Parse(token, func(t *jwt.Token)(interface{}, error){
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok{
			return nil, fmt.Errorf("unexpected method: %s", t.Header["alg"])
		}
		return key, nil
	})

	if err != nil{
		log.Println(err.Error())
		return nil, errs.NewUnexpectedError(err.Error())
	}
	
	if claims, valid := parsedToken.Claims.(jwt.MapClaims); valid && parsedToken.Valid{
		return claims, nil
	}else if ve, ok := err.(*jwt.ValidationError); ok{
		if ve.Errors&jwt.ValidationErrorMalformed != 0{
			return nil, errs.NewMalformedTokenError()
		}else if ve.Errors&jwt.ValidationErrorExpired!= 0{
			return nil, errs.NewExpiredTokenError()
		}else {
			return nil, errs.NewUnexpectedError(ve.Error())
		}
	}else{
		return nil, errs.NewUnexpectedError(ve.Error())
	}
	
}

// ValidateGoogleIdToken receives a Google ID token as input and validates
// it using Google's provided jwk url and returns a jwt.MapClaims response with its claims
func ValidateGoogleIdToken(token string) (jwt.MapClaims, error){

	jwks, err := keyfunc.Get(GoogleJwkUrl,keyfunc.Options{})
	if err != nil{
		log.Fatalf("Failed to create JWK from url" + err.Error())
		log.Println("Failed to create JWK from urk " + err.Error())
	}

	//Flag to check for clock skew issue due to jwt library not implementing tolerance now
	iatFlag := 0

	parsedToken, parseErr := jwt.Parse(token, jwks.Keyfunc)
	if parseErr != nil {

		// If the exact error is not encountered, proceed with returning an error normally
		if parseErr.Error() != "Token used before issued"{

			validationErr, isValidationError := parseErr.(*jwt.ValidationError)
			if !isValidationError {
				log.Println("JWT parsing failed and is not a ValidationError")
				return nil, errs.NewTokenValidationError(parseErr.Error())
			}

			hasIssuedAtValidationError := validationErr.Errors&jwt.ValidationErrorIssuedAt != 0
			if !hasIssuedAtValidationError {
				log.Println("JWT parsing failed, but is not a ValidationErrorIssuedAt")
				return nil, errs.NewTokenValidationError(parseErr.Error())
			}

			// toggle ValidationErrorIssuedAt and check if it was the only validation error
			remainingErrors := validationErr.Errors ^ jwt.ValidationErrorIssuedAt
			if remainingErrors > 0 {
				log.Println("JWT parsing failed, but has other errors besides ValidationErrorIssuedAt")
				return nil, errs.NewTokenValidationError(parseErr.Error())
			}

		}else{
			iatFlag = 1
		}

	}
	if claims, valid := parsedToken.Claims.(jwt.MapClaims); valid && (parsedToken.Valid || iatFlag == 1){
		return claims, nil
	}else if ve, ok := err.(*jwt.ValidationError); ok{
		if ve.Errors&jwt.ValidationErrorMalformed != 0{
			return nil, errs.NewMalformedTokenError()
		}else if ve.Errors&jwt.ValidationErrorExpired!= 0{
			return nil, errs.NewExpiredTokenError()
		}else {
			return nil, errs.NewUnexpectedError(ve.Error())
		}
	}else{
		return nil, errs.NewUnexpectedError(ve.Error())
	}
}
