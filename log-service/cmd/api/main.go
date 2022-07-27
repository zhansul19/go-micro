package main

import (
	"context"
	"log"
	"time"

	"github.com/zhansul19/log-service/database"

	"go.mongodb.org/mongo-driver/mongo"
)

const webPort = "80"

var client *mongo.Client

func main() {
	//connect to db
	conn, err := database.ConnectToMongo()
	if err != nil {
		log.Panic("could't connect mongo")
	}
	client = conn

	//context in order to disconnect
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	defer func() {
		if client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	config:= &Config{
		Models: database.New(client),
	}
	go config.serve()
}
