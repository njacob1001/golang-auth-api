package creating

import (
	"context"
	rumm "rumm-api/internal/client"
)

type ClientService struct {
	clientRepository rumm.ClientRepository
}

func NewClientService(clientRepository rumm.ClientRepository) ClientService {
	return ClientService{
		clientRepository: clientRepository,
	}
}

func (service ClientService) CreateClient(ctx context.Context, uuid, name, lastName, birthday, email, city, address, cellphone, password string) error {
	client, err := rumm.NewClient(uuid, name, lastName, birthday, email, city, address, cellphone, password)
	if err != nil {
		return err
	}
	return service.clientRepository.Save(ctx, client)
}
