package handlers

import (
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/robesmi/MSISDNApp/model/errs"
	"github.com/robesmi/MSISDNApp/service"
	"github.com/robesmi/MSISDNApp/utils"
)

type AuthApiHandler struct {
	Service service.AuthService
}

type RefreshRequest struct{
	RefreshToken string `json:"refresh_token"`
}
 
var	(
	validateAccessToken = utils.ValidateAccessToken
	validateRefreshToken = utils.ValidateRefreshToken
)

// HandleNativeRegister gets a username/password combination from a form, performs needed validation and creates
// a new user, returning a pair of access/refresh tokens and a success json
func (a AuthApiHandler) HandleNativeRegisterCall(c *gin.Context){
	var login LoginForm
	if err := c.Bind(&login); err != nil{
		return
	}

	// Checking for any valid email address
	emailRegex := regexp.MustCompile("[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*@(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?")
	if !emailRegex.MatchString(login.Username){
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Enter a valid email address",
		})
		return
	}

	// Using a negative password regex because golang regex does not support lookahead
	// At least 8 characters, must contain one uppercase character, 1 lowercase and 1 number
	passwordRegex := regexp.MustCompile(`^(.{0,7}|[^0-9]*|[^A-Z]*|[^a-z]*|[a-zA-Z0-9]*)$`)
		if passwordRegex.MatchString(login.Password){
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Password must have at least 8 characters, contain at least 1 uppercase letter, 1 lower case letter,1 number and a special character.",
		})
		return
	}

	loginResp, err := a.Service.RegisterNativeUser(login.Username, login.Password)
	if err != nil{
		if _,ok := err.(*errs.UserAlreadyExists); ok{
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Email already in use",
			})
			return
		}else{
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal error, please try again",
			})
			return
		}
	}

	c.SetCookie("access_token", loginResp.AccessToken, int(60 * 15),"/","localhost",false,true)
	c.SetCookie("refresh_token", loginResp.RefreshToken, int(60 * 60 * 24),"/","localhost",false,true)

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"access_token" : loginResp.AccessToken,
		"refresh_token" : loginResp.RefreshToken,
	})
}


// HandleNativeLogin will log in the users that choose to use a local account
func (a AuthApiHandler) HandleNativeLoginCall(c *gin.Context){
	var login LoginForm
	if err := c.Bind(&login); err != nil{
		log.Println("Api error binding login form: " + err.Error())
		return
	}
	
	emailRegex := regexp.MustCompile("[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*@(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?")
	if !emailRegex.MatchString(login.Username){
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Enter a valid email address",
		})
		return
	}

	passwordRegex := regexp.MustCompile(`^(.{0,7}|[^0-9]*|[^A-Z]*|[^a-z]*|[a-zA-Z0-9]*)$`)
	if passwordRegex.MatchString(login.Password){
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Password must have at least 8 characters, contain at least 1 uppercase letter, 1 lower case letter,1 number and a special character.",
		})
		return
	}

	loginResp, err := a.Service.LoginNativeUser(login.Username, login.Password)
	if err != nil{
		if _,ok := err.(*errs.InvalidCredentials); ok{
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Email or password is incorrect",
			})
			return
		}else{
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal error: " + err.Error(),
			})
			return
		}
	}

	c.SetCookie("access_token", loginResp.AccessToken, int(60 * 15),"/","localhost",false,true)
	c.SetCookie("refresh_token", loginResp.RefreshToken, int(60 * 60 * 24),"/","localhost",false,true)

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"access_token" : loginResp.AccessToken,
		"refresh_token" : loginResp.RefreshToken,
	})
}


// RefreshAccessToken takes a refresh token, checks the validity and responds with new tokens on successful authentication
func (a AuthApiHandler) RefreshAccessTokenCall(c *gin.Context){
	
	var refToken RefreshRequest
	err := c.ShouldBind(&refToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"error": err.Error(),
		})
		return
	}

	refClaims, valErr := validateAccessToken(refToken.RefreshToken)
	if valErr != nil{
		log.Println("Error validating refresh token:" + valErr.Error())
		c.SetCookie("access_token", "", 0,"/","localhost",false,true)
		c.SetCookie("refresh_token", "", 0,"/","localhost",false,true)
		c.Redirect(http.StatusTemporaryRedirect, "/login")
		return
	}
	resp, err := a.Service.RefreshTokens(fmt.Sprint(refClaims["id"]),refToken.RefreshToken)
	if err != nil{
		log.Println("Error refreshing access token: " + err.Error())
		c.Redirect(http.StatusTemporaryRedirect, "/login")
		return
	}
	c.SetCookie("access_token", resp.AccessToken, int(60 * 15),"/","localhost",false,true)
	c.SetCookie("refresh_token", resp.RefreshToken, int(60 * 60 * 24),"/","localhost",false,true)
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"access_token" : resp.AccessToken,
		"refresh_token" : resp.RefreshToken,
	})
}

func (a AuthApiHandler) LogOutCall(c *gin.Context) {

	var refToken RefreshRequest
	err := c.ShouldBind(&refToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"error": err.Error(),
		})
		return
	}
	if err != nil {
		log.Println("Error logging user out: " + err.Error())
		c.Redirect(http.StatusTemporaryRedirect, "/login")
		return
	}
	refClaims, valErr := validateRefreshToken(refToken.RefreshToken)
	if valErr != nil{
		c.SetCookie("access_token", "", 0,"/","localhost",false,true)
		c.SetCookie("refresh_token", "", 0,"/","localhost",false,true)
		c.JSON(http.StatusOK, gin.H{
			"status": "failed",
			"error": valErr.Error(),
		})
		return
	}
	erro := a.Service.LogOutUser(fmt.Sprint(refClaims["id"]))
	if erro != nil{
		c.SetCookie("access_token", "", 0,"/","localhost",false,true)
		c.SetCookie("refresh_token", "", 0,"/","localhost",false,true)
		c.JSON(http.StatusOK, gin.H{
			"status": "failed",
			"error": erro.Error(),
		})
		return
	}
	c.SetCookie("access_token", "", 0,"/","localhost",false,true)
	c.SetCookie("refresh_token", "", 0,"/","localhost",false,true)
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}