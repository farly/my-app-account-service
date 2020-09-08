package accounts

import (
	schema "accounts/datastore/schema"
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

/**
 * Pass mongo client here
 */
type UserModel struct {
	collection *mongo.Collection
}

func NewUserModel(collection *mongo.Collection) *UserModel {
	return &UserModel{collection}
}

func (model *UserModel) List() schema.Users {
	return accountList
}

func (model *UserModel) Create(user schema.User) error {

	_, err := model.collection.InsertOne(context.Background(), user)

	return err
}

var accountList = schema.Users{
	schema.User{
		ID:        "testId",
		Firstname: "Firstname",
		Lastname:  "lastname",
	},
}
