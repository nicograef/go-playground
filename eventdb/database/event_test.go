package database

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewEvent(t *testing.T) {
	user := user{"ID": "1", "Name": "John Doe", "Email": "john.doe@example.com"}
	event := NewEvent("user.new", "1", user)

	if err := uuid.Validate(event.ID.String()); err != nil {
		t.Fatal("Event ID is not valid")
	}

	if event.Type != "user.new" {
		t.Fatal("Event type is not the same as the one created")
	}

	if event.EntityID != "1" {
		t.Fatal("Event entity ID is not the same as the one created")
	}

	if event.Timestamp.IsZero() {
		t.Fatal("Event timestamp is not set")
	}

	if !equalUser(event.Data, user) {
		t.Fatal("Event retrieved is not the same as the one created")
	}

}
