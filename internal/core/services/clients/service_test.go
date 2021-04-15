package accountservice

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"rumm-api/internal/core/domain"
	"rumm-api/mocks/mockups"
	"testing"
)

func Test_ClientService_CreateClient_RepositoryError(t *testing.T) {
	UUID, Name, LastName, Birthday, Email, City, Address, Cellphone := "66021013-a0ce-4104-b29f-329686825aeb", "test", "test", "2020-01-01", "test", "test", "test", "testing"

	client, err := domain.NewClient(
		UUID,
		domain.WithPersonalInformation(Name, LastName, Birthday),
		domain.WithLocation(City, Address),
		domain.WithAccount(Email, Cellphone))
	require.NoError(t, err)

	clientRepositoryMock := new(storagemocks.ClientRepository)
	accountRepositoryMock := new(storagemocks.AccountRepository)
	clientRepositoryMock.On("Create", mock.Anything, client).Return(errors.New("something unexpected happened"))

	clientService := NewAccountService(accountRepositoryMock, clientRepositoryMock)

	err = clientService.CreateClient(context.Background(), UUID, Name, LastName, Birthday, Email, City, Address, Cellphone)

	clientRepositoryMock.AssertExpectations(t)
	assert.Error(t, err)
}

func Test_ClientService_CreateClient_Succeed(t *testing.T) {
	UUID, Name, LastName, Birthday, Email, City, Address, Cellphone := "66021013-a0ce-4104-b29f-329686825aeb", "test", "test", "2020-01-01", "test", "test", "test", "testing"

	course, err := domain.NewClient(UUID,
		domain.WithPersonalInformation(Name, LastName, Birthday),
		domain.WithLocation(City, Address),
		domain.WithAccount(Email, Cellphone))
	require.NoError(t, err)

	clientRepositoryMock := new(storagemocks.ClientRepository)
	accountRepositoryMock := new(storagemocks.AccountRepository)
	clientRepositoryMock.On("Create", mock.Anything, course).Return(nil)

	clientService := NewAccountService(accountRepositoryMock, clientRepositoryMock)

	err = clientService.CreateClient(context.Background(), UUID, Name, LastName, Birthday, Email, City, Address, Cellphone)

	clientRepositoryMock.AssertExpectations(t)
	assert.NoError(t, err)
}
