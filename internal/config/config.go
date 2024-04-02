package config

import (
	"fmt"
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
)

type Config struct {
	DbConnection string
	AppPort      string
}

func LoadConfig() *Config {
	return &Config{
		DbConnection: loadDBConfig(),
		AppPort:      os.Getenv(Port),
	}
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
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", dbUsername, dbPassword, dbHost, dbName)
}
