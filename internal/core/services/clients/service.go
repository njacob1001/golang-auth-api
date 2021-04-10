package clients

import (
	"context"
	"rumm-api/internal/core/domain"
	"rumm-api/internal/core/ports"
)

type ClientService struct {
	clientRepository ports.ClientRepository
}

func NewClientService(clientRepository ports.ClientRepository) ClientService {
	return ClientService{
		clientRepository: clientRepository,
	}
}

func (service ClientService) CreateClient(ctx context.Context, uuid, name, lastName, birthday, email, city, address, cellphone, password string) error {
	client, err := domain.NewClient(
		uuid,
		domain.WithAccount(email, password, cellphone),
		domain.WithLocation(city, address),
		domain.WithPersonalInformation(name, lastName, birthday),
	)
	if err != nil {
		return err
	}
	return service.clientRepository.Save(ctx, client)
}


