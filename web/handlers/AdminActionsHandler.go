package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/robesmi/MSISDNApp/model"
	"github.com/robesmi/MSISDNApp/model/errs"
	"github.com/robesmi/MSISDNApp/service"
)

type AdminActionsHandler struct {
	AuthService service.AuthService
	MSISDNService service.MSISDNService
}

type AccountRequest struct {
	Username 	string	`form:"username"`
	Password 	string	`form:"password"`
	Role 		string	`form:"role"`
}

type AccountList struct {
	Accounts model.User
}

func (adh AdminActionsHandler) GetAdminPanelPage(c *gin.Context){
	c.HTML(http.StatusOK, "adminpanel.html", nil)
}

func (adh AdminActionsHandler) InsertNewUser(c *gin.Context){

	acReq := AccountRequest{}
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

	c.HTML(http.StatusOK, "adminpanel.html", gin.H{
		"status" : "Success",
	})

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