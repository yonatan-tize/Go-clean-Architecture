// Code generated by mockery v2.44.1. DO NOT EDIT.

package mocks

import (
	context "context"
	Domain "example/go-clean-architecture/Domain"

	mock "github.com/stretchr/testify/mock"
)

// UserUseCase is an autogenerated mock type for the UserUseCase type
type UserUseCase struct {
	mock.Mock
}

// AuthenticateUser provides a mock function with given fields: ctx, userName, password
func (_m *UserUseCase) AuthenticateUser(ctx context.Context, userName string, password string) (Domain.User, string, error) {
	ret := _m.Called(ctx, userName, password)

	if len(ret) == 0 {
		panic("no return value specified for AuthenticateUser")
	}

	var r0 Domain.User
	var r1 string
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) (Domain.User, string, error)); ok {
		return rf(ctx, userName, password)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) Domain.User); ok {
		r0 = rf(ctx, userName, password)
	} else {
		r0 = ret.Get(0).(Domain.User)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) string); ok {
		r1 = rf(ctx, userName, password)
	} else {
		r1 = ret.Get(1).(string)
	}

	if rf, ok := ret.Get(2).(func(context.Context, string, string) error); ok {
		r2 = rf(ctx, userName, password)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// CreateAccount provides a mock function with given fields: ctx, user
func (_m *UserUseCase) CreateAccount(ctx context.Context, user *Domain.User) (Domain.User, error) {
	ret := _m.Called(ctx, user)

	if len(ret) == 0 {
		panic("no return value specified for CreateAccount")
	}

	var r0 Domain.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *Domain.User) (Domain.User, error)); ok {
		return rf(ctx, user)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *Domain.User) Domain.User); ok {
		r0 = rf(ctx, user)
	} else {
		r0 = ret.Get(0).(Domain.User)
	}

	if rf, ok := ret.Get(1).(func(context.Context, *Domain.User) error); ok {
		r1 = rf(ctx, user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateUserRole provides a mock function with given fields: ctx, id
func (_m *UserUseCase) UpdateUserRole(ctx context.Context, id string) error {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for UpdateUserRole")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewUserUseCase creates a new instance of UserUseCase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUserUseCase(t interface {
	mock.TestingT
	Cleanup(func())
}) *UserUseCase {
	mock := &UserUseCase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
