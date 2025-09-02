package database

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID        uuid.UUID      `json:"id"`
	Timestamp time.Time      `json:"timestamp"`
	Type      string         `json:"type"`
	EntityID  string         `json:"entity_id"`
	Data      map[string]any `json:"data"`
}

func NewEvent(eventType string, entityID string, data map[string]any) Event {
	return Event{
		ID:        uuid.New(),
		Timestamp: time.Now(),
		Type:      eventType,
		EntityID:  entityID,
		Data:      data,
	}
}
