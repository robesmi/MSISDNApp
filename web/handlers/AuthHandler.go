package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/robesmi/MSISDNApp/model"
	"github.com/robesmi/MSISDNApp/service"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

type AuthHandler struct {
	Service service.AuthService
}

// The response token received from the Identity Provider's OAuth token endpoint
type OAuthAccessResponse struct {
	AccessToken string `json:"access_token"`
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
	file,err := os.Open("docker/oauthcred.txt")
	if err != nil {
		log.Fatalf(err.Error())
	}
	defer file.Close()
	var credentials = make(map[string]model.OauthCredentials,2)
	json.NewDecoder(file).Decode(&credentials)
	googleCreds := credentials["google"]
	githubCreds := credentials["github"]

	googleConfig.ClientID = googleCreds.ClientId
	googleConfig.ClientSecret = googleCreds.ClientSecret

	githubConfig.ClientID = githubCreds.ClientId
	githubConfig.ClientSecret = githubCreds.ClientSecret

}

// GetLoginPage returns the login page that will be used for presenting the
// available authentication methods
func (a AuthHandler)GetLoginPage(c *gin.Context){
	c.HTML(http.StatusOK, "index.html", nil)
}

// HandleNativeLogin will log in the users that choose to use a local account
func (a AuthHandler) HandleNativeLogin(c *gin.Context){
	panic("panic!!")
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
// to receive the appropriate access token and refresh token
func (a AuthHandler) HandleGoogleCode(c *gin.Context){

	session := sessions.Default(c)
	responseState := session.Get("state")
	if responseState != c.Query("state"){
		c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("Invalid session state %s",responseState))
		return
	}
	token, err := googleConfig.Exchange(oauth2.NoContext, c.Query("code"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest,err)
		return
	}
	client := googleConfig.Client(oauth2.NoContext,token)

	resp,err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	defer resp.Body.Close()

	panic("panic!!")
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
		c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("Invalid session state %s",responseState))
		return
	}
	token, err := googleConfig.Exchange(oauth2.NoContext, c.Query("code"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest,err)
		return
	}
	client := googleConfig.Client(oauth2.NoContext,token)

	resp,err := client.Get("https://api.github.com/user/emails")
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	defer resp.Body.Close()

	panic("panic!!")
}

func randToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

