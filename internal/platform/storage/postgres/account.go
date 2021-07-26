package postgres

import "time"

const (
	sqlAccountTable = "accounts"
)

type Profile struct {
	ID        string `db:"id"`
	Code      string `db:"code"`
	IsActive  bool   `db:"is_active"`
	DarkMode  bool   `db:"dark_mode"`
	AccountID string `db:"account_id"`
}

type CreatorClientPerson struct {
	ID        string    `db:"id"`
	Name      string    `db:"name"`
	IDType    string    `db:"id_type"`
	IDNumber  string    `db:"id_number"`
	LastName  string    `db:"last_name"`
	Cellphone string    `db:"cellphone"`
	BirthDate time.Time `db:"birth_date"`
	Email     string    `db:"email"`
	Country   string    `db:"country"`
	City      string    `db:"city"`
	Address   string    `db:"address"`
	Photo     string    `db:"photo"`
}

type Person struct {
	ID        string    `db:"id"`
	Name      string    `db:"name"`
	IDType    string    `db:"id_type"`
	IDNumber  string    `db:"id_number"`
	LastName  string    `db:"last_name"`
	Cellphone string    `db:"cellphone"`
	BirthDate time.Time `db:"birth_date"`
	Email     string    `db:"email"`
	Country   string    `db:"country"`
	City      string    `db:"city"`
	Address   string    `db:"address"`
	Photo     string    `db:"photo"`
	CompanyID string    `db:"company_id"`
}

type CreatorAccount struct {
	ID         string `db:"id"`
	Identifier string `db:"identifier"`
	Password   string `db:"password"`
	PersonID   string `db:"person_id"`
	TypeID     string `db:"type_id"`
}

type Account struct {
	ID            string    `db:"id"`
	Identifier    string    `db:"identifier"`
	Password      string    `db:"password"`
	RequestLogin  bool      `db:"request_login"`
	ResetPassword bool      `db:"reset_password"`
	PersonID      string    `db:"person_id"`
	TypeID        string    `db:"type_id"`
	LastLogin     time.Time `db:"last_login"`
}

type sqlAccount struct {
	ID          string `db:"id"`
	Password    string `db:"password"`
	Identifier  string `db:"identifier"`
	AccountType string `db:"type_id"`
	ClientID    string `db:"client_id"`
}
type accountInfo struct {
	ID         string `db:"id"`
	Identifier string `db:"identifier"`
	Password   string `db:"password"`
	//CreatedAt    time.Time `db:"created_at"`
	//UpdatedAt    time.Time `db:"updated_at"`
	//LastLogin    time.Time `db:"last_login"`
	//RequestLogin bool      `db:"request_login"`
	AccountID string `db:"type_id"`
}
