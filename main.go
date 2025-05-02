package main

import (
	"backend-challenge/config"
	"backend-challenge/handlers"
	"backend-challenge/middleware"
	"backend-challenge/repository"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	cfg := config.Load()

	clientOpts := options.Client().ApplyURI(cfg.MongoURI)
	client, err := mongo.Connect(context.Background(), clientOpts)
	if err != nil {
		log.Fatal(err)
	}

	repo := repository.NewUserRepository(client.Database(cfg.DBName))

	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			count, err := repo.Count(context.Background())
			if err != nil {
				log.Printf("Error counting users: %v", err)
				continue
			}
			log.Printf("User count: %d", count)
		}
	}()

	r := mux.NewRouter()
	r.Use(middleware.Logger)

	uh := handlers.NewUserHandler(repo, cfg.JWTSecret)
	r.HandleFunc("/register", uh.Register).Methods("POST")
	r.HandleFunc("/login", uh.Login).Methods("POST")

	auth := r.PathPrefix("/users").Subrouter()
	auth.Use(middleware.Auth(cfg.JWTSecret))
	auth.HandleFunc("", uh.List).Methods("GET")
	auth.HandleFunc("/{id}", uh.GetByID).Methods("GET")
	auth.HandleFunc("/{id}", uh.Update).Methods("PUT")
	auth.HandleFunc("/{id}", uh.Delete).Methods("DELETE")

	port := cfg.Port
	fmt.Printf("Listening on :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
