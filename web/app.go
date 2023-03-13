package web

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/robesmi/MSISDNApp/config"
	"github.com/robesmi/MSISDNApp/middleware"
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
	aurepo := repository.NewAuthRepository(dbClient)
	mh := handlers.MSISDNLookupHandler{Service: service.NewMSISDNService(msrepo)}
	ah := handlers.AuthHandler{Service: service.ReturnAuthService(aurepo)}

	//Wiring
	router.LoadHTMLGlob("templates/*.html")
	
	router.GET("/", func(c *gin.Context){
		c.Redirect(http.StatusMovedPermanently, "/login")
	})
	
	
	router.GET("/register", ah.GetRegisterPage)
	router.POST("/register", ah.HandleNativeRegister)
	
	router.GET("/login", ah.GetLoginPage)
	router.POST("/login/", ah.HandleNativeLogin)

	router.GET("/refresh", ah.RefreshAccessToken)
	router.GET("/logout", ah.LogOut)

	

	router.GET("/oauth/google", ah.HandleGoogleLogin)
	router.GET("/oauth/google/callback", ah.HandleGoogleCode)
	router.GET("/oauth/github", ah.HandleGithubLogin)
	router.GET("/oauth/github/callback", ah.HandleGithubCode)


	authorized := router.Group("/service")
	authorized.Use(middleware.ValidateTokenUserSection)
	{
		authorized.GET("/lookup", mh.GetLookupPage)
		authorized.POST("/api/lookup", mh.NumberLookup)
	}

	adminGroup := router.Group("/admin")
	adminGroup.Use(middleware.ValidateTokenAdminSection)
	{
		
	}

	//Starting up server
	config, _ := config.LoadConfig()
	router.Run(":" + config.Port)
}

// getStubDbClient initializes the db connection and returns it to Start
// Using mysql as a placeholder until a solid solution is decided on
func getStubDbClient() *sqlx.DB{

	config, appErr := config.LoadConfig()
	if appErr != nil{
		panic(appErr)
	}

	client, err := sqlx.Open(config.MySqlDriver,config.MySqlSource)
	if err != nil {
		panic(err)
	}
	
	client.SetMaxOpenConns(10)
	client.SetMaxIdleConns(10)
	client.SetConnMaxLifetime(time.Hour)

	return client
}