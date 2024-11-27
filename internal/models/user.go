package models

import (
	"auth-service/internal/common"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type User struct {
	UserID                string     `json:"userID"`
	Login                 string     `json:"login"`
	VisibleID             string     `json:"visibleID"`
	HashedPassword        string     `json:"hashedPassword"`
	Person                Person     `json:"person"`
	Roles                 []*Role    `json:"roles"`
	CreatedAt             time.Time  `json:"createdAt"`
	UpdatedAt             time.Time  `json:"updatedAt"`
	DeletedAt             *time.Time `json:"deletedAt"`
	LastPasswordRestoreAt *time.Time `json:"lastPassword"`
	SearchIndex           *string    `json:"searchIndex"`
}

type Person struct {
	Firstname string    `json:"firstname"`
	Lastname  string    `json:"lastname"`
	Birthdate time.Time `json:"birthdate,omitempty"`
	Email     string    `json:"email,omitempty"`
	Gender    string    `json:"gender,omitempty"`
}

func (a *Person) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *Person) Scan(value interface{}) error {
	b, ok := value.(string)
	if !ok {
		return errors.New("type assertion to string failed")
	}

	return json.Unmarshal([]byte(b), &a)
}

type AuthResponse struct {
	PublicToken       string
	PublicTokenExpiry time.Time
}

type UserInput struct {
	Login     string
	Password  string
	Firstname string
	Lastname  string
	Birthdate time.Time
	Email     string
	Gender    string
	RoleIDs   []string
}

type UserUpdateInput struct {
	Login     *string
	Password  *string
	Firstname *string
	Lastname  *string
	Birthdate *time.Time
	Email     *string
	Gender    *string
	RoleIDs   []string
}

func (u *UserUpdateInput) ToUpdatedModel(oldUser *User) {
	if u.Login != nil {
		oldUser.Login = *u.Login
	}
	if u.Password != nil {
		oldUser.HashedPassword = common.HashPassword(*u.Password)
	}
	if u.Firstname != nil {
		oldUser.Person.Firstname = *u.Firstname
	}
	if u.Lastname != nil {
		oldUser.Person.Lastname = *u.Lastname
	}
	if u.Birthdate != nil {
		oldUser.Person.Birthdate = *u.Birthdate
	}
	if u.Email != nil {
		oldUser.Person.Email = *u.Email
	}
	if u.Gender != nil {
		oldUser.Person.Gender = *u.Gender
	}
}

type UserFilter struct {
	UserID *[]string
	Login  *[]string
	Email  *[]string
}
