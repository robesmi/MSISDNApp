package handlers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
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

type GithubEmail struct{
	Email string		`json:"email"`
	Primary bool		`json:"primary"`
	Verified bool		`json:"verified"`
	Visibility string	`json:"visibility"`
}

//Configuring credentials for OAuth websites that will be implemented
var googleConfig = &oauth2.Config{
	Endpoint: google.Endpoint,
	Scopes: []string{
		"https://www.googleapis.com/auth/userinfo.email",
	},
}
var githubConfig = &oauth2.Config{
	Endpoint: github.Endpoint,
	Scopes: []string{
		"user:email",
	},
}

func init(){

	//Get the clientId/clientSecret pairs from a non published file and set them in the configs
	config, _ := config.LoadConfig()

	googleConfig.ClientID = config.GoogleClientID
	googleConfig.ClientSecret = config.GoogleClientSecret
	googleConfig.RedirectURL = config.GoogleRedirect

	githubConfig.ClientID = config.GithubClientID
	githubConfig.ClientSecret = config.GithubClientSecret
	githubConfig.RedirectURL = config.GithubRedirect

}

// GetRegisterPage returns the registration page that will be used for presenting
// the available registration methods
func(a AuthHandler) GetRegisterPage(c *gin.Context){
	c.HTML(http.StatusOK, "register.html", nil)
}

// GetLoginPage returns the login page that will be used for presenting the
// available authentication methods
func (a AuthHandler)GetLoginPage(c *gin.Context){
	c.HTML(http.StatusOK, "login.html", nil)
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
		return
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
		if _,ok := err.(*errs.UserAlreadyExists); ok{
			c.HTML(http.StatusBadRequest, "register.html", gin.H{
				"error": "Email already in use",
				"prevUsername": login.Username,
				"prevPassword": login.Password,
			})
			return
		}else{
			c.HTML(http.StatusInternalServerError, "register.html", gin.H{
				"error": "Internal error, please try again",
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
		return
	}

	passwordRegex := regexp.MustCompile(`^(.{0,7}|[^0-9]*|[^A-Z]*|[^a-z]*|[a-zA-Z0-9]*)$`)
	if passwordRegex.MatchString(login.Password){
		c.HTML(http.StatusBadRequest, "login.html", gin.H{
			"error": "Password must have at least 8 characters, contain at least 1 uppercase letter, 1 lower case letter and a number.",
			"prevUsername": login.Username,
			"prevPassword": login.Password,
		})
		return
	}

	loginResp, err := a.Service.LoginNativeUser(login.Username, login.Password)
	if err != nil{
		if _,ok := err.(*errs.UserAlreadyExists); ok{
			c.HTML(http.StatusBadRequest, "login.html", gin.H{
				"error": "Email already in use",
				"prevUsername": login.Username,
				"prevPassword": login.Password,
			})
			return
		}
		if _,ok := err.(*errs.InvalidCredentials); ok{
			c.HTML(http.StatusBadRequest, "login.html", gin.H{
				"error": "Email or password is incorrect",
				"prevUsername": login.Username,
				"prevPassword": login.Password,
			})
			return
		}else{
			c.HTML(http.StatusBadRequest, "login.html", gin.H{
				"error": "Internal error " + err.Error(),
				"prevUsername": login.Username,
				"prevPassword": login.Password,
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

// HandleGoogleCode receives the ID token from google's redirect and registers/logs in
// the user and returns the appropriate tokens
func (a AuthHandler) HandleGoogleCode(c *gin.Context){

	// Get the data from google's response
	respData,readErr := io.ReadAll(c.Request.Body)
	if readErr != nil{
		c.Redirect(http.StatusTemporaryRedirect, "/login")
		return
	}

	// Get the relevant fields from the response
	config, _ := config.LoadConfig()
	idTokenFields := strings.Split(string(respData), "&")
	clientId, idJWT ,csrfToken := strings.Split(idTokenFields[1],"=")[1], strings.Split(idTokenFields[2],"=")[1], strings.Split(idTokenFields[4],"=")[1]
	csrfCookie,err := c.Request.Cookie("g_csrf_token")

	// Verify the fields
	if err != nil{
		log.Println("Error getting csrf protection cookie" + err.Error())
		c.Redirect(http.StatusTemporaryRedirect, "/login")
	}
	if csrfToken != csrfCookie.Value{
		log.Println("CSRF protection cookie and token do not match ")
		c.Redirect(http.StatusTemporaryRedirect, "/login")
	}
	if config.GoogleClientID != clientId{
		log.Println("ID Token client id does not match")
		c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	// Extract the email from the ID token claims and use it to register/login the user
	tokenClaims, valErr := utils.ValidateGoogleIdToken(idJWT)
	if valErr != nil{
		log.Println(valErr.Error())
	}

	if fmt.Sprint(tokenClaims["iss"]) != "https://accounts.google.com"{
		log.Println("Unauthorized google id token issuer, received: " + fmt.Sprint((tokenClaims["iss"])))
	}
	login, appErr := a.Service.RegisterImportedUser(fmt.Sprint(tokenClaims["email"]))
	if _,ok := appErr.(*errs.UserAlreadyExists); ok{
		var newErr error
		login, newErr = a.Service.LoginImportedUser(fmt.Sprint(tokenClaims["email"]))
		if newErr != nil{
			c.Abort()
			return
		}
	}else if appErr != nil {
		log.Println("Error with registering/logging a google user" + appErr.Error())
		return
	}

	c.SetCookie("access_token", login.AccessToken, int(60 * 15),"/","localhost",false,true)
	c.SetCookie("refresh_token", login.RefreshToken, int(60 * 60 * 24),"/","localhost",false,true)

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"access_token" : login.AccessToken,
		"refresh_token" : login.RefreshToken,
	})
	
}


// HandleGithubLogin redirects user to the github oauth2 authorization page with a state and the required parameters
func (a AuthHandler) HandleGithubLogin(c *gin.Context){
	//Set a random state parameter
	state := randToken()
	session := sessions.Default(c)
	session.Set("state",state)
	session.Save()
	// Redirect user to github's consent page with the appropriate parameters
	url := githubConfig.AuthCodeURL(state, oauth2.SetAuthURLParam("client_id", githubConfig.ClientID), oauth2.SetAuthURLParam("redirect_uri",githubConfig.RedirectURL), oauth2.SetAuthURLParam("scope", githubConfig.Scopes[0]))
	c.Redirect(http.StatusTemporaryRedirect,url)
}
// HandleGithubCode receives the access code from google's redirect and makes a post request
// to receive the appropriate access token and refresh token
func (a AuthHandler) HandleGithubCode(c *gin.Context){

	// Check whether the states match
	session := sessions.Default(c)
	responseState := session.Get("state")
	if responseState != c.Query("state"){
		c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("invalid session state %s",responseState))
		return
	}

	//Exchange the received code for an access token to github
	token, err := githubConfig.Exchange(context.Background(), c.Query("code"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest,err)
		return
	}
	client := githubConfig.Client(context.Background(),token)

	// Use token to get the email needed for registering user
	resp,err := client.Get("https://api.github.com/user/emails")
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	defer resp.Body.Close()
	username, errr := io.ReadAll(resp.Body)
	if errr != nil {
		c.Redirect(http.StatusTemporaryRedirect, "/register")
		return
	}
	var userEmails []GithubEmail
	json.Unmarshal(username,&userEmails)

	var primaryEmail string
	for k := range userEmails {		
		if userEmails[k].Primary{
			primaryEmail = userEmails[k].Email
		}
	}
	
	login, appErr := a.Service.RegisterImportedUser(primaryEmail)
	if _,ok := appErr.(*errs.UserAlreadyExists); ok{
		var newErr error
		login, newErr = a.Service.LoginImportedUser(primaryEmail)
		if newErr != nil{
			log.Println("Error " + newErr.Error())
			return
		}
	}else if appErr != nil{
		log.Println("Error with github authentication")
		c.Redirect(http.StatusTemporaryRedirect, "/login")
		return
	}

	c.SetCookie("access_token", login.AccessToken, int(60 * 15),"/","localhost",false,true)
	c.SetCookie("refresh_token", login.RefreshToken, int(60 * 60 * 24),"/","localhost",false,true)

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
		c.Redirect(http.StatusTemporaryRedirect, "/login")
		return
	}
	refClaims, valErr := utils.ValidateRefreshToken(refToken)
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
	refClaims, valErr := utils.ValidateRefreshToken(refToken)
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

