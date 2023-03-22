package web

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/robesmi/MSISDNApp/config"
	"github.com/robesmi/MSISDNApp/middleware"
	"github.com/robesmi/MSISDNApp/repository"
	"github.com/robesmi/MSISDNApp/service"
	"github.com/robesmi/MSISDNApp/web/handlers"
	"github.com/rs/zerolog"
)

// Start initializes the needed route handling, connections between layers and starts the server
func Start(){

	//Setup
	router := gin.New()
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d \"%s\" %s\"\n",
				param.ClientIP,
				param.TimeStamp.Format(time.RFC1123),
				param.Method,
				param.Path,
				param.Request.Proto,
				param.StatusCode,
				param.Request.UserAgent(),
				param.ErrorMessage,
		)
	}))
	router.Use(gin.Recovery())


	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	dbClient := getStubDbClient()
	msrepo := repository.NewMSISDNRepository(dbClient)
	aurepo := repository.NewAuthRepository(dbClient)
	mh := handlers.MSISDNLookupHandler{Service: service.NewMSISDNService(msrepo), Logger: logger}
	ah := handlers.AuthHandler{Service: service.ReturnAuthService(aurepo), Logger: logger}
	aph := handlers.AuthApiHandler{Service: service.ReturnAuthService(aurepo)}
	adh := handlers.AdminActionsHandler{AuthService: service.ReturnAuthService(aurepo), 
		MSISDNService: service.NewMSISDNService(msrepo), 
		Logger: logger,}

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

	userSection := router.Group("/service")
	userSection.Use(middleware.ValidateTokenUserSection)
	
	{
		userSection.GET("/lookup", mh.GetLookupPage)
		userSection.POST("/lookup", mh.NumberLookup)
	}

	adminSection := router.Group("/admin")
	adminSection.Use(middleware.ValidateTokenAdminSection)
	{
		adminSection.GET("/panel", adh.GetAdminPanelPage)

		adminSection.POST("/adduser", adh.InsertNewUser)
		adminSection.POST("/edituserpanel", adh.EditUserPage)
		adminSection.POST("/edituser", adh.EditUser)
		adminSection.POST("/removeuser", adh.RemoveUser)

		adminSection.POST("/addcountry", adh.InsertNewCountry)
		adminSection.POST("/removecountry", adh.RemoveCountry)

		adminSection.POST("/addoperator", adh.InsertNewMobileOperator)
		adminSection.POST("/removeoperator", adh.RemoveOperator)
	
		adminSection.POST("/getusers", adh.GetAllUsers)
		adminSection.POST("/getcountries", adh.GetAllCountries)
		adminSection.POST("/getoperators", adh.GetAllMobileOperators)
	}

	router.NoRoute( func(c *gin.Context){
		c.HTML(http.StatusNotFound, "notfound.html",nil)
	})

	
	//Starting up server
	config, _ := config.LoadConfig()

	_, err := ah.Service.RegisterNativeUser(config.AdminUsername, config.AdminPassword, "admin")
	if err != nil{
		logger.Err(err).Str("package","web").Str("context","init").Msg("Error during init")
	}
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