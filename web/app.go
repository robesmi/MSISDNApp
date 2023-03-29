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
	"github.com/robesmi/MSISDNApp/middleware"
	"github.com/robesmi/MSISDNApp/repository"
	"github.com/robesmi/MSISDNApp/service"
	"github.com/robesmi/MSISDNApp/web/handlers"
	"github.com/rs/zerolog"
	vaultapi "github.com/hashicorp/vault/api"
	"github.com/robesmi/MSISDNApp/vault"

)

// Start initializes the needed route handling, connections between layers and starts the server
func Start(){

	//Setting up gin router, recovery and logging middleware
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

	// Setup the client for interacting with the vault
	vault_config := vaultapi.DefaultConfig()

	vAddr, set := os.LookupEnv("VAULT_ADDR")
	if !set{
		logger.Error().Msg("VAULT_ADDR is not set in env")
		os.Exit(1)
	}else{
		vault_config.Address = vAddr
	}
	
	token,ok := os.LookupEnv("MY_VAULT_TOKEN")
	if !ok {
		logger.Error().Msg("MY_VAULT_TOKEN is not set in env")
		os.Exit(1)
	}

	client, vaultErr  := vault.New(vault_config, token)
	if vaultErr != nil {
		logger.Error().Err(vaultErr).Str("package","web").Str("context","Start").Msg("Error starting vault client")
	}

	// Immediately get some variables that will be needed for setup
	startupVars, fetchErr := client.Fetch("appvars", "Secret", "PORT")
	if fetchErr != nil{
		logger.Error().Err(fetchErr).Str("package","web").Str("context","Start").Msg("Error getting startup variables")
	}
	
	// Setup the db connection along with initializing the layers
	dbClient := getDbClient(client, &logger)
	msrepo := repository.NewMSISDNRepository(dbClient)
	aurepo := repository.NewAuthRepository(dbClient)
	mh := handlers.MSISDNLookupHandler{Service: service.NewMSISDNService(msrepo), Logger: logger}
	//ah := handlers.AuthHandler{Service: service.ReturnAuthService(aurepo), Logger: logger, Vault: client}
	ah := handlers.NewAuthHandler(service.ReturnAuthService(aurepo, client), logger, client)
	aph := handlers.AuthApiHandler{Service: service.ReturnAuthService(aurepo, client), Vault: client}
	adh := handlers.AdminActionsHandler{AuthService: service.ReturnAuthService(aurepo, client), MSISDNService: service.NewMSISDNService(msrepo), Logger: logger, Vault: client}

	//Wiring
	router.LoadHTMLGlob("templates/*.html")

	store := cookie.NewStore([]byte(startupVars["Secret"]))
  	router.Use(sessions.Sessions("mysession", store))
	
	router.GET("/", mh.GetMainPage)

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

	router.POST("/service/api/lookup", middleware.ValidateApiTokenUserSection(client), mh.NumberLookupApi)

	userSection := router.Group("/service")
	userSection.Use(middleware.ValidateTokenUserSection(client))
	
	{
		userSection.GET("/lookup", mh.GetLookupPage)
		userSection.POST("/lookup", mh.NumberLookup)
	}

	adminSection := router.Group("/admin")
	adminSection.Use(middleware.ValidateTokenAdminSection(client))
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

		adminSection.POST("/getsecrets", adh.GetAllSecrets)
	}

	router.NoRoute( func(c *gin.Context){
		c.HTML(http.StatusNotFound, "notfound.html", nil)
	})

	// Initialize an admin user
	user, userErr := client.Fetch("superuser", "AdminUsername", "AdminPassword")
	if userErr != nil{
		logger.Err(userErr).Str("package","web").Str("context","init").Msg("Error Error fetching admin credentials from vault")
	}
	_, regErr := ah.Service.RegisterNativeUser(user["AdminUsername"], user["AdminPassword"], "admin")
	if regErr != nil{
		logger.Err(regErr).Str("package","web").Str("context","init").Msg("Error during init")
	}

	//Starting up server
	router.Run(":" + startupVars["PORT"])
}

// getDbClient initializes the db connection and returns it to Start
func getDbClient(vault vault.VaultInterface, logger *zerolog.Logger) *sqlx.DB{

	dbCreds, fetchErr := vault.Fetch("appvars", "MYSQL_DRIVER", "MYSQL_SOURCE")
	if fetchErr != nil {
		logger.Error().Err(fetchErr).Str("package","web").Str("context","getDbClient").Msg("Error getting db details from vault")
	}


	client, err := sqlx.Open(dbCreds["MYSQL_DRIVER"],dbCreds["MYSQL_SOURCE"])
	if err != nil {
		logger.Error().Err(err).Str("package","web").Str("context","getDbClient").Msg("Error opening db connection")
		
	}
	
	client.SetMaxOpenConns(10)
	client.SetMaxIdleConns(10)
	client.SetConnMaxLifetime(time.Hour)

	return client
}