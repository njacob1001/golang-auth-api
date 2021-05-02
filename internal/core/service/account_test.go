package service

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"rumm-api/mocks/mockups"
	"testing"
)

func TestAccountService(t *testing.T) {
	t.Run("Test ClientService CreateClient Repository Error", func(t *testing.T) {
		UUID, Name, LastName, Birthday, Email, City, Address, Cellphone := "66021013-a0ce-4104-b29f-329686825aeb", "test", "test", "2020-01-01", "test", "test", "test", "testing"

		clientRepositoryMock := new(storagemocks.ClientRepository)
		accountRepositoryMock := new(storagemocks.AccountRepository)
		clientRepositoryMock.On("Create", mock.Anything, mock.Anything).Return(errors.New("something unexpected happened"))

		clientService := NewAccountService(accountRepositoryMock, clientRepositoryMock)

		err := clientService.CreateClient(context.Background(), UUID, Name, LastName, Birthday, Email, City, Address, Cellphone)

		clientRepositoryMock.AssertExpectations(t)
		assert.Error(t, err)
	})

	t.Run("Test ClientService CreateClient Succeed", func(t *testing.T) {
		UUID, Name, LastName, Birthday, Email, City, Address, Cellphone := "66021013-a0ce-4104-b29f-329686825aeb", "test", "test", "2020-01-01", "test", "test", "test", "testing"

		clientRepositoryMock := new(storagemocks.ClientRepository)
		accountRepositoryMock := new(storagemocks.AccountRepository)
		clientRepositoryMock.On("Create", mock.Anything, mock.Anything).Return(nil)

		clientService := NewAccountService(accountRepositoryMock, clientRepositoryMock)

		err := clientService.CreateClient(context.Background(), UUID, Name, LastName, Birthday, Email, City, Address, Cellphone)

		clientRepositoryMock.AssertExpectations(t)
		assert.NoError(t, err)
	})
}
