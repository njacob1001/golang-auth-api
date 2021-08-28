package port

import (
	"context"
	"rumm-api/internal/core/domain"
	"rumm-api/kit/security"
)

type AccountRepository interface {
	Create(ctx context.Context, account domain.Account, profile domain.Profile, person domain.Person) (*security.TokenDetails, error)
	Authenticate(ctx context.Context, accIdentifier, password string) (*security.TokenDetails, error)
	Logout(ctx context.Context, accessUuid string) error
	Refresh(ctx context.Context, token string) (*security.TokenDetails, error)
	ValidateRegister(ctx context.Context, account domain.Account, profile domain.Profile, person domain.Person) error
}

//go:generate mockery --case=snake --outpkg=storagemocks --output=../../../mocks/mockups --name=AccountRepository
