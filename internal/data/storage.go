package data

import (
	"polymarket/internal/models"
	"sync"
	"time"
)

type Storage struct {
	timeStamp time.Time
	mu        sync.RWMutex
	idMap     map[string]models.Event
}

func NewStorage() *Storage {
	return &Storage{
		timeStamp: time.Now().UTC(),
		idMap:     make(map[string]models.Event),
	}
}
