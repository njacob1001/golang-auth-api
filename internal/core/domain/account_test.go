package domain

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"rumm-api/kit/security"
	"testing"
)

func TestAccountDomain(t *testing.T) {
	t.Run("Validate password should be true", func(t *testing.T) {

		pass, err := security.GetHash("testingpass")
		require.NoError(t, err)

		acc := Account{
			Password: string(pass),
		}

		result, err := acc.ValidatePassword("testingpass")
		require.NoError(t, err)

		assert.True(t, result)
	})

	t.Run("Validate password should be false", func(t *testing.T) {
		pass, err := security.GetHash("testingpass")
		require.NoError(t, err)

		acc := Account{
			Password: string(pass),
		}

		result, err := acc.ValidatePassword("testingpass ")

		assert.False(t, result)
	})
}
