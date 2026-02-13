package main

import (
	"log"
	"net/http"
	"time"

	"github.com/generic/selectel-craas-web/internal/api"
	"github.com/generic/selectel-craas-web/internal/auth"
	"github.com/generic/selectel-craas-web/internal/config"
	"github.com/generic/selectel-craas-web/internal/craas"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	authClient := auth.New(cfg)
	craasService := craas.New()

	router := api.New(authClient, craasService)

	srv := &http.Server{
		Addr:         ":" + cfg.WebPort,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	log.Printf("Starting server on port %s", cfg.WebPort)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
