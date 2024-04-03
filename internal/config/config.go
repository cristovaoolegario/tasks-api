package config

import (
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

const (
	Port             = "PORT"
	DBPassword       = "DB_PASSWORD"
	DBUser           = "DB_USER"
	DBHost           = "DB_HOST"
	DBName           = "DB_NAME"
	IsLocalContainer = "LOCAL_CONTAINER"
	AuthSecret       = "JWT_SECRET"
)

type Config struct {
	DbConnection string
	AppPort      string
	AuthSecret   string
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load(".env")

	secret, err := loadAuthConfig()
	if err != nil {
		return nil, err
	}

	return &Config{
		DbConnection: loadDBConfig(),
		AppPort:      os.Getenv(Port),
		AuthSecret:   secret,
	}, nil
}

func loadDBConfig() string {
	dbPassword := os.Getenv(DBPassword)
	if dbPassword == "" {
		log.Fatal("DB_PASSWORD environment variable is required but was not set")
	}
	dbUsername := os.Getenv(DBUser)
	if dbUsername == "" {
		log.Fatal("DB_USER environment variable is required but was not set")
	}
	dbHost := os.Getenv(DBHost)
	if dbHost == "" {
		localDockerEnv := os.Getenv(IsLocalContainer)
		if localDockerEnv == "" || localDockerEnv == "false" {
			dbHost = "localhost"
		} else {
			dbHost = "host.docker.internal"
		}
	}

	dbName := os.Getenv(DBName)
	if dbName == "" {
		dbName = "tasks"
	}
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true", dbUsername, dbPassword, dbHost, dbName)
}

func loadAuthConfig() (string, error) {
	secret := os.Getenv(AuthSecret)
	if secret == "" {
		return "", errors.New("JWT_SECRET environment variable needs to be set")
	}
	return secret, nil
}
