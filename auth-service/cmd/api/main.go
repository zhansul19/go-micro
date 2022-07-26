package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/zhansul19/go-micro/auth-service/data"
	"github.com/zhansul19/go-micro/auth-service/database"
)

const webPort = "8001"

func main() {
	log.Printf("Connecting to auth service on port: %s \n", webPort)

	//connect to postgres
	conn := database.ConnectToDB()
	if conn == nil {
		log.Panic("can't connect to Postgres...")
	}
	app := &Config{
		DB:     conn,
		Models: data.New(conn),
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Panic(err)
	}

}
