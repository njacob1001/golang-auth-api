package domain

import "rumm-api/kit/identifier"

type Account struct {
	id           string
	accountType  string
	identifier   string
	password     string
}

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

func WithAccountID(id, accIdentifier string) AccountOption {
	return func(a *Account) error {
		safeUUID, err := identifier.ValidateIdentifier(id)
		if err != nil {
			return err
		}
		a.id = safeUUID.String
		a.identifier = accIdentifier
		return nil
	}
}

func WithAccountPass(password string) AccountOption {
	return func(a *Account) error {
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

func (a Account) Identifier() string {
	return a.identifier
}

func (a Account) Password() string {
	return a.password
}

func (a Account) AccountType() string{
	return a.accountType
}