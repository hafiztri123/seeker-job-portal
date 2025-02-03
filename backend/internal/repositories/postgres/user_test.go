package postgres

import (
	"testing"

	"github.com/hafiztri123/internal/core/domain"
	"github.com/stretchr/testify/assert"
)

func TestUserRepositry(t *testing.T) {
	db := SetupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)

	t.Run("CreateAndFind", func(t *testing.T) {
		user, _ := domain.NewUser("test@example.com", "password123", "Test user")

		err := repo.Create(user)
		assert.NoError(t, err)

		found, err := repo.FindByEmail(user.Email)
		assert.NoError(t, err)
		assert.Equal(t, user.Email, found.Email)
		assert.Equal(t, user.FullName, found.FullName)
	})

	t.Run("NotFound", func(t *testing.T) {
		_, err := repo.FindByEmail("nonexistent@example.com")
		assert.Error(t, err)
	})
}