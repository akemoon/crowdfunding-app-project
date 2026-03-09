package project

import (
	"github.com/akemoon/crowdfunding-app-project/domain"
	"github.com/google/uuid"
)

type EventType int

const (
	EventTypeUnknown EventType = iota
	EventTypeFinished
)

type Event struct {
	EventID uuid.UUID      `json:"eventID"`
	Type    EventType      `json:"type"`
	Project domain.Project `json:"project"`
}
