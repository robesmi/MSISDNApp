package handlers

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/robesmi/MSISDNApp/service"
)

type MSISDNLookupHandler struct {
	Service service.MSISDNService
}

func (msh MSISDNLookupHandler) NumberLookup(c *gin.Context){


	// Check for empty input, trim and validate number
	input := c.Query("number")
	if input == ""{
		//If no input, just serve up the page
		c.HTML(http.StatusOK,"index.html",nil)
	}else{
		
		number := strings.Trim(input,"0")
		number = strings.Trim(number,"-")
		number = strings.TrimSpace(number)
		number = strings.ReplaceAll(number," ","")
		var validNumberRegex = regexp.MustCompile(`^[0-9]{7,15}$`)
		if !validNumberRegex.MatchString(number){
			writeResponse(c, http.StatusBadRequest, "Invalid Number Entered")

		}else{

			// Execute service layer logic and receive a response
			response, error := msh.Service.LookupMSISDN(number)
			if error != nil{
				writeResponse(c,error.Code, error.AsMessage())
			}
			
			// Send response back
			writeResponse(c, http.StatusOK, response)
		}
	}
}


func writeResponse(c *gin.Context,code int, data interface{}){
	c.JSON(code,data)
}
