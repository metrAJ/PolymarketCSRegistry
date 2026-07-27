package transport

import (
	"context"
	"log/slog"
	"net/http"
)

type ScraperService interface {
	ScrapeActiveEvents(ctx context.Context) error
}

type ScraperHandler struct {
	service ScraperService
	logger  *slog.Logger
}

func NewScraperHandler(service ScraperService, logger *slog.Logger) *ScraperHandler {
	return &ScraperHandler{
		service: service,
		logger:  logger,
	}
}

func (h *ScraperHandler) ScrapeCSEvents(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	err := h.service.ScrapeActiveEvents(ctx)
	if err != nil {
		h.logger.Error("service/scraper/transport failed to scrape events from service", "error", http.StatusInternalServerError)
		http.Error(w, "Could not scrape events", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "success", "message": "events successfully scraped"}`))
}
