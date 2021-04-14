// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package storagemocks

import (
	context "context"
	domain "rumm-api/internal/core/domain"

	mock "github.com/stretchr/testify/mock"
)

// AccountRepository is an autogenerated mock type for the AccountRepository type
type AccountRepository struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, clientID, account
func (_m *AccountRepository) Create(ctx context.Context, clientID string, account domain.Account) error {
	ret := _m.Called(ctx, clientID, account)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, domain.Account) error); ok {
		r0 = rf(ctx, clientID, account)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
