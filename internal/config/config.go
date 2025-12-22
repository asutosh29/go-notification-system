package config

import (
	"log"
	"strconv"

	"github.com/asutosh29/go-gin/internal/env"
)

type DbConfig struct {
	DbName     string
	DbUser     string
	DbPassword string
	DbHost     string
	DbPort     int
}

type Config struct {
	DbConfig
}

func InitConfig() Config {
	env.InitEnv() // Loads env before use

	DbPort, err := strconv.Atoi(env.GetEnvString("DB_PORT", "5432"))
	if err != nil {
		log.Fatal("Error parsing port: ", err)
	}

	return Config{
		DbConfig: DbConfig{
			DbName:     env.GetEnvString("DB_NAME", "notification"),
			DbUser:     env.GetEnvString("DB_USER", "amx"),
			DbPassword: env.GetEnvString("DB_PASSWORD", "amx"),
			DbHost:     env.GetEnvString("DB_HOST", "localhost"),
			DbPort:     DbPort,
		},
	}
}
