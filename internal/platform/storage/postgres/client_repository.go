package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/huandu/go-sqlbuilder"
	"rumm-api/internal/core/domain"
	"time"
)

var clientSQLStruck = sqlbuilder.NewStruct(new(sqlClient)).For(sqlbuilder.PostgreSQL)
var clientInfoSQLStruck = sqlbuilder.NewStruct(new(clientInfo)).For(sqlbuilder.PostgreSQL)
var updateClientSQLStruck = sqlbuilder.NewStruct(new(sqlUpdateClient)).For(sqlbuilder.PostgreSQL)

type ClientRepository struct {
	db        *sql.DB
	dbTimeout time.Duration
	rdb       *redis.Client
}

func NewClientRepository(db *sql.DB, dbTimeout time.Duration, rdb *redis.Client) *ClientRepository {
	return &ClientRepository{
		db:        db,
		dbTimeout: dbTimeout,
		rdb:       rdb,
	}
}

func (r *ClientRepository) Create(ctx context.Context, client domain.Client) error {
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

	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	_, err := r.db.ExecContext(ctxTimeout, query, args...)

	if err != nil {
		return fmt.Errorf("error trying to persist client on database: %v", err)
	}

	return nil
}

func (r *ClientRepository) Find(ctx context.Context, clientID string) (domain.Client, error) {

	sb := clientInfoSQLStruck.SelectFrom(sqlClientTable)
	query, args := sb.Where(sb.Equal("id", clientID)).Build()

	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	row := r.db.QueryRowContext(ctxTimeout, query, args...)

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

func (r *ClientRepository) Delete(ctx context.Context, clientID string) error {
	sb := clientInfoSQLStruck.DeleteFrom(sqlClientTable)
	query, args := sb.Where(sb.Equal("id", clientID)).Build()
	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()
	_, err := r.db.ExecContext(ctxTimeout, query, args...)
	if err != nil {
		return fmt.Errorf("error trying to delete client: %v", err)
	}

	return nil
}

func (r *ClientRepository) Update(ctx context.Context, clientID string, client domain.Client) error {

	newClient := sqlUpdateClient{
		Name:      client.Name(),
		LastName:  client.LastName(),
		Birthday:  client.BirthDay(),
		Cellphone: client.Cellphone(),
		Address:   client.Address(),
		Email:     client.Email(),
		City:      client.City(),
	}
	sb := updateClientSQLStruck.Update(sqlClientTable, newClient)
	query, args := sb.Where(sb.Equal("id", clientID)).Build()

	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	_, err := r.db.ExecContext(ctxTimeout, query, args...)
	if err != nil {
		return fmt.Errorf("error trying to update client: %v", err)
	}
	return nil
}

func (r *ClientRepository) CreateTemporal(ctx context.Context, client domain.Client) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()
	c := sqlClient{
		ID:        client.ID(),
		Email:     client.Email(),
		City:      client.City(),
		Cellphone: client.Cellphone(),
		Name:      client.Name(),
		LastName:  client.LastName(),
		Address:   client.Address(),
		Birthday:  client.BirthDay(),
	}

	j, err := json.Marshal(c)

	if err != nil {
		return fmt.Errorf("error trying to persisr client on cache: %v", err)
	}

	infoDuration := 10 * time.Minute

	if err := r.rdb.Set(ctxTimeout, client.ID(), j, infoDuration).Err(); err != nil {
		return fmt.Errorf("error trying to persisr client on cache: %v", err)
	}

	return nil
}
