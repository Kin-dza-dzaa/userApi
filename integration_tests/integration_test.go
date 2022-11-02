// +build integration
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
	"github.com/Kin-dza-dzaa/userApi/internal/dto"
	"github.com/Kin-dza-dzaa/userApi/pkg/handlers"
	"github.com/Kin-dza-dzaa/userApi/pkg/repositories"
	"github.com/Kin-dza-dzaa/userApi/pkg/service"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
)

const (
	queryDeleteUsers = "DELETE FROM USERS;"
)

type IntegrationTestSuite struct {
	suite.Suite
	server *http.Server
	pool *pgxpool.Pool
}

func (suite *IntegrationTestSuite) SetupSuite() {
	logger := zerolog.New(nil).With().Timestamp().Caller().Logger()
	config, err := config.ReadConfig(&logger)
	if err != nil {
		suite.FailNow(err.Error())
	}
	config.TemplateLocation = "./../internal/templates/response_template.html"
	pool ,err := pgxpool.Connect(context.TODO(), config.LocalDbUrlTest)
	if err != nil {
		suite.FailNow(err.Error())
	}
	suite.pool = pool
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
		suite.FailNow(err.Error())
	}
}

func (suite *IntegrationTestSuite) TearDownTest() {
	if _, err := suite.pool.Exec(context.TODO(), queryDeleteUsers); err != nil {
		suite.FailNow(err.Error())
	}
}

func (suite *IntegrationTestSuite) TestSignUp() {
	testCases := []struct{
		name 		string
		input 		*dto.UserSignUpDto
		result 		apierror.ErrorStruct
	}{
		{
			name: "good_user_data",
			input: &dto.UserSignUpDto{
				Email: "testEmail@gmail.com",
				UserName: "testUserValid",
				Password: "12345ValidPassword",
			},
			result: apierror.ErrorStruct{
				Result: "ok",
				Message: "email was sent",
				Code: 200,
			},
		},
		{
			name: "nil_input",
			input: nil,
			result: apierror.ErrorStruct{
				Result: "error",
				Message: dto.ErrInvalidCredentials.Error(),
				Code: 400,
			},
		},
		{
			name: "empty_user",
			input: &dto.UserSignUpDto{},
			result: apierror.ErrorStruct{
				Result: "error",
				Message: dto.ErrInvalidCredentials.Error(),
				Code: 400,
			},
		},
	}
	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			byteData, err := json.Marshal(tc.input)
			if err != nil {
				suite.FailNow(err.Error())
			}
			res, err := http.Post("http://localhost:8001/user", "application/json", bytes.NewReader(byteData))
			if err != nil {
				suite.FailNow(err.Error())
			}
			var response apierror.ErrorStruct
			if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
				suite.FailNow(err.Error())
			}
			suite.Equal(tc.result, response)
		})
	}
}

func TestHandlers(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
