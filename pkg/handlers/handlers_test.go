package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	config "github.com/Kin-dza-dzaa/userApi/configs"
	"github.com/Kin-dza-dzaa/userApi/internal/models"
	"github.com/Kin-dza-dzaa/userApi/pkg/mocks"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type TestData struct {
	user models.User
	result string
	message string
	status int
}

type response struct {
	Message string		`json:"message,omitempty"`
	Result string		`json:"result"`
}

type HandlersSuite struct {
	suite.Suite
	service *mocks.Service
	handlers *Handlers
}

func (suite *HandlersSuite) SetupSuite() {
	logger := zerolog.New(zerolog.NewConsoleWriter()).With().Timestamp().Caller().Logger()
	config, err := config.ReadConfig(&logger)
	if err != nil {
		suite.FailNow(err.Error())
	}
	suite.service = mocks.NewService(suite.T())
	suite.handlers = NewHandlers(suite.service, *config, &logger)
}

func (suite *HandlersSuite) TestSignUpHandler() {
	var testSlice []TestData = []TestData{
		{
			user: models.User{
				UserName: "testUser",
				Email: "testEmail@gmail.com",
				Password: "12345Qwerty",
			},
			result: "ok",
			message: "email was sent",
			status: http.StatusOK,
		},
	}
	for _, v := range testSlice {
		byteData, err := json.Marshal(v.user)
		if err != nil {
			suite.FailNow(err.Error())
		}
		r := httptest.NewRequest("POST", "/user",bytes.NewReader(byteData))	
		w := httptest.NewRecorder()
		suite.service.On("SignUpUser", mock.Anything).Return(nil).Once()
		suite.handlers.SignUpHandler().ServeHTTP(w, r)
		var response response
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			suite.FailNow(err.Error())
		}
		suite.Equal(v.message, response.Message)
		suite.Equal(v.result, response.Result)
		suite.Equal(v.status, w.Result().StatusCode)
	}
}

func (suite *HandlersSuite) TestSignInHandler() {
	var testSlice []TestData = []TestData{
		{
			user: models.User{
				UserName: "testUser",
				Email: "testEmail@gmail.com",
				Password: "12345Qwerty",
			},
			result: "ok",
			status: http.StatusOK,
		},
	}
	for _, v := range testSlice {
		byteData, err := json.Marshal(v.user)
		if err != nil {
			suite.FailNow(err.Error())
		}
		r := httptest.NewRequest("POST", "/user",bytes.NewReader(byteData))	
		w := httptest.NewRecorder()
		suite.service.On("SignInUser", mock.Anything).Return(nil).Once()
		suite.handlers.SignInHandler().ServeHTTP(w, r)
		var response response
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			suite.FailNow(err.Error())
		}
		suite.Equal(v.message, response.Message)
		suite.Equal(v.result, response.Result)
		suite.Equal(v.status, w.Result().StatusCode)
	}
}

func TestHandlers(t *testing.T) {
	suite.Run(t, new(HandlersSuite))
}