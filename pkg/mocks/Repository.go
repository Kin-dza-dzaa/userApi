// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"

	models "github.com/Kin-dza-dzaa/userApi/internal/models"
	mock "github.com/stretchr/testify/mock"
)

// Repository is an autogenerated mock type for the Repository type
type Repository struct {
	mock.Mock
}

// AddUser provides a mock function with given fields: ctx, user
func (_m *Repository) AddUser(ctx context.Context, user *models.User) error {
	ret := _m.Called(ctx, user)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.User) error); ok {
		r0 = rf(ctx, user)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetUUid provides a mock function with given fields: ctx, user
func (_m *Repository) GetUUid(ctx context.Context, user *models.User) error {
	ret := _m.Called(ctx, user)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.User) error); ok {
		r0 = rf(ctx, user)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetVerifiedUser provides a mock function with given fields: ctx, user
func (_m *Repository) GetVerifiedUser(ctx context.Context, user *models.User) (*models.User, error) {
	ret := _m.Called(ctx, user)

	var r0 *models.User
	if rf, ok := ret.Get(0).(func(context.Context, *models.User) *models.User); ok {
		r0 = rf(ctx, user)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *models.User) error); ok {
		r1 = rf(ctx, user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IfUnverifiedUserExists provides a mock function with given fields: ctx, user, result
func (_m *Repository) IfUnverifiedUserExists(ctx context.Context, user *models.User, result *bool) error {
	ret := _m.Called(ctx, user, result)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.User, *bool) error); ok {
		r0 = rf(ctx, user, result)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateCredentials provides a mock function with given fields: ctx, user
func (_m *Repository) UpdateCredentials(ctx context.Context, user *models.User) error {
	ret := _m.Called(ctx, user)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.User) error); ok {
		r0 = rf(ctx, user)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateRefreshToken provides a mock function with given fields: ctx, user
func (_m *Repository) UpdateRefreshToken(ctx context.Context, user *models.User) error {
	ret := _m.Called(ctx, user)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.User) error); ok {
		r0 = rf(ctx, user)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// VerifyUser provides a mock function with given fields: ctx, user
func (_m *Repository) VerifyUser(ctx context.Context, user *models.User) error {
	ret := _m.Called(ctx, user)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.User) error); ok {
		r0 = rf(ctx, user)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewRepository creates a new instance of Repository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewRepository(t mockConstructorTestingTNewRepository) *Repository {
	mock := &Repository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
