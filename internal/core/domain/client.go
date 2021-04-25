package domain

import (
	"rumm-api/kit/identifier"
)

type Client struct {
	id        string `json: "id"`
	name      string `json:"name"`
	lastName  string `json:"last_name"`
	birthday  string `json:"birthday"`
	email     string `json:"email"`
	city      string `json:"city"`
	address   string `json:"address"`
	cellphone string `json:"cellphone"`
}

type ClientOption func(*Client) error

func NewClient(uuid string, options ...ClientOption) (Client, error) {

	safeUUID, err := identifier.ValidateIdentifier(uuid)
	if err != nil {
		return Client{}, err
	}
	client := Client{}
	client.id = safeUUID.String

	for _, option := range options {
		err := option(&client)
		if err != nil {
			return Client{}, err
		}
	}

	return client, nil
}

func WithAccount(email, cellphone string) ClientOption {
	return func(client *Client) error {
		client.email = email
		client.cellphone = cellphone
		return nil
	}
}

func WithPersonalInformation(name, lastName, birthDay string) ClientOption {
	return func(client *Client) error {
		client.name = name
		client.lastName = lastName
		client.birthday = birthDay
		return nil
	}
}

func WithLocation(city, address string) ClientOption {
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
