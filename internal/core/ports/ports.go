package ports

import (
	"context"
	"rumm-api/internal/core/domain"
)

type ClientRepository interface {
	Create(ctx context.Context, client domain.Client) error
	Find(ctx context.Context, clientID string) (domain.Client, error)
	Delete(ctx context.Context, clientID string) error
	Update(ctx context.Context, clientID string, client domain.Client) error
}

//go:generate mockery --case=snake --outpkg=storagemocks --output=../../../mocks/mockups --name=ClientRepository

type AccountRepository interface {
	Create(ctx context.Context, clientID string, account domain.Account) error
}

//go:generate mockery --case=snake --outpkg=storagemocks --output=../../../mocks/mockups --name=AccountRepository
