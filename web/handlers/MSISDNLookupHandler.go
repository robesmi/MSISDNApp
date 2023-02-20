package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/robesmi/MSISDNApp/service"
)

type MSISDNLookupHandler struct {
	Service service.MSISDNService
}

func (msh MSISDNLookupHandler) NumberLookup(c *gin.Context){
	panic("panic!!")
}

func writeResponse(c *gin.Context,code int, data interface{}){
	c.JSON(code,data)
}

