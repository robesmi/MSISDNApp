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

func (msh MSISDNLookupHandler) GetLookupPage(c *gin.Context){
	c.HTML(http.StatusOK, "index.html",nil)
}

func (msh MSISDNLookupHandler) NumberLookup(c *gin.Context){

	// Check for empty input, trim and validate number
	var req LookupRequest
	if err := c.ShouldBind(&req); err != nil{
		writeResponse(c, http.StatusBadRequest, map[string]string{ "error":"Enter a proper MSISDN ex: 38977123456"})
		return
	}
	if req.Number == ""{
		//If no input
		writeResponse(c,http.StatusBadRequest,map[string]string{ "error":"Please enter a MSISDN"})
		return
	}

		number := strings.Trim(req.Number,"0")
		m1 := regexp.MustCompile(`\D`)
		number = m1.ReplaceAllString(number,"")
		var validNumberRegex = regexp.MustCompile(`^[0-9]{7,15}$`)
		if !validNumberRegex.MatchString(number){
			writeResponse(c, http.StatusBadRequest, map[string]string{ "error":"The MSISDN must only contain digits and be 7-15 digits long"})
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
}


func writeResponse(c *gin.Context,code int, data interface{}){
	c.JSON(code,data)
}
