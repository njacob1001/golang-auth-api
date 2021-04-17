// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package storagemocks

import (
	context "context"
	domain "rumm-api/internal/core/domain"

	mock "github.com/stretchr/testify/mock"

	security "rumm-api/kit/security"
)

// AccountRepository is an autogenerated mock type for the AccountRepository type
type AccountRepository struct {
	mock.Mock
}

// Authenticate provides a mock function with given fields: ctx, accIdentifier, password
func (_m *AccountRepository) Authenticate(ctx context.Context, accIdentifier string, password string) (domain.Account, *security.TokenDetails, error) {
	ret := _m.Called(ctx, accIdentifier, password)

	var r0 domain.Account
	if rf, ok := ret.Get(0).(func(context.Context, string, string) domain.Account); ok {
		r0 = rf(ctx, accIdentifier, password)
	} else {
		r0 = ret.Get(0).(domain.Account)
	}

	var r1 *security.TokenDetails
	if rf, ok := ret.Get(1).(func(context.Context, string, string) *security.TokenDetails); ok {
		r1 = rf(ctx, accIdentifier, password)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*security.TokenDetails)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, string, string) error); ok {
		r2 = rf(ctx, accIdentifier, password)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// Create provides a mock function with given fields: ctx, account
func (_m *AccountRepository) Create(ctx context.Context, account domain.Account) (*security.TokenDetails, error) {
	ret := _m.Called(ctx, account)

	var r0 *security.TokenDetails
	if rf, ok := ret.Get(0).(func(context.Context, domain.Account) *security.TokenDetails); ok {
		r0 = rf(ctx, account)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*security.TokenDetails)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, domain.Account) error); ok {
		r1 = rf(ctx, account)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Logout provides a mock function with given fields: ctx, accessUuid
func (_m *AccountRepository) Logout(ctx context.Context, accessUuid string) error {
	ret := _m.Called(ctx, accessUuid)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, accessUuid)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
