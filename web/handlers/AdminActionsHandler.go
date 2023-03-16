package handlers

import (
	"log"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/robesmi/MSISDNApp/model/dto"
	"github.com/robesmi/MSISDNApp/model/errs"
	"github.com/robesmi/MSISDNApp/service"
)

type AdminActionsHandler struct {
	AuthService service.AuthService
	MSISDNService service.MSISDNService
}

func (adh AdminActionsHandler) GetAdminPanelPage(c *gin.Context){
	c.HTML(http.StatusOK, "adminpanel.html", nil)
}

func (adh AdminActionsHandler) InsertNewUser(c *gin.Context){

	acReq := dto.AccountRequest{}
	err := c.ShouldBind(&acReq)
	if err != nil{
		log.Printf("Error adding new user from admin panel: " + err.Error())
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

func (adh AdminActionsHandler) InsertNewCountry(c *gin.Context){

	cReq := dto.CountryRequest{}
	err := c.ShouldBind(&cReq)
	if err != nil{
		log.Printf("Error adding new country from admin panel: " + err.Error())
		c.HTML(http.StatusBadRequest, "adminpanel.html", gin.H{
			"error": "Error adding new country" + err.Error(),
		})
		return
	}
	
	numberRegex := regexp.MustCompile(`[0-9]{1}`)
	if !numberRegex.MatchString(cReq.CountryCodeLength){
		c.HTML(http.StatusBadRequest, "adminpanel.html", gin.H{
			"error": "The Country Code Length must be a number",
		})
		return
	}
	codeRegex := regexp.MustCompile(`[0-9]{6}`)
	if !codeRegex.MatchString(cReq.CountryCode){
		c.HTML(http.StatusBadRequest, "adminpanel.html", gin.H{
			"error": "The Country Code can't be longer than 6 digits",
		})
		return
	}
	ciRegex := regexp.MustCompile(`[a-zA-Z]{2}`)
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

func (adh AdminActionsHandler) InsertNewMobileOperator(c *gin.Context){

	mnoReq := dto.OperatorRequest{}
	err := c.ShouldBind(&mnoReq)
	if err != nil{
		log.Printf("Error adding new operator from admin panel: " + err.Error())
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