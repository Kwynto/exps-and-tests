package main

import (
	"net/http"
	"os"
	"rest-url-short/internal/config"
	"rest-url-short/internal/http-server/handlers/redirect"
	"rest-url-short/internal/http-server/handlers/url/save"
	mwLogger "rest-url-short/internal/http-server/middleware/logger"
	"rest-url-short/internal/lib/helpers/loghelper"
	"rest-url-short/internal/storage/sqlite"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/exp/slog"
)

func main() {
	// Init config: cleanenv
	configPath := os.Getenv("CONFIG_PATH")
	cfg := config.MustLoad(configPath)

	// Init logger: slog
	log := loghelper.SetupLogger(cfg.Env)
	log.Info("starting rest-url-short", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	// Init storage: sqlite3
	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", loghelper.Err(err))
		os.Exit(1)
	}

	// Init router: chi, chi render
	router := chi.NewRouter()
	// middleware
	router.Use(middleware.RequestID)
	if cfg.Env == "local" { // || cfg.Env == "dev"
		router.Use(middleware.Logger)
	}
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat) // для шаблонизации ссылок, например /{alias}

	// router for "/url"
	router.Route("/url", func(r chi.Router) {
		r.Use(middleware.BasicAuth("url-shortener", map[string]string{
			cfg.HTTPServer.User: cfg.HTTPServer.Password,
		}))

		r.Post("/", save.New(log, storage))
		//TODO: add DELETE /url/{id}
	})

	// root router
	router.Get("/{alias}", redirect.New(log, storage))

	// Run server
	log.Info("starting server", slog.String("address", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	log.Error("server stopped")
}
