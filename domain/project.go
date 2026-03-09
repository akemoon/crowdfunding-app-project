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

type Currency string

const (
	CurrencyRUB Currency = "RUB"
	CurrencyUSD Currency = "USD"
)

type CreateProjectReq struct {
	Category     Category `json:"category"`
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Currency     Currency `json:"currency"`
	GoalAmount   int64    `json:"goalAmount"`
	DurationDays int      `json:"durationDays"`
}

type Project struct {
	ID            uuid.UUID `json:"id"`
	UserID        uuid.UUID `json:"userID"`
	Category      Category  `json:"category"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	Currency      Currency  `json:"currency"`
	GoalAmount    int64     `json:"goalAmount"`
	CurrentAmount int64     `json:"currentAmount"`
	StartedAt     time.Time `json:"startedAt"`
	DurationDays  int       `json:"durationDays"`
	Status        Status    `json:"status"`
	IsBoosted     bool      `json:"isBoosted"`
}

type Status string

const (
	StatusActive   Status = "active"
	StatusFinished Status = "finished"
)

func ValidateStatus(s Status) error {
	switch s {
	case StatusActive, StatusFinished:
		return nil
	default:
		return ErrUnknownStatus
	}
}

type Sort string

const (
	SortDefault Sort = "default"
	SortDate    Sort = "date"
)

func ValidateSort(s Sort) error {
	switch s {
	case SortDefault, SortDate:
		return nil
	default:
		return ErrUnknownSort
	}
}

type GetProjectsReq struct {
	Status   Status
	Sort     Sort
	Category *Category
	Search   *string
	Limit    int
	Offset   int
}

type GetProjectsResp struct {
	Items []Project `json:"items"`
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

func ValidateCurrency(c Currency) error {
	switch c {
	case CurrencyRUB, CurrencyUSD:
		return nil
	default:
		return ErrUnknownCurrency
	}
}

func ValidateCategory(c Category) error {
	switch c {
	case CategoryScience, CategoryTech, CategoryArchitectureAndUrban,
		CategoryMusic, CategorySport:
		return nil
	default:
		return ErrUnknownCategory
	}
}
