package creating

import (
	"context"
	"errors"
	rumm "rumm-api/internal/client"
	"rumm-api/internal/platform/storage/storagemocks"

	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func Test_ClientService_CreateClient_RepositoryError(t *testing.T) {
	UUID, Name, LastName, Birthday, Email, City, Address, password, Cellphone := "66021013-a0ce-4104-b29f-329686825aeb", "test", "test", "2020-01-01", "test", "test", "test", "some", "testing"

	client, err := rumm.NewClient(UUID, Name, LastName, Birthday, Email, City, Address, Cellphone, password)
	require.NoError(t, err)

	clientRepositoryMock := new(storagemocks.ClientRepository)
	clientRepositoryMock.On("Save", mock.Anything, client).Return(errors.New("something unexpected happened"))

	clientService := NewClientService(clientRepositoryMock)

	err = clientService.CreateClient(context.Background(), UUID, Name, LastName, Birthday, Email, City, Address, Cellphone, password)

	clientRepositoryMock.AssertExpectations(t)
	assert.Error(t, err)
}

func Test_ClientService_CreateClient_Succeed(t *testing.T) {
	UUID, Name, LastName, Birthday, Email, City, Address, password, Cellphone := "66021013-a0ce-4104-b29f-329686825aeb", "test", "test", "2020-01-01", "test", "test", "test", "test", "testing"

	course, err := rumm.NewClient(UUID, Name, LastName, Birthday, Email, City, Address, Cellphone, password)
	require.NoError(t, err)

	clientRepositoryMock := new(storagemocks.ClientRepository)
	clientRepositoryMock.On("Save", mock.Anything, course).Return(nil)

	clientService := NewClientService(clientRepositoryMock)

	err = clientService.CreateClient(context.Background(), UUID, Name, LastName, Birthday, Email, City, Address, Cellphone, password)

	clientRepositoryMock.AssertExpectations(t)
	assert.NoError(t, err)
}
