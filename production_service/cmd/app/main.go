package main

import (
	"log"

	"github.com/Kwynto/production_service/internal/app"
	"github.com/Kwynto/production_service/internal/config"
	"github.com/Kwynto/production_service/pkg/logging"
)

func main() {
	log.Print("Config initializing.")
	cfg := config.GetConfig()

	log.Print("Logger initializing")
	logger := logging.GetLogger()

	app, err := app.NewApp(cfg, &logger)
	if err != nil {
		logger.Fatal(err)
	}
}
