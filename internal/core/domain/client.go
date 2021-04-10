package domain



import (
	"context"
	"rumm-api/kit/identifier"
)

type ClientRepository interface {
	Save(ctx context.Context, client Client) error
}

//go:generate mockery --case=snake --outpkg=storagemocks --output=../../../mocks/mockups --name=ClientRepository

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

type Option func(*Client) error

func NewClient(uuid string, options ...Option) (Client, error) {

	safeUUID, err := identifier.ValidateIdentifier(uuid)
	if err != nil {
		return Client{}, err
	}
	client := Client{}
	client.id = safeUUID.String

	for _, option := range options {
		err := option(&client)
		if err != nil {
			return  Client{}, err
		}
	}

	return client, nil
}

func WithAccount(email, password, cellphone string) Option {
	return func(client *Client) error {
		client.password = password
		client.email = email
		client.cellphone = cellphone
		return nil
	}
}

func WithPersonalInformation(name, lastName, birthDay string) Option {
	return func(client *Client) error {
		client.name = name
		client.lastName = lastName
		client.birthday = birthDay
		return nil
	}
}

func WithLocation(city, address string) Option {
	return func(client *Client) error {
		client.city = city
		client.address = address
		return nil
	}
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
