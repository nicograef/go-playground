package database

import (
	"os"
	"testing"

	"github.com/google/uuid"
)

func TestPersistToJsonFile(t *testing.T) {
	defer os.Remove("database.json") // Clean up after test

	db := New()
	if db == nil {
		t.Fatal("Failed to create database")
	}

	data := User{
		"ID":   "1",
		"Name": "John Doe",
	}
	db.AddEvent("user.new", "1", data)

	if err := db.PersistToJsonFile(); err != nil {
		t.Fatal("Failed to persist to JSON file:", err)
	}
}

func TestLoadDatabaseFromJsonFile(t *testing.T) {
	err := os.WriteFile("database.json", []byte(`{"created":"2025-09-01T17:09:53.684081205Z","events":{"f8ceae97-5a98-473d-a075-c1f0a530da2c":{"id":"f8ceae97-5a98-473d-a075-c1f0a530da2c","timestamp":"2025-09-01T17:09:53.68409515Z","type":"user.new","entity_id":"1","data":{"ID":"1","Name":"John Doe"}}}}`), 0644)
	if err != nil {
		t.Fatal("Failed to write database.json:", err)
	}
	defer os.Remove("database.json")

	db, err := LoadDatabaseFromJsonFile()
	if err != nil {
		t.Fatal("Failed to load database from JSON file:", err)
	}

	if db == nil {
		t.Fatal("Loaded database is nil")
	}

	if len(db.events) != 1 {
		t.Fatal("Failed to load events")
	}

	id, err := uuid.Parse("f8ceae97-5a98-473d-a075-c1f0a530da2c")
	if err != nil {
		t.Fatal("Failed to parse UUID:", err)
	}
	event := db.GetEventByID(id)
	if event == nil {
		t.Fatal("Failed to get event by ID")
	}

	expectedUser := User{
		"ID":   "1",
		"Name": "John Doe",
	}
	if !equalUser(event.Data, expectedUser) {
		t.Fatal("Event data is not the same as the one created")
	}
}
