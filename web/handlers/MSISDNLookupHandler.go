package handlers

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/robesmi/MSISDNApp/service"
	"github.com/rs/zerolog"
)

type MSISDNLookupHandler struct {
	Service service.MSISDNService
	Logger zerolog.Logger
}
type LookupRequest struct {
	Number string `form:"number"`
}
type ApiLookupRequest struct{
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
		c.HTML(http.StatusBadRequest, "index.html", gin.H{
			"error" : "Unexpected error",
		})
		return
	}
	if req.Number == ""{
		//If no input
		c.HTML(http.StatusBadRequest, "index.html", gin.H{
			"error" : "Please enter a MSISDN",
		})
		return
	}

		number := strings.TrimLeft(req.Number,"0")
		m1 := regexp.MustCompile(`\D`)
		number = m1.ReplaceAllString(number,"")
		var validNumberRegex = regexp.MustCompile(`^[0-9]{7,15}$`)
		if !validNumberRegex.MatchString(number){
			c.HTML(http.StatusBadRequest, "index.html", gin.H{
				"error" : "The MSISDN must only contain digits and be 7-15 digits long",
			})
			return
		}

			// Execute service layer logic and receive a response
			response, lookupErr := msh.Service.LookupMSISDN(number)
			if lookupErr != nil{
				msh.Logger.Error().Err(lookupErr).Str("package","handlers").Str("context","NumberLookupApi").Msg("Error making lookup")
				c.HTML(http.StatusBadRequest, "index.html", gin.H{
					"error" : lookupErr.Error(),
				})
				return
			}
			
			// Send response back
			c.HTML(http.StatusOK, "index.html", gin.H{
				"mno" : response.MNO,
				"cc" :	response.CC,
				"sn":	response.SN,
				"ci": response.CI,
			})
}

func (msh MSISDNLookupHandler) NumberLookupApi(c *gin.Context){

	// Check for empty input, trim and validate number
	var req ApiLookupRequest
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
			response, lookupErr := msh.Service.LookupMSISDN(number)
			if lookupErr != nil{
				msh.Logger.Error().Err(lookupErr).Str("package","handlers").Str("context","NumberLookupApi").Msg("Error making lookup")
				writeResponse(c,http.StatusBadRequest, map[string]string{ "error": lookupErr.Error()})
				return
			}
			
			// Send response back
			writeResponse(c, http.StatusOK, response)
}


func writeResponse(c *gin.Context,code int, data interface{}){
	c.JSON(code,data)
}
