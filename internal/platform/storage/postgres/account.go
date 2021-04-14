package postgres

import "time"

const (
	sqlAccountTable = "accounts"
)

type sqlAccount struct {
	ID           string    `db:"id"`
	Password     string    `db:"password"`
	Identifier   string    `db:"password"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
	LastLogin    time.Time `db:"last_login"`
	RequestLogin bool      `db:"request_login"`
}
type accountInfo struct {
	ID           string    `db:"id"`
	Identifier   string    `db:"password"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
	LastLogin    time.Time `db:"last_login"`
	RequestLogin bool      `db:"request_login"`
}
