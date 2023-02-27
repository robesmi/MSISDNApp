package handlers

import (
	"bytes"
	"encoding/json"
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
	router.POST("/lookup", lh.NumberLookup)


	return func() {
		ctx = nil
		router = nil
		defer ctrl.Finish()
	}

}

func TestNumberLookup(t *testing.T) {

	tt := []struct{
		Name string
		Input string
		TestErrorMessage string
		ExpectedReturnCode int
		CallsService bool

	}{
		{
			Name: 				"Test Empty Input",
			Input: 				"",
			TestErrorMessage: 	"Failed while testing empty input value",
			ExpectedReturnCode:	http.StatusBadRequest,
			CallsService: 		false,
		},
		{
			Name: 				"Test Negative input",
			Input:				"-212315231",
			TestErrorMessage: 	"Failed while testing negative input value",
			ExpectedReturnCode:	http.StatusInternalServerError,
			CallsService:		true,
		},
		{
			Name:				"Test Invalid Number",
			Input:				"lorem ipsum",
			TestErrorMessage: 	"Failed while testing invalid input",
			ExpectedReturnCode: http.StatusBadRequest,
			CallsService:		false,

		},
		{
			Name: 				"Test Massive Number",
			Input:				"237128937019023213123121232",
			TestErrorMessage:	"Failed while testing massive input value",
			ExpectedReturnCode: http.StatusBadRequest,
			CallsService:		false,
		},
		{
			Name:				"Test Sending a base 64 image as input",
			Input:				"/9j/4AAQSkZJRgABAQAAAQABAAD/2wBDAAUDBAQEAwUEBAQFBQUGBwwIBwcHBw8LCwkMEQ8SEhEPERETFhwXExQaFRERGCEYGh0dHx8fExciJCIeJBweHx7/2wBDAQUFBQcGBw4ICA4eFBEUHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh7/wAARCAFUAgoDASIAAhEBAxEB/8QAHAAAAQUBAQEAAAAAAAAAAAAABAABAgMFBgcI/8QATRAAAQMCBAQCBwUECAMFCQEAAQACAwQRBRIhYQYTMVFBcQcUIjKBkaEjQlKxwUNigtEVJDNTcpKi4QjC8CY0Y6OyFiUnNTZEZIOz8f/EABsBAAMBAQEBAQAAAAAAAAAAAAABAgMEBQYH/8QAMBEBAQACAgIBAgQDCQEBAAAAAAECEQMSBCExIkEFEzJRcZGhBiMzQmGBscHwFOH/2gAMAwEAAhEDEQA/APl2yVtlOyVl7TnQ1SU7JiD2TIyVk9krJgk4SsEkyJOkkgqcJJJIBkkkyYOnTJBMJBOFEdU4TCScFRBToJJOFEFOgkkk10kge6dMkgHumSSQRBOmTpgkkkkAkk6VkEZJPZKyNgySeyVtkgilZTtsmtsgI2CVgp2KSRI2TWU8qfKkELbJ7KeXZOGpbCuyfKrAxSDNkthTlT5VdkSyJBQWqJbsiSxRLNktmHypsqILE2RI1BaolqJLFEsQYYsUS1FFiiWI2YbKmyoksUciNmHypi1EFijkQew5amLUQWJi1Bhy1NlPZX5Usp7I2FmVLKrcpSypltVlTZVdZNZMKi1NZWkbJiEBVZKysITWTG0LJ7KVk1tkyMknslZMIplKxTJgyQSSTB06ZOmDhOCmCdAOnCYJ0JOnTJJUHTqKcII6SSVkyJOkE4CQIJwlZSAQEbJwFIN/6KcBAQsntuphqcBBK7JZfNWZUrJBWAkRsFZZNZGzRslZTyqQakFYapBqmGqYYlslYYpBmyta1SDFOwqDFIMVwYrAxLZbDcvZS5aJDFIM2S2NhOWmMaM5aXKS2NgTHsm5aP5WyYwpbMBy0xjR5hUTFsls9gTGomNHcpMYtkbPYEx7KJj2Rxi2UTEjZ7A8tNy0cY1ExI2NgjGomNGmJRMaNnsEY03LRpiUeWjY2hkSLFflSy7LQth8h7KJYisqbIg9hciiWosxqDo9kxsNlTFqILFEtQagtTWVxamLUwqITWVhamLU4FdkxCsskQmFRCVtlMtSsmEEgpWSsmDBOlZPbsmCThKydCSST2SskVIBOkAnsgEAnATgKQCCRAUgFIBSDUiQAUgFMNUg1AQDU4arMuycNQFeVPlVgbspBuyAqypFuiuypZNkgoyp8hV4YpCNIBwwqQYiBEeymIktgMGKxrNl1nBPAmPcWiaXDIqeKlhOV9TUvLI8/wCBtgS53S9hYX1IQ3FXC2L8L4i2ixenax0jC+GWRU1YOqsZ1SSU1NTCSSSRUiolJJBIO6qtySSCQKiUkkGiU3ikkgHTDqkkpBDqpBJJSDhOEkkESSSSAiVA9UklUCDuqrckkrgVvVLkklpFxS5VOSSWkWqd1VZSSWkVDFR8UkkzMUkkkKIJXKSSk4fqmPRJJIRByrKSSDVuKrckkkuKnqpySSqLil/ZdHwv70XmkkuLzv0w8/01v+mDXAsLHhn/AEXmB8Eklt4n+Gnx/wBEIJvBJJdLYkkkkGSQSSQCSSSQH//Z",
			TestErrorMessage: 	"Failed while testing base64 image input value",
			ExpectedReturnCode: http.StatusBadRequest,
			CallsService: 		false,
		},
		{
			Name: 				"Valid Number",
			Input:				"38977123456",
			TestErrorMessage: 	"Failed while testing valid number",
			ExpectedReturnCode:	http.StatusOK,
			CallsService: 		true,
		},
		{
			Name:				"Valid number with whitespace",
			Input: 				"389 77 123 456",
			TestErrorMessage: 	"Failed while testing valid number with whitespace",
			ExpectedReturnCode: http.StatusOK,
			CallsService: 		true,
		},
	}

	for _, test := range tt{
		fn := func(t *testing.T){

			//Arrange
			recorder := httptest.NewRecorder()
			teardown := setup(t,recorder)
			defer teardown()

			if test.CallsService{
				mockService.EXPECT().LookupMSISDN(gomock.Any()).Return(nil, &errs.AppError{ Code: test.ExpectedReturnCode})
			}
			jsonReq := LookupRequest{
				Number: test.Input,
			}
			jsonVal,_ := json.Marshal(jsonReq)	

			//Act
			req, err := http.NewRequest(http.MethodPost,"/lookup",bytes.NewBuffer(jsonVal))
			req.Header.Set("Content-Type","application/json")
			if err!= nil{
				t.Fatalf("Could not make request in " + test.Name)
			}
			router.ServeHTTP(recorder,req)

			//Assert
			if recorder.Code != test.ExpectedReturnCode{
				t.Error(test.TestErrorMessage)
			}
		}
		t.Run(test.Name, fn)
	}
}