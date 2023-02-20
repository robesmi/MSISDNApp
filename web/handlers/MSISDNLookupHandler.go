package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strconv"
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
		writeResponse(c, http.StatusBadRequest,"No Number Entered")
	}else{

		number := strings.Trim(input,"0")
		number = strings.Trim(number,"-")
		number = strings.ReplaceAll(number," ","")
		var validNumber = regexp.MustCompile(`[0-9]{7,15}`)
		if !validNumber.MatchString(number){
			writeResponse(c, http.StatusBadRequest, "Invalid Number Entered")
		}else{
		
			convNumber, err := strconv.Atoi(number)
			if err != nil{
				log.Print(err)
				writeResponse(c, http.StatusInternalServerError, "There was an unexpected error")
			}else{
		
				// Execute service layer logic and receive a response
				response, error := msh.Service.LookupMSISDN(uint64(convNumber))
				if error != nil{
					writeResponse(c,error.Code, error.AsMessage())
				}
				// Send response back
				writeResponse(c, http.StatusOK, response)
			}
		}
	}

}

func writeResponse(c *gin.Context,code int, data interface{}){
	c.JSON(code,data)
}
