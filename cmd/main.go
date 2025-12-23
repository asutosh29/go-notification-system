package main

import (
	"fmt"

	"github.com/asutosh29/go-gin/internal/api"
	"github.com/asutosh29/go-gin/internal/config"
	"github.com/asutosh29/go-gin/internal/database"
	"github.com/asutosh29/go-gin/internal/env"
)

func main() {
	// Configs
	config := config.InitConfig()

	// DB Setup
	database.InitDatabase(config)
	database.Migrate()

	// Web server Setup
	PORT := env.GetEnvString("SERVER_PORT", "8080")
	router := api.InitRouter()

	router.Run(fmt.Sprintf("localhost:%s", PORT))
}
