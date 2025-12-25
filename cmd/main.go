package main

import (
	"fmt"

	"github.com/asutosh29/go-gin/internal/api"
	"github.com/asutosh29/go-gin/internal/config"
	"github.com/asutosh29/go-gin/internal/database"
	"github.com/asutosh29/go-gin/internal/env"
	"github.com/asutosh29/go-gin/internal/hub"
)

func main() {
	// Configs
	config := config.InitConfig()
	// Notification hub
	hub := hub.NewHub()
	go hub.Listen()
	// DB Setup
	database.InitDatabase(config)
	database.Migrate()

	// Web server Setup
	PORT := env.GetEnvString("SERVER_PORT", "8080")
	router := api.InitRouter(hub)

	router.Run(fmt.Sprintf("localhost:%s", PORT))
}
