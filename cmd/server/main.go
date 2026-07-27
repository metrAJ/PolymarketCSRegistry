package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"polymarket/internal/config"
	"polymarket/internal/data"
	event_service "polymarket/internal/services/event"
	event_handler "polymarket/internal/services/event/transport"
	scraper_service "polymarket/internal/services/scraper"
	scraper_handler "polymarket/internal/services/scraper/transport"
	gamma "polymarket/pkg/gammaapi"
)

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	if cfg.Port == "" {
		cfg.Port = "3000"
	}

	client := gamma.NewClient(logger)
	storage := data.NewStorage()
	storageRepo := data.NewStorageRepository(storage)
	scraperService := scraper_service.NewScraperService(storageRepo, client)
	eventService := event_service.NewEventService(storageRepo, logger)
	eventHandler := event_handler.NewEventHandler(eventService, logger)
	scraperHandler := scraper_handler.NewScraperHandler(scraperService, logger)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/events", eventHandler.GetAllEvents)
	mux.HandleFunc("POST /api/scrape", scraperHandler.ScrapeCSEvents)

	logger.Info("Starting server", "port", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, mux); err != nil {
		logger.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}
