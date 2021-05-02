package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"rumm-api/internal/core/domain"
	"testing"
	"time"
)

func TestAccountRepository(t *testing.T) {
	t.Run("Should not create account correctly", func(t *testing.T) {
		id, accountType, identifier, password := "b157588e-75fe-40a4-b405-a7eed0c21663", "822d9c35-e508-4b23-ab9d-58336df8df19", "3213432233", "123345"
		account, err := domain.NewAccount(domain.WithAccountID(id), domain.WithAccountIdentifier(identifier), domain.WithAccountPass(password), domain.WithAccountType(accountType))
		require.NoError(t, err)
		cid, name, lastName, birthday, email, city, address, cellphone := "66021013-a0ce-4104-b29f-329686825aeb", "test", "test", "2020-01-01", "test", "test", "test", "testing"
		c, err := domain.NewClient(cid, domain.WithPersonalInformation(name, lastName, birthday), domain.WithLocation(city, address), domain.WithAccount(email, cellphone))
		require.NoError(t, err)

		db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		redisServer = mockRedis()
		redisClient := redis.NewClient(&redis.Options{
			Addr: redisServer.Addr(),
		})

		sqlMock.ExpectExec(
			"INSERT INTO clients (id, name, account_type_id, identifier, password) VALUES ($1, $2, $3, $4)",
		).WithArgs(id, accountType, identifier, password).WillReturnError(errors.New("something-failed"))

		repo := NewAccountRepository(db, 5*time.Second, "", redisClient)

		_, err = repo.Create(context.Background(), account, c)

		assert.Error(t, sqlMock.ExpectationsWereMet())
		assert.Error(t, err)

	})

	t.Run("Should create account correctly", func(t *testing.T) {
		id, accountType, identifier, password := "b157588e-75fe-40a4-b405-a7eed0c21663", "822d9c35-e508-4b23-ab9d-58336df8df19", "3213432233", "123345"
		account, err := domain.NewAccount(domain.WithAccountID(id), domain.WithAccountIdentifier(identifier), domain.WithAccountPass(password), domain.WithAccountType(accountType))
		require.NoError(t, err)
		cid, name, lastName, birthday, email, city, address, cellphone := "66021013-a0ce-4104-b29f-329686825aeb", "test", "test", "2020-01-01", "test", "test", "test", "testing"
		c, err := domain.NewClient(cid, domain.WithPersonalInformation(name, lastName, birthday), domain.WithLocation(city, address), domain.WithAccount(email, cellphone))

		require.NoError(t, err)

		db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		redisServer = mockRedis()
		redisClient := redis.NewClient(&redis.Options{
			Addr: redisServer.Addr(),
		})

		createClientQuery := "INSERT INTO clients (id, name, last_name, birth_day, email, city, address, cellphone) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"
		createAccountQuery := "INSERT INTO accounts (id, identifier, password, type_id, client_id) VALUES ($9, $10, $11, $12, $13)"
		query := fmt.Sprintf("WITH new_client as (%s) %s", createClientQuery, createAccountQuery)

		sqlMock.ExpectExec(query).WithArgs(
			c.ID(),
			c.Name(),
			c.LastName(),
			c.BirthDay(),
			c.Email(),
			c.City(),
			c.Address(),
			c.Cellphone(),
			account.ID(),
			account.Identifier(),
			account.Password(),
			account.AccountType(),
			c.ID(),
		).WillReturnResult(sqlmock.NewResult(0, 1))

		repo := NewAccountRepository(db, 5*time.Second, "", redisClient)

		_, err = repo.Create(context.Background(), account, c)

		assert.NoError(t, sqlMock.ExpectationsWereMet())
		assert.NoError(t, err)

	})

}
