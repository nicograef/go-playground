package database

import (
	"encoding/json"
	"os"
	"time"

	"github.com/google/uuid"
)

// DatabaseMarshal is a helper struct for JSON marshalling/unmarshalling
type DatabaseMarshal struct {
	Created time.Time           `json:"created"`
	Events  map[uuid.UUID]Event `json:"events"`
}

func (db *Database) PersistToJsonFile() error {
	data, err := json.MarshalIndent(DatabaseMarshal{
		Events:  db.events,
		Created: db.created,
	}, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile("database.json", data, 0644)
}

func LoadDatabaseFromJsonFile() (*Database, error) {
	data, err := os.ReadFile("database.json")
	if err != nil {
		return nil, err
	}

	var marshal DatabaseMarshal
	if err := json.Unmarshal(data, &marshal); err != nil {
		return nil, err
	}

	return &Database{
		created: marshal.Created,
		events:  marshal.Events,
	}, nil
}
