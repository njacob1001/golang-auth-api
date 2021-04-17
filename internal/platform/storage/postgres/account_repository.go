package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/huandu/go-sqlbuilder"
	"rumm-api/internal/core/domain"
	"rumm-api/kit/security"
	"time"
)

type AccountRepository struct {
	db        *sql.DB
	dbTimeout time.Duration
	jwtSecret string
	rdb       *redis.Client
}

var accountSQLStruck = sqlbuilder.NewStruct(new(sqlAccount)).For(sqlbuilder.PostgreSQL)
var accountInfoSQLStruck = sqlbuilder.NewStruct(new(accountInfo)).For(sqlbuilder.PostgreSQL)

func NewAccountRepository(db *sql.DB, dbTimeout time.Duration, jwtSecret string, rdb *redis.Client) *AccountRepository {
	return &AccountRepository{
		db:        db,
		dbTimeout: dbTimeout,
		jwtSecret: jwtSecret,
		rdb:       rdb,
	}
}

func (r *AccountRepository) Create(ctx context.Context, account domain.Account) (*security.TokenDetails, error) {
	query, args := accountSQLStruck.InsertInto(sqlAccountTable, sqlAccount{
		ID:          account.ID(),
		Identifier:  account.Identifier(),
		Password:    account.Password(),
		AccountType: account.AccountType(),
	}).Build()

	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	_, err := r.db.ExecContext(ctxTimeout, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error trying to persist account on database: %v", err)
	}

	td, err := security.CreateToken(r.jwtSecret, account.ID())
	if err != nil {
		return nil, err
	}

	return td, nil
}
func (r *AccountRepository) Authenticate(ctx context.Context, accIdentifier, password string) (domain.Account, *security.TokenDetails, error) {

	sb := accountInfoSQLStruck.SelectFrom(sqlAccountTable)
	query, args := sb.Where(sb.Equal("identifier", accIdentifier)).Build()

	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	row := r.db.QueryRowContext(ctxTimeout, query, args...)

	var acc accountInfo
	if err := row.Scan(accountInfoSQLStruck.Addr(&acc)...); err != nil {
		return domain.Account{}, nil, fmt.Errorf("error trying to find account on database, account doesn't exist: %v", err)
	}

	account, err := domain.NewAccount(
		domain.WithAccountID(acc.ID),
		domain.WithAccountIdentifier(acc.Identifier),
		domain.WithAccountHashedPass(acc.Password),
		domain.WithAccountType(acc.AccountID))
	if err != nil {
		return domain.Account{}, nil, err
	}

	isValid, err := account.ValidatePassword(password)

	if err != nil {
		return domain.Account{}, nil, err
	}
	if isValid {
		td, err := security.CreateToken(r.jwtSecret, account.ID())
		if err != nil {
			return domain.Account{}, nil, err
		}

		if err := security.CreateAuth(ctxTimeout, account.ID(), td, r.rdb); err != nil {
			return domain.Account{}, nil, err
		}

		return account, td, nil
	}

	return domain.Account{}, nil, domain.ErrAccountValidation
}

func (r *AccountRepository) Logout(ctx context.Context, accessUUID string) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	_, err := security.DeleteAuth(ctxTimeout, r.rdb,accessUUID)

	if err != nil {
		return err
	}
	return nil



}