package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/robesmi/MSISDNApp/model/errs"
	"github.com/robesmi/MSISDNApp/utils"
)

func ValidateTokenUserSection() gin.HandlerFunc{

	return func(c *gin.Context){

		//Get the token either from authorization header or cookie
		var access_token string
		cookie, erro := c.Cookie("access_token")
		authorizationHeader := c.Request.Header.Get("Authorization")
		fields := strings.Fields(authorizationHeader)

		if len(fields) != 0 && fields[0] == "Bearer" {
			access_token = fields[1]
		} else if erro == nil {
			access_token = cookie
		}
		// If there's no token in either, kick user back
		if access_token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "You are not logged in"})
			return
		}
		
		// Check whether the token is valid
		var claims jwt.MapClaims
		var err error
		claims,err = utils.ValidateToken(access_token)
		if err != nil{
			if err == err.(errs.ExpiredTokenError){
				//Do token refresh logic here
				var refresh_token string
				refresh_token,err = c.Cookie("refresh_token")
				if err != nil{
					c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{ "status": "fail", "message": "How do you not have a cookie"})
				}
				_, valErr := utils.ValidateToken(refresh_token)
				if valErr != nil{
					c.SetCookie("access_token", "", 0,"/","localhost",false,true)
					c.SetCookie("refresh_token", "", 0,"/","localhost",false,true)
					c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{ "status": "fail", "message": "Error with refresh token" + err.Error()})
				}
				c.Redirect(http.StatusTemporaryRedirect,"/refresh?redirect=" + c.FullPath())

			}else{
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{ "status":"fail", "message": err.Error()})
			}
		}
		
		// After that we gotta check if the token has the permissions for the route smh...
		role := claims["role"]
		if role == "user"{
			c.Next()
		}else{
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{ "status":"fail", "message": "Unauthoirized"})
		}
		
	}
}

func ValidateTokenAdminSection() gin.HandlerFunc{

	return func(c *gin.Context){

		//Get the token either from authorization header or cookie
		var access_token string
		cookie, erro := c.Cookie("access_token")
		authorizationHeader := c.Request.Header.Get("Authorization")
		fields := strings.Fields(authorizationHeader)

		if len(fields) != 0 && fields[0] == "Bearer" {
			access_token = fields[1]
		} else if erro == nil {
			access_token = cookie
		}
		// If there's no token in either, kick user back
		if access_token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "You are not logged in"})
			return
		}
		
		// Check whether the token is valid
		var claims jwt.MapClaims
		var err error
		claims,err = utils.ValidateToken(access_token)
		if err != nil{
			if err == err.(errs.ExpiredTokenError){
				//Do token refresh logic here
				var refresh_token string
				refresh_token,err = c.Cookie("refresh_token")
				if err != nil{
					c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{ "status": "fail", "message": "How do you not have a cookie"})
				}
				_, valErr := utils.ValidateToken(refresh_token)
				if valErr != nil{
					c.SetCookie("access_token", "", 0,"/","localhost",false,true)
					c.SetCookie("refresh_token", "", 0,"/","localhost",false,true)
					c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{ "status": "fail", "message": "Error with refresh token" + err.Error()})
				}
				c.Redirect(http.StatusTemporaryRedirect,"/refresh?redirect=" + c.FullPath())

			}else{
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{ "status":"fail", "message": err.Error()})
			}
		}
		
		// After that we gotta check if the token has the permissions for the route smh...
		role := claims["role"]
		if role == "user"{
			c.Next()
		}else{
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{ "status":"fail", "message": "Unauthoirized"})
		}
		
	}
}