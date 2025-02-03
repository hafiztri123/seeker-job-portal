package domain

import (
	"errors"
	"regexp"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID uuid.UUID `json:"id"`
	Email string `json:"email"`
	HashedPassword string `json:"-"`
	FullName string `json:"full_name"`
}

var (
	ErrRegisterInvalidInput = errors.New("register invalid input")
)

const (
	EmailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
)

var emailRegex = regexp.MustCompile(EmailRegex)


func NewUser(email, fullname, password string) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}


	return &User{
		ID: uuid.New(),
		Email: email,
		HashedPassword: string(hashedPassword),
		FullName: fullname,
	}, nil
}

func (u *User) ValidatePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.HashedPassword), []byte(password))
}

