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

func (msh MSISDNLookupHandler) GetMainPage(c *gin.Context){
	c.HTML(http.StatusOK, "home.html", nil)
}

func (msh MSISDNLookupHandler) GetLookupPage(c *gin.Context){
	c.HTML(http.StatusOK, "index.html",nil)
}

func (msh MSISDNLookupHandler) NumberLookup(c *gin.Context){

	// Check for empty input, trim and validate number
	var req LookupRequest
	if err := c.ShouldBind(&req); err != nil{
		writeResponse(c, http.StatusBadRequest, map[string]string{ "error":"API call type should be string"})
		return
	}
	if req.Number == ""{
		//If no input
		writeResponse(c,http.StatusBadRequest,map[string]string{ "error":"Please enter a MSISDN"})
		return
	}

		number := strings.TrimLeft(req.Number,"0")
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
				writeResponse(c,http.StatusBadRequest, map[string]string{ "error": error.Error()})
				return
			}
			
			// Send response back
			writeResponse(c, http.StatusOK, response)
}


func writeResponse(c *gin.Context,code int, data interface{}){
	c.JSON(code,data)
}
