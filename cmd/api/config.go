package main

import (
	"os"
	"github.com/joho/godotenv"
)

type Config struct {
	database databaseConfig
	server   serverConfig
}

type databaseConfig struct {
	dsn      string
	dbDriver string
}

type serverConfig struct {
	port string
}

func LoadConfig() (*Config, error) {
	godotenv.Load()

	dsn := os.Getenv("DATA_SOURCE_NAME")
	if dsn == "" {
		dsn = "appuser:apppass@tcp(db:3306)/urldb?parseTime=true"
	}

	driver := os.Getenv("DATABASE_DRIVER")
	if driver == "" {
		driver = "mysql"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = ":4000"
	}

	return &Config{
		database: databaseConfig{
			dsn:      dsn,
			dbDriver: driver,
		},
		server: serverConfig{
			port: port,
		},
	}, nil
}