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
		"Name": "John Doe",
	}
	db.AddEvent("user.new", data)

	if err := db.PersistToJsonFile(); err != nil {
		t.Fatal("Failed to persist to JSON file:", err)
	}
}

func TestLoadDatabaseFromJsonFile(t *testing.T) {
	err := os.WriteFile("database.json", []byte(`{"events":{"67099372-1947-4c9d-bc68-423c52d00b76":{"id":"67099372-1947-4c9d-bc68-423c52d00b76","timestamp":"2025-08-31T11:38:51.517259326Z","type":"user.new","data":{"Name":"John Doe"}}}}`), 0644)
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

	id, err := uuid.Parse("67099372-1947-4c9d-bc68-423c52d00b76")
	if err != nil {
		t.Fatal("Failed to parse UUID:", err)
	}
	event := db.GetEventByID(id)
	if event == nil {
		t.Fatal("Failed to get event by ID")
	}

	expectedUser := User{
		"Name": "John Doe",
	}
	if !equalUser(event.Data, expectedUser) {
		t.Fatal("Event data is not the same as the one created")
	}
}
