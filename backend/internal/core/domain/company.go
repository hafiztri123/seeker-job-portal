package domain

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Company struct {
	Id uuid.UUID `json:"id"`
	Email string `json:"email"`
	HashedPassword string `json:"-"`
	Name string `json:"name"`
	PhoneNumber *string `json:"phone_number,omitempty"`
	Location *json.RawMessage `json:"location,omitempty"`
	About *string `json:"about,omitempty"`
	BussinessCategory *string `json:"bussiness_category,omitempty"`
	Size *string `json:"company_size,omitempty"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	ProfilePicture *string `json:"profile_picture,omitempty"`
}

func NewCompany(name, email, password string) (*Company, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}


	return &Company{
		Id: uuid.New(),
		Email: email,
		HashedPassword: string(hashedPassword),
		Name: name,
	}, nil
}

func (c *Company) ValidatePassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(c.HashedPassword), []byte(password))
	if err != nil {
		return errors.New("invalid password")
	}

	return nil
}