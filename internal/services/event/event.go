package services

import (
	"context"
	"polymarket/internal/models"
)

type StorageRepository interface {
	GetEvents(ctx context.Context) ([]models.Event, error)
}

type Service struct {
	repo StorageRepository
}

func NewEventService(repo StorageRepository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetAllEvents(ctx context.Context) ([]models.Event, error) {
	events, err := s.repo.GetEvents(ctx)
	if err != nil {
		return nil, err
	}

	return events, nil
}
