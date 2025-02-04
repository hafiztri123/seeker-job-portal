package domain

import (
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID uuid.UUID `json:"id"`
	Email string `json:"email"`
	HashedPassword string `json:"-"`
	FullName string `json:"full_name"`
	PhoneNumber string `json:"phone_number"`
	About string `json:"about"`
	ProfilePicture string `json:"profile_picture"`
}




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
	err := bcrypt.CompareHashAndPassword([]byte(u.HashedPassword), []byte(password))
	if err != nil {
		return errors.New("invalid password")
	}

	return nil
}

