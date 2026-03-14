package project

import (
	"github.com/akemoon/crowdfunding-app-project/domain"
)

type EventType int

const (
	EventTypeUnknown EventType = iota
	EventTypeFinished
)

type Event struct {
	Type    EventType      `json:"type"`
	Project domain.Project `json:"project"`
}
