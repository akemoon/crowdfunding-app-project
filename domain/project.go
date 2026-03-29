package domain

import (
	"time"

	"github.com/google/uuid"
)

type Category string

const (
	CategoryScience              Category = "science"
	CategoryTech                 Category = "tech"
	CategoryArchitectureAndUrban Category = "architecture_and_urban"
	CategorySport                Category = "sport"
	CategoryMusic                Category = "music"
)

func ValidateCategory(c Category) error {
	switch c {
	case CategoryScience, CategoryTech, CategoryArchitectureAndUrban, CategorySport, CategoryMusic:
		return nil
	default:
		return ErrUnknownCategory
	}
}

type Currency string

const (
	CurrencyRUB Currency = "RUB"
	CurrencyUSD Currency = "USD"
)

func ValidateCurrency(c Currency) error {
	switch c {
	case CurrencyRUB, CurrencyUSD:
		return nil
	default:
		return ErrUnknownCurrency
	}
}

type Status string

const (
	StatusDraft    Status = "draft"
	StatusReview   Status = "review"
	StatusActive   Status = "active"
	StatusFinished Status = "finished"
)

type CreateProjectReq struct {
	UserID       uuid.UUID `json:"user_id"`
	Category     Category  `json:"category"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Currency     Currency  `json:"currency"`
	GoalAmount   int64     `json:"goal_amount"`
	DurationDays int       `json:"duration_days"`
}

type Project struct {
	ID           uuid.UUID  `json:"id"`
	UserID       uuid.UUID  `json:"user_id"`
	Category     Category   `json:"category"`
	Name         string     `json:"name"`
	Description  string     `json:"description"`
	Currency     Currency   `json:"currency"`
	GoalAmount   int64      `json:"goal_amount"`
	StartedAt    *time.Time `json:"started_at,omitempty"`
	DurationDays int        `json:"duration_days"`
	Status       Status     `json:"status"`
}

const (
	MinNameLen = 1
	MaxNameLen = 100
)

func ValidateName(name string) error {
	if len(name) < MinNameLen || len(name) > MaxNameLen {
		return ErrInvalidNameLen
	}
	return nil
}

const (
	MinDescriptionLen = 0
	MaxDescriptionLen = 1000
)

func ValidateDescription(description string) error {
	if len(description) < MinDescriptionLen || len(description) > MaxDescriptionLen {
		return ErrInvalidDescriptionLen
	}
	return nil
}

const MinGoalAmount = 1

func ValidateGoalAmount(amount int64) error {
	if amount < MinGoalAmount {
		return ErrInvalidGoalAmount
	}
	return nil
}

const (
	MinDurationDays = 1
	MaxDurationDays = 60
)

func ValidateDurationDays(daysNum int) error {
	if daysNum < MinDurationDays || daysNum > MaxDurationDays {
		return ErrInvalidDurationDays
	}
	return nil
}
