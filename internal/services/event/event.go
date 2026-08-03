package services

import (
	"context"
	"log/slog"
	"polymarket/internal/models"
)

type StorageRepository interface {
	GetEvents(ctx context.Context) ([]models.Event, error)
}

type Service struct {
	repo   StorageRepository
	logger *slog.Logger
}

func NewEventService(repo StorageRepository, logger *slog.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

func (s *Service) GetAllEvents(ctx context.Context) ([]models.Event, error) {
	events, err := s.repo.GetEvents(ctx)
	if err != nil {
		s.logger.Error("service/event failed to fetch events", "error", err)
		return nil, err
	}

	return events, nil
}
