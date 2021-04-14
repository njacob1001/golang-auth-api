package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/huandu/go-sqlbuilder"
	"rumm-api/internal/core/domain"
	"rumm-api/kit/encryption"
	"time"
)

type AccountRepository struct {
	db        *sql.DB
	dbTimeout time.Duration
}

var accountSQLStruck = sqlbuilder.NewStruct(new(sqlAccount)).For(sqlbuilder.PostgreSQL)
var accountInfoSQLStruck = sqlbuilder.NewStruct(new(accountInfo)).For(sqlbuilder.PostgreSQL)

func NewAccountRepository(db *sql.DB, dbTimeout time.Duration) *AccountRepository {
	return &AccountRepository{
		db:        db,
		dbTimeout: dbTimeout,
	}
}

func (r *AccountRepository) Create(ctx context.Context, account domain.Account) error {
	password, err := encryption.GetHash(account.Password())
	if err != nil {
		return fmt.Errorf("error trying to persist account on database: %v", err)
	}
	query, args := accountSQLStruck.InsertInto(sqlAccountTable, sqlAccount{
		ID:          account.ID(),
		Identifier:  account.Identifier(),
		Password:    password,
		AccountType: account.AccountType(),
	}).Build()

	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	_, err = r.db.ExecContext(ctxTimeout, query, args...)
	if err != nil {
		return fmt.Errorf("error trying to persist account on database: %v", err)
	}

	return nil
}
func (r *AccountRepository) Authenticate(ctx context.Context, accIdentifier, password string) (domain.Account, error) {

	sb := accountInfoSQLStruck.SelectFrom(sqlAccountTable)
	query, args := sb.Where(sb.Equal("identifier", accIdentifier)).Build()

	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	row := r.db.QueryRowContext(ctxTimeout, query, args...)

	var acc accountInfo
	if err := row.Scan(accountInfoSQLStruck.Addr(&acc)...); err != nil {
		return domain.Account{}, fmt.Errorf("error trying to find account on database, account doesn't exist: %v", err)
	}

	isValid, err := encryption.ValidatePassword(acc.Password, password)

	if err != nil {
		return domain.Account{}, fmt.Errorf("error in accout validation: %v: ", err)
	}
	if isValid {
		account, err :=  domain.NewAccount(
			domain.WithAccountID(acc.ID, acc.Identifier),
			domain.WithAccountType(acc.AccountID))

		if err != nil {
			return domain.Account{}, fmt.Errorf("error in accout format: %v: ", err)
		}
		return account, nil
	}

	return domain.Account{}, nil
}
