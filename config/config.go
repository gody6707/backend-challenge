package config

import (
	"log"
	"os"
)

type Config struct {
	MongoURI  string
	DBName    string
	JWTSecret string
	Port      string
}

func Load() *Config {
	uri := os.Getenv("MONGO_URI")
	db := os.Getenv("DB_NAME")
	secret := os.Getenv("JWT_SECRET")
	port := os.Getenv("PORT")
	if uri == "" || db == "" || secret == "" {
		log.Fatal("MONGO_URI, DB_NAME, and JWT_SECRET must be set")
	}
	if port == "" {
		port = "8080"
	}
	return &Config{MongoURI: uri, DBName: db, JWTSecret: secret, Port: port}
}
