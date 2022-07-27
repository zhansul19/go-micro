package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/zhansul19/log-service/database"
)

type Config struct {
	Models database.Models
}

func (c *Config) serve() {
	srv := http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: c.routes(),
	}
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
func (c *Config) routes() http.Handler {
	mux := chi.NewMux()

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	mux.Use(middleware.Heartbeat("/ping"))
	mux.Post("/log", c.Writelog)

	return mux
}
