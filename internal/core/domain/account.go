package domain

import (
	"errors"
	"fmt"
	"rumm-api/kit/security"
	"time"
)

type Company struct {
	ID                  string `json:"id" db:"id" validate:"required,uuid4"`
	Name                string `json:"name" db:"name" validate:"required"`
	Nit                 string `json:"nit" db:"nit" validate:"required"`
	NumberOfEmployees   int    `json:"number_of_employees" db:"number_of_employees" validate:"required"`
	Phone               string `json:"phone" db:"phone" validate:"required"`
	LegalRepresentative string `json:"legal_representative" db:"legal_representative" validate:"required"`
	DocumentNumber      string `json:"document_number" db:"document_number" validate:"required"`
	Email               string `json:"email" db:"email" validate:"required"`
	Logo                string `json:"logo" db:"logo" validate:"required"`
}

type Profile struct {
	ID        string `json:"id" db:"id" validate:"required,uuid4"`
	Code      string `json:"code" db:"code" validate:"required"`
	IsActive  bool   `json:"is_active" db:"is_active" validate:"required"`
	DarkMode  bool   `json:"dark_mode" db:"dark_mode" validate:"required"`
	AccountID string `json:"account_id" db:"account_id",validate:"required,uuid4" validate:"required"`
}

type Person struct {
	ID        string    `json:"id" db:"id" validate:"required,uuid4"`
	Name      string    `json:"name" db:"name" validate:"required"`
	IDType    string    `json:"id_type" db:"id_type" validate:"required"`
	IDNumber  string    `json:"id_number" db:"id_number" validate:"required"`
	LastName  string    `json:"last_name" db:"last_name" validate:"required"`
	Cellphone string    `json:"cellphone" db:"cellphone" validate:"required"`
	BirthDate time.Time `json:"birth_date" db:"birth_date" validate:"required"`
	Email     string    `json:"email" db:"email" validate:"required"`
	Country   string    `json:"country" db:"country" validate:"required"`
	City      string    `json:"city" db:"city" validate:"required"`
	Address   string    `json:"address" db:"address" validate:"required"`
	Photo     string    `json:"photo" db:"photo" validate:"required"`
	CompanyID string    `json:"company_id" db:"company_id",validate:"required,uuid4"`
}

type Account struct {
	ID            string    `json:"id" db:"id" validate:"required,uuid4"`
	Identifier    string    `json:"identifier" db:"identifier" validate:"required"`
	Password      []byte    `json:"password" db:"password" validate:"required"`
	RequestLogin  bool      `json:"request_login" db:"request_login"`
	ResetPassword bool      `json:"request_reset_password" db:"request_reset_password"`
	PersonID      string    `json:"person_id" db:"person_id" validate:"required"`
	TypeID        string    `json:"type_id" db:"type_id" validate:"required,uuid4"`
	LastLogin     time.Time `json:"last_login" db:"last_login"`
}

var ErrAccountValidation = errors.New("account validation error")

func (a Account) ValidatePassword(password string) (bool, error) {
	isValid, err := security.ValidatePassword(a.Password, password)

	if err != nil {
		return false, fmt.Errorf("%w", ErrAccountValidation)
	}

	return isValid, nil
}
