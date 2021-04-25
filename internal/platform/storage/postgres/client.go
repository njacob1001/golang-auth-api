package postgres

import "time"

const (
	sqlClientTable = "clients"
)


type sqlClient struct {
	ID        string `db:"id" json:"id"`
	Name      string `db:"name" json:"name"`
	LastName  string `db:"last_name" json:"last_name"`
	Birthday  string `db:"birth_day" json:"birthday"`
	Email     string `db:"email" json:"email"`
	City      string `db:"city" json:"city"`
	Address   string `db:"address" json:"address"`
	Cellphone string `db:"cellphone" json:"cellphone"`
}
type clientInfo struct {
	ID        string `db:"id"`
	Name      string `db:"name"`
	LastName  string `db:"last_name"`
	Birthday  time.Time `db:"birth_day"`
	Email     string `db:"email"`
	City      string `db:"city"`
	Address   string `db:"address"`
	Cellphone string `db:"cellphone"`
}

type sqlUpdateClient struct {
	Name      string `db:"name"`
	LastName  string `db:"last_name"`
	Birthday  string `db:"birth_day"`
	Email     string `db:"email"`
	City      string `db:"city"`
	Address   string `db:"address"`
	Cellphone string `db:"cellphone"`
}