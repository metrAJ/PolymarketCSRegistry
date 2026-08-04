package transport

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"polymarket/internal/models"
	"time"
)

type EventResponse struct {
	Title   string           `json:"title"`
	EndDate time.Time        `json:"endDate"`
	Tags    []string         `json:"tags"`
	Markets []MarketResponse `json:"markets"`
}

type MarketResponse struct {
	Question     string            `json:"question"`
	VolumeNum    float64           `json:"volumeNum"`
	LiquidityNum float64           `json:"liquidityNum"`
	Outcomes     []OutcomeResponse `json:"outcomes"`
}

type OutcomeResponse struct {
	Outcome       string  `json:"outcome"`
	OutcomePrices float64 `json:"outcomePrices"`
}

func toEventResponceDTO(e models.Event) EventResponse {
	var markets []MarketResponse

	for _, m := range e.Markets {
		var outcomes []OutcomeResponse
		for _, o := range m.Outcomes {
			outcomes = append(outcomes, OutcomeResponse{
				Outcome:       o.Outcome,
				OutcomePrices: o.OutcomePrices,
			})
		}

		markets = append(markets, MarketResponse{
			Question:     m.Question,
			VolumeNum:    m.VolumeNum,
			LiquidityNum: m.LiquidityNum,
			Outcomes:     outcomes,
		})
	}

	return EventResponse{
		Title:   e.Title,
		EndDate: e.EndDate,
		Tags:    e.Tags,
		Markets: markets,
	}
}

type Service interface {
	GetAllEvents(ctx context.Context) ([]models.Event, error)
}

type EventHandler struct {
	service Service
	logger  *slog.Logger
}

func NewEventHandler(service Service, logger *slog.Logger) *EventHandler {
	return &EventHandler{
		service: service,
		logger:  logger,
	}
}

func (h *EventHandler) GetAllEvents(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var responseDTOs []EventResponse

	events, err := h.service.GetAllEvents(ctx)
	if err != nil {
		http.Error(w, "Could not get events", http.StatusInternalServerError)
		h.logger.Error("service/event/transport failed to fetch events from service", "error", err)
		return
	}

	for _, event := range events {
		responseDTOs = append(responseDTOs, toEventResponceDTO(event))
	}
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(responseDTOs); err != nil {
		h.logger.Error("service/event/transport failed to encode events in response", "error", err)
	}
}
