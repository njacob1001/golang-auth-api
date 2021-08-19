package service

import (
	"context"
	"rumm-api/internal/core/domain"
	"rumm-api/internal/core/port"
	"rumm-api/kit/security"
)

type AccountService struct {
	accountRepository port.AccountRepository
}

func NewAccountService(accountRepository port.AccountRepository) AccountService {
	return AccountService{
		accountRepository: accountRepository,
	}
}

func (s AccountService) CreateAccount(ctx context.Context, person domain.Person, account domain.Account, profile domain.Profile) (*security.TokenDetails, error) {

	hash, err := security.GetHash(account.Password)

	if err != nil {
		return nil, err
	}

	newAcc := domain.Account{
		ID:         account.ID,
		Identifier: account.Identifier,
		TypeID:     account.TypeID,
		Password:   string(hash),
		PersonID:   account.PersonID,
	}

	if err != nil {
		return nil, err
	}

	return s.accountRepository.Create(ctx, newAcc, profile, person)
}

func (s AccountService) Authenticate(ctx context.Context, accIdentifier, password string) (*security.TokenDetails, error) {
	return s.accountRepository.Authenticate(ctx, accIdentifier, password)
}

func (s AccountService) Logout(ctx context.Context, accessUUID string) error {
	return s.accountRepository.Logout(ctx, accessUUID)
}

func (s AccountService) Refresh(ctx context.Context, refreshToken string) (*security.TokenDetails, error) {
	return s.accountRepository.Refresh(ctx, refreshToken)
}
