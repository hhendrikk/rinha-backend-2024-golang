package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

const (
	DEFAULT_DATABASE_MAX_CONNS = 4
)

type RinhaBackendConfig struct {
	TimeZone    int
	Port        int
	Environment string
	Database    struct {
		Host        string
		Port        string
		User        string
		Pass        string
		DbName      string
		Application string
		MaxConns    int
	}
}

func LoadRinhaBackendEnv() {
	env := os.Getenv("RINHA_BACKEND_ENV")

	if env == "" {
		env = "development"
	}

	_ = godotenv.Load(fmt.Sprintf(".env.%s", env))
	_ = godotenv.Load()
}

func NewRinhaBackendConfig() *RinhaBackendConfig {
	LoadRinhaBackendEnv()

	maxConns, err := strconv.Atoi(os.Getenv("RINHA_BACKEND_DB_MAX_CONNS"))
	if err != nil {
		maxConns = DEFAULT_DATABASE_MAX_CONNS
	}

	timeZone, err := strconv.Atoi(os.Getenv("RINHA_BACKEND_TIME_ZONE"))
	if err != nil {
		timeZone = 0
	}

	port, err := strconv.Atoi(os.Getenv("RINHA_BACKEND_PORT"))
	if err != nil {
		port = 8080
	}

	return &RinhaBackendConfig{
		TimeZone:    timeZone,
		Port:        port,
		Environment: os.Getenv("RINHA_BACKEND_ENV"),
		Database: struct {
			Host        string
			Port        string
			User        string
			Pass        string
			DbName      string
			Application string
			MaxConns    int
		}{
			Host:        os.Getenv("RINHA_BACKEND_DB_HOST"),
			Port:        os.Getenv("RINHA_BACKEND_DB_PORT"),
			User:        os.Getenv("RINHA_BACKEND_DB_USER"),
			Pass:        os.Getenv("RINHA_BACKEND_DB_PASS"),
			DbName:      os.Getenv("RINHA_BACKEND_DB_DBNAME"),
			Application: os.Getenv("RINHA_BACKEND_DB_APPLICATION"),
			MaxConns:    maxConns,
		},
	}
}
