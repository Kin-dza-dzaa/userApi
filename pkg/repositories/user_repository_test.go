package repository

import (
	"context"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/Kin-dza-dzaa/userApi/internal/apierror"
	"github.com/Kin-dza-dzaa/userApi/internal/models"
	"github.com/chrisyxlee/pgxpoolmock"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
)

func TestAddUser(t *testing.T) {
	// just to make sure that no one will mess with an order of the func paramaters
	testCases := []struct{
		name string
		args struct{
			ctx context.Context
			user *models.User
		}
		beforeTest func(mockPool *pgxpoolmock.MockPgxIface, user *models.User)
	}{
		{
			name: "AddUser",
			args: struct{ctx context.Context; user *models.User}{
				ctx: context.TODO(),
				user: new(models.User),
			},
			beforeTest: func(mockPool *pgxpoolmock.MockPgxIface, user *models.User) {
				mockPool.EXPECT().Exec(gomock.Any(), queryCreateUser, user.UserId, user.UserName, user.Email, user.Password, user.RegistrationTime, user.VerificationCode, user.Verified).Return(nil, nil).Times(1)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			MockPool := pgxpoolmock.NewMockPgxIface(ctrl)
			repository := NewUserRepository(MockPool)
			tc.beforeTest(MockPool, tc.args.user)
			if err := repository.AddUser(tc.args.ctx, tc.args.user); err != nil {
				t.FailNow()
			}
		})
	}
}

func TestUpdateCredentials(t *testing.T) {
	// just to make sure that no one will mess with an order of the func paramaters
	testCases := []struct{
		name string
		args struct{
			ctx context.Context
			user *models.User
		}
		beforeTest func(mockPool *pgxpoolmock.MockPgxIface, user *models.User)
	}{
		{
			name: "UpdateCredentials",
			args: struct{ctx context.Context; user *models.User}{
				ctx: context.TODO(),
				user: new(models.User),
			},
			beforeTest: func(mockPool *pgxpoolmock.MockPgxIface, user *models.User) {
				mockPool.EXPECT().Exec(gomock.Any(), queryUpdateCreditnails, user.UserName, user.Password, user.VerificationCode, user.Email).Return(nil, nil).Times(1)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			MockPool := pgxpoolmock.NewMockPgxIface(ctrl)
			repository := NewUserRepository(MockPool)
			tc.beforeTest(MockPool, tc.args.user)
			if err := repository.UpdateCredentials(tc.args.ctx, tc.args.user); err != nil {
				t.FailNow()
			}
		})
	}
}

func TestVerifyUser(t *testing.T) {
	// just to make sure that no one will mess with an order of the func paramaters
	testCases := []struct{
		name string
		args struct{
			ctx context.Context
			user *models.User
		}
		beforeTest func(mockPool *pgxpoolmock.MockPgxIface, user *models.User)
		err *apierror.ErrorStruct
	}{
		{
			name: "VerifyUser",
			args: struct{ctx context.Context; user *models.User}{
				ctx: context.TODO(),
				user: new(models.User),
			},
			beforeTest: func(mockPool *pgxpoolmock.MockPgxIface, user *models.User) {
				mockPool.EXPECT().Exec(gomock.Any(), queryVerifyUser, user.RefreshToken, user.ExpirationTime, user.VerificationCode).Return(nil, nil).Times(1)
			},
			err: apierror.NewErrorStruct(ErrWrongVerificationCode.Error(), "error", http.StatusBadRequest),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			MockPool := pgxpoolmock.NewMockPgxIface(ctrl)
			repository := NewUserRepository(MockPool)
			tc.beforeTest(MockPool, tc.args.user)
			err := repository.VerifyUser(tc.args.ctx, tc.args.user)
			Err, ok := err.(*apierror.ErrorStruct)
			if !ok {
				t.FailNow()
			} else {
				if !reflect.DeepEqual(*tc.err, *Err) {
					t.FailNow()
				}
			}
		})
	}
}

func TestGetUUid(t *testing.T) {
	// just to make sure that no one will mess with an order of the func paramaters
	testCases := []struct{
		name string
		args struct{
			ctx context.Context
			user *models.User
		}
		beforeTest func(mockPool *pgxpoolmock.MockPgxIface, user *models.User)
	}{
		{
			name: "GetUUid",
			args: struct{ctx context.Context; user *models.User}{
				ctx: context.TODO(),
				user: new(models.User),
			},
			beforeTest: func(mockPool *pgxpoolmock.MockPgxIface, user *models.User) {
				mockPool.EXPECT().QueryRow(gomock.Any(), queryGetUUid, user.RefreshToken, gomock.Any()).Return(pgxpoolmock.NewRow(uuid.New().String())).Times(1)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			MockPool := pgxpoolmock.NewMockPgxIface(ctrl)
			repository := NewUserRepository(MockPool)
			tc.beforeTest(MockPool, tc.args.user)
			if err := repository.GetUUid(tc.args.ctx, tc.args.user); err != nil {
				t.FailNow()
			}
		})
	}
}

func TestGetVerifiedUser(t *testing.T) {
	// just to make sure that no one will mess with an order of the func paramaters
	testCases := []struct{
		name string
		args struct{
			ctx context.Context
			user *models.User
		}
		beforeTest func(mockPool *pgxpoolmock.MockPgxIface, user *models.User)
	}{
		{
			name: "GetVerifiedUser",
			args: struct{ctx context.Context; user *models.User}{
				ctx: context.TODO(),
				user: &models.User{
					UserId: uuid.New(),
					Password: "testPassword",
					RefreshToken: "testToken",
					ExpirationTime: time.Now().UTC(),
				},
			},
			beforeTest: func(mockPool *pgxpoolmock.MockPgxIface, user *models.User) {
				mockPool.EXPECT().QueryRow(gomock.Any(), queryGetVerifiedUser, user.Email).Return(pgxpoolmock.NewRow(user.UserId, user.Password, user.RefreshToken, user.ExpirationTime)).Times(1)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			MockPool := pgxpoolmock.NewMockPgxIface(ctrl)
			repository := NewUserRepository(MockPool)
			tc.beforeTest(MockPool, tc.args.user)
			user, err := repository.GetVerifiedUser(tc.args.ctx, tc.args.user)
			if !reflect.DeepEqual(user, tc.args.user) || err != nil {
				t.FailNow()
			}
		})
	}
}

func TestUpdateRefreshToken(t *testing.T) {
	// just to make sure that no one will mess with an order of the func paramaters
	testCases := []struct{
		name string
		args struct{
			ctx context.Context
			user *models.User
		}
		beforeTest func(mockPool *pgxpoolmock.MockPgxIface, user *models.User)
	}{
		{
			name: "UpdateRefreshToken",
			args: struct{ctx context.Context; user *models.User}{
				ctx: context.TODO(),
				user: new(models.User),
			},
			beforeTest: func(mockPool *pgxpoolmock.MockPgxIface, user *models.User) {
				mockPool.EXPECT().Exec(gomock.Any(), queryUpdateRefreshToken, user.RefreshToken, user.ExpirationTime, user.Email).Return(nil, nil).Times(1)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			MockPool := pgxpoolmock.NewMockPgxIface(ctrl)
			repository := NewUserRepository(MockPool)
			tc.beforeTest(MockPool, tc.args.user)
			if err := repository.UpdateRefreshToken(tc.args.ctx, tc.args.user); err != nil {
				t.FailNow()
			}
		})
	}
}

func TestIfUnverifiedUserExists(t *testing.T) {
	// just to make sure that no one will mess with an order of the func paramaters
	testCases := []struct{
		name string
		args struct{
			ctx context.Context
			user *models.User
		}
		beforeTest func(mockPool *pgxpoolmock.MockPgxIface, user *models.User)
	}{
		{
			name: "IfUnverifiedUserExists",
			args: struct{ctx context.Context; user *models.User}{
				ctx: context.TODO(),
				user: new(models.User),
			},
			beforeTest: func(mockPool *pgxpoolmock.MockPgxIface, user *models.User) {
				mockPool.EXPECT().QueryRow(gomock.Any(), queryIfUnverifiedUserExists, user.Email).Return(pgxpoolmock.NewRow(true)).Times(1)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			MockPool := pgxpoolmock.NewMockPgxIface(ctrl)
			repository := NewUserRepository(MockPool)
			tc.beforeTest(MockPool, tc.args.user)
			err := repository.IfUnverifiedUserExists(tc.args.ctx, tc.args.user, new(bool))
			if err != nil {
				t.Fail()
			}
		})
	}
}