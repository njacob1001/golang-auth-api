package service

import (
	"context"
	"rumm-api/internal/core/domain"
	"rumm-api/internal/core/port"
	"rumm-api/kit/security"
)

type AccountService struct {
	clientRepository  port.ClientRepository
	accountRepository port.AccountRepository
}

func NewAccountService(accountRepository port.AccountRepository, clientRepository port.ClientRepository) AccountService {
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

func(s AccountService) CreateTemporalClient(ctx context.Context, uuid, name, lastName, birthday, email, city, address, cellphone string) error {
	client, err := domain.NewClient(
		uuid,
		domain.WithAccount(email, cellphone),
		domain.WithLocation(city, address),
		domain.WithPersonalInformation(name, lastName, birthday),
	)
	if err != nil {
		return err
	}
	return s.clientRepository.CreateTemporal(ctx, client)
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

func (s AccountService) CreateAccount(ctx context.Context, id, identifier, password, accountType, clientID string) (*security.TokenDetails, error) {
	account, err := domain.NewAccount(
		domain.WithAccountID(id),
		domain.WithAccountPass(password),
		domain.WithAccountIdentifier(identifier),
		domain.WithAccountType(accountType))

	if err != nil {
		return nil, err
	}

	client, err := s.accountRepository.GetTemporalClient(ctx, clientID)
	if err != nil {
		return nil, err
	}

	return s.accountRepository.Create(ctx, account, client)
}

func (s AccountService) Authenticate(ctx context.Context, accIdentifier, password string) (domain.Account, *security.TokenDetails, error) {
	return s.accountRepository.Authenticate(ctx, accIdentifier, password)
}

func (s AccountService) Logout(ctx context.Context, accessUUID string) error {
	return s.accountRepository.Logout(ctx, accessUUID)
}

func (s AccountService) Refresh(ctx context.Context, refreshToken string) (*security.TokenDetails, error) {
	return s.accountRepository.Refresh(ctx, refreshToken)
}
