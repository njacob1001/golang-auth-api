package ports

import (
	"context"
	"rumm-api/internal/core/domain"
)

type ClientRepository interface {
	Save(ctx context.Context, client domain.Client) error
	FindByID(ctx context.Context, clientID string) (domain.Client, error)
}
//go:generate mockery --case=snake --outpkg=storagemocks --output=../../../mocks/mockups --name=ClientRepository
