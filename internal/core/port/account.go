package port

import (
	"context"
	"rumm-api/internal/core/domain"
	"rumm-api/kit/security"
)

type ClientRepository interface {
	Create(ctx context.Context, client domain.Client) error
	Find(ctx context.Context, clientID string) (domain.Client, error)
	Delete(ctx context.Context, clientID string) error
	Update(ctx context.Context, clientID string, client domain.Client) error
}

//go:generate mockery --case=snake --outpkg=storagemocks --output=../../../mocks/mockups --name=ClientRepository

type AccountRepository interface {
	Create(ctx context.Context, account domain.Account) (*security.TokenDetails, error)
	Authenticate(ctx context.Context, accIdentifier, password string) (domain.Account, *security.TokenDetails,error)
	Logout(ctx context.Context, accessUuid string) error
}

//go:generate mockery --case=snake --outpkg=storagemocks --output=../../../mocks/mockups --name=AccountRepository
