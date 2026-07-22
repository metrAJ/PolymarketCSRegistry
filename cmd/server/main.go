package main

import (
	"log"
	"net/http"
	"polymarket/internal/config"

	"go.uber.org/zap"
)

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}

	if cfg.Port == "" {
		cfg.Port = "3000"
	}
	logger.Info("Starting server on port:" + cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, nil))
}
