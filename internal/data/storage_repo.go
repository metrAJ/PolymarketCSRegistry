package data

import (
	"context"
	"polymarket/internal/models"
	"time"
)

type StorageRepository struct {
	storage *Storage
}

func NewStorageRepository(storage *Storage) *StorageRepository {
	return &StorageRepository{
		storage: storage,
	}
}

func (r *StorageRepository) SaveEvents(_ context.Context, events []models.Event) error {
	newEvents := make(map[string]models.Event, len(events))
	for _, event := range events {
		newEvents[event.ID] = event
	}
	r.storage.mu.Lock()
	defer r.storage.mu.Unlock()
	r.storage.idMap = newEvents
	r.storage.timeStamp = time.Now().UTC()
	return nil
}

func (r *StorageRepository) GetEvents(_ context.Context) ([]models.Event, error) {
	r.storage.mu.RLock()
	defer r.storage.mu.RUnlock()
	events := make([]models.Event, 0, len(r.storage.idMap))
	for _, event := range r.storage.idMap {
		events = append(events, event)
	}
	return events, nil
}
