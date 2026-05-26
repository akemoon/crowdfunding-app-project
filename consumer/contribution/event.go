package contribution

import "github.com/google/uuid"

type EventType int

const (
	EventTypeUnknown EventType = iota
	EventTypeCreated EventType = 1
)

type Event struct {
	EventID   uuid.UUID `json:"eventID"`
	Type      EventType `json:"type"`
	ProjectID uuid.UUID `json:"projectID"`
	Amount    int64     `json:"amount"`
}
