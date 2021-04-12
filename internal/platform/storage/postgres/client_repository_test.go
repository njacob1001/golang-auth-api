package postgres

import (
	"context"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"rumm-api/internal/core/domain"
	"testing"
	"time"
)

func Test_ClientRepository_Save_RepositoryError(t *testing.T) {
	UUID, Name, LastName, Birthday, Email, City, Address, Cellphone := "66021013-a0ce-4104-b29f-329686825aeb", "test", "test", "2020-01-01", "test", "test", "test", "testing"

	client, err := domain.NewClient(UUID,
		domain.WithPersonalInformation(Name, LastName, Birthday),
		domain.WithLocation(City, Address),
		domain.WithAccount(Email, Cellphone))
	require.NoError(t, err)

	db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)

	sqlMock.ExpectExec(
		"INSERT INTO clientsHandler (id, name, last_name, birth_day, email, city, address, password, cellphone) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
	).WithArgs(UUID, Name, LastName, Birthday, Email, City, Address, Cellphone).WillReturnError(errors.New("something-failed"))

	repo := NewClientRepository(db, 5*time.Second)

	err = repo.Save(context.Background(), client)

	assert.NoError(t, sqlMock.ExpectationsWereMet())
	assert.Error(t, err)
}

func Test_ClientRepository_Save_RepositorySucceed(t *testing.T) {
	UUID, Name, LastName, Birthday, Email, City, Address, Cellphone := "66021013-a0ce-4104-b29f-329686825aeb", "test", "test", "2020-01-01", "test", "test", "test", "testing"

	client, err := domain.NewClient(UUID,
		domain.WithPersonalInformation(Name, LastName, Birthday),
		domain.WithLocation(City, Address),
		domain.WithAccount(Email, Cellphone))
	require.NoError(t, err)

	db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)

	sqlMock.ExpectExec(
		"INSERT INTO clientsHandler (id, name, last_name, birth_day, email, city, address, password, cellphone) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
	).WithArgs(UUID, Name, LastName, Birthday, Email, City, Address, Cellphone).WillReturnResult(sqlmock.NewResult(0, 1))

	repo := NewClientRepository(db, 5*time.Second)

	err = repo.Save(context.Background(), client)

	assert.NoError(t, sqlMock.ExpectationsWereMet())
	assert.NoError(t, err)
}
