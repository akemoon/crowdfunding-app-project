package domain

import (
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/google/uuid"
)

type Category string

const (
	CategoryScience              Category = "science"
	CategoryTech                 Category = "tech"
	CategoryArchitectureAndUrban Category = "architecture_and_urban"
	CategorySport                Category = "sport"
	CategoryMusic                Category = "music"
	CategoryArt                  Category = "art"
	CategoryFilm                 Category = "film"
	CategoryGames                Category = "games"
	CategoryEducation            Category = "education"
	CategoryFood                 Category = "food"
	CategoryFashion              Category = "fashion"
	CategoryHealth               Category = "health"
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

type ProjectImage struct {
	ID         uuid.UUID `json:"id"`
	URL        string    `json:"url"`
	StorageKey string    `json:"-"`
}

type Project struct {
	ID            uuid.UUID      `json:"id"`
	UserID        uuid.UUID      `json:"userID"`
	Category      Category       `json:"category"`
	Name          string         `json:"name"`
	Description   string         `json:"description"`
	Currency      Currency       `json:"currency"`
	GoalAmount    int64          `json:"goalAmount"`
	CurrentAmount int64          `json:"currentAmount"`
	StartedAt     *time.Time     `json:"startedAt,omitempty"`
	DurationDays  int            `json:"durationDays"`
	Status        Status         `json:"status"`
	IsBoosted     bool           `json:"isBoosted"`
	BoostedUntil  *time.Time     `json:"boostedUntil,omitempty"`
	CoverURL      string         `json:"coverURL,omitempty"`
	CoverKey      *string        `json:"-"`
	Images        []ProjectImage `json:"images,omitempty"`
}

type Status string

const (
	StatusDraft    Status = "draft"
	StatusReview   Status = "review"
	StatusActive   Status = "active"
	StatusFinished Status = "finished"
)


type ApplicationStatus string

const (
	ApplicationStatusPending  ApplicationStatus = "pending"
	ApplicationStatusInReview ApplicationStatus = "review"
	ApplicationStatusRejected ApplicationStatus = "rejected"
	ApplicationStatusApproved ApplicationStatus = "approved"
)

type Application struct {
	Status       ApplicationStatus `json:"status"`
	RejectReason string            `json:"rejectReason,omitempty"`
	CreatedAt    time.Time         `json:"createdAt"`
	AssignedAt   *time.Time        `json:"assignedAt,omitempty"`
	ProcessedAt  *time.Time        `json:"processedAt,omitempty"`
	Project      *Project          `json:"project,omitempty"`
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

var allowedNamePunct = map[rune]bool{
	'.': true, ',': true, '!': true, '?': true,
	'-': true, '\'': true, '"': true, ':': true,
}

func ValidateName(name string) error {
	n := utf8.RuneCountInString(name)
	if n < MinNameLen || n > MaxNameLen {
		return ErrInvalidNameLen
	}
	if name != strings.TrimSpace(name) {
		return ErrInvalidNameLen
	}
	if strings.Contains(name, "  ") {
		return ErrInvalidNameLen
	}
	for _, r := range name {
		if unicode.Is(unicode.Latin, r) || unicode.Is(unicode.Cyrillic, r) || unicode.IsDigit(r) || r == ' ' || allowedNamePunct[r] {
			continue
		}
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
		CategorySport, CategoryMusic, CategoryArt, CategoryFilm,
		CategoryGames, CategoryEducation, CategoryFood, CategoryFashion,
		CategoryHealth:
		return nil
	default:
		return ErrUnknownCategory
	}
}
