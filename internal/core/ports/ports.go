package ports

import (
	"context"
	"rumm-api/internal/core/domain"
)

type ClientRepository interface {
	Save(ctx context.Context, client domain.Client) error
}
