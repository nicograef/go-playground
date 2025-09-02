package database

import "maps"

// GetEntity retrieves all events for one entity
// and builds the aggregate entity by applying all events sorted by their timestamp
// to the initial state.
func (db *Database) GetEntity(entityID string) map[string]any {
	events := db.GetEventsByEntity(entityID)

	if len(events) == 0 {
		return nil
	}

	aggregate := make(map[string]any)

	for _, event := range events {
		maps.Copy(aggregate, event.Data)
	}

	return aggregate
}

// DeleteEntity deletes all events of an entity.
func (db *Database) DeleteEntity(entityID string) {
	events := db.GetEventsByEntity(entityID)

	for _, event := range events {
		db.DeleteEvent(event.ID)
	}

	delete(db.entityIndex, entityID)
}
