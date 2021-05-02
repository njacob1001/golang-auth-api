package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestClientDomain(t *testing.T) {
	t.Run("Should create user successes", func(t *testing.T) {
		id, name, lastName, birthday, email, city, address, cellphone := "66021013-a0ce-4104-b29f-329686825aeb", "test", "test", "2020-01-01", "test", "test", "test", "testing"
		_, err := NewClient(id, WithPersonalInformation(name, lastName, birthday), WithLocation(city, address), WithAccount(email, cellphone))

		assert.NoError(t, err)

	})
	t.Run("Should not create user with invalid ID", func(t *testing.T) {
		_, err := NewClient("wrong-id")
		assert.Error(t, err)
	})
}
