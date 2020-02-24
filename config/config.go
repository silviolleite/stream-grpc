package config

import (
	"github.com/joho/godotenv"
	"os"
)

type configurations struct {
	BrockerUrl           string
	BrockerPort string
	Topic        string
	GroupID string
	SqlitePath string
}

var Config configurations

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

func init() {
	_ = godotenv.Load()

	Config = configurations{
		getEnv("BROCKER_URL", "localhost"),
		getEnv("BROCKER_PORT", "9092"),
		"transactions",
		"machines",
		getEnv("SQLITE_PATH", "sqlite.db"),
	}
}

