package domain

import (
	"errors"
	"fmt"
	"rumm-api/kit/security"
	"time"
)

type Profile struct {
	ID        string `json:"id" validate:"required,uuid4"`
	Code      string `json:"code" validate:"required"`
	IsActive  bool   `json:"is_active" validate:"required"`
	DarkMode  bool   `json:"dark_mode" validate:"required"`
	AccountID string `json:"account_id" validate:"required,uuid4"`
}

type Person struct {
	ID        string    `json:"id" validate:"required,uuid4"`
	Name      string    `json:"name" validate:"required"`
	IDType    string    `json:"id_type" validate:"required"`
	IDNumber  string    `json:"id_number" validate:"required"`
	LastName  string    `json:"last_name" validate:"required"`
	Cellphone string    `json:"cellphone" validate:"required,e164"`
	BirthDate time.Time `json:"birth_date" validate:"required"`
	Email     string    `json:"email" validate:"required"`
	Country   string    `json:"country"`
	City      string    `json:"city"`
	Address   string    `json:"address"`
	Photo     string    `json:"photo"`
	CompanyID string    `json:"company_id"`
}

type Account struct {
	ID            string    `json:"id" db:"id" validate:"required,uuid4"`
	Identifier    string    `json:"identifier" validate:"required"`
	Password      string    `json:"password" validate:"required"`
	RequestLogin  bool      `json:"request_login"`
	ResetPassword bool      `json:"request_reset_password"`
	PersonID      string    `json:"person_id" validate:"required,uuid4"`
	TypeID        string    `json:"type_id" validate:"required,uuid4"`
	LastLogin     time.Time `json:"last_login"`
}

var ErrAccountValidation = errors.New("account validation error")

func (a Account) ValidatePassword(password string) (bool, error) {
	isValid, err := security.ValidatePassword([]byte(a.Password), password)

	if err != nil {
		return false, fmt.Errorf("%w", ErrAccountValidation)
	}

	return isValid, nil
}
