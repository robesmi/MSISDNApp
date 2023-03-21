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

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/robesmi/MSISDNApp/config"
	"github.com/robesmi/MSISDNApp/model/errs"
	"github.com/robesmi/MSISDNApp/service"
	"github.com/robesmi/MSISDNApp/utils"
	"github.com/rs/zerolog"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

type AuthHandler struct {
	Service service.AuthService
	Logger zerolog.Logger
}

type LoginForm struct {
	Username string `form:"username" json:"username"`
	Password string `form:"password" json:"password"`
}

type GithubEmail struct{
	Email string		`json:"email"`
	Primary bool		`json:"primary"`
	Verified bool		`json:"verified"`
	Visibility string	`json:"visibility"`
}

//Configuring credentials for OAuth websites that will be implemented
var githubConfig = &oauth2.Config{
	Endpoint: github.Endpoint,
	Scopes: []string{
		"user:email",
	},
}

func init(){

	//Get the identity provider information from a non published file and set them in the configs
	config, err := config.LoadConfig()
	if err != nil{
		log.Println("Error loading config: " + err.Error())
	}

	githubConfig.ClientID = config.GithubClientID
	githubConfig.ClientSecret = config.GithubClientSecret
	githubConfig.RedirectURL = config.GithubRedirect

}

// GetRegisterPage returns the registration page that will be used for presenting
// the available registration methods
func(a AuthHandler) GetRegisterPage(c *gin.Context){
	redirectErr := c.Query("error")
	if redirectErr == ""{
		c.HTML(http.StatusOK, "register.html", nil)
	}else if redirectErr == "AuthError"{
		c.HTML(http.StatusOK, "register.html", gin.H{
			"error" : "There was an error registering, please try again",
		})
	}
}

// GetLoginPage returns the login page that will be used for presenting the
// available authentication methods
func (a AuthHandler)GetLoginPage(c *gin.Context){

	redirectErr := c.Query("error")
	if redirectErr == ""{
		c.HTML(http.StatusOK, "login.html", nil)
	}else if redirectErr == "AuthError"{
		c.HTML(http.StatusOK, "login.html", gin.H{
			"error" : "There was an error logging you in, please try again",
		})
	}
	
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
			"error": "Password must have at least 8 characters, contain at least 1 uppercase letter, 1 lower case letter,1 number and a special character.",
			"prevUsername": login.Username,
			"prevPassword": login.Password,
		})
		return
	}

	loginResp, err := a.Service.RegisterNativeUser(login.Username, login.Password, "user")
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

	c.SetCookie("access_token", loginResp.AccessToken, int(60 * 60 * 24),"/","localhost",false,true)
	c.SetCookie("refresh_token", loginResp.RefreshToken, int(60 * 60 * 24),"/","localhost",false,true)

	c.Redirect(http.StatusFound,"/")
}

// HandleNativeLogin will log in the users that choose to use a local account
func (a AuthHandler) HandleNativeLogin(c *gin.Context){
	var login LoginForm
	if err := c.Bind(&login); err != nil{
		a.Logger.Error().Err(err).Str("package","handlers").Str("context","HandleNativeLogin").Msg("Error binding login form")
		log.Println("Error binding login form: " + err.Error())
		return
	}
	
	emailRegex := regexp.MustCompile("[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*@(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?")
	if !emailRegex.MatchString(login.Username){
		c.HTML(http.StatusBadRequest, "login.html", gin.H{
			"error": "Enter a valid email address",
			"prevUsername": login.Username,
			"prevPassword": login.Password,
		})
		return
	}

	passwordRegex := regexp.MustCompile(`^(.{0,7}|[^0-9]*|[^A-Z]*|[^a-z]*|[a-zA-Z0-9]*)$`)
	if passwordRegex.MatchString(login.Password){
		c.HTML(http.StatusBadRequest, "login.html", gin.H{
			"error": "Password must have at least 8 characters, contain at least 1 uppercase letter, 1 lower case letter,1 number and a special character.",
			"prevUsername": login.Username,
			"prevPassword": login.Password,
		})
		return
	}

	loginResp, err := a.Service.LoginNativeUser(login.Username, login.Password)
	if err != nil{
		if _,ok := err.(*errs.InvalidCredentials); ok{
			c.HTML(http.StatusBadRequest, "login.html", gin.H{
				"error": "Email or password is incorrect",
				"prevUsername": login.Username,
				"prevPassword": login.Password,
			})
			return
		}else if _,ok := err.(*errs.UserNotFoundError); ok{
			c.HTML(http.StatusBadRequest, "login.html", gin.H{
				"error": "Email does not exist",
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

	c.SetCookie("access_token", loginResp.AccessToken, int(60 * 60 * 24),"/","localhost",false,true)
	c.SetCookie("refresh_token", loginResp.RefreshToken, int(60 * 60 * 24),"/","localhost",false,true)

	c.Redirect(http.StatusFound, "/")
}

// HandleGoogleCode receives the ID token from google's redirect and registers/logs in
// the user and returns the appropriate tokens
func (a AuthHandler) HandleGoogleCode(c *gin.Context){

	// Get the data from google's response
	respData,readErr := io.ReadAll(c.Request.Body)
	if readErr != nil{
		a.Logger.Error().Err(readErr).Str("package","handlers").Str("context","HandleGoogleCode").Msg("Error reading the google id token response")
		c.Redirect(http.StatusFound, "/register?error=AuthError")
		return
	}

	// Get the relevant fields from the response
	config, _ := config.LoadConfig()
	idTokenFields := strings.Split(string(respData), "&")
	clientId, idJWT ,csrfToken := strings.Split(idTokenFields[1],"=")[1], strings.Split(idTokenFields[2],"=")[1], strings.Split(idTokenFields[4],"=")[1]

	// Verify the fields
	csrfCookie,err := c.Request.Cookie("g_csrf_token")
	if err != nil{
		a.Logger.Error().Err(err).Str("package","handlers").Str("context","HandleGoogleCode").Msg("Error getting csrf protection cookie")
		c.Redirect(http.StatusFound, "/register?error=AuthError")
	}
	if csrfToken != csrfCookie.Value{
		a.Logger.Error().Err(err).Str("package","handlers").Str("context","HandleGoogleCode").Msg("CSRF protection cookie and token do not match")
		c.Redirect(http.StatusFound, "/register?error=AuthError")
	}
	if config.GoogleClientID != clientId{
		a.Logger.Error().Err(err).Str("package","handlers").Str("context","HandleGoogleCode").Msg("ID Token client id does not match")
		c.Redirect(http.StatusFound, "/register?error=AuthError")
	}

	// Extract the email from the ID token claims and use it to register/login the user
	tokenClaims, valErr := utils.ValidateGoogleIdToken(idJWT)
	if valErr != nil{
		a.Logger.Error().Err(err).Str("package","handlers").Str("context","HandleGoogleCode").Msg("Error validating google id token")
		c.Redirect(http.StatusFound, "/register?error=AuthError")
	}

	if fmt.Sprint(tokenClaims["iss"]) != "https://accounts.google.com"{
		a.Logger.Error().Err(err).Str("package","handlers").Str("context","HandleGoogleCode").Msg("Unauthorized google id token issuer: received: " + fmt.Sprint((tokenClaims["iss"])))
		c.Redirect(http.StatusFound, "/register?error=AuthError")
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
		a.Logger.Error().Err(err).Str("package","handlers").Str("context","HandleGoogleCode").Msg("Error with registering/logging a google user")
		c.Redirect(http.StatusFound, "/register")
		return
	}

	c.SetCookie("access_token", login.AccessToken, int(60 * 15),"/","localhost",false,true)
	c.SetCookie("refresh_token", login.RefreshToken, int(60 * 60 * 24),"/","localhost",false,true)

	c.Redirect(http.StatusFound, "/")
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
		a.Logger.Error().Str("package","handlers").Str("context","HandleGithubCode").Msg(fmt.Sprintf("Error validating state in github login, expected %s got %s",responseState, c.Query("state")))
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

	userEmailsJson, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		c.Redirect(http.StatusTemporaryRedirect, "/register")
		a.Logger.Error().Err(readErr).Str("package","handlers").Str("context","HandleGithubCode").Msg("Error reading github user emails json")
		return
	}
	var userEmails []GithubEmail
	json.Unmarshal(userEmailsJson,&userEmails)

	var primaryEmail string
	for k := range userEmails {		
		if userEmails[k].Primary{
			primaryEmail = userEmails[k].Email
			break
		}
	}
	
	login, appErr := a.Service.RegisterImportedUser(primaryEmail)
	if _,ok := appErr.(*errs.UserAlreadyExists); ok{
		var newErr error
		login, newErr = a.Service.LoginImportedUser(primaryEmail)
		if newErr != nil{
			a.Logger.Error().Err(newErr).Str("package","handlers").Str("context","HandleGithubCode").Msg("Error logging in imported user")
			return
		}
	}else if appErr != nil{
		a.Logger.Error().Err(appErr).Str("package","handlers").Str("context","HandleGithubCode").Msg("Error with github authentication")
		c.Redirect(http.StatusTemporaryRedirect, "/login")
		return
	}

	c.SetCookie("access_token", login.AccessToken, int(60 * 60 * 15),"/","localhost",false,true)
	c.SetCookie("refresh_token", login.RefreshToken, int(60 * 60 * 24),"/","localhost",false,true)

	c.Redirect(http.StatusFound, "/")
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
		a.Logger.Error().Err(valErr).Str("package","handlers").Str("context","RefreshAccessToken").Msg("Error validating refresh token")
		c.SetCookie("access_token", "", 0,"/","localhost",false,true)
		c.SetCookie("refresh_token", "", 0,"/","localhost",false,true)
		c.Redirect(http.StatusTemporaryRedirect, "/login")
		return
	}
	resp, refErr := a.Service.RefreshTokens(fmt.Sprint(refClaims["id"]),refToken)
	if err != nil{
		a.Logger.Error().Err(refErr).Str("package","handlers").Str("context","RefreshAccessToken").Msg("Error refreshing access token")
		c.Redirect(http.StatusTemporaryRedirect, "/login")
		return
	}
	c.SetCookie("access_token", resp.AccessToken, int(60 * 15),"/","localhost",false,true)
	c.SetCookie("refresh_token", resp.RefreshToken, int(60 * 60 * 24),"/","localhost",false,true)
	c.Redirect(http.StatusTemporaryRedirect, c.Query("redirect"))
}

// LogOut invalidates the request's auth tokens and cookies and removes the account's refresh token from db
func (a AuthHandler) LogOut(c *gin.Context) {
	refToken, err := c.Cookie("refresh_token")
	if refToken == ""{
		c.Redirect(http.StatusTemporaryRedirect, "/login")
		return
	}
	if err != nil {
		a.Logger.Error().Err(err).Str("package","handlers").Str("context","LogOut").Msg("Error logging user out")
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
	session := sessions.Default(c)
	session.Clear()
	c.Redirect(http.StatusFound, c.Query("redirect"))
}

func randToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

