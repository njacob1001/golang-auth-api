package domain

import (
	"errors"
	"fmt"
	"rumm-api/kit/security"
	"time"
)


type Account struct {
	ID            string
	Type   string
	Identifier    string
	Password      []byte
	RequestLogin  bool
	ResetPassword bool
	PersonID      string
	LastLogin     time.Time
	Created       time.Time
	Updated       time.Time
}


var ErrAccountValidation = errors.New("account validation error")


func (a Account) ValidatePassword(password string) (bool, error) {
	isValid, err := security.ValidatePassword(a.Password, password)

	if err != nil {
		return false, fmt.Errorf("%w", ErrAccountValidation)
	}

	return isValid, nil
}
