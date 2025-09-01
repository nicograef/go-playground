package database

import (
	"testing"

	"github.com/google/uuid"
)

type User map[string]any

func equalUser(a, b User) bool {
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
	if db == nil {
		t.Fatal("Failed to create database")
	}

	db.AddEvent("user.new", "1", User{"ID": "1", "Name": "John Doe"})
	db.AddEvent("user.update", "1", User{"ID": "1", "Email": "john.doe@example.com"})

	events := db.GetEvents()
	if len(events) != 2 {
		t.Fatal("Failed to get events")
	}

	if !equalUser(events[0].Data, User{"ID": "1", "Name": "John Doe"}) {
		t.Fatal("Event retrieved is not the same as the one created")
	}

	if !equalUser(events[1].Data, User{"ID": "1", "Email": "john.doe@example.com"}) {
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
	if db == nil {
		t.Fatal("Failed to create database")
	}

	user := User{
		"ID":   "1",
		"Name": "John Doe",
	}
	db.AddEvent("user.new", "1", user)

	if nonExistingEvent := db.GetEventByID(uuid.New()); nonExistingEvent != nil {
		t.Fatal("Expected no event to be found")
	}

	// Get the first event from the map
	var firstEvent Event
	for _, event := range db.events {
		firstEvent = event
		break
	}

	event := db.GetEventByID(firstEvent.ID)
	if event == nil {
		t.Fatal("Failed to get event by ID")
	}

	if !equalUser(event.Data, user) {
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
	if db == nil {
		t.Fatal("Failed to create database")
	}

	db.AddEvent("user.new", "1", User{"ID": "1", "Name": "John Doe"})
	db.AddEvent("user.update", "1", User{"ID": "1", "Email": "john.doe@example.com"})
	db.AddEvent("user.new", "2", User{"ID": "2", "Name": "Max Mustermann"})

	if nonExistingEvents := db.GetEventsByType("non-existing-type"); nonExistingEvents != nil {
		t.Fatal("Expected no event to be found")
	}

	events := db.GetEventsByType("user.new")
	if len(events) != 2 {
		t.Fatal("Failed to get events by type")
	}

	if !equalUser(events[0].Data, User{"ID": "1", "Name": "John Doe"}) {
		t.Fatal("Event retrieved is not the same as the one created")
	}

	if !equalUser(events[1].Data, User{"ID": "2", "Name": "Max Mustermann"}) {
		t.Fatal("Event retrieved is not the same as the one created")
	}

	if events[0].Type != "user.new" {
		t.Fatal("Event type is not the same as the one created")
	}

	if events[1].Type != "user.new" {
		t.Fatal("Event type is not the same as the one created")
	}
}

func TestGetEventsByEntityID(t *testing.T) {
	db := New()
	if db == nil {
		t.Fatal("Failed to create database")
	}

	db.AddEvent("user.new", "1", User{"ID": "1", "Name": "John Doe"})
	db.AddEvent("user.update", "1", User{"ID": "1", "Email": "john.doe@example.com"})
	db.AddEvent("user.new", "2", User{"ID": "2", "Name": "Max Mustermann"})

	if nonExistingEvents := db.GetEventsByEntityID("non-existing-id"); nonExistingEvents != nil {
		t.Fatal("Expected no event to be found")
	}

	events := db.GetEventsByEntityID("1")
	if len(events) != 2 {
		t.Fatal("Failed to get events by entity ID")
	}

	if !equalUser(events[0].Data, User{"ID": "1", "Name": "John Doe"}) {
		t.Fatal("Event retrieved is not the same as the one created")
	}

	if !equalUser(events[1].Data, User{"ID": "1", "Email": "john.doe@example.com"}) {
		t.Fatal("Event retrieved is not the same as the one created")
	}

	if events[0].Type != "user.new" {
		t.Fatal("Event type is not the same as the one created")
	}

	if events[1].Type != "user.update" {
		t.Fatal("Event type is not the same as the one created")
	}
}

func TestGetEntity(t *testing.T) {
	db := New()
	if db == nil {
		t.Fatal("Failed to create database")
	}

	db.AddEvent("user.new", "1", User{"ID": "1", "Name": "John Doe"})
	db.AddEvent("user.update", "1", User{"ID": "1", "Email": "john.doe@example.com"})
	db.AddEvent("user.update", "1", User{"ID": "1", "Name": "johnny donny"})
	db.AddEvent("user.new", "2", User{"ID": "2", "Name": "Max Mustermann"})

	if entity := db.GetEntity("non-existing-id"); len(entity) != 0 {
		t.Fatal("Expected no entity to be found")
	}

	entity := db.GetEntity("1")
	if entity == nil {
		t.Fatal("Failed to get entity by ID")
	}

	if !equalUser(entity, User{"ID": "1", "Name": "johnny donny", "Email": "john.doe@example.com"}) {
		t.Fatal("Entity retrieved is not the same as the one created")
	}

}
