package middleware

import (
	"log"
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
				c.Redirect(http.StatusFound, "/login")
				c.Abort()
				return
			}
			
			// Check whether the token is valid
			var claims jwt.MapClaims
			var err error
			claims, err = utils.ValidateAccessToken(access_token)
			if err != nil{
				if _,ok := err.(*errs.ExpiredTokenError); ok{

					//Check for presence and validity of refresh token
					var refresh_token string
					refresh_token,err = c.Cookie("refresh_token")
					if err != nil{
						c.SetCookie("access_token", "", 0,"/","localhost",false,true)
						c.SetCookie("refresh_token", "", 0,"/","localhost",false,true)
						c.Redirect(http.StatusFound, "/login")
						c.Abort()
						return
					}
					_, valErr := utils.ValidateRefreshToken(refresh_token)
					if valErr != nil{
						c.SetCookie("access_token", "", 0,"/","localhost",false,true)
						c.SetCookie("refresh_token", "", 0,"/","localhost",false,true)
						log.Println("Error with validating refresh token: " + valErr.Error())
						c.Redirect(http.StatusFound, "/login")
						c.Abort()
						return
					}
					// Redirect to refresh handler
					c.Redirect(http.StatusFound,"/refresh?redirect=" + c.FullPath())
					c.Abort()
					return
	
				}else{
					log.Println("Token error " + err.Error())
					c.Redirect(http.StatusFound, "/login")
					c.Abort()
					return
				}
			}
			
			// Check if token contains appropriate role
			role := claims["role"]
			if role == "user" || role == "admin"{
				c.Next()
				return
			}else{
				c.Redirect(http.StatusFound, "/?error=Unauthorized")
				c.Abort()
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
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}
		
		// Check whether the token is valid
		var claims jwt.MapClaims
		var err error
		claims, err = utils.ValidateAccessToken(access_token)
		if err != nil{
			if _,ok := err.(*errs.ExpiredTokenError); ok{

				//Check for presence and validity of refresh token
				var refresh_token string
				refresh_token,err = c.Cookie("refresh_token")
				if err != nil{
					c.Redirect(http.StatusFound, "login")
					c.Abort()
					return
				}
				_, valErr := utils.ValidateRefreshToken(refresh_token)
				if valErr != nil{
					c.SetCookie("access_token", "", 0,"/","localhost",false,true)
					c.SetCookie("refresh_token", "", 0,"/","localhost",false,true)
					log.Println("Error with validating refresh token: " + valErr.Error())
					c.Redirect(http.StatusFound, "/login")
					c.Abort()
					return
				}
				// Redirect to refresh handler
				c.Redirect(http.StatusTemporaryRedirect,"/refresh?redirect=" + c.FullPath())
				c.Abort()
				return

			}else{
				c.Redirect(http.StatusFound, "/login")
				c.Abort()
				return
			}
		}
	
		
		// Check if token has appropriate role
		role := claims["role"]
		if role == "admin"{
			c.Next()
		}else{
			c.Redirect(http.StatusTemporaryRedirect, "/?error=Unauthorized")
			c.Abort()
			return
		}
		
}

func ValidateApiTokenUserSection(c *gin.Context){
	//Get the token either from authorization header or cookie
	var access_token string
	authorizationHeader := c.Request.Header.Get("Authorization")
	fields := strings.Fields(authorizationHeader)

	if len(fields) != 0 && fields[0] == "Bearer" {
		access_token = fields[1]
	}
	// If there's no token in either, kick user back
	if access_token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "No auth token",
		})
		c.Abort()
		return
	}
	
	// Check whether the token is valid
	var claims jwt.MapClaims
	var err error
	claims, err = utils.ValidateAccessToken(access_token)
	if err != nil{
		if _,ok := err.(*errs.ExpiredTokenError); ok{
			// Redirect to refresh handler
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization token expired",
			})
			c.Abort()
			return

		}else{
			log.Println("Token error " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal error",
			})
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
		c.JSON(http.StatusBadRequest,gin.H{
			"error":"Invalid access token, reauthorize.",
		})
		c.Abort()
		return
	}
	
}