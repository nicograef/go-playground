package database

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID        uuid.UUID      `json:"id"`
	Timestamp time.Time      `json:"timestamp"`
	Type      string         `json:"type"`
	Data      map[string]any `json:"data"`
}

type Database struct {
	created   time.Time
	events    map[uuid.UUID]Event
	typeIndex map[string][]uuid.UUID
}

func New() *Database {
	return &Database{
		events:    make(map[uuid.UUID]Event),
		typeIndex: make(map[string][]uuid.UUID),
	}
}

func (db *Database) AddEvent(eventType string, data map[string]any) {
	event := Event{
		ID:        uuid.New(),
		Timestamp: time.Now(),
		Type:      eventType,
		Data:      data,
	}

	db.events[event.ID] = event
	db.typeIndex[eventType] = append(db.typeIndex[eventType], event.ID)
}

func (db *Database) GetEvents() []Event {
	events := make([]Event, 0, len(db.events))
	for _, event := range db.events {
		events = append(events, event)
	}

	return events
}

func (db *Database) GetEventByID(id uuid.UUID) *Event {
	if event, exists := db.events[id]; exists {
		return &event
	}

	return nil
}

func (db *Database) GetEventsByType(eventType string) []Event {
	if eventIDs, exists := db.typeIndex[eventType]; exists {
		events := make([]Event, 0, len(eventIDs))
		for _, id := range eventIDs {
			if event, exists := db.events[id]; exists {
				events = append(events, event)
			}
		}

		return events
	}

	return nil
}
