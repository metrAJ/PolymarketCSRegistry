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
	service  ScraperService
	logger   *slog.Logger
	interval time.Duration
	stop     chan struct{}
}

func NewScraperWorker(service ScraperService, logger *slog.Logger, interval time.Duration) *ScraperWorker {
	return &ScraperWorker{
		service:  service,
		logger:   logger,
		interval: interval,
		stop:     make(chan struct{}),
	}
}

func (w *ScraperWorker) scrape(ctx context.Context) {
	if err := w.service.ScrapeActiveEvents(ctx); err != nil {
		w.logger.Error("worker/cron failed to scrape active events", "error", err)
	} else {
		w.logger.Info("worker/cron successfully updated events")
	}
}

func (w *ScraperWorker) Start(ctx context.Context) {
	w.logger.Info("woker/cron started", "interval", w.interval)
	w.scrape(ctx)

	ticker := time.NewTicker(w.interval)

	go func() {
		for {
			select {
			case <-ticker.C: // Another scrape
				w.scrape(ctx)
			case <-w.stop: // Stop by receiving signal from closing chan
				ticker.Stop()
				w.logger.Info("worker/cron scraper cron stopped")

				return
			case <-ctx.Done(): // Stop by context
				ticker.Stop()
				w.logger.Info("worker/cron scraper cron ctx canceled")

				return
			}
		}
	}()
}

func (w *ScraperWorker) Stop() {
	close(w.stop)
}
