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
	StartDate     time.Time
	DurationDays  int
	IsBoosted     bool
}
