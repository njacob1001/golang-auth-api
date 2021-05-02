package postgres

import (
	"context"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"rumm-api/internal/core/domain"
	"testing"
	"time"
)

var redisServer *miniredis.Miniredis

func mockRedis() *miniredis.Miniredis {
	s, err := miniredis.Run()

	if err != nil {
		panic(err)
	}

	return s
}

func TestClientRepository(t *testing.T) {

	t.Run("Test ClientRepository Save RepositoryError", func(t *testing.T) {
		UUID, Name, LastName, Birthday, Email, City, Address, Cellphone := "66021013-a0ce-4104-b29f-329686825aeb", "test", "test", "2020-01-01", "test", "test", "test", "testing"

		client, err := domain.NewClient(UUID,
			domain.WithPersonalInformation(Name, LastName, Birthday),
			domain.WithLocation(City, Address),
			domain.WithAccount(Email, Cellphone))
		require.NoError(t, err)

		db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		redisServer = mockRedis()
		redisClient := redis.NewClient(&redis.Options{
			Addr: redisServer.Addr(),
		})

		sqlMock.ExpectExec(
			"INSERT INTO clients (id, name, last_name, birth_day, email, city, address, cellphone) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		).WithArgs(UUID, Name, LastName, Birthday, Email, City, Address, Cellphone).WillReturnError(errors.New("something-failed"))

		repo := NewClientRepository(db, 5*time.Second, redisClient)

		err = repo.Create(context.Background(), client)

		assert.NoError(t, sqlMock.ExpectationsWereMet())
		assert.Error(t, err)
	})

	t.Run("Test ClientRepository Save RepositorySucceed", func(t *testing.T) {
		UUID, Name, LastName, Birthday, Email, City, Address, Cellphone := "66021013-a0ce-4104-b29f-329686825aeb", "test", "test", "2020-01-01", "test", "test", "test", "testing"

		client, err := domain.NewClient(UUID,
			domain.WithPersonalInformation(Name, LastName, Birthday),
			domain.WithLocation(City, Address),
			domain.WithAccount(Email, Cellphone))
		require.NoError(t, err)

		db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		redisServer = mockRedis()
		redisClient := redis.NewClient(&redis.Options{
			Addr: redisServer.Addr(),
		})

		sqlMock.ExpectExec(
			"INSERT INTO clients (id, name, last_name, birth_day, email, city, address, cellphone) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		).WithArgs(UUID, Name, LastName, Birthday, Email, City, Address, Cellphone).WillReturnResult(sqlmock.NewResult(0, 1))

		repo := NewClientRepository(db, 5*time.Second, redisClient)

		err = repo.Create(context.Background(), client)

		assert.NoError(t, sqlMock.ExpectationsWereMet())
		assert.NoError(t, err)
	})

}
