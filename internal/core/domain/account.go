package domain

import (
	"errors"
	"fmt"
	"rumm-api/kit/identifier"
	"rumm-api/kit/security"
)

type Account struct {
	id           string
	accountType  string
	identifier   string
	password     string
}

var ErrInvalidClientUUID = errors.New("invalid password")
var ErrAccountValidation = errors.New("account validation error")

type AccountOption func(*Account) error

func NewAccount(options ...AccountOption) (Account, error) {
	account := Account{}
	for _, option := range options {
		err := option(&account)
		if err != nil {
			return Account{}, err
		}
	}

	return account, nil

}

func WithAccountID(id string) AccountOption {
	return func(a *Account) error {
		safeUUID, err := identifier.ValidateIdentifier(id)
		if err != nil {
			return err
		}
		a.id = safeUUID.String
		return nil
	}
}

func WithAccountPass(password string) AccountOption {
	return func(a *Account) error {
		hashPassword, err := security.GetHash(password)
		if err != nil {
			return fmt.Errorf("%w", ErrInvalidClientUUID)
		}
		a.password = hashPassword
		return nil
	}
}

func WithAccountHashedPass(password string) AccountOption {
	return func(a *Account) error {
		if password == "" {
			return ErrInvalidClientUUID
		}
		a.password = password
		return nil
	}
}

func WithAccountIdentifier(accIdentifier string) AccountOption {
	return func(a *Account) error {
		a.identifier = accIdentifier
		return nil
	}
}


func WithAccountType(accountType string) AccountOption {
	return func(a *Account) error {
		a.accountType = accountType
		return nil
	}
}


func (a Account) ID() string {
	return a.id
}

func (a Account) Password() string {
	return a.password
}

func (a Account) Identifier() string {
	return a.identifier
}

func (a Account) AccountType() string{
	return a.accountType
}

func (a Account) ValidatePassword(password string) (bool, error) {
	isValid, err := security.ValidatePassword(a.password, password)

	if err != nil {
		return false, fmt.Errorf("%w", ErrAccountValidation)
	}

	return isValid, nil
}
