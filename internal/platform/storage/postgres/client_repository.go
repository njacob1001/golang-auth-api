package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/huandu/go-sqlbuilder"
	rumm "rumm-api/internal/client"
	"time"
)

type ClientRepository struct {
	db *sql.DB
	dbTimeout time.Duration
}

func NewClientRepository(db *sql.DB, dbTimeout time.Duration) *ClientRepository {
	return &ClientRepository{
		db: db,
		dbTimeout: dbTimeout,
	}
}

func (repository *ClientRepository) Save(ctx context.Context, client rumm.Client) error {

	clientSQLStruck := sqlbuilder.NewStruct(new(sqlClient))

	query, args := clientSQLStruck.InsertInto(sqlClientTable, sqlClient{
		ID:        client.ID(),
		Name:      client.Name(),
		LastName:  client.LastName(),
		Birthday:  client.BirthDay(),
		Email:     client.Email(),
		City:      client.City(),
		Address:   client.Address(),
		Password:  client.Password(),
		Cellphone: client.Cellphone(),
	}).Build()

	ctxTimeout, cancel := context.WithTimeout(ctx, repository.dbTimeout)
	defer cancel()

	_, err := repository.db.ExecContext(ctxTimeout, query, args...)

	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("error trying to persist client on database: %v", err)
	}

	return nil
}
