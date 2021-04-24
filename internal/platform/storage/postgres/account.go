package postgres

const (
	sqlAccountTable = "accounts"
)

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
