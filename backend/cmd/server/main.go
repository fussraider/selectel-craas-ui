package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/generic/selectel-craas-web/internal/api"
	"github.com/generic/selectel-craas-web/internal/auth"
	"github.com/generic/selectel-craas-web/internal/config"
	"github.com/generic/selectel-craas-web/internal/craas"
	"github.com/generic/selectel-craas-web/pkg/logger"
)

var Version = "dev"

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	appLogger := logger.New(cfg.LogLevel, cfg.LogFormat)
	appLogger.Info("starting application", "version", Version, "port", cfg.WebPort, "log_level", cfg.LogLevel)

	// Auth validation
	if cfg.AuthEnabled {
		if cfg.AuthLogin == "" || cfg.AuthPassword == "" {
			log.Fatal("Authentication is ENABLED but AUTH_LOGIN or AUTH_PASSWORD is not set. Please set these environment variables.")
		}
		appLogger.Info("Authentication: ENABLED")
	} else {
		appLogger.Warn("Authentication: DISABLED (Anyone can access the application)")
	}

	if cfg.CORSAllowedOrigin == "*" {
		appLogger.Warn("CORS: ALLOWED_ORIGIN is set to '*' (INSECURE). Do not use this in production.")
	} else if cfg.CORSAllowedOrigin == "" {
		appLogger.Info("CORS: ALLOWED_ORIGIN is empty (CORS disabled). Requests from other origins will be blocked.")
	} else {
		appLogger.Info("CORS: ALLOWED_ORIGIN is set", "origin", cfg.CORSAllowedOrigin)
	}

	authClient := auth.New(cfg, appLogger)
	craasService := craas.New(cfg, appLogger)

	router := api.New(authClient, craasService, appLogger, cfg)

	srv := &http.Server{
		Addr:         ":" + cfg.WebPort,
		Handler:      router,
		ReadTimeout:  300 * time.Second,
		WriteTimeout: 300 * time.Second,
	}

	// Server run context
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-sig

		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, cancel := context.WithTimeout(serverCtx, 30*time.Second)
		defer cancel()

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		// Trigger graceful shutdown
		appLogger.Info("shutting down server...")
		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Fatal(err)
		}
		serverStopCtx()
	}()

	// Run the server
	appLogger.Info("server listening", "addr", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		appLogger.Error("server failed", "error", err)
		log.Fatal(err)
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()
	appLogger.Info("server exited")
}
