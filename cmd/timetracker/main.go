package main

import (
	"log/slog"
	"net/http"
	"os"
	"timetracker_effective_mobile/internal/config"
	mwLogger "timetracker_effective_mobile/internal/handlers/mlwr/logger"
	"timetracker_effective_mobile/internal/handlers/urlpath/adduser"
	"timetracker_effective_mobile/internal/handlers/urlpath/deleteuser"
	"timetracker_effective_mobile/internal/lib/logger/sl"
	sqls "timetracker_effective_mobile/internal/storage/sqls"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	cfg := config.MustLoad()
	log := setupSlogger(cfg.Env)

	storage, err := sqls.New(cfg.ConnStr)
	if err != nil {
		log.Error("can't init storage", sl.Err(err))
	}

	_ = storage

	log.Info("start program", slog.String("env", cfg.Env))

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/user", adduser.New(log, storage))
	router.Post("/deleteuser", deleteuser.New(log, storage))

	server := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}
	log.Error("server stopped")
}

func setupSlogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case "local":
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case "prod":
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn}),
		)
	}
	return log
}
