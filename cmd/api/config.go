package main

import (
	"os"

	"github.com/joho/godotenv"
)

type DatabaseConfig struct {
	dsn      string
	dbDriver string
}

type Config struct {
	port    string
	database DatabaseConfig
}

func LoadConfig() (*Config, error){
	err := godotenv.Load()
	if err != nil {	
		return nil, err
	}

	return &Config{
		port: os.Getenv("PORT"),
		database: DatabaseConfig{
			dsn: os.Getenv("DATA_SOURCE_NAME"),
			dbDriver: os.Getenv("DATABASE_DRIVER"),
		},
	}, nil
}