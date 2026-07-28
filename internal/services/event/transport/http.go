package transport

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"polymarket/internal/models"
)

type EventService interface {
	GetAllEvents(ctx context.Context) ([]models.Event, error)
}

type EventHandler struct {
	service EventService
	logger  *slog.Logger
}

func NewEventHandler(service EventService, logger *slog.Logger) *EventHandler {
	return &EventHandler{
		service: service,
		logger:  logger,
	}
}

func (h *EventHandler) GetAllEvents(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	events, err := h.service.GetAllEvents(ctx)
	if err != nil {
		http.Error(w, "Could not get events", http.StatusInternalServerError)
		h.logger.Error("service/event/transport failed to fetch events from service", "error", err)

		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(events); err != nil {
		h.logger.Error("service/event/transport failed to encode events in response", "error", err)
	}
}
