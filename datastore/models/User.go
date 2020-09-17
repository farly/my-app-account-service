package accounts

import (
	schema "accounts/datastore/schema"
	"context"

	"go.mongodb.org/mongo-driver/bson"
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

	_, err := model.collection.InsertOne(context.Background(), user.HashPassword())

	return err
}

func (model *UserModel) FindOneByUsername(email string) (schema.User, error) {
	user := schema.User{}

	err := model.collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)

	return user, err
}

var accountList = schema.Users{
	schema.User{
		ID:        "testId",
		Firstname: "Firstname",
		Lastname:  "lastname",
	},
}
