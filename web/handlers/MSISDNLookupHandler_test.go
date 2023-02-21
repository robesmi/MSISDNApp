package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/robesmi/MSISDNApp/mocks/service"
	"github.com/robesmi/MSISDNApp/model/errs"
)


var ctx *gin.Context
var router *gin.Engine
var lh MSISDNLookupHandler
var mockService *service.MockMSISDNService

func setup(t *testing.T, w *httptest.ResponseRecorder) func(){
	
	ctrl := gomock.NewController(t)
	mockService = service.NewMockMSISDNService(ctrl)
	lh = MSISDNLookupHandler{mockService}

	gin.SetMode(gin.TestMode)
	ctx, router = gin.CreateTestContext(w)
	router.GET("/lookup", lh.NumberLookup)


	return func() {
		ctx = nil
		router = nil
		defer ctrl.Finish()
	}

}

func TestEmptyInput(t *testing.T) {

	//Arrange
	recorder := httptest.NewRecorder()
	teardown := setup(t,recorder)
	defer teardown()

	//Act
	req, err := http.NewRequest(http.MethodGet,"/lookup",nil)
	if err != nil{
		t.Fatalf("Could not make request")
	}
	router.ServeHTTP(recorder,req)

	//Assert
	if recorder.Code != http.StatusBadRequest{
		t.Error("Failed while testing invalid input value")
	}
}

func TestSendingNegativeNumberAsInput(t *testing.T) {

	//Arrange
	recorder := httptest.NewRecorder()
	teardown := setup(t,recorder)
	defer teardown()
	input := "212315231"
	mockService.EXPECT().LookupMSISDN(input).Return(nil, errs.UnexpectedError("Invalid input"))

	//Act
	req, err := http.NewRequest(http.MethodGet,fmt.Sprintf("/lookup?number=-%v",input),nil)
	if err != nil{
		t.Fatalf("Could not make request")
	}
	router.ServeHTTP(recorder,req)

	//Assert
	if recorder.Code != http.StatusInternalServerError{
		t.Error("Failed while testing invalid input value")
	}
}

func TestInvalidNumber(t *testing.T) {

	//Arrange
	recorder := httptest.NewRecorder()
	teardown := setup(t,recorder)
	defer teardown()
	input := "lorem ipsum"

	//Act
	req, err := http.NewRequest(http.MethodGet,fmt.Sprintf("/lookup?number=-%v",input),nil)
	if err != nil{
		t.Fatalf("Could not make request")
	}
	router.ServeHTTP(recorder,req)

	//Assert
	if recorder.Code != http.StatusBadRequest{
		t.Error("Failed while testing invalid input value")
	}
}
func TestMassiveNumber(t *testing.T) {

	//Arrange
	recorder := httptest.NewRecorder()
	teardown := setup(t,recorder)
	defer teardown()
	var input  = "2371289370190232132"
	mockService.EXPECT().LookupMSISDN(input).Return(nil, errs.UnexpectedError("Invalid input"))

	//Act
	req, err := http.NewRequest(http.MethodGet,fmt.Sprintf("/lookup?number=-%v",input),nil)
	if err != nil{
		t.Fatalf("Could not make request")
	}
	router.ServeHTTP(recorder,req)

	//Assert
	if recorder.Code != http.StatusInternalServerError{
		t.Error("Failed while testing invalid input value")
	}
}

func TestSendingABase64ImageAsInput(t *testing.T) {

	//Arrange
	recorder := httptest.NewRecorder()
	teardown := setup(t,recorder)
	defer teardown()
	input := "/9j/4AAQSkZJRgABAQAAAQABAAD/2wBDAAUDBAQEAwUEBAQFBQUGBwwIBwcHBw8LCwkMEQ8SEhEPERETFhwXExQaFRERGCEYGh0dHx8fExciJCIeJBweHx7/2wBDAQUFBQcGBw4ICA4eFBEUHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh7/wAARCAFUAgoDASIAAhEBAxEB/8QAHAAAAQUBAQEAAAAAAAAAAAAABAABAgMFBgcI/8QATRAAAQMCBAQCBwUECAMFCQEAAQACAwQRBRIhYQYTMVFBcQcUIjKBkaEjQlKxwUNigtEVJDNTcpKi4QjC8CY0Y6OyFiUnNTZEZIOz8f/EABsBAAMBAQEBAQAAAAAAAAAAAAABAgMEBQYH/8QAMBEBAQACAgIBAgQDCQEBAAAAAAECEQMSBCExIkEFEzJRcZGhBiMzQmGBscHwFOH/2gAMAwEAAhEDEQA/APl2yVtlOyVl7TnQ1SU7JiD2TIyVk9krJgk4SsEkyJOkkgqcJJJIBkkkyYOnTJBMJBOFEdU4TCScFRBToJJOFEFOgkkk10kge6dMkgHumSSQRBOmTpgkkkkAkk6VkEZJPZKyNgySeyVtkgilZTtsmtsgI2CVgp2KSRI2TWU8qfKkELbJ7KeXZOGpbCuyfKrAxSDNkthTlT5VdkSyJBQWqJbsiSxRLNktmHypsqILE2RI1BaolqJLFEsQYYsUS1FFiiWI2YbKmyoksUciNmHypi1EFijkQew5amLUQWJi1Bhy1NlPZX5Usp7I2FmVLKrcpSypltVlTZVdZNZMKi1NZWkbJiEBVZKysITWTG0LJ7KVk1tkyMknslZMIplKxTJgyQSSTB06ZOmDhOCmCdAOnCYJ0JOnTJJUHTqKcII6SSVkyJOkE4CQIJwlZSAQEbJwFIN/6KcBAQsntuphqcBBK7JZfNWZUrJBWAkRsFZZNZGzRslZTyqQakFYapBqmGqYYlslYYpBmyta1SDFOwqDFIMVwYrAxLZbDcvZS5aJDFIM2S2NhOWmMaM5aXKS2NgTHsm5aP5WyYwpbMBy0xjR5hUTFsls9gTGomNHcpMYtkbPYEx7KJj2Rxi2UTEjZ7A8tNy0cY1ExI2NgjGomNGmJRMaNnsEY03LRpiUeWjY2hkSLFflSy7LQth8h7KJYisqbIg9hciiWosxqDo9kxsNlTFqILFEtQagtTWVxamLUwqITWVhamLU4FdkxCsskQmFRCVtlMtSsmEEgpWSsmDBOlZPbsmCThKydCSST2SskVIBOkAnsgEAnATgKQCCRAUgFIBSDUiQAUgFMNUg1AQDU4arMuycNQFeVPlVgbspBuyAqypFuiuypZNkgoyp8hV4YpCNIBwwqQYiBEeymIktgMGKxrNl1nBPAmPcWiaXDIqeKlhOV9TUvLI8/wCBtgS53S9hYX1IQ3FXC2L8L4i2ixenax0jC+GWRU1YOqsZ1SSU1NTCSSSRUiolJJBIO6qtySSCQKiUkkGiU3ikkgHTDqkkpBDqpBJJSDhOEkkESSSSAiVA9UklUCDuqrckkrgVvVLkklpFxS5VOSSWkWqd1VZSSWkVDFR8UkkzMUkkkKIJXKSSk4fqmPRJJIRByrKSSDVuKrckkkuKnqpySSqLil/ZdHwv70XmkkuLzv0w8/01v+mDXAsLHhn/AEXmB8Eklt4n+Gnx/wBEIJvBJJdLYkkkkGSQSSQCSSSQH//Z"

	//Act
	req, err := http.NewRequest(http.MethodGet,fmt.Sprintf("/lookup?number=-%v",input),nil)
	if err != nil{
		t.Fatalf("Could not make request")
	}
	router.ServeHTTP(recorder,req)

	//Assert
	if recorder.Code != http.StatusBadRequest{
		t.Error("Failed while testing invalid input value")
	}
}