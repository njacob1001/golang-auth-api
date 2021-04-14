package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/huandu/go-sqlbuilder"
	"rumm-api/internal/core/domain"
	"time"
)

type AccountRepository struct {
	db        *sql.DB
	dbTimeout time.Duration
}


var accountSQLStruck = sqlbuilder.NewStruct(new(sqlAccount)).For(sqlbuilder.PostgreSQL)

func NewAccountRepository(db *sql.DB, dbTimeout time.Duration) *AccountRepository {
	return &AccountRepository{
		db:        db,
		dbTimeout: dbTimeout,
	}
}

func (r *AccountRepository) Create(ctx context.Context, clientID string, account domain.Account) error {
	query, args := accountSQLStruck.InsertInto(sqlAccountTable, sqlAccount{
		ID:         account.ID(),
		Identifier: account.Identifier(),
		Password:   account.Password(),
	}).Build()

	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	_, err := r.db.ExecContext(ctxTimeout, query, args...)
	if err != nil {
		return fmt.Errorf("error trying to persist account on database: %v", err)
	}

	return nil
}