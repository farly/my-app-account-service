package accounts

import (
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        string `json:"id,omitempty" bson:"_id,omitempty"`
	Firstname string `json:"firstname" validate:"required"`
	Lastname  string `json:"lastname" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required"`
}

func (user *User) HashPassword() *User {
	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)

	if err != nil {
		return nil
	}

	user.Password = string(hashed)

	return user
}

type Users []User
