package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/robesmi/MSISDNApp/model/errs"
	"github.com/robesmi/MSISDNApp/utils"
)

func ValidateTokenUserSection(c *gin.Context){
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
				c.Redirect(http.StatusTemporaryRedirect, "/login")
				c.Abort()
				return
			}
			
			// Check whether the token is valid
			var claims jwt.MapClaims
			var err error
			claims, err = utils.ValidateAccessToken(access_token)
			if err != nil{
				if _,ok := err.(*errs.ExpiredTokenError); ok{
					//Do token refresh logic here
					var refresh_token string
					refresh_token,err = c.Cookie("refresh_token")
					if err != nil{
						c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{ "status": "fail", "message": "Access cookie with no refresh cookie case"})
						return
					}
					_, valErr := utils.ValidateRefreshToken(refresh_token)
					if valErr != nil{
						c.SetCookie("access_token", "", 0,"/","localhost",false,true)
						c.SetCookie("refresh_token", "", 0,"/","localhost",false,true)
						c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{ "status": "fail", "message": "Error with refresh token" + err.Error()})
						return
					}
					c.Redirect(http.StatusTemporaryRedirect,"/refresh?redirect=" + c.FullPath())
					c.Abort()
					return
	
				}else{
					c.Redirect(http.StatusTemporaryRedirect, "/login")
					c.Abort()
					return
				}
			}
			
			// Check if token contains appropriate role
			role := claims["role"]
			if role == "user"{
				c.Next()
				return
			}else{
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{ "status":"fail", "message": "Unauthorized"})
				return
			}
			
}

func ValidateTokenAdminSection(c *gin.Context){

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
		claims,err = utils.ValidateAccessToken(access_token)
		if err != nil{
			if _,ok := err.(errs.ExpiredTokenError);ok{
				//Do token refresh logic here
				var refresh_token string
				refresh_token,err = c.Cookie("refresh_token")
				if err != nil{
					c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{ "status": "fail", "message": "Access cookie with no refresh cookie case"})
				}
				_, valErr := utils.ValidateRefreshToken(refresh_token)
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
		
		// Check if token has appropriate role
		role := claims["role"]
		if role == "user"{
			c.Next()
		}else{
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{ "status":"fail", "message": "Unauthorized"})
		}
		
}