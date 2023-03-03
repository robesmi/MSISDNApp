package web

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/robesmi/MSISDNApp/repository"
	"github.com/robesmi/MSISDNApp/service"
	"github.com/robesmi/MSISDNApp/web/handlers"
)

// Start initializes the needed route handling, connections between layers and starts the server
func Start(){

	//Setup
	router := gin.Default()
	dbClient := getStubDbClient()
	msrepo := repository.NewMSISDNRepository(dbClient)
	mh := handlers.MSISDNLookupHandler{Service: service.NewMSISDNService(msrepo)}
	ah := handlers.AuthHandler{}

	//Wiring
	router.LoadHTMLGlob("templates/*.html")
	
	router.GET("/", func(c *gin.Context){
		c.Redirect(http.StatusMovedPermanently, "/lookup")
	})
	router.GET("/lookup", mh.GetLookupPage)
	router.POST("/api/lookup", mh.NumberLookup)
	
	router.GET("/login", ah.GetLoginPage)
	router.GET("/oauth/google", ah.HandleGoogleLogin)
	router.GET("/oauth/google/callback", ah.HandleGoogleCode)
	router.GET("/oauth/github", ah.HandleGithubLogin)
	router.GET("/oauth/github/callback", ah.HandleGithubCode)
	


	//Starting up server
	router.Run(":" + os.Getenv("APP_PORT"))
}

// getStubDbClient initializes the db connection and returns it to Start
// Using mysql as a placeholder until a solid solution is decided on
func getStubDbClient() *sqlx.DB{

	client, err := sqlx.Open("mysql","docker:password@tcp(godockerDB)/msisdn")
	if err != nil {
		panic(err)
	}
	
	client.SetMaxOpenConns(10)
	client.SetMaxIdleConns(10)
	client.SetConnMaxLifetime(time.Hour)

	return client
}