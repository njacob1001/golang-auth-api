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

	client := domain.Client{
		ID:        uuid,
		Name:      email,
		LastName:  name,
		Birthday:  birthday,
		City:      city,
		Address:   address,
		Email:     email,
		Cellphone: cellphone,
	}
	return s.clientRepository.Create(ctx, client)
}

func (s AccountService) CreateTemporalClient(ctx context.Context, uuid, name, lastName, birthday, email, city, address, cellphone string) error {
	client := domain.Client{
		ID:        uuid,
		Name:      email,
		LastName:  name,
		Birthday:  birthday,
		City:      city,
		Address:   address,
		Email:     email,
		Cellphone: cellphone,
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
	client := domain.Client{
		ID:        uuid,
		Name:      email,
		LastName:  name,
		Birthday:  birthday,
		City:      city,
		Address:   address,
		Email:     email,
		Cellphone: cellphone,
	}
	return s.clientRepository.Update(ctx, uuid, client)
}

func (s AccountService) CreateAccount(ctx context.Context, person domain.Person, account domain.Account, profile domain.Profile) (*security.TokenDetails, error) {

	hash, err := security.GetHash(string(account.Password))

	if err != nil {
		return nil, err
	}

	newAcc := domain.Account{
		ID:         account.ID,
		Identifier: account.Identifier,
		TypeID:       account.Identifier,
		Password:   hash,
	}

	if err != nil {
		return nil, err
	}

	return s.accountRepository.Create(ctx, newAcc, profile, person)
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
