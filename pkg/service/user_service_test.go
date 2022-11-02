package service

import (
	"context"
	"errors"
	"testing"
	config "github.com/Kin-dza-dzaa/userApi/configs"
	"github.com/Kin-dza-dzaa/userApi/internal/models"
	"github.com/Kin-dza-dzaa/userApi/pkg/mocks"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type UserServiceSuite struct {
	service Service
	suite.Suite
	repository *mocks.Repository
}

func (suite *UserServiceSuite) SetupSuite() {
	logger := zerolog.New(nil).With().Timestamp().Caller().Logger()
	config, err := config.ReadConfig(&logger)
	if err != nil {
		suite.FailNow(err.Error())
	}
	suite.repository = mocks.NewRepository(suite.T())
	suite.service = NewService(suite.repository, config)
}

func (suite *UserServiceSuite) TestSignUp() {
	testSlice := []struct{
			expectedError bool
			err string
		}{
		{
			expectedError: true,
			err:           "",
		},
		{
			expectedError: true,
			err:           "",
		},
		{
			expectedError: true,
			err:           "",
		},
		{
			expectedError: true,
			err:           "",
		},
		
	}
	for _, tc := range testSlice {
		suite.T().Run("SignUp", func(t *testing.T) {
			suite.repository.On("IfUnverifiedUserExists", mock.Anything, mock.Anything, mock.Anything).Return(errors.New(""))
			err := suite.service.SignUpUser(context.TODO(), &models.User{})
			suite.Equal(tc.err, err.Error())
		})
	}
}

func TestUserService(t *testing.T) {
	suite.Run(t, new(UserServiceSuite))
}
