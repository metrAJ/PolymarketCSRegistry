package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"polymarket/internal/config"
	"polymarket/internal/data"
	scraper "polymarket/internal/services/scraper"
	gamma "polymarket/pkg/gammaapi"

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
	// logger.Info("Starting server on port:" + cfg.Port)
	// log.Fatal(http.ListenAndServe(":"+cfg.Port, nil))

	client := gamma.NewClient(logger)                 // GammaAPIClient
	storage := data.NewStorage()                      // Storage
	storageRepo := data.NewStorageRepository(storage) // StorageRepository
	scraperService := scraper.NewScraperService(storageRepo, client)

	scraperCtx := context.Background()

	// Raw functions use to check before making http endpoints

	err = scraperService.ScrapeActiveEvents(scraperCtx)
	if err != nil {
		logger.Fatal("scraper service failed", zap.Error(err))
	}

	events, err := storageRepo.GetEvents(scraperCtx)
	if err != nil {
		logger.Fatal("failed to read from repository", zap.Error(err))
	}
	if len(events) > 0 {
		eventJSON, marshalErr := json.MarshalIndent(events[0], "", "  ")
		if marshalErr != nil {
			logger.Error("failed to marshal event for debugging", zap.Error(marshalErr))
		} else {
			fmt.Println(string(eventJSON))
		}
	} else {
		logger.Info("no active events found to print.")
	}
}
