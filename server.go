package main

import (
	"context"
	"log"
	"net/http"
	"time"

	models "accounts/models"
	routes "accounts/routes"
	server "accounts/server"
	util "accounts/utils"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	var router = mux.NewRouter()

	clientOptions := options.Client().ApplyURI("mongodb://admin:password@0.0.0.0:27020")

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// ensure connected to server
	err = client.Ping(ctx, nil)

	if err != nil {
		log.Fatal(err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     "0.0.0.0:6379", // use default Addr
		Password: "",             // no password set
		DB:       0,              // use default DB
	})

	_, err = rdb.Ping().Result()

	if err != nil {
		log.Fatal(err)
	}

	userModel := models.NewUserModel(client.Database("accounts").Collection("users"))

	context := &server.Context{
		UserModel: userModel,
		Rdb:       rdb,
	}

	// Set response content type to application/json
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			w.Header().Add("Content-Type", "application/json")
			next.ServeHTTP(w, r)
		})
	})

	publisher := util.NewPublisher(rdb)
	publisher.Publish("create.user", "server started")
	// merge router and models into one context?
	routes.SetUpRoutes(router, context)

	http.ListenAndServe(":6060", router)
}
