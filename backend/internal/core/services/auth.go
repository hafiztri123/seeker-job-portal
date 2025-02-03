package services

import (
	"errors"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/hafiztri123/internal/core/domain"
	"github.com/hafiztri123/internal/repositories/postgres"
)

type authService struct {
   userRepo *postgres.UserRepository
   jwtSecret string
}

func NewAuthService(userRepo *postgres.UserRepository) *authService {
   return &authService{
       userRepo: userRepo,
       jwtSecret: "your-secret-key", // Should come from config
   }
}

func (s *authService) Login(email, password string) (string, error) {
   user, err := s.userRepo.FindByEmail(email)
   if err != nil {
       return "", err
   }

   if err := user.ValidatePassword(password); err != nil {
       return "", err
   }

   token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
       "sub": user.ID,
       "exp": time.Now().Add(24 * time.Hour).Unix(),
   })

   return token.SignedString([]byte(s.jwtSecret))
}

func (s *authService) ValidateToken(tokenString string) (*domain.User, error) {
   token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
       return []byte(s.jwtSecret), nil
   })

   if err != nil || !token.Valid {
       return nil, err
   }

   claims := token.Claims.(jwt.MapClaims)
   userID := claims["sub"].(string)

   return s.userRepo.FindByID(userID)
}

func (s *authService) Register(user *domain.User) error {
   if err := validateRegister(user); err != nil {
       return err
   }

   return s.userRepo.Create(user)
}

func validateRegister(user *domain.User) error {
   emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
   if !emailRegex.MatchString(user.Email) {
       return errors.New("invalid email format")
   }
   if user.FullName == "" {
       return errors.New("full name is required")
   }

   if user.HashedPassword == "" {
       return errors.New("password is required")
   }

   return nil
}