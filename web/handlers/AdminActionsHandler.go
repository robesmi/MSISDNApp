package handlers

import (
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/robesmi/MSISDNApp/model"
	"github.com/robesmi/MSISDNApp/model/dto"
	"github.com/robesmi/MSISDNApp/model/errs"
	"github.com/robesmi/MSISDNApp/service"
	"github.com/robesmi/MSISDNApp/vault"
	"github.com/rs/zerolog"
)

type AdminActionsHandler struct {
	AuthService service.AuthService
	MSISDNService service.MSISDNService
	Logger zerolog.Logger
	Vault vault.VaultInterface
}

func (adh AdminActionsHandler) GetAdminPanelPage(c *gin.Context){
	c.HTML(http.StatusOK, "adminpanel.html", nil)
}

func (adh AdminActionsHandler) InsertNewUser(c *gin.Context){

	acReq := dto.AccountRequest{}
	err := c.ShouldBind(&acReq)
	if err != nil{
		adh.Logger.Error().Err(err).Str("package","handlers").Str("context","InsertNewUser").Msg("Error adding new user from admin panel")
		c.HTML(http.StatusBadRequest, "adminpanel.html", gin.H{
			"error": "Error adding new user" + err.Error(),
		})
		return
	}
	
	_, addErr := adh.AuthService.RegisterNativeUser(acReq.Username, acReq.Password, acReq.Role)
	if addErr != nil{
		if _,ok := err.(*errs.UserAlreadyExists); ok{
			c.HTML(http.StatusBadRequest, "adminpanel.html", gin.H{
				"error": "Email already in use",
				"prevUsername": acReq.Username,
				"prevPassword": acReq.Password,
				"prevRole":		acReq.Role,
			})
			return
		}else{
			c.HTML(http.StatusInternalServerError, "adminpanel.html", gin.H{
				"error": "Internal error, please try again " + addErr.Error(),
				"prevUsername": acReq.Username,
				"prevPassword": acReq.Password,
				"prevRole":		acReq.Role,
			})
			return
		}
	}

	c.Redirect(http.StatusFound, "/admin/panel")

}

func (adh AdminActionsHandler) GetAllUsers(c *gin.Context){

	usersList, err := adh.AuthService.GetAllUsers()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "adminpanel.html", gin.H{
			"error": "Internal error: " + err.Error(),
		})
		return
	}
	c.HTML(http.StatusOK, "adminpanel.html", gin.H{
		"users": usersList,
	})
}

func (adh AdminActionsHandler) EditUserPage(c *gin.Context){
	
	userIdParam,err := io.ReadAll(c.Request.Body)
	if err != nil {
		adh.Logger.Error().Err(err).Str("package","handlers").Str("context","EditUserPage").Msg("Error reading POST body")
		c.HTML(http.StatusBadRequest, "adminpanel.html", gin.H{
			"error": "Internal Error: " + err.Error(),
		})
		return
	}
	userId := strings.Split(string(userIdParam),"=")[1]
	user, getErr := adh.AuthService.GetUserById(string(userId))
	if getErr != nil {
		c.HTML(http.StatusBadRequest, "adminpanel.html", gin.H{
			"error": "Internal Error: " + getErr.Error(),
		})
		return
	}

	c.HTML(http.StatusOK, "edituser.html", gin.H{
		"user" : user,
	})
}

func (adh AdminActionsHandler) EditUser(c *gin.Context){

	var editUser model.User
	err := c.ShouldBind(&editUser)
	if err != nil{
		adh.Logger.Error().Err(err).Str("package","handlers").Str("context","EditUser").Msg("Error reading POST body")
		c.Redirect(http.StatusInternalServerError, "/admin/panel")
		return
	}
	// Do the rest of the edit user logic here
	editErr := adh.AuthService.EditUserById(editUser.UUID, editUser.Username, editUser.Password, editUser.Role)
	if editErr != nil{
		adh.Logger.Error().Err(editErr).Str("package","handlers").Str("context","EditUser").Msg("Error editing user")
		c.HTML(http.StatusBadRequest, "edituser.html", gin.H{
			"error": "Internal Error: " + editErr.Error(),
		})
		return
	}
	c.Redirect(http.StatusFound, "/admin/panel")
	
}

func (adh AdminActionsHandler) RemoveUser(c *gin.Context){

	userIdParam,err := io.ReadAll(c.Request.Body)
	if err != nil {
		adh.Logger.Error().Err(err).Str("package","handlers").Str("context","RemoveUser").Msg("Error reading POST body")
		c.HTML(http.StatusBadRequest, "adminpanel.html", gin.H{
			"error": "Internal Error: " + err.Error(),
		})
		return
	}
	
	userId := strings.Split(string(userIdParam),"=")[1]
	decodedId, escErr := url.QueryUnescape(userId)
	if escErr != nil{
		adh.Logger.Error().Err(err).Str("package","handlers").Str("context","RemoveOperator").Msg("Error decoding uri characters")
		c.HTML(http.StatusBadRequest, "adminpanel.html", gin.H{
			"error": "Internal Error: " + err.Error(),
		})
		return
	}
	
	rmErr := adh.AuthService.RemoveUserById(decodedId)
	if rmErr != nil {
		c.HTML(http.StatusBadRequest, "adminpanel.html", gin.H{
			"error": "Internal Error: " + rmErr.Error(),
		})
		return
	}

	c.Redirect( http.StatusFound, "/admin/panel")
}

func (adh AdminActionsHandler) InsertNewCountry(c *gin.Context){

	cReq := dto.CountryRequest{}
	err := c.ShouldBind(&cReq)
	if err != nil{
		adh.Logger.Error().Err(err).Str("package","handlers").Str("context","InsertNewCountry").Msg("Error adding new country from admin panel")
		c.HTML(http.StatusBadRequest, "adminpanel.html", gin.H{
			"error": "Error adding new country" + err.Error(),
		})
		return
	}
	
	numberRegex := regexp.MustCompile(`^\d{1}$`)
	if !numberRegex.MatchString(cReq.CountryCodeLength){
		c.HTML(http.StatusBadRequest, "adminpanel.html", gin.H{
			"error": "The Country Code Length must be a number",
		})
		return
	}
	codeRegex := regexp.MustCompile(`^\d{1,6}$`)
	if !codeRegex.MatchString(cReq.CountryCode){
		c.HTML(http.StatusBadRequest, "adminpanel.html", gin.H{
			"error": "The Country Code can't be longer than 6 digits",
		})
		return
	}
	ciRegex := regexp.MustCompile(`^[a-zA-Z]{2}$`)
	if !ciRegex.MatchString(cReq.CountryIdentifier){
		c.HTML(http.StatusBadRequest, "adminpanel.html", gin.H{
			"error": "The Country identifier must be 2 letters",
		})
		return
	}

	addErr := adh.MSISDNService.AddNewCountry(&cReq)
	if addErr != nil{
		c.HTML(http.StatusInternalServerError, "adminpanel.html", gin.H{
			"error": "Internal error adding new country, please try again " + addErr.Error(),
			"prevCountryRequest" : cReq,
		})
		return
	}

	c.Redirect(http.StatusOK, "/admin/panel")

}

func (adh AdminActionsHandler) GetAllCountries(c *gin.Context){

	countriesList, err := adh.MSISDNService.GetAllCountries()
	if err != nil{
		c.HTML(http.StatusInternalServerError, "adminpanel.html", gin.H{
			"error": "Internal error: " + err.Error(),
		})
		return
	}

	c.HTML(http.StatusOK, "adminpanel.html", gin.H{
		"countries" : countriesList,
	})
}

func (adh AdminActionsHandler) RemoveCountry(c *gin.Context){

	countryPrefixParam, err := io.ReadAll(c.Request.Body)
	if err != nil {
		adh.Logger.Error().Err(err).Str("package","handlers").Str("context","RemoveCountry").Msg("Error reading POST body")
		c.HTML(http.StatusBadRequest, "adminpanel.html", gin.H{
			"error": "Internal Error: " + err.Error(),
		})
		return
	}
	countryPrefix := strings.Split(string(countryPrefixParam),"=")[1]
	decodedPrefix, escErr := url.QueryUnescape(countryPrefix)
	if escErr != nil{
		adh.Logger.Error().Err(err).Str("package","handlers").Str("context","RemoveOperator").Msg("Error decoding uri characters")
		c.HTML(http.StatusBadRequest, "adminpanel.html", gin.H{
			"error": "Internal Error: " + err.Error(),
		})
		return
	}
	
	rmErr := adh.MSISDNService.RemoveCountry(decodedPrefix)
	if rmErr != nil {
		c.HTML(http.StatusBadRequest, "adminpanel.html", gin.H{
			"error": "Internal Error: " + rmErr.Error(),
		})
		return
	}

	c.Redirect( http.StatusFound, "/admin/panel")
}

func (adh AdminActionsHandler) InsertNewMobileOperator(c *gin.Context){

	mnoReq := dto.OperatorRequest{}
	err := c.ShouldBind(&mnoReq)
	if err != nil{
		adh.Logger.Error().Err(err).Str("package","handlers").Str("context","InsertNewMobileOperator").Msg("Error adding new operator from admin panel")
		c.HTML(http.StatusBadRequest, "adminpanel.html", gin.H{
			"error": "Error adding new operator" + err.Error(),
		})
		return
	}


	numberRegex := regexp.MustCompile(`[0-9]{1}`)
	if !numberRegex.MatchString(mnoReq.PrefixLength){
		c.HTML(http.StatusBadRequest, "adminpanel.html", gin.H{
			"error": "The Prefix length must be a number",
		})
		return
	}
	ciRegex := regexp.MustCompile(`[a-zA-Z]{2}`)
	if !ciRegex.MatchString(mnoReq.CountryIdentifier){
		c.HTML(http.StatusBadRequest, "adminpanel.html", gin.H{
			"error": "The Country identifier must be 2 letters",
		})
		return
	}

	addErr := adh.MSISDNService.AddNewMobileOperator(&mnoReq)
	if addErr != nil{
		c.HTML(http.StatusInternalServerError, "adminpanel.html", gin.H{
			"error": "Internal error adding operator, please try again " + addErr.Error(),
			"prevCountryRequest" : mnoReq,
		})
		return
	}

	c.Redirect(http.StatusFound, "/admin/panel")

}

func (adh AdminActionsHandler) GetAllMobileOperators(c *gin.Context){

	operatorsList, err := adh.MSISDNService.GetAllMobileOperators()
	if err != nil{
		c.HTML(http.StatusInternalServerError, "adminpanel.html", gin.H{
			"error": "Internal error: " + err.Error(),
		})
		return
	}

	c.HTML(http.StatusOK, "adminpanel.html", gin.H{
		"operators" : operatorsList,
	})
}

func (adh AdminActionsHandler) RemoveOperator(c *gin.Context){

	operatorPrefixParam, err := io.ReadAll(c.Request.Body)
	if err != nil {
		adh.Logger.Error().Err(err).Str("package","handlers").Str("context","RemoveOperator").Msg("Error reading POST body")
		c.HTML(http.StatusBadRequest, "adminpanel.html", gin.H{
			"error": "Internal Error: " + err.Error(),
		})
		return
	}
	operatorPrefix := strings.Split(string(operatorPrefixParam),"=")[1]
	decodedPrefix, escErr := url.QueryUnescape(operatorPrefix)
	if escErr != nil{
		adh.Logger.Error().Err(err).Str("package","handlers").Str("context","RemoveOperator").Msg("Error decoding uri characters")
		c.HTML(http.StatusBadRequest, "adminpanel.html", gin.H{
			"error": "Internal Error: " + err.Error(),
		})
		return
	}
	rmErr := adh.MSISDNService.RemoveOperator(decodedPrefix)
	if rmErr != nil {
		c.HTML(http.StatusBadRequest, "adminpanel.html", gin.H{
			"error": "Internal Error: " + rmErr.Error(),
		})
		return
	}

	c.Redirect( http.StatusFound, "/admin/panel")
}