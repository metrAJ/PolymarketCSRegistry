package scraper

import (
	"context"
	"fmt"
	"polymarket/internal/models"
	"time"
)

type StorageRepository interface {
	SaveEvents(ctx context.Context, events []models.Event) error
}

type GammaAPIClient interface {
	GetAllCSEvents(ctx context.Context) ([]models.Event, error)
}

type ScraperService struct {
	repo StorageRepository
	api  GammaAPIClient
}

func NewService(repo StorageRepository, api GammaAPIClient) *ScraperService {
	return &ScraperService{
		repo: repo,
		api:  api,
	}
}

func (s *ScraperService) ScrapeActiveEvents(ctx context.Context) error {
	events, err := s.api.GetAllCSEvents(ctx)
	if err != nil {
		return err
	}

	activeEvents := s.filterActiveEvents(ctx, events)

	if err := s.repo.SaveEvents(ctx, activeEvents); err != nil {
		return fmt.Errorf("service/scraper: failed to save active events: %w", err)
	}

	return nil
}

func (s *ScraperService) filterActiveEvents(_ context.Context, events []models.Event) []models.Event {
	activeEvents := make([]models.Event, 0, len(events))

	for _, event := range events {
		activeMarkets := make([]models.Market, 0, len(event.Markets))

		for _, market := range event.Markets {
			if market.AcceptingOrders && market.EndDate.After(time.Now().UTC()) && market.UmaResolutionStatus == "" {
				activeMarkets = append(activeMarkets, market)
			}
		}

		if len(activeMarkets) > 0 {
			event.Markets = activeMarkets
			activeEvents = append(activeEvents, event)
		}
	}

	return activeEvents
}
