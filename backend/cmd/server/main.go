package main

import (
	"log"
	"net/http"
	"time"

	"github.com/generic/selectel-craas-web/internal/api"
	"github.com/generic/selectel-craas-web/internal/auth"
	"github.com/generic/selectel-craas-web/internal/config"
	"github.com/generic/selectel-craas-web/internal/craas"
	"github.com/generic/selectel-craas-web/pkg/logger"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	appLogger := logger.New(cfg.LogLevel, cfg.LogFormat)
	appLogger.Info("starting application", "port", cfg.WebPort, "log_level", cfg.LogLevel)

	authClient := auth.New(cfg, appLogger)
	craasService := craas.New(appLogger)

	router := api.New(authClient, craasService, appLogger)

	srv := &http.Server{
		Addr:         ":" + cfg.WebPort,
		Handler:      router,
		ReadTimeout:  300 * time.Second,
		WriteTimeout: 300 * time.Second,
	}

	appLogger.Info("server listening", "addr", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		appLogger.Error("server failed", "error", err)
		log.Fatal(err)
	}
}
