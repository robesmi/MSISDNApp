package web

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"
  	"github.com/gin-contrib/sessions/cookie"
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
	aph := handlers.AuthApiHandler{Service: service.ReturnAuthService(aurepo)}
	adh := handlers.AdminActionsHandler{AuthService: service.ReturnAuthService(aurepo), 
		MSISDNService: service.NewMSISDNService(msrepo)}

	//Wiring
	router.LoadHTMLGlob("templates/*.html")
	
	router.GET("/", mh.GetMainPage)
	
	store := cookie.NewStore([]byte("secret"))
  	router.Use(sessions.Sessions("mysession", store))

	router.GET("/register", ah.GetRegisterPage)
	router.POST("/register", ah.HandleNativeRegister)

	router.GET("/login", ah.GetLoginPage)
	router.POST("/login", ah.HandleNativeLogin)

	router.GET("/refresh", ah.RefreshAccessToken)
	router.POST("/refresh", func(c *gin.Context){
		c.Redirect(http.StatusTemporaryRedirect, "/refresh")
	})

	router.GET("/logout", ah.LogOut)
	router.POST("/logout", func(c *gin.Context){
		c.Redirect(http.StatusTemporaryRedirect, "/logout")
	})

	router.POST("/oauth/google/callback", ah.HandleGoogleCode)
	router.GET("/oauth/google/callback", func(c *gin.Context){
		c.Redirect(http.StatusTemporaryRedirect,"/login")
	})

	router.GET("/oauth/github", ah.HandleGithubLogin)
	router.POST("/oauth/github", func(c *gin.Context){
		c.Redirect(http.StatusTemporaryRedirect,"/oauth/github")
	})
	router.GET("/oauth/github/callback", ah.HandleGithubCode)
	router.POST("/oauth/github/callback", func(c *gin.Context){
		c.Redirect(http.StatusTemporaryRedirect,"/login")
	})

	router.POST("/api/register", aph.HandleNativeRegisterCall)
	router.POST("/api/login", aph.HandleNativeLoginCall)
	router.POST("/api/refresh", aph.RefreshAccessTokenCall)
	router.POST("/api/logout", aph.LogOutCall)

	router.POST("/service/api/lookup", middleware.ValidateApiTokenUserSection, mh.NumberLookupApi)

	authorized := router.Group("/service")
	authorized.Use(middleware.ValidateTokenUserSection)
	
	{
		authorized.GET("/lookup", mh.GetLookupPage)
		authorized.POST("/lookup", mh.NumberLookup)
		
	}

	router.GET("/admin/panel", adh.GetAdminPanelPage)

	router.POST("/admin/adduser", adh.InsertNewUser)
	router.POST("/admin/addcountry", adh.InsertNewCountry)
	router.POST("/admin/addoperator", adh.InsertNewMobileOperator)
	
	router.POST("/admin/getusers", adh.GetAllUsers)
	router.POST("/admin/getcountries", adh.GetAllCountries)
	router.POST("/admin/getoperators", adh.GetAllMobileOperators)

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