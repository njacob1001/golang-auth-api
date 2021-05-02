package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAccountDomain(t *testing.T) {
	id, accountType, identifier, password := "b157588e-75fe-40a4-b405-a7eed0c21663", "822d9c35-e508-4b23-ab9d-58336df8df19", "3213432233", "123345"
	t.Run("Should create an account", func(t *testing.T) {
		_, err := NewAccount(WithAccountID(id), WithAccountIdentifier(identifier), WithAccountPass(password), WithAccountType(accountType))
		assert.NoError(t, err)
	})

	t.Run("Should not create an account with invalid id", func(t *testing.T) {
		_, err := NewAccount(WithAccountID("invalid id"))
		assert.Error(t, err)
	})
}
