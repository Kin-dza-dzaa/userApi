package repository

import (
	"context"
	"testing"
	"time"

	config "github.com/Kin-dza-dzaa/userApi/configs"
	"github.com/Kin-dza-dzaa/userApi/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
)

type TestStruct struct {
    user models.User
	method string
	errorExpected bool
	err string
	ifExistsResult bool
}

var testSlice []TestStruct = []TestStruct{
	{
		user: models.User{UserId: uuid.New(), UserName: "TestUser", Password: "12345", Email: "testuser@gmail.com", RegistrationTime: time.Now().UTC(), VerificationCode: "verify", Verified: false},
		method: "AddUser",
		errorExpected: false,
	},
	{
		user: models.User{UserId: uuid.New(), UserName: "TestUser", Password: "12345", Email: "testuser@gmail.com", RegistrationTime: time.Now().UTC(), VerificationCode: "newVerify", Verified: false},
		method: "UpdateCredentials",
		errorExpected: false,
	},
	{
		user: models.User{Email: "testuser@gmail.com"},
		method: "IfUnverifiedUserExists",
		ifExistsResult: true,
		errorExpected: false,
	},
	{
		user: models.User{Email: "user@gmail.com"},
		method: "IfUnverifiedUserExists",
		ifExistsResult: false,
		errorExpected: false,
	},
	{
		user: models.User{VerificationCode: "newVerify", ExpirationTime: time.Now().UTC().Add(time.Minute * 10), RefreshToken: "refresh"},
		method: "VerifyUser",
		errorExpected: false,
	},
	{
		user: models.User{VerificationCode: "verify", ExpirationTime: time.Now().UTC().Add(time.Minute * 10), RefreshToken: "refresh"},
		method: "VerifyUser",
		errorExpected: true,
		err: "wrong verification code",
	},
	{
		user: models.User{UserId: uuid.New(), UserName: "TestUser", Email: "testuser@gmail.com", Password: "12345", RegistrationTime: time.Now().UTC(), VerificationCode: "verify1", Verified: false},
		method: "AddUser",
		errorExpected: true,
		err: "user already exists",
	},
	{
		user: models.User{RefreshToken: "refresh"},
		method: "GetUUid",
		errorExpected: false,
	},
	{
		user: models.User{RefreshToken: "wrong code"},
		method: "GetUUid",
		errorExpected: true,
		err: "user doesn't exist",
	},
	{
		user: models.User{Email: "testuser@gmail.com"},
		method: "GetVerifiedUser",
		errorExpected: false,
	},	
	{
		user: models.User{Email: "wrongemail@gmail.com"},
		method: "GetVerifiedUser",
		errorExpected: true,
		err: "wrong email",
	},
	{
		user: models.User{ExpirationTime: time.Now().UTC().Add(time.Minute * 10), Email: "testuser@gmail.com", RefreshToken: "refresh"},
		method: "UpdateRefreshToken",
		errorExpected: false,
	},
	{
		user: models.User{VerificationCode: "verify", ExpirationTime: time.Now().UTC().Add(time.Minute * 10), Email: "user@gmail.com", RefreshToken: "refresh"},
		method: "UpdateRefreshToken",
		errorExpected: false,
	},
}

type RepositorySuite struct {
	pool *pgxpool.Pool
	suite.Suite
	repository Repository
}

func (suite *RepositorySuite) SetupSuite() {
	logger := zerolog.New(zerolog.NewConsoleWriter()).With().Timestamp().Caller().Logger()
	config, err := config.ReadConfig(&logger)
	if err != nil {
		suite.FailNow(err.Error())
	}
	suite.pool, err = pgxpool.Connect(context.TODO(), config.DbUrl)
	if err != nil {
		suite.FailNow(err.Error())
	}
	suite.repository = NewRepository(suite.pool)
}

func (suite *RepositorySuite) TearDownSuite() {
	for _, v := range testSlice {
		if v.method == "AddUser" {
			if _, err := suite.pool.Exec(context.TODO(), "DELETE FROM USERS WHERE email = $1;", v.user.Email); err != nil {
				suite.FailNow(err.Error())
			}
		}
	}
}

func (suite *RepositorySuite) Test() {
	for _, v := range testSlice {
		if v.method == "AddUser" {
			err := suite.repository.AddUser(context.TODO(), &v.user)
			if v.errorExpected {
				suite.EqualError(err, v.err)
			} else {
				suite.Nil(err)
			}
		}
		if v.method == "UpdateCredentials" {
			err := suite.repository.UpdateCredentials(context.TODO(), &v.user)
			if v.errorExpected {
				suite.EqualError(err, v.err)
			} else {
				suite.Nil(err)
			}
		}
		if v.method == "IfUnverifiedUserExists" {
			res, err := suite.repository.IfUnverifiedUserExists(context.TODO(), &v.user)
			if v.errorExpected {
				suite.EqualError(err, v.err)
			} else {
				suite.Nil(err)
				suite.Equal(v.ifExistsResult, res)
			}
		}
		if v.method == "VerifyUser" {
			err := suite.repository.VerifyUser(context.TODO(), &v.user)
			if v.errorExpected {
				suite.EqualError(err, v.err)
			} else {
				suite.Nil(err)
			}
		}
		if v.method == "GetUUid" {
			err := suite.repository.GetUUid(context.TODO(), &v.user)
			if v.errorExpected {
				suite.EqualError(err, v.err)
			} else {
				suite.Nil(err)
			}
		}
		if v.method == "GetVerifiedUser" {
			_, err := suite.repository.GetVerifiedUser(context.TODO(), &v.user)
			if v.errorExpected {
				suite.EqualError(err, v.err)
			} else {
				suite.Nil(err)
			}
		}
		if v.method == "UpdateRefreshToken" {
			err := suite.repository.UpdateRefreshToken(context.TODO(), &v.user)
			if v.errorExpected {
				suite.EqualError(err, v.err)
			} else {
				suite.Nil(err)
			}
		}
	}
}

func TestUserRepo(t *testing.T) {
	suite.Run(t, new(RepositorySuite))
}



