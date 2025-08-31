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

func TestAddEvent(t *testing.T) {
	db := New()
	if db == nil {
		t.Fatal("Failed to create database")
	}

	data := User{
		"Name": "John Doe",
	}
	db.AddEvent("user.new", data)

	// Get the first event from the map
	var firstEvent Event
	for _, event := range db.events {
		firstEvent = event
		break
	}

	if !equalUser(firstEvent.Data, data) {
		t.Fatal("Event added is not the same as the one created")
	}

	if firstEvent.Type != "user.new" {
		t.Fatal("Event type is not the same as the one created")
	}

	if err := uuid.Validate(firstEvent.ID.String()); err != nil {
		t.Fatal("Event ID is not valid")
	}
}

func TestGetEvents(t *testing.T) {
	db := New()
	if db == nil {
		t.Fatal("Failed to create database")
	}

	data := User{
		"Name": "John Doe",
	}
	db.AddEvent("user.new", data)

	events := db.GetEvents()
	if len(events) != 1 {
		t.Fatal("Failed to get events")
	}

	if !equalUser(events[0].Data, data) {
		t.Fatal("Event retrieved is not the same as the one created")
	}

	if events[0].Type != "user.new" {
		t.Fatal("Event type is not the same as the one created")
	}

	if err := uuid.Validate(events[0].ID.String()); err != nil {
		t.Fatal("Event ID is not valid")
	}
}

func TestGetEventByID(t *testing.T) {
	db := New()
	if db == nil {
		t.Fatal("Failed to create database")
	}

	data := User{
		"Name": "John Doe",
	}
	db.AddEvent("user.new", data)

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

	if !equalUser(event.Data, data) {
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

	data := User{
		"Name": "John Doe",
	}
	db.AddEvent("user.new", data)

	events := db.GetEventsByType("user.new")
	if len(events) != 1 {
		t.Fatal("Failed to get events by type")
	}

	if !equalUser(events[0].Data, data) {
		t.Fatal("Event retrieved is not the same as the one created")
	}

	if events[0].Type != "user.new" {
		t.Fatal("Event type is not the same as the one created")
	}

	if err := uuid.Validate(events[0].ID.String()); err != nil {
		t.Fatal("Event ID is not valid")
	}
}
