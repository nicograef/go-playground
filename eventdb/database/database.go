package database

import (
	"sort"
	"time"

	"github.com/google/uuid"
)

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

func (db *Database) AddEvent(eventType string, entityID string, data map[string]any) uuid.UUID {
	event := NewEvent(eventType, entityID, data)

	db.events[event.ID] = event
	db.typeIndex[eventType] = append(db.typeIndex[eventType], event.ID)
	db.entityIndex[entityID] = append(db.entityIndex[entityID], event.ID)

	return event.ID
}

func (db *Database) GetEvent(id uuid.UUID) *Event {
	event, exists := db.events[id]

	if !exists {
		return nil
	}

	return &event
}

func (db *Database) DeleteEvent(id uuid.UUID) {
	if event, exists := db.events[id]; exists {
		delete(db.events, id)
		db.typeIndex[event.Type] = removeIDFromSlice(db.typeIndex[event.Type], id)
		db.entityIndex[event.EntityID] = removeIDFromSlice(db.entityIndex[event.EntityID], id)
	}
}

func (db *Database) GetEvents() []Event {
	events := make([]Event, 0, len(db.events))

	for _, event := range db.events {
		events = append(events, event)
	}

	sortEventsByTimestamp(events)

	return events
}

func (db *Database) GetEventsByType(eventType string) []Event {
	eventIDs, exists := db.typeIndex[eventType]

	if !exists || len(eventIDs) == 0 {
		return []Event{}
	}

	events := make([]Event, 0, len(eventIDs))
	for _, id := range eventIDs {
		if event, exists := db.events[id]; exists {
			events = append(events, event)
		}
	}

	sortEventsByTimestamp(events)

	return events

}

// GetEventsByEntity returns all events for a specific entity sorted by their timestamp
func (db *Database) GetEventsByEntity(entityID string) []Event {
	eventIDs, exists := db.entityIndex[entityID]

	if !exists || len(eventIDs) == 0 {
		return []Event{}
	}

	events := make([]Event, 0, len(eventIDs))
	for _, id := range eventIDs {
		if event, exists := db.events[id]; exists {
			events = append(events, event)
		}
	}

	sortEventsByTimestamp(events)

	return events
}

// removeIDFromSlice removes an ID from a slice of UUIDs
func removeIDFromSlice(originalSlice []uuid.UUID, id uuid.UUID) []uuid.UUID {
	result := make([]uuid.UUID, 0, len(originalSlice))

	for _, currentId := range originalSlice {
		if currentId == id {
			continue
		} else {
			result = append(result, currentId)
		}
	}

	return result
}

// sortEventsByTimestamp sorts events by their timestamp
func sortEventsByTimestamp(events []Event) {
	sort.Slice(events, func(i, j int) bool {
		return events[i].Timestamp.Before(events[j].Timestamp)
	})
}
