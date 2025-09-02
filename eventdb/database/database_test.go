package database

import (
	"testing"

	"github.com/google/uuid"
)

type user map[string]any

func equalUser(a, b user) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if bv, ok := b[k]; !ok || bv != v {
			return false
		}
	}
	return true
}

func TestDatabaseCreation(t *testing.T) {
	db := New()

	if db == nil {
		t.Fatal("Failed to create database")
	}
	if db.events == nil {
		t.Fatal("Failed to create events")
	}
	if len(db.events) != 0 {
		t.Fatal("Failed to create events")
	}
}

func TestAddAndGetEvents(t *testing.T) {
	db := New()
	id1 := db.AddEvent("user.new", "1", user{"ID": "1", "Name": "John Doe"})
	id2 := db.AddEvent("user.update", "1", user{"ID": "1", "Email": "john.doe@example.com"})

	if len(db.events) != 2 {
		t.Fatal("Failed to add events")
	}
	if len(db.typeIndex) != 2 {
		t.Fatal("Failed to create type index")
	}
	if len(db.entityIndex) != 1 {
		t.Fatal("Failed to create entity index")
	}

	events := db.GetEvents()
	if len(events) != 2 {
		t.Fatal("Failed to get events")
	}

	if id1 != events[0].ID {
		t.Fatal("Event ID is not the same as the one created")
	}
	if id2 != events[1].ID {
		t.Fatal("Event ID is not the same as the one created")
	}

	if !equalUser(events[0].Data, user{"ID": "1", "Name": "John Doe"}) {
		t.Fatal("Event retrieved is not the same as the one created")
	}
	if !equalUser(events[1].Data, user{"ID": "1", "Email": "john.doe@example.com"}) {
		t.Fatal("Event retrieved is not the same as the one created")
	}

	if events[0].Type != "user.new" {
		t.Fatal("Event type is not the same as the one created")
	}
	if events[1].Type != "user.update" {
		t.Fatal("Event type is not the same as the one created")
	}
}

func TestGetEventByID(t *testing.T) {
	db := New()
	eventId := db.AddEvent("user.new", "1", user{"ID": "1", "Name": "John Doe"})

	if nonExistingEvent := db.GetEvent(uuid.New()); nonExistingEvent != nil {
		t.Fatal("Expected no event to be found")
	}

	event := db.GetEvent(eventId)
	if event == nil {
		t.Fatal("Failed to get event by ID")
	}

	if !equalUser(event.Data, user{"ID": "1", "Name": "John Doe"}) {
		t.Fatal("Event retrieved is not the same as the one created")
	}
	if event.Type != "user.new" {
		t.Fatal("Event type is not the same as the one created")
	}
	if err := uuid.Validate(event.ID.String()); err != nil {
		t.Fatal("Event ID is not valid")
	}
}

func TestGetEventsByType(t *testing.T) {
	db := New()
	db.AddEvent("user.new", "1", user{"ID": "1", "Name": "John Doe"})
	db.AddEvent("user.update", "1", user{"ID": "1", "Email": "john.doe@example.com"})
	db.AddEvent("user.new", "2", user{"ID": "2", "Name": "Max Mustermann"})

	if nonExistingEvents := db.GetEventsByType("non-existing-type"); len(nonExistingEvents) != 0 {
		t.Fatal("Expected no event to be found")
	}

	events := db.GetEventsByType("user.new")
	if len(events) != 2 {
		t.Fatal("Failed to get events by type")
	}

	if events[0].Type != "user.new" {
		t.Fatal("Event type is not the same as the one created")
	}
	if events[1].Type != "user.new" {
		t.Fatal("Event type is not the same as the one created")
	}

	if !equalUser(events[0].Data, user{"ID": "1", "Name": "John Doe"}) {
		t.Fatal("Event retrieved is not the same as the one created")
	}
	if !equalUser(events[1].Data, user{"ID": "2", "Name": "Max Mustermann"}) {
		t.Fatal("Event retrieved is not the same as the one created")
	}
}

func TestGetEventsByEntityID(t *testing.T) {
	db := New()
	db.AddEvent("user.new", "1", user{"ID": "1", "Name": "John Doe"})
	db.AddEvent("user.update", "1", user{"ID": "1", "Email": "john.doe@example.com"})
	db.AddEvent("user.new", "2", user{"ID": "2", "Name": "Max Mustermann"})

	if nonExistingEvents := db.GetEventsByEntity("non-existing-id"); len(nonExistingEvents) != 0 {
		t.Fatal("Expected no event to be found")
	}

	events := db.GetEventsByEntity("1")
	if len(events) != 2 {
		t.Fatal("Failed to get events by entity ID")
	}

	if events[0].EntityID != "1" {
		t.Fatal("Event entity ID is not the same as the one created")
	}
	if events[1].EntityID != "1" {
		t.Fatal("Event entity ID is not the same as the one created")
	}

	if !equalUser(events[0].Data, user{"ID": "1", "Name": "John Doe"}) {
		t.Fatal("Event retrieved is not the same as the one created")
	}
	if !equalUser(events[1].Data, user{"ID": "1", "Email": "john.doe@example.com"}) {
		t.Fatal("Event retrieved is not the same as the one created")
	}
}

func TestDeleteEvent(t *testing.T) {
	db := New()

	// adding 3 events
	db.AddEvent("user.new", "1", user{"ID": "1", "Name": "John Doe"})
	eventId := db.AddEvent("user.update", "1", user{"ID": "1", "Email": "john.doe@example.com"})
	db.AddEvent("user.new", "2", user{"ID": "2", "Name": "Max Mustermann"})

	if len(db.events) != 3 {
		t.Fatal("Expected 3 events to be present")
	}
	if updateEvents := db.GetEventsByType("user.update"); len(updateEvents) != 1 {
		t.Fatal("Expected 1 update event to be present")
	}
	if entityEvents := db.GetEventsByEntity("1"); len(entityEvents) != 2 {
		t.Fatal("Expected 2 entity events to be present")
	}

	// deleting the second event
	db.DeleteEvent(eventId)

	if len(db.events) != 2 {
		t.Fatal("Expected 2 events to be present")
	}
	if event := db.GetEvent(eventId); event != nil {
		t.Fatal("Expected event to be deleted")
	}
	if updateEvents := db.GetEventsByType("user.update"); len(updateEvents) != 0 {
		t.Fatal("Expected no update events to be present")
	}
	if entityEvents := db.GetEventsByEntity("1"); len(entityEvents) != 1 {
		t.Fatal("Expected only 1 entity event to be present")
	}
	if len(db.typeIndex["user.update"]) != 0 {
		t.Fatal("Expected no update events to be present")
	}
	if len(db.entityIndex["1"]) != 1 {
		t.Fatal("Expected only 1 entity event to be present")
	}
}
