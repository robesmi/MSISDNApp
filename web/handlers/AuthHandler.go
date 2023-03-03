package handlers

import (
	"encoding/json"
	"log"
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
func (a AuthHandler)GetLoginPage(c *gin.Context){
	panic("panic!!")
}

// GetAccessCode receives the authorization code from the OAuth2 workflow and
// makes a post request to the Identity Provider's oauth token endpoint
// to receive an access token and refresh token
func (a AuthHandler)GetAccessCode(c *gin.Context){
	panic("panic!!")
}