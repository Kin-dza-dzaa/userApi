package apierror

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/Kin-dza-dzaa/userApi/internal/dto"
	repository "github.com/Kin-dza-dzaa/userApi/pkg/repositories"
	"github.com/Kin-dza-dzaa/userApi/pkg/service"
	"github.com/jackc/puddle"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	apiError *ApiError
}

func (suite *TestSuite) SetupSuite() {
	logger := zerolog.New(nil).With().Timestamp().Caller().Logger()
	suite.apiError = NewApiError(&logger)
}

func (suite *TestSuite) TestMiddleware() {
	testCases := []struct {
		expectedResponse *ErrorStruct
		raiseError error
	}{
		{
			expectedResponse: NewErrorStruct("internal server error", "error"),
			raiseError: puddle.ErrClosedPool,
		},
		{
			expectedResponse: NewErrorStruct("too many requests", "error"),
			raiseError: puddle.ErrNotAvailable,
		},
		{
			expectedResponse: NewErrorStruct(dto.ErrInvalidCredentials.Error(), "error"),
			raiseError: dto.ErrInvalidCredentials,
		},
		{
			expectedResponse: NewErrorStruct(service.ErrWrongPassowrd.Error(), "error"),
			raiseError: service.ErrWrongPassowrd,
		},
		{
			expectedResponse: NewErrorStruct(repository.ErrUserDoesntExists.Error(), "error"),
			raiseError: repository.ErrUserDoesntExists,
		},
		{
			expectedResponse: NewErrorStruct(repository.ErrWrongEmail.Error(), "error"),
			raiseError: repository.ErrWrongEmail,
		},
		{
			expectedResponse: NewErrorStruct(repository.ErrWrongVerificationCode.Error(), "error"),
			raiseError: repository.ErrWrongVerificationCode,
		},
		{
			expectedResponse: NewErrorStruct("unexpected error", "error"),
			raiseError: errors.New("not validated error"),
		},
	}	
	for _, tc := range testCases {
		r := httptest.NewRequest("POST", "/test", nil)
		w := httptest.NewRecorder()
		suite.apiError.ErrorMiddleWare(func(w http.ResponseWriter, r *http.Request) error {return tc.raiseError}).ServeHTTP(w, r)
		var response ErrorStruct
		err := json.NewDecoder(w.Body).Decode(&response)
		if err != nil {
			suite.FailNow(err.Error())
		}
		suite.Equal(*tc.expectedResponse, response)
	}
}

func TestMiddleware(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func DummyHandler(w http.ResponseWriter, r *http.Request) error {
	return nil
}