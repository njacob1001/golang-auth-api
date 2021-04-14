package service

import (
	"context"
	"rumm-api/internal/core/domain"
	"rumm-api/internal/core/ports"
)

type AccountService struct {
	clientRepository ports.ClientRepository
	accountRepository ports.AccountRepository
}

func NewAccountService(accountRepository ports.AccountRepository,clientRepository ports.ClientRepository) AccountService {
	return AccountService{
		clientRepository: clientRepository,
		accountRepository: accountRepository,
	}
}

func (service AccountService) CreateClient(ctx context.Context, uuid, name, lastName, birthday, email, city, address, cellphone string) error {
	client, err := domain.NewClient(
		uuid,
		domain.WithAccount(email, cellphone),
		domain.WithLocation(city, address),
		domain.WithPersonalInformation(name, lastName, birthday),
	)
	if err != nil {
		return err
	}
	return service.clientRepository.Save(ctx, client)
}

func (service AccountService) FindClientByID(ctx context.Context, id string) (domain.Client, error) {
	return service.clientRepository.FindByID(ctx, id)
}

func (service AccountService) DeleteClientByID(ctx context.Context, clientID string) error {
	return service.clientRepository.DeleteByID(ctx, clientID)
}