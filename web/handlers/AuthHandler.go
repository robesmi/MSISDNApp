package handlers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/robesmi/MSISDNApp/config"
	"github.com/robesmi/MSISDNApp/model/errs"
	"github.com/robesmi/MSISDNApp/service"
	"github.com/robesmi/MSISDNApp/utils"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

type AuthHandler struct {
	Service service.AuthService
}

type LoginForm struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

//Configuring credentials for OAuth websites that will be implemented
var googleRedirect = "http://127.0.0.1/8080/oauth/google/callback"
var githubRedirect = "http://127.0.0.1/8080/oauth/github/callback"

var googleConfig = &oauth2.Config{
	Endpoint: google.Endpoint,
	RedirectURL: googleRedirect,
	Scopes: []string{
		"https://www.googleapis.com/auth/userinfo.email",
	},
}
var githubConfig = &oauth2.Config{
	Endpoint: github.Endpoint,
	RedirectURL: githubRedirect,
	Scopes: []string{
		"user:email",
	},
}

func init(){

	//Get the clientId/clientSecret pairs from a non published file and set them in the configs
	config, _ := config.LoadConfig()

	googleConfig.ClientID = config.GoogleClientID
	googleConfig.ClientSecret = config.GoogleClientSecret

	githubConfig.ClientID = config.GithubClientID
	githubConfig.ClientSecret = config.GithubClientSecret

}

// GetRegisterPage returns the registration page that will be used for presenting
// the available registration methods
func(a AuthHandler) GetRegisterPage(c *gin.Context){
	c.HTML(http.StatusOK, "register.html", nil)
}

// GetLoginPage returns the login page that will be used for presenting the
// available authentication methods
func (a AuthHandler)GetLoginPage(c *gin.Context){
	c.HTML(http.StatusOK, "index.html", nil)
}

// HandleNativeRegister gets a username/password combination from a form, performs needed validation and creates
// a new user, returning a pair of access/refresh tokens and a success json
func (a AuthHandler) HandleNativeRegister(c *gin.Context){
	var login LoginForm
	if err := c.Bind(&login); err != nil{
		return
	}

	// Checking for any valid email address
	emailRegex := regexp.MustCompile("[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*@(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?")
	if !emailRegex.MatchString(login.Username){
		c.HTML(http.StatusBadRequest, "register.html", gin.H{
			"error": "Enter a valid email address",
			"prevUsername": login.Username,
			"prevPassword": login.Password,
		})
	}

	// Using a negative password regex because golang regex does not support lookahead
	// At least 8 characters, must contain one uppercase character, 1 lowercase and 1 number
	passwordRegex := regexp.MustCompile(`^(.{0,7}|[^0-9]*|[^A-Z]*|[^a-z]*|[a-zA-Z0-9]*)$`)
	if passwordRegex.MatchString(login.Password){
		c.HTML(http.StatusBadRequest, "register.html", gin.H{
			"error": "Password must have at least 8 characters, contain at least 1 uppercase letter, 1 lower case letter and a number.",
			"prevUsername": login.Username,
			"prevPassword": login.Password,
		})
		return
	}

	loginResp, err := a.Service.RegisterNativeUser(login.Username, login.Password)
	if err != nil{
		if _,ok := err.(*errs.UserNotFoundError); ok{
			c.HTML(http.StatusBadRequest, "register.html", gin.H{
				"error": "Email already in use",
				"prevUsername": login.Username,
				"prevPassword": login.Password,
			})
			return
		}else if _,ok := err.(*errs.UserAlreadyExists); ok{
			c.HTML(http.StatusBadRequest, "register.html", gin.H{
				"error": "User already exists",
				"prevUsername": login.Username,
				"prevPassword": login.Password,
			})
			return
		}else{
			c.HTML(http.StatusInternalServerError, "register.html", gin.H{
				"error": "Internal error " + err.Error(),
				"prevUsername": login.Username,
				"prevPassword": login.Password,
			})
			return
		}
	}

	c.SetCookie("access_token", loginResp.AccessToken, int(time.Minute * 15),"/","localhost",false,true)
	c.SetCookie("refresh_token", loginResp.RefreshToken, int(time.Hour * 24),"/","localhost",false,true)

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"access_token" : loginResp.AccessToken,
		"refresh_token" : loginResp.RefreshToken,
	})
}

// HandleNativeLogin will log in the users that choose to use a local account
func (a AuthHandler) HandleNativeLogin(c *gin.Context){
	var login LoginForm
	if err := c.Bind(&login); err != nil{
		return
	}
	
	emailRegex := regexp.MustCompile("[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*@(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?")
	if !emailRegex.MatchString(login.Username){
		c.HTML(http.StatusOK, "login.html", gin.H{
			"error": "Enter a valid email address",
			"prevUsername": login.Username,
			"prevPassword": login.Password,
		})
	}

	passwordRegex := regexp.MustCompile(`^(.{0,7}|[^0-9]*|[^A-Z]*|[^a-z]*|[a-zA-Z0-9]*)$`)
	if passwordRegex.MatchString(login.Password){
		c.HTML(http.StatusBadRequest, "login.html", gin.H{
			"error": "Password must have at least 8 characters, contain at least 1 uppercase letter, 1 lower case letter and a number.",
			"prevUsername": login.Username,
			"prevPassword": login.Password,
		})
	}

	loginResp, err := a.Service.LoginNativeUser(login.Username, login.Password)
	if err != nil{
		if err == err.(*errs.UserAlreadyExists){
			c.HTML(http.StatusBadRequest, "login.html", gin.H{
				"error": "Email already in use",
				"prevUsername": login.Username,
				"prevPassword": login.Password,
			})
		}else{
			c.HTML(http.StatusBadRequest, "login.html", gin.H{
				"error": "Internal error " + err.Error(),
				"prevUsername": login.Username,
				"prevPassword": login.Password,
			})
		}
	}

	c.SetCookie("access_token", loginResp.AccessToken, int(time.Minute * 15),"/","localhost",false,true)
	c.SetCookie("refresh_token", loginResp.RefreshToken, int(time.Hour * 24),"/","localhost",false,true)

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"access_token" : loginResp.AccessToken,
		"refresh_token" : loginResp.RefreshToken,
	})
}


// HandleGoogleLogin redirects user to the google oauth2 authorization page
func (a AuthHandler) HandleGoogleLogin(c *gin.Context){
	state := randToken()
	session := sessions.Default(c)
	session.Set("state",state)
	session.Save()

	url := googleConfig.AuthCodeURL(state)
	c.Redirect(http.StatusTemporaryRedirect,url)
}
// HandleGoogleCode receives the access code from google's redirect and makes a post request
// to receive the appropriate access token and refresh token, used to obtain a username
// to register in the database
func (a AuthHandler) HandleGoogleCode(c *gin.Context){

	session := sessions.Default(c)
	responseState := session.Get("state")
	if responseState != c.Query("state"){
		c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("invalid session state %s",responseState))
		return
	}
	token, err := googleConfig.Exchange(context.Background(), c.Query("code"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest,err)
		return
	}
	client := googleConfig.Client(context.Background(),token)

	resp,err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	defer resp.Body.Close()
	username, errr := io.ReadAll(resp.Body)
	if errr != nil {
		c.AbortWithError(http.StatusInternalServerError, errr)
	}
	login, appErr := a.Service.RegisterImportedUser(string(username))
	if appErr == err.(*errs.UserAlreadyExists){
		var newErr error
		login, newErr = a.Service.LoginImportedUser(string(username))
		if newErr != nil{
			c.Abort()
		}
	}

	c.SetCookie("access_token", login.AccessToken, int(time.Minute * 15),"/","localhost",false,true)
	c.SetCookie("refresh_token", login.RefreshToken, int(time.Hour * 24),"/","localhost",false,true)

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"access_token" : login.AccessToken,
		"refresh_token" : login.RefreshToken,
	})
}


// HandleGithubLogin redirects user to the github oauth2 authorization page
func (a AuthHandler) HandleGithubLogin(c *gin.Context){
	state := randToken()
	session := sessions.Default(c)
	session.Set("state",state)
	session.Save()
	
	url := githubConfig.AuthCodeURL(state)
	c.Redirect(http.StatusTemporaryRedirect,url)
}
// HandleGithubCode receives the access code from google's redirect and makes a post request
// to receive the appropriate access token and refresh token
func (a AuthHandler) HandleGithubCode(c *gin.Context){
	session := sessions.Default(c)
	responseState := session.Get("state")
	if responseState != c.Query("state"){
		c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("invalid session state %s",responseState))
		return
	}
	token, err := googleConfig.Exchange(context.Background(), c.Query("code"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest,err)
		return
	}
	client := googleConfig.Client(context.Background(),token)

	resp,err := client.Get("https://api.github.com/user/emails")
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	defer resp.Body.Close()
	username, errr := io.ReadAll(resp.Body)
	if errr != nil {
		c.AbortWithError(http.StatusInternalServerError, errr)
	}
	login, appErr := a.Service.RegisterImportedUser(string(username))
	if appErr == err.(*errs.UserAlreadyExists){
		var newErr error
		login, newErr = a.Service.LoginImportedUser(string(username))
		if newErr != nil{
			c.Abort()
		}
	}

	c.SetCookie("access_token", login.AccessToken, int(time.Minute * 15),"/","localhost",false,true)
	c.SetCookie("refresh_token", login.RefreshToken, int(time.Hour * 24),"/","localhost",false,true)

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"access_token" : login.AccessToken,
		"refresh_token" : login.RefreshToken,
	})
}

// RefreshAccessToken takes a refresh token, checks the validity and responds with new tokens on successful authentication
func (a AuthHandler) RefreshAccessToken(c *gin.Context){
	refToken, err := c.Cookie("refresh_token")

	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": "Could not refresh token"})
		return
	}
	refClaims, valErr := utils.ValidateToken(refToken)
	if valErr != nil{
		c.SetCookie("access_token", "", 0,"/","localhost",false,true)
		c.SetCookie("refresh_token", "", 0,"/","localhost",false,true)
		c.Redirect(http.StatusTemporaryRedirect, "/login")
		return
	}
	resp, err := a.Service.RefreshTokens(fmt.Sprint(refClaims["id"]),refToken)
	if err != nil{
		c.Redirect(http.StatusTemporaryRedirect, "/login")
		return
	}
	c.SetCookie("access_token", resp.AccessToken, int(time.Minute * 15),"/","localhost",false,true)
	c.SetCookie("refresh_token", resp.RefreshToken, int(time.Hour * 24),"/","localhost",false,true)
	c.Redirect(http.StatusTemporaryRedirect, c.Query("redirect"))
}
func (a AuthHandler) LogOut(c *gin.Context) {
	refToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, "/login")
		return
	}
	refClaims, valErr := utils.ValidateToken(refToken)
	if valErr != nil{
		c.SetCookie("access_token", "", 0,"/","localhost",false,true)
		c.SetCookie("refresh_token", "", 0,"/","localhost",false,true)
		c.Redirect(http.StatusTemporaryRedirect, "/login")
		return
	}
	erro := a.Service.LogOutUser(fmt.Sprint(refClaims["id"]))
	if erro != nil{
		c.SetCookie("access_token", "", 0,"/","localhost",false,true)
		c.SetCookie("refresh_token", "", 0,"/","localhost",false,true)
		c.Redirect(http.StatusTemporaryRedirect, "/login")
		return
	}
	c.SetCookie("access_token", "", 0,"/","localhost",false,true)
	c.SetCookie("refresh_token", "", 0,"/","localhost",false,true)
	c.Redirect(http.StatusTemporaryRedirect, c.Query("redirect"))
}

func randToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

