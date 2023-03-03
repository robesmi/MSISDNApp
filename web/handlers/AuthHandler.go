package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/robesmi/MSISDNApp/model"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

type AuthHandler struct {
}

// The response token received from the Identity Provider's OAuth token endpoint
type OAuthAccessResponse struct {
	AccessToken string `json:"access_token"`
}

//Configuring credentials for OAuth websites that will be implemented
var redirectUrl = "http://127.0.0.1/8080/oauth/redirect"

var googleConfig = &oauth2.Config{
	Endpoint: google.Endpoint,
	RedirectURL: redirectUrl,
	Scopes: []string{
		"https://www.googleapis.com/auth/userinfo.email",
	},
}
var githubConfig = &oauth2.Config{
	Endpoint: github.Endpoint,
	RedirectURL: redirectUrl,
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
//TODO: generate state dynamically and save it in session
func (a AuthHandler)GetLoginPage(c *gin.Context){
	c.HTML(http.StatusOK, "index.html", nil)
}

// HandleGoogleLogin redirects user to the google oauth2 authorization page
func (a AuthHandler) HandleGoogleLogin(c *gin.Context){
	url := googleConfig.AuthCodeURL("state")
	c.Redirect(http.StatusTemporaryRedirect,url)
}
// HandleGoogleCode receives the access code from google's redirect and makes a post request
// to receive the appropriate access token and refresh token
func (a AuthHandler) HandleGoogleCode(c *gin.Context){
	panic("panic!!")
}
// HandleGithubLogin redirects user to the github oauth2 authorization page
//TODO: generate state dynamically and save it in session
func (a AuthHandler) HandleGithubLogin(c *gin.Context){
	url := githubConfig.AuthCodeURL("state")
	c.Redirect(http.StatusTemporaryRedirect,url)
}
// HandleGithubCode receives the access code from google's redirect and makes a post request
// to receive the appropriate access token and refresh token
func (a AuthHandler) HandleGithubCode(c *gin.Context){
	panic("panic!!")
}

