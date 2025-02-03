package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestUser(t *testing.T) {
	t.Run("ValidatePassword", func(t *testing.T) {
		password := "correctpass"
		hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		assert.NoError(t, err)
		user := &User{
			Email: "test@example.com",
			HashedPassword: string(hashedBytes),
		}

		assert.NoError(t, user.ValidatePassword("correctpass"))
		assert.Error(t, user.ValidatePassword("wrongpass"))
	})
}

