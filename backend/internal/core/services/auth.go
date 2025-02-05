package services

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/hafiztri123/config"
	"github.com/hafiztri123/internal/core/domain"
	"github.com/hafiztri123/internal/repositories/postgres"
	"github.com/redis/go-redis/v9"
)

type AuthService struct {
   userRepo *postgres.UserRepository
   redisClient *redis.Client
   jwtSecret string
   refreshSecret string
   refreshTTL time.Duration
   accessTTL time.Duration
}

var (
    ErrInvalidRefreshToken = errors.New("invalid refresh token")
    ErrInvalidTokenType = errors.New("invalid token type")
    ErrTokenNotFound = errors.New("token not found")
    
)

func NewAuthService(userRepo *postgres.UserRepository, redisClient *redis.Client, config *config.Config) *AuthService {
   return &AuthService{
       userRepo: userRepo,
       redisClient: redisClient,
       jwtSecret: config.JWT.Secret,
       refreshSecret: config.JWT.RefreshSecret,
       refreshTTL: config.JWT.RefreshTTL,
       accessTTL: config.JWT.AccessTTL,
   }
}

func (s *AuthService) Login(email, password string) (*domain.TokenPair, error) {
    fmt.Printf("Attempting to find user with email: %s\n", email)
    user, err := s.userRepo.FindByEmail(email)
    if err != nil {
        fmt.Printf("Database error finding user: %v\n", err)
        return nil, err
    }
   if err := user.ValidatePassword(password); err != nil {

       return nil, err
   }

//    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
//        "sub": user.ID,
//        "exp": time.Now().Add(24 * time.Hour).Unix(),
//    })

   return s.generateTokenPair(user.ID)
}

func (s *AuthService) ValidateToken(tokenString string) (*domain.User, error) {
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

func (s *AuthService) Register(user *domain.User) error {
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

func (s *AuthService) generateTokenPair(userID uuid.UUID) (*domain.TokenPair, error) {
    accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "sub": userID,
        "exp": time.Now().Add(s.accessTTL).Unix(),
        "type": "access", 
    })

    refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "sub": userID,
        "exp": time.Now().Add(s.refreshTTL).Unix(),
        "type": "refresh", 
    })

    accessTokenString, err := accessToken.SignedString([]byte(s.jwtSecret))
    if err != nil {
        return nil, err
    }

    refreshTokenString, err := refreshToken.SignedString([]byte(s.refreshSecret))
    if err != nil {
        return nil, err
    }

    err = s.redisClient.Set(context.Background(), refreshTokenString, userID.String(), s.refreshTTL).Err()
    if err != nil {
        return nil, err
    }

    return &domain.TokenPair{
        AccessToken: accessTokenString,
        RefreshToken: refreshTokenString,
    }, nil
}

func (s *AuthService) RefreshToken(refreshToken string) (*domain.TokenPair, error) {
    claims, err := s.validateRefreshToken(refreshToken)
    if err != nil {
        return nil, err
    }

    userID := claims["sub"].(string)
    storedUserID, err := s.redisClient.Get(context.Background(), refreshToken).Result()
    if err != nil || storedUserID != userID {
        return nil, ErrInvalidRefreshToken
    }

    err = s.RevokeRefreshToken(refreshToken)
    if err != nil {
        return nil, err
    }



    return s.generateTokenPair(uuid.MustParse(userID))
}

func (s *AuthService) validateRefreshToken(tokenString string) (jwt.MapClaims, error) {
    token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
        return []byte(s.refreshSecret), nil
    })

    if err != nil || !token.Valid {
        fmt.Print(token)
        return nil, ErrInvalidRefreshToken
    }

    claims := token.Claims.(jwt.MapClaims)

    if claims["type"] != "refresh" {
        return nil, ErrInvalidTokenType 
    }
    
    return claims, nil
}

func (s *AuthService) RevokeRefreshToken(token string) error {
    result := s.redisClient.Del(context.Background(), token)

    if result.Err() != nil {
        return result.Err()
    }

    if result.Val() == 0 {
        return ErrTokenNotFound
    }

    return nil
}