package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"
	config "github.com/Kin-dza-dzaa/userApi/configs"
	"github.com/Kin-dza-dzaa/userApi/internal/apierror"
	"github.com/Kin-dza-dzaa/userApi/internal/models"
	"github.com/Kin-dza-dzaa/userApi/pkg/handlers"
	"github.com/Kin-dza-dzaa/userApi/pkg/repositories"
	"github.com/Kin-dza-dzaa/userApi/pkg/service"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
)

type TestData struct {
	user models.User
	expected response
	status int
}

type response struct {
	Message string		`json:"message,omitempty"`
	Result string		`json:"result"`
}

type IntegrationTestSuite struct {
	suite.Suite
	server *http.Server
}

func (suite *IntegrationTestSuite) SetupSuite() {
	logger := zerolog.New(nil).With().Timestamp().Caller().Logger()
	config, err := config.ReadConfig(&logger)
	if err != nil {
		suite.FailNow(err.Error())
	}
	pool ,err := pgxpool.Connect(context.TODO(), config.DbUrl)
	if err != nil {
		suite.FailNow(err.Error())
	}
	ApiError := apierror.NewApiError(&logger)
	MyRepository := repository.NewRepository(pool)
	MyService := service.NewService(MyRepository, config)
	MyHandlers := handlers.NewHandlers(MyService, config, ApiError)
	suite.server = &http.Server{
		Addr: config.Adress,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		Handler: MyHandlers.Cors.Handler(MyHandlers.Router),
	}
	go func() {
		suite.server.ListenAndServe()
	}()
}

func (suite *IntegrationTestSuite) TearDownSuite() {
	if err := suite.server.Shutdown(context.TODO()); err != nil {
		suite.Fail(err.Error())
	}
}

func (suite *IntegrationTestSuite) TestSignUp() {
	var testSlice []TestData = []TestData{
		{
			user: models.User{
				UserName: "testUser",
				Email: "testuser@gmail.com",
				Password: "12345Qwerty",
			},
			expected: response{
				Result: "ok",
				Message: "email was sent",
			},
			status: http.StatusOK,
		},		
		{
			user: models.User{
				UserName: "testUser",
				Email: "testuser@gmail.com",
				Password: "12345Qwerty",
			},
			expected: response{
				Result: "ok",
				Message: "email was sent",
			},
			status: http.StatusOK,
		},
		{
			user: models.User{
				Password: "12345Qwerty",
			},
			expected: response{
				Result: "error",
				Message: "invalid credentials",
			},
			status: http.StatusBadRequest,
		},
		{
			user: models.User{
			},
			expected: response{
				Result: "error",
				Message: "invalid credentials",
			},
			status: http.StatusBadRequest,
		},
		{
			user: models.User{		
				UserName: "testUser",
				Password: "12345Qwerty",
			},
			expected: response{
				Result: "error",
				Message: "invalid credentials",
			},
			status: http.StatusBadRequest,
		},
		{
			user: models.User{		
				Email: "testuser@gmail.com",
				Password: "12345Qwerty",
			},
			expected: response{
				Result: "error",
				Message: "invalid credentials",
			},
			status: http.StatusBadRequest,
		},
		{
			user: models.User{		
				UserName: "",
				Email: "",
				Password: "12345Qwerty",
			},
			expected: response{
				Result: "error",
				Message: "invalid credentials",
			},
			status: http.StatusBadRequest,
		},
		{
			user: models.User{		
				UserName: "",
				Email: "testuser@gmail.com",
				Password: "12345Qwerty",
			},
			expected: response{
				Result: "error",
				Message: "invalid credentials",
			},
			status: http.StatusBadRequest,
		},
		{
			user: models.User{		
				UserName: "testUser",
				Email: "",
				Password: "12345Qwerty",
			},
			expected: response{
				Result: "error",
				Message: "invalid credentials",
			},
			status: http.StatusBadRequest,
		},
		{
			user: models.User{		
				UserName: "testUser",
				Email: "testuser@gmail.com",
				Password: "",
			},
			expected: response{
				Result: "error",
				Message: "invalid credentials",
			},
			status: http.StatusBadRequest,
		},
		{
			user: models.User{		
				UserName: "",
				Email: "",
				Password: "",
			},
			expected: response{
				Result: "error",
				Message: "invalid credentials",
			},
			status: http.StatusBadRequest,
		},
	}
	for _, v := range testSlice {
		byteData, err := json.Marshal(v.user)
		if err != nil {
			suite.FailNow(err.Error())
		}
		res, err := http.Post("http://localhost:8001/user", "application/json", bytes.NewReader(byteData))
		if err != nil {
			suite.FailNow(err.Error())
		}
		var responseStruct response
		if err := json.NewDecoder(res.Body).Decode(&responseStruct); err != nil {
			suite.FailNow(err.Error())
		}
		defer res.Body.Close()
		suite.Equal(v.expected, responseStruct)
		suite.Equal(v.status, res.StatusCode)
	}
}

func TestHandlers(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}