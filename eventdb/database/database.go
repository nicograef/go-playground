package database

import (
	"sort"
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

type Database struct {
	created     time.Time
	events      map[uuid.UUID]Event
	typeIndex   map[string][]uuid.UUID
	entityIndex map[string][]uuid.UUID
}

func New() *Database {
	return &Database{
		created:     time.Now().UTC(),
		events:      make(map[uuid.UUID]Event),
		typeIndex:   make(map[string][]uuid.UUID),
		entityIndex: make(map[string][]uuid.UUID),
	}
}

func (db *Database) AddEvent(eventType string, entityID string, data map[string]any) {
	event := Event{
		ID:        uuid.New(),
		Timestamp: time.Now(),
		Type:      eventType,
		EntityID:  entityID,
		Data:      data,
	}

	db.events[event.ID] = event
	db.typeIndex[eventType] = append(db.typeIndex[eventType], event.ID)
	db.entityIndex[entityID] = append(db.entityIndex[entityID], event.ID)
}

func (db *Database) GetEvents() []Event {
	events := make([]Event, 0, len(db.events))
	for _, event := range db.events {
		events = append(events, event)
	}

	sort.Slice(events, func(i, j int) bool {
		return events[i].Timestamp.Before(events[j].Timestamp)
	})

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

		sort.Slice(events, func(i, j int) bool {
			return events[i].Timestamp.Before(events[j].Timestamp)
		})

		return events
	}

	return nil
}

// GetEventsByEntityID returns all events for a specific entity sorted by their timestamp
func (db *Database) GetEventsByEntityID(entityID string) []Event {
	if eventIDs, exists := db.entityIndex[entityID]; exists {
		events := make([]Event, 0, len(eventIDs))
		for _, id := range eventIDs {
			if event, exists := db.events[id]; exists {
				events = append(events, event)
			}
		}

		sort.Slice(events, func(i, j int) bool {
			return events[i].Timestamp.Before(events[j].Timestamp)
		})

		return events
	}

	return nil
}

// GetEntity retrieves all events for one entity
// and builds the aggregate entity by applying all events sorted by their timestamp
// to the initial state.
func (db *Database) GetEntity(entityID string) map[string]any {
	events := db.GetEventsByEntityID(entityID)
	aggregate := make(map[string]any)

	for _, event := range events {
		for k, v := range event.Data {
			aggregate[k] = v
		}
	}

	return aggregate
}
