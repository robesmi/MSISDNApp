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

type LookupRequest struct{
	Number string `json:"number" xml:"number"`
}

func (msh MSISDNLookupHandler) NumberLookup(c *gin.Context){


	// Check for empty input, trim and validate number
	var req LookupRequest
	if err := c.ShouldBind(&req); err != nil{
		writeResponse(c, http.StatusBadRequest, "Bad Request")
		return
	}
	if req.Number == ""{
		//If no input, just serve up the page
		c.HTML(http.StatusOK,"index.html",nil)
		return
	}
		
		number := strings.Trim(req.Number,"0")
		m1 := regexp.MustCompile(`\D`)
		number = m1.ReplaceAllString(number,"")
		var validNumberRegex = regexp.MustCompile(`^[0-9]{7,15}$`)
		if !validNumberRegex.MatchString(number){
			writeResponse(c, http.StatusBadRequest, "Invalid Number Entered")
			return

		}

			// Execute service layer logic and receive a response
			response, error := msh.Service.LookupMSISDN(number)
			if error != nil{
				writeResponse(c,error.Code, error.AsMessage())
				return
			}
			
			// Send response back
			writeResponse(c, http.StatusOK, response)
			return
		
	

}


func writeResponse(c *gin.Context,code int, data interface{}){
	c.JSON(code,data)
}
