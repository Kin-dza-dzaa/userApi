package service

import (
	"context"
	"errors"
	"testing"
	"time"
	config "github.com/Kin-dza-dzaa/userApi/configs"
	"github.com/Kin-dza-dzaa/userApi/internal/models"
	"github.com/Kin-dza-dzaa/userApi/pkg/mocks"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
)

type TestStruct struct {
	user          models.User
	errorExpected bool
	err           string
}

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
	var testSlice []TestStruct = []TestStruct{
		{
			user:          models.User{UserName: "TestUser", Password: "12345", Email: "testuser@gmail.com"},
			errorExpected: true,
			err:           "invalid credentials",
		},
		{
			user:          models.User{UserName: "TestUser@!", Password: "12345Qwerty", Email: "testuser@gmail.com"},
			errorExpected: true,
			err:           "invalid credentials",
		},
		{
			user:          models.User{UserName: "", Password: "", Email: "testuser@gmail.com"},
			errorExpected: true,
			err:           "invalid credentials",
		},
		{
			user:          models.User{Password: "12345Qwerty", Email: "testuser@gmail.com"},
			errorExpected: true,
			err:           "invalid credentials",
		},
		{
			user:          models.User{UserName: "TestUser", Password: "12345Qwerty", Email: "testusergmail.com"},
			errorExpected: true,
			err:           "invalid credentials",
		},
		{
			user:          models.User{UserName: "", Password: "", Email: ""},
			errorExpected: true,
			err:           "invalid credentials",
		},
		{
			errorExpected: true,
			err:           "invalid credentials",
		},
		{
			user:          models.User{UserName: "TestUser", Password: "12345Qwerty", Email: "testuser@gmail.com"},
			errorExpected: true,
			err:           "service logic ok",
		},
	}
	for _, v := range testSlice {
		suite.repository.On("IfUnverifiedUserExists", mock.Anything, &v.user).Return(false, errors.New("service logic ok"))
		err := suite.service.SignUpUser(context.TODO(), &v.user)
		if v.errorExpected {
			suite.EqualError(err, v.err)
		} else {
			suite.Nil(err)
		}
	}
}

func (suite *UserServiceSuite) TestSignIn() {
	var testSlice []TestStruct = []TestStruct{
		{
			user:          models.User{UserName: "TestUser", Password: "12345Qwerty", Email: "testuser@gmail.com"},
			errorExpected: true,
			err:           "wrong password",
		},
		{
			user:          models.User{UserName: "TestUser", Password: "12345Qwerty", Email: "testuser@gmail.com"},
			errorExpected: true,
			err:           "wrong password",
		},
		{
			user:          models.User{UserName: "TestUser", Password: "12345Qwerty", Email: "testuser@gmail.com"},
			errorExpected: false,
		},
	}
	for _, v := range testSlice {
		hash, err := bcrypt.GenerateFromPassword([]byte(v.user.Password), 14)
		if err != nil {
			suite.FailNow(err.Error())
		}
		if v.errorExpected {
			suite.repository.On("GetVerifiedUser", mock.Anything, &v.user).Return(&models.User{Password: "", ExpirationTime: time.Now().Add(time.Minute * 5).UTC()}, nil).Once()
			err = suite.service.SignInUser(context.TODO(), &v.user)
			suite.EqualError(err, v.err)
		} else {
			suite.repository.On("GetVerifiedUser", mock.Anything, &v.user).Return(&models.User{Password: string(hash), ExpirationTime: time.Now().Add(time.Minute * 5).UTC()}, nil).Once()
			err = suite.service.SignInUser(context.TODO(), &v.user)
			suite.Nil(err)
		}
	}
}

func TestUserService(t *testing.T) {
	suite.Run(t, new(UserServiceSuite))
}
