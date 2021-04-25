// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package storagemocks

import (
	context "context"
	domain "rumm-api/internal/core/domain"

	mock "github.com/stretchr/testify/mock"
)

// ClientRepository is an autogenerated mock type for the ClientRepository type
type ClientRepository struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, client
func (_m *ClientRepository) Create(ctx context.Context, client domain.Client) error {
	ret := _m.Called(ctx, client)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.Client) error); ok {
		r0 = rf(ctx, client)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreateTemporal provides a mock function with given fields: ctx, client
func (_m *ClientRepository) CreateTemporal(ctx context.Context, client domain.Client) error {
	ret := _m.Called(ctx, client)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.Client) error); ok {
		r0 = rf(ctx, client)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: ctx, clientID
func (_m *ClientRepository) Delete(ctx context.Context, clientID string) error {
	ret := _m.Called(ctx, clientID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, clientID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Find provides a mock function with given fields: ctx, clientID
func (_m *ClientRepository) Find(ctx context.Context, clientID string) (domain.Client, error) {
	ret := _m.Called(ctx, clientID)

	var r0 domain.Client
	if rf, ok := ret.Get(0).(func(context.Context, string) domain.Client); ok {
		r0 = rf(ctx, clientID)
	} else {
		r0 = ret.Get(0).(domain.Client)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, clientID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: ctx, clientID, client
func (_m *ClientRepository) Update(ctx context.Context, clientID string, client domain.Client) error {
	ret := _m.Called(ctx, clientID, client)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, domain.Client) error); ok {
		r0 = rf(ctx, clientID, client)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
