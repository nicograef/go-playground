package database

import (
	"testing"
)

func TestGetEntity(t *testing.T) {
	db := New()
	if db == nil {
		t.Fatal("Failed to create database")
	}

	db.AddEvent("user.new", "1", user{"ID": "1", "Name": "John Doe"})
	db.AddEvent("user.update", "1", user{"ID": "1", "Email": "john.doe@example.com"})
	db.AddEvent("user.update", "1", user{"ID": "1", "Name": "johnny donny"})
	db.AddEvent("user.new", "2", user{"ID": "2", "Name": "Max Mustermann"})

	if entity := db.GetEntity("non-existing-id"); len(entity) != 0 {
		t.Fatal("Expected no entity to be found")
	}

	entity := db.GetEntity("1")
	if entity == nil {
		t.Fatal("Failed to get entity by ID")
	}

	if !equalUser(entity, user{"ID": "1", "Name": "johnny donny", "Email": "john.doe@example.com"}) {
		t.Fatal("Entity retrieved is not the same as the one created")
	}
}

func TestDeleteEntity(t *testing.T) {
	db := New()

	// adding 3 events
	db.AddEvent("user.new", "1", user{"ID": "1", "Name": "John Doe"})
	db.AddEvent("user.update", "1", user{"ID": "1", "Email": "john.doe@example.com"})
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

	// deleting the whole entity
	db.DeleteEntity("1")

	if len(db.events) != 1 {
		t.Fatal("Expected only 1 event to be present")
	}
	if events := db.GetEventsByEntity("1"); len(events) != 0 {
		t.Fatal("Expected event to be deleted")
	}
	if entity := db.GetEntity("1"); entity != nil {
		t.Fatal("Expected entity to be deleted")
	}
	if len(db.typeIndex["user.update"]) != 0 {
		t.Fatal("Expected no update events to be present")
	}
	if len(db.entityIndex) != 1 {
		t.Fatal("Expected only 1 entity to be present in the index")
	}
}
