package main

import (
	"context"
	"log"
	"net/http"
	"time"

	models "accounts/models"
	routes "accounts/routes"
	server "accounts/server"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	var router = mux.NewRouter()

	clientOptions := options.Client().ApplyURI("mongodb://admin:password@0.0.0.0:27020")

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)

	// ensure connected to server
	err = client.Ping(ctx, nil)

	if err != nil {
		log.Fatal(err)
	}

	userModel := models.NewUserModel(client.Database("accounts").Collection("users"))

	context := &server.Context{
		UserModel: userModel,
	}

	// merge router and models into one context?
	routes.SetUpRoutes(router, context)

	http.ListenAndServe(":6060", router)
}
