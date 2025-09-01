package database

import (
	"encoding/json"
	"os"
	"time"

	"github.com/google/uuid"
)

// DatabaseMarshal is a helper struct for JSON marshalling/unmarshalling
type DatabaseMarshal struct {
	Created     time.Time              `json:"created"`
	Events      map[uuid.UUID]Event    `json:"events"`
	TypeIndex   map[string][]uuid.UUID `json:"type_index"`
	EntityIndex map[string][]uuid.UUID `json:"entity_index"`
}

func (db *Database) PersistToJsonFile() error {
	data, err := json.MarshalIndent(DatabaseMarshal{
		Created:     db.created,
		Events:      db.events,
		TypeIndex:   db.typeIndex,
		EntityIndex: db.entityIndex,
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
		created:     marshal.Created,
		events:      marshal.Events,
		typeIndex:   marshal.TypeIndex,
		entityIndex: marshal.EntityIndex,
	}, nil
}
