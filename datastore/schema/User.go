package accounts

type User struct {
	ID        string `json:"id,omitempty" bson:"_id,omitempty"`
	Firstname string `json:"firstname" validate:"required"`
	Lastname  string `json:"lastname" validate:"required"`
}

type Users []User
