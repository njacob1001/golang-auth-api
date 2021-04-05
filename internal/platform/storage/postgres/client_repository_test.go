package postgres

import (
	"context"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	rumm "rumm-api/internal/client"
	"testing"
)

func Test_ClientRepository_Save_RepositoryError(t *testing.T) {
	UUID, Name, LastName, Birthday, Email, City, Address, password, Cellphone := "66021013-a0ce-4104-b29f-329686825aeb", "test", "test", "2020-01-01", "test", "test", "test", "test", "testing"

	client, err := rumm.NewClient(UUID, Name, LastName, Birthday, Email, City, Address, Cellphone, password)
	require.NoError(t, err)

	db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)

	sqlMock.ExpectExec(
		"INSERT INTO clients (id, name, last_name, birth_day, email, city, address, password, cellphone) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
	).WithArgs(UUID, Name, LastName, Birthday, Email, City, Address, password, Cellphone).WillReturnError(errors.New("something-failed"))

	repo := NewClientRepository(db)

	err = repo.Save(context.Background(), client)

	assert.NoError(t, sqlMock.ExpectationsWereMet())
	assert.Error(t, err)
}

func Test_ClientRepository_Save_RepositorySucceed(t *testing.T){
	UUID, Name, LastName, Birthday, Email, City, Address, password, Cellphone := "66021013-a0ce-4104-b29f-329686825aeb", "test", "test", "2020-01-01", "test", "test", "test", "test", "testing"

	client, err := rumm.NewClient(UUID, Name, LastName, Birthday, Email, City, Address, Cellphone, password)
	require.NoError(t, err)

	db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)

	sqlMock.ExpectExec(
		"INSERT INTO clients (id, name, last_name, birth_day, email, city, address, password, cellphone) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
	).WithArgs(UUID, Name, LastName, Birthday, Email, City, Address, password, Cellphone).WillReturnResult(sqlmock.NewResult(0,1))

	repo := NewClientRepository(db)

	err = repo.Save(context.Background(), client)

	assert.NoError(t, sqlMock.ExpectationsWereMet())
	assert.NoError(t, err)
}
