package rumm

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
)

type ClientRepository interface {
	Save(ctx context.Context, client Client) error
}

//go:generate mockery --case=snake --outpkg=storagemocks --output=../platform/storage/storagemocks --name=ClientRepository

type Client struct {
	id        string
	name      string
	lastName  string
	birthday  string
	email     string
	city      string
	address   string
	cellphone string
	password  string
}

var ErrInvalidClientUUID = errors.New("invalid Client UUID")

type ClientUUID struct {
	value string
}

func NewClientUUID(value string) (ClientUUID, error) {
	v, err := uuid.Parse(value)
	if err != nil {
		return ClientUUID{}, fmt.Errorf("%w: %s", ErrInvalidClientUUID, value)
	}

	return ClientUUID{
		value: v.String(),
	}, nil
}

func NewClient(uuid, name, lastName, birthday, email, city, address, cellphone, password string) (Client, error) {
	safeUUID, err := NewClientUUID(uuid)
	if err != nil {
		return Client{}, err
	}

	return Client{
		id:        safeUUID.value,
		name:      name,
		lastName:  lastName,
		birthday:  birthday,
		email:     email,
		city:      city,
		address:   address,
		cellphone: cellphone,
		password:  password,
	}, nil
}

func (c Client) ID() string {
	return c.id
}

func (c Client) Name() string {
	return c.name
}
func (c Client) LastName() string {
	return c.lastName
}
func (c Client) BirthDay() string {
	return c.birthday
}
func (c Client) Email() string {
	return c.email
}
func (c Client) City() string {
	return c.city
}
func (c Client) Address() string {
	return c.address
}
func (c Client) Cellphone() string {
	return c.cellphone
}

func (c Client) Password() string {
	return c.password
}
