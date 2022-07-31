package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/zhansul19/go-micro/log-service/database"
	"go.mongodb.org/mongo-driver/mongo"
)

const webPort = "80"

var client *mongo.Client


func main() {
	// connect to mongo
	mongoClient, err := database.ConnectToMongo()
	if err != nil {
		log.Panic(err)
	}
	client = mongoClient

	// create a context in order to disconnect
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// close connection
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	app := Config{
		Models: database.New(client),
	}

	// start web server
	// go app.serve()
	log.Println("Starting service on port", webPort)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Panic()
	}

}
