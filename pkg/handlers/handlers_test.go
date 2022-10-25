package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	config "github.com/Kin-dza-dzaa/userApi/configs"
	"github.com/Kin-dza-dzaa/userApi/internal/apierror"
	"github.com/Kin-dza-dzaa/userApi/internal/dto"
	"github.com/Kin-dza-dzaa/userApi/pkg/mocks"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type HandlersSuite struct {
	suite.Suite
	service  *mocks.Service
	handlers *Handlers
}

func (suite *HandlersSuite) SetupSuite() {
	logger := zerolog.New(nil).With().Timestamp().Caller().Logger()
	config, err := config.ReadConfig(&logger)
	if err != nil {
		suite.FailNow(err.Error())
	}
	suite.service = mocks.NewService(suite.T())
	ApiError := apierror.NewApiError(&logger)
	suite.handlers = NewHandlers(suite.service, config, ApiError)
}

func (suite *HandlersSuite) TestSignUpHandler() {
	testCases := []struct {
		user        dto.UserSignUpDto
		response    apierror.ErrorStruct
		expectedErr bool
	}{
		{
			user: dto.UserSignUpDto{
				UserName: "testuser",
				Email:    "testEmail@gmail.com",
				Password: "12345Qwerty",
			},
			response: apierror.ErrorStruct{
				Message: "email was sent",
				Result:  "ok",
			},
			expectedErr: false,
		},
		{
			user: dto.UserSignUpDto{
				Email:    "testEmail@gmail.com",
				Password: "12345Qwerty",
			},
			response: apierror.ErrorStruct{
				Message: dto.ErrInvalidCredentials.Error(),
				Result:  "error",
			},
			expectedErr: true,
		},
		{
			user: dto.UserSignUpDto{
				Email:    "testEmail@gmail.com",
				UserName: "ValidUserName",
				Password: "",
			},
			response: apierror.ErrorStruct{
				Message: dto.ErrInvalidCredentials.Error(),
				Result:  "error",
			},
			expectedErr: true,
		},
		{
			user: dto.UserSignUpDto{},
			response: apierror.ErrorStruct{
				Message: dto.ErrInvalidCredentials.Error(),
				Result:  "error",
			},
			expectedErr: true,
		},
	}
	for _, tc := range testCases {
		suite.T().Run("SignUpHandler", func(t *testing.T) {
			byteData, err := json.Marshal(tc.user)
			if err != nil {
				suite.FailNow(err.Error())
			}
			r := httptest.NewRequest("POST", "/user", bytes.NewReader(byteData))
			w := httptest.NewRecorder()
			if !tc.expectedErr {
				suite.service.On("SignUpUser", mock.Anything, mock.Anything).Return(nil).Once()
			}
			suite.handlers.ApiError.ErrorMiddleWare(suite.handlers.SignUpHandler()).ServeHTTP(w, r)
			var response apierror.ErrorStruct
			if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
				suite.FailNow(err.Error())
			}
			suite.Equal(tc.response, response)
		})
	}
}

func (suite *HandlersSuite) TestSignInHandler() {
	testSlice := []struct {
		user        dto.UserSignUpDto
		response    apierror.ErrorStruct
		expectedErr bool
	}{
		{
			user: dto.UserSignUpDto{
				Email:    "testEmail@gmail.com",
				Password: "12345Qwerty",
			},
			response: apierror.ErrorStruct{
				Result: "ok",
			},
			expectedErr: false,
		},
		{
			user: dto.UserSignUpDto{
				Email:    "testEmail@gmail.com",
				Password: "",
			},
			response: apierror.ErrorStruct{
				Result:  "error",
				Message: dto.ErrInvalidCredentials.Error(),
			},
			expectedErr: true,
		},
		{
			user: dto.UserSignUpDto{
				Email:    "",
				Password: "123456789",
			},
			response: apierror.ErrorStruct{
				Result:  "error",
				Message: dto.ErrInvalidCredentials.Error(),
			},
			expectedErr: true,
		},
		{
			user: dto.UserSignUpDto{},
			response: apierror.ErrorStruct{
				Result:  "error",
				Message: dto.ErrInvalidCredentials.Error(),
			},
			expectedErr: true,
		},
	}
	for _, tc := range testSlice {
		suite.T().Run("SignInHandler", func(t *testing.T) {
			byteData, err := json.Marshal(tc.user)
			if err != nil {
				suite.FailNow(err.Error())
			}
			r := httptest.NewRequest("POST", "/user", bytes.NewReader(byteData))
			w := httptest.NewRecorder()
			if !tc.expectedErr {
				suite.service.On("SignInUser", mock.Anything, mock.Anything).Return(nil).Once()
			}
			suite.handlers.ApiError.ErrorMiddleWare(suite.handlers.SignInHandler()).ServeHTTP(w, r)
			var response apierror.ErrorStruct
			if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
				suite.FailNow(err.Error())
			}
			suite.Equal(tc.response, response)
		})
	}
}

func (suite *HandlersSuite) TestGetTokenHandler() {
	testCases := []struct {
		result      apierror.ErrorStruct
		expectedErr bool
	}{
		{
			result: apierror.ErrorStruct{
				Result: "ok",
			},
			expectedErr: false,
		},
		{
			result: apierror.ErrorStruct{
				Result:  "error",
				Message: "unexpected error",
			},
			expectedErr: true,
		},
	}
	for _, tc := range testCases {
		suite.T().Run("GetTokenHandler", func(t *testing.T) {
			r := httptest.NewRequest("POST", "/user/token", nil)
			w := httptest.NewRecorder()
			if !tc.expectedErr {
				r.AddCookie(&http.Cookie{
					Name: "Refresh-token",
				})
				suite.service.Mock.On("GetAccessToken", mock.Anything, mock.Anything).Return(nil).Once()
			}
			suite.handlers.ApiError.ErrorMiddleWare(suite.handlers.GetTokenHandler()).ServeHTTP(w, r)
			var response apierror.ErrorStruct
			if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
				suite.T().FailNow()
			}
			suite.Equal(tc.result, response)
		})
	}

}

func (suite *HandlersSuite) TestLogOutHandler() {
	testCases := []struct {
		result        apierror.ErrorStruct
		expectedError bool
	}{
		{
			result: apierror.ErrorStruct{
				Result: "ok",
			},
			expectedError: false,
		},
	}
	for _, tc := range testCases {
		suite.T().Run("LogOutHandler", func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("PoST", "/user/logout", nil)
			suite.handlers.ApiError.ErrorMiddleWare(suite.handlers.LogOutHandler()).ServeHTTP(w, r)
			var result apierror.ErrorStruct
			if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
				t.FailNow()
			}
			suite.Equal(tc.result, result)
		})
	}
}

func TestHandlers(t *testing.T) {
	suite.Run(t, new(HandlersSuite))
}

