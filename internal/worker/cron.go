package worker

import (
	"context"
	"log/slog"
	"time"
)

type ScraperService interface {
	ScrapeActiveEvents(ctx context.Context) error
}

type ScraperWorker struct {
	service ScraperService
	logger  *slog.Logger
}

func NewScraperWorker(service ScraperService, logger *slog.Logger) *ScraperWorker {
	return &ScraperWorker{
		service: service,
		logger:  logger,
	}
}

func (w *ScraperWorker) scrape(ctx context.Context) {
	if err := w.service.ScrapeActiveEvents(ctx); err != nil {
		w.logger.Error("worker/cron failed to scrape active events", "error", err)
		return
	}

	w.logger.Info("worker/cron successfully updated events")
}

func (w *ScraperWorker) Start(ctx context.Context, interval time.Duration) {
	w.logger.Info("woker/cron started", "interval", interval)
	w.scrape(ctx)

	ticker := time.NewTicker(interval)

	go func() {
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C: // Another scrape
				w.scrape(ctx)

				return
			case <-ctx.Done(): // Stop by context
				ticker.Stop()
				w.logger.Info("worker/cron scraper cron ctx canceled")

				return
			}
		}
	}()
}
