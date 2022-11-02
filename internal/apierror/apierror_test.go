package apierror

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

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
		name             string
		expectedResponse *ErrorStruct
	}{
		{
			name:             "closed_pool",
			expectedResponse: NewErrorStruct("internal server error", "error", http.StatusInternalServerError),
		},
		{
			name:             "busy_pool",
			expectedResponse: NewErrorStruct("too many requests", "error", http.StatusTooManyRequests),
		},
		{
			name:             "ivalid_credentials",
			expectedResponse: NewErrorStruct("invalid credentials", "error", http.StatusBadRequest),
		},
		{
			name:             "wrong_password",
			expectedResponse: NewErrorStruct("wrong password", "error", http.StatusBadRequest),
		},
		{
			name:             "user_doesn't_exists",
			expectedResponse: NewErrorStruct("user doesn't exist", "error", http.StatusBadRequest),
		},
		{
			name:             "wrong_email",
			expectedResponse: NewErrorStruct("wrong email", "error", http.StatusBadRequest),
		},
		{
			name:             "wrong_verif_code",
			expectedResponse: NewErrorStruct("wrong verification code", "error", http.StatusBadRequest),
		},
		{
			name:             "not_validated_error",
			expectedResponse: NewErrorStruct("unexpected error", "error", http.StatusInternalServerError),
		},
	}
	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest("POST", "/test", nil)
			w := httptest.NewRecorder()
			suite.apiError.ErrorMiddleWare(func(w http.ResponseWriter, r *http.Request) error { return tc.expectedResponse }).ServeHTTP(w, r)
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
