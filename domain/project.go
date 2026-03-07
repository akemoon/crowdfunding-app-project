package domain

import (
	"errors"
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

type ProjectsListResp struct {
	Items []Project `json:"items"`
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
}

const (
	MinNameLen = 1
	MaxNameLen = 100
)

// NOTE: maybe my unique field for testing
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

func ValidateCreateProjectReq(req CreateProjectReq) error {
	return validateAll(
		func() error { return ValidateCategory(req.Category) },
		func() error { return ValidateName(req.Name) },
		func() error { return ValidateDescription(req.Description) },
		func() error { return ValidateCurrency(req.Currency) },
		func() error { return ValidateGoalAmount(req.GoalAmount) },
		func() error { return ValidateDurationDays(req.DurationDays) },
	)
}

func validateAll(checks ...func() error) error {
	var errs []error
	for _, check := range checks {
		err := check()
		if err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return errors.Join(errs...)
}
