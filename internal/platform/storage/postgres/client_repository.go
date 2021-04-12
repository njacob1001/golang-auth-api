package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/huandu/go-sqlbuilder"
	"rumm-api/internal/core/domain"
	"time"
)

var clientSQLStruck = sqlbuilder.NewStruct(new(sqlClient)).For(sqlbuilder.PostgreSQL)
var clientInfoSQLStruck = sqlbuilder.NewStruct(new(clientInfo)).For(sqlbuilder.PostgreSQL)

type ClientRepository struct {
	db        *sql.DB
	dbTimeout time.Duration
}

func NewClientRepository(db *sql.DB, dbTimeout time.Duration) *ClientRepository {
	return &ClientRepository{
		db:        db,
		dbTimeout: dbTimeout,
	}
}

func (repository *ClientRepository) Save(ctx context.Context, client domain.Client) error {
	query, args := clientSQLStruck.InsertInto(sqlClientTable, sqlClient{
		ID:        client.ID(),
		Name:      client.Name(),
		LastName:  client.LastName(),
		Birthday:  client.BirthDay(),
		Email:     client.Email(),
		City:      client.City(),
		Address:   client.Address(),
		Cellphone: client.Cellphone(),
	}).Build()

	ctxTimeout, cancel := context.WithTimeout(ctx, repository.dbTimeout)
	defer cancel()

	_, err := repository.db.ExecContext(ctxTimeout, query, args...)

	if err != nil {
		return fmt.Errorf("error trying to persist client on database: %v", err)
	}

	return nil
}

func (repository *ClientRepository) FindByID(ctx context.Context, clientID string) (domain.Client, error) {

	sb := clientInfoSQLStruck.SelectFrom(sqlClientTable)
	query, args := sb.Where(sb.Equal("id", clientID)).Build()

	ctxTimeout, cancel := context.WithTimeout(ctx, repository.dbTimeout)
	defer cancel()

	row := repository.db.QueryRowContext(ctxTimeout, query, args...)

	var dbClient clientInfo

	if err := row.Scan(clientInfoSQLStruck.Addr(&dbClient)...); err != nil {
		return domain.Client{}, fmt.Errorf("error trying to find client on database, client doesn't exist: %v", err)
	}

	client, err := domain.NewClient(
		dbClient.ID,
		domain.WithAccount(dbClient.Email, dbClient.Cellphone),
		domain.WithLocation(dbClient.City, dbClient.Address),
		domain.WithPersonalInformation(dbClient.Name, dbClient.LastName, dbClient.Birthday.Format("2006-01-02")))
	if err != nil {
		return domain.Client{}, err
	}

	return client, nil
}
