package database

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const mongoURL = "mongodb://mongo:27017"

func ConnectToMongo() (*mongo.Client, error) {
	//connect options
	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "admin",
	})
	conn, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("error: connecting mongo")
		return nil, err
	}

	return conn, nil
}
