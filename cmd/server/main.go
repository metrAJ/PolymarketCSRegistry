package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"polymarket/internal/config"
	"polymarket/internal/data"
	event_service "polymarket/internal/services/event"
	event_handler "polymarket/internal/services/event/transport"
	scraper_service "polymarket/internal/services/scraper"
	scraper_handler "polymarket/internal/services/scraper/transport"
	"polymarket/internal/worker"
	gamma "polymarket/pkg/gammaapi"
	"syscall"
	"time"
)

const ScraperIntervalSec = 15

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	client := gamma.NewClient(logger)
	storage := data.NewStorage()
	storageRepo := data.NewStorageRepository(storage)
	scraperService := scraper_service.NewService(storageRepo, client)
	eventService := event_service.NewEventService(storageRepo)
	eventHandler := event_handler.NewEventHandler(eventService, logger)
	scraperHandler := scraper_handler.NewScraperHandler(scraperService, logger)

	cronJob := worker.NewScraperWorker(scraperService, logger)

	cronJob.Start(ctx, ScraperIntervalSec*time.Second)
	defer cronJob.Stop()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/events", eventHandler.GetAllEvents)
	mux.HandleFunc("POST /api/scrape", scraperHandler.ScrapeCSEvents)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: mux,
	}

	go func() {
		logger.Info("Starting server", "port", cfg.Port)

		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("Server failed to start", "error", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()

	stop()
	logger.Info("Shutting down gracefully,Ctrl + C to force")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("Server forced to shutdown", "error", err)
	}

	logger.Info("Server exiting")
}
