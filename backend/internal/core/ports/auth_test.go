package ports

import (
	"testing"

	"github.com/hafiztri123/internal/core/domain"
	"github.com/hafiztri123/internal/core/services"
	"github.com/hafiztri123/internal/repositories/postgres"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthService(t *testing.T){
	db := postgres.SetupTestDB(t)
	defer db.Close()
	username := "test@example.com"
	password := "password123"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatal(err)
	}
	userRepo := postgres.NewUserRepository(db)
	authService := services.NewAuthService(userRepo)
	
	
	userRepo.Create(&domain.User{
		Email: username,
		HashedPassword: string(hashedPassword),
		FullName: "Test User",
	})






	t.Run("Login", func(t *testing.T) {
		

		token, err := authService.Login(username, password)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

	})

	t.Run("InvalidLogin", func(t *testing.T) {

		_, err := authService.Login("invalid_email", password)
		assert.Error(t, err)
	})

	t.Run("Register", func(t *testing.T) {
		user, _ := domain.NewUser("test1@example.com", "password123", "Test User")

		err := authService.Register(user)

		assert.NoError(t, err)
	})

	t.Run("InvalidRegister", func(t *testing.T) {
		user, _ := domain.NewUser("invalid_email", "password123", "Test User")

		err := authService.Register(user)

		assert.Error(t, err)
	})

	
	
}