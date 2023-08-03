package app

import (
	"github.com/Kwynto/production_service/internal/config"
	"github.com/Kwynto/production_service/pkg/logging"
)

type App struct {
	config *config.Config
	logger *logging.Logger
}

func NewApp(config *config.Config, logger *logging.Logger) (App, error) {
	return App{
		config,
		logger,
	}, nil
}
