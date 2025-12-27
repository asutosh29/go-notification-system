package env

import (
	"os"

	"github.com/joho/godotenv"
)

func InitEnv() {
	_ = godotenv.Load() // While pushing prod env is injected from docker compose!
}

func GetEnvString(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
