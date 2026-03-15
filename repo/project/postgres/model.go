package postgres

import (
	"time"

	"github.com/google/uuid"
)

type ProjectDB struct {
	ID            uuid.UUID
	UserID        uuid.UUID
	CategoryID    int
	Name          string
	Description   string
	CurrencyID    int
	GoalAmount    int64
	CurrentAmount int64
	StartedAt     *time.Time
	DurationDays  int
	StatusID      int
	IsBoosted     bool
	BoostedUntil  *time.Time
}

type ApplicationDB struct {
	StatusID     int
	AssignedTo   *uuid.UUID
	RejectReason string
	CreatedAt    time.Time
	AssignedAt   *time.Time
	ProcessedAt  *time.Time
}
