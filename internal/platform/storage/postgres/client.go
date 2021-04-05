package postgres

const (
	sqlClientTable = "clients"
)


type sqlClient struct {
	ID        string `db:"id"`
	Name      string `db:"name"`
	LastName  string `db:"last_name"`
	Birthday  string `db:"birth_day"`
	Email     string `db:"email"`
	City      string `db:"city"`
	Address   string `db:"address"`
	Password  string `db:"password"`
	Cellphone string `db:"cellphone"`
}
