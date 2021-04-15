package accountservice

import (
	"context"
	"rumm-api/internal/core/domain"
	"rumm-api/internal/core/ports"
)

type AccountService struct {
	clientRepository  ports.ClientRepository
	accountRepository ports.AccountRepository
}

func NewAccountService(accountRepository ports.AccountRepository, clientRepository ports.ClientRepository) AccountService {
	return AccountService{
		clientRepository:  clientRepository,
		accountRepository: accountRepository,
	}
}

func (s AccountService) CreateClient(ctx context.Context, uuid, name, lastName, birthday, email, city, address, cellphone string) error {
	client, err := domain.NewClient(
		uuid,
		domain.WithAccount(email, cellphone),
		domain.WithLocation(city, address),
		domain.WithPersonalInformation(name, lastName, birthday),
	)
	if err != nil {
		return err
	}
	return s.clientRepository.Create(ctx, client)
}

func (s AccountService) FindClientByID(ctx context.Context, id string) (domain.Client, error) {
	return s.clientRepository.Find(ctx, id)
}

func (s AccountService) DeleteClientByID(ctx context.Context, clientID string) error {
	return s.clientRepository.Delete(ctx, clientID)
}

func (s AccountService) UpdateClientByID(ctx context.Context, uuid, name, lastName, birthday, email, city, address, cellphone string) error {
	client, err := domain.NewClient(
		uuid,
		domain.WithAccount(email, cellphone),
		domain.WithLocation(city, address),
		domain.WithPersonalInformation(name, lastName, birthday),
	)
	if err != nil {
		return err
	}
	return s.clientRepository.Update(ctx, uuid, client)
}
func (s AccountService) CreateAccount(ctx context.Context, id, identifier, password, accountType string) error {
	account, err := domain.NewAccount(
		domain.WithAccountID(id),
		domain.WithAccountPass(password),
		domain.WithAccountIdentifier(identifier),
		domain.WithAccountType(accountType))

	if err != nil {
		return err
	}

	return s.accountRepository.Create(ctx, account)
}

func (s AccountService) Authenticate(ctx context.Context, accIdentifier, password string) (domain.Account, error) {
	return s.accountRepository.Authenticate(ctx, accIdentifier, password)
}
