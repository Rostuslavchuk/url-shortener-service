package main

import (
	"log/slog"
	"net/http"
	"os"

	cfg "url_shortener/internal/config"
	del "url_shortener/internal/http-server/handlers/delete"
	"url_shortener/internal/http-server/handlers/redirect"
	"url_shortener/internal/http-server/handlers/save"
	"url_shortener/internal/lib/sl"
	"url_shortener/internal/storage/postgres"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	if err := godotenv.Load(); err != nil {
		slog.Error("No .env file found", sl.Err(err))
		return
	}

	config := cfg.MustLoad()

	log := loggerSetup(config.Env)
	log = log.With(slog.String("env", config.Env))

	db, err := postgres.New(config.StoragePath)
	if err != nil {
		log.Error("faild to init storage", sl.Err(err))
		return
	}
	_ = db

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Route("/url", func(r chi.Router) {
		r.Use(middleware.BasicAuth("url_shortener", map[string]string{
			config.User: config.Password,
		}))

		r.Post("/", save.New(log, db))
		r.Delete("/{alias}", del.New(log, db))
	})
	router.Get("/{alias}", redirect.New(log, db))

	log.Info("server starting", slog.String("address", config.Address))

	srv := &http.Server{
		Addr:         config.Address,
		Handler:      router,
		ReadTimeout:  config.Timeout,
		WriteTimeout: config.Timeout,
		IdleTimeout:  config.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("faild to start server", sl.Err(err))
	}
	log.Error("server stopped")
}

func loggerSetup(env string) *slog.Logger {
	// local - text log
	// dev, prod - json log (dev - debug, prod - info)
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
