package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/huandu/go-sqlbuilder"
	"gorm.io/gorm"
	"rumm-api/internal/core/domain"
	"time"
)

var clientSQLStruck = sqlbuilder.NewStruct(new(sqlClient)).For(sqlbuilder.PostgreSQL)
var clientInfoSQLStruck = sqlbuilder.NewStruct(new(clientInfo)).For(sqlbuilder.PostgreSQL)
var updateClientSQLStruck = sqlbuilder.NewStruct(new(sqlUpdateClient)).For(sqlbuilder.PostgreSQL)

type ClientRepository struct {
	db        *gorm.DB
	dbTimeout time.Duration
	rdb       *redis.Client
}

func NewClientRepository(db *gorm.DB, dbTimeout time.Duration, rdb *redis.Client) *ClientRepository {
	return &ClientRepository{
		db:        db,
		dbTimeout: dbTimeout,
		rdb:       rdb,
	}
}

func (r *ClientRepository) Create(_ context.Context, _ domain.Client) error {


	return nil


}

func (r *ClientRepository) Find(_ context.Context, _ string) (domain.Client, error) {
	return domain.Client{}, nil
}

func (r *ClientRepository) Delete(_ context.Context, _ string) error {

	return nil
}

func (r *ClientRepository) Update(_ context.Context, _ string, _ domain.Client) error {


	return nil
}

func (r *ClientRepository) CreateTemporal(ctx context.Context, client domain.Client) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()
	c := sqlClient{
		ID:        client.ID,
		Email:     client.Email,
		City:      client.City,
		Cellphone: client.Cellphone,
		Name:      client.Name,
		LastName:  client.LastName,
		Address:   client.Address,
		Birthday:  client.Birthday,
	}

	j, err := json.Marshal(c)

	if err != nil {
		return fmt.Errorf("error trying to persisr client on cache: %v", err)
	}

	infoDuration := 10 * time.Minute

	if err := r.rdb.Set(ctxTimeout, client.ID, j, infoDuration).Err(); err != nil {
		return fmt.Errorf("error trying to persisr client on cache: %v", err)
	}

	return nil
}
