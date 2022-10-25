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
		name string
		expectedResponse *ErrorStruct
		raiseError error
	}{
		{
			name: "closed_pool",
			expectedResponse: NewErrorStruct("internal server error", "error"),
			raiseError: puddle.ErrClosedPool,
		},
		{
			name: "busy_pool",
			expectedResponse: NewErrorStruct("too many requests", "error"),
			raiseError: puddle.ErrNotAvailable,
		},
		{
			name: "ivalid_credentials",
			expectedResponse: NewErrorStruct(dto.ErrInvalidCredentials.Error(), "error"),
			raiseError: dto.ErrInvalidCredentials,
		},
		{
			name: "wrong_password",
			expectedResponse: NewErrorStruct(service.ErrWrongPassowrd.Error(), "error"),
			raiseError: service.ErrWrongPassowrd,
		},
		{
			name: "user_doesn't_exists",
			expectedResponse: NewErrorStruct(repository.ErrUserDoesntExists.Error(), "error"),
			raiseError: repository.ErrUserDoesntExists,
		},
		{
			name: "wrong_email",
			expectedResponse: NewErrorStruct(repository.ErrWrongEmail.Error(), "error"),
			raiseError: repository.ErrWrongEmail,
		},
		{
			name: "wrong_verif_code",
			expectedResponse: NewErrorStruct(repository.ErrWrongVerificationCode.Error(), "error"),
			raiseError: repository.ErrWrongVerificationCode,
		},
		{
			name: "not_validated_error",
			expectedResponse: NewErrorStruct("unexpected error", "error"),
			raiseError: errors.New("not validated error"),
		},
	}	
	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest("POST", "/test", nil)
			w := httptest.NewRecorder()
			suite.apiError.ErrorMiddleWare(func(w http.ResponseWriter, r *http.Request) error {return tc.raiseError}).ServeHTTP(w, r)
			var response ErrorStruct
			err := json.NewDecoder(w.Body).Decode(&response)
			if err != nil {
				suite.FailNow(err.Error())
			}
			suite.Equal(*tc.expectedResponse, response)
		})
	}
}

func TestMiddleware(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
