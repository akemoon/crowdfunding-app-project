package postgres

import (
	"fmt"

	"github.com/akemoon/crowdfunding-app-project/domain"
)

func MapStatusToDB(s domain.Status) (int, error) {
	switch s {
	case domain.StatusActive:
		return 1, nil
	case domain.StatusFinished:
		return 2, nil
	default:
		return 0, fmt.Errorf("%w: map status to db err", domain.ErrInternal)
	}
}

func MapStatusFromDB(id int) (domain.Status, error) {
	switch id {
	case 1:
		return domain.StatusActive, nil
	case 2:
		return domain.StatusFinished, nil
	default:
		return "", fmt.Errorf("%w: map status from db err", domain.ErrInternal)
	}
}

func MapCategoryToDB(c domain.Category) (int, error) {
	switch c {
	case domain.CategoryScience:
		return 1, nil
	case domain.CategoryTech:
		return 2, nil
	case domain.CategoryArchitectureAndUrban:
		return 3, nil
	case domain.CategorySport:
		return 4, nil
	case domain.CategoryMusic:
		return 5, nil
	default:
		return 0, fmt.Errorf("%w: map category to db err", domain.ErrInternal)
	}
}

func MapCurrencyToDB(c domain.Currency) (int, error) {
	switch c {
	case domain.CurrencyRUB:
		return 1, nil
	case domain.CurrencyUSD:
		return 2, nil
	default:
		return 0, fmt.Errorf("%w: map currency to db err", domain.ErrInternal)
	}
}

func MapCategoryFromDB(id int) (domain.Category, error) {
	switch id {
	case 1:
		return domain.CategoryScience, nil
	case 2:
		return domain.CategoryTech, nil
	case 3:
		return domain.CategoryArchitectureAndUrban, nil
	case 4:
		return domain.CategorySport, nil
	case 5:
		return domain.CategoryMusic, nil
	default:
		return "", fmt.Errorf("%w: map category from db err", domain.ErrInternal)
	}
}

func MapCurrencyFromDB(id int) (domain.Currency, error) {
	switch id {
	case 1:
		return domain.CurrencyRUB, nil
	case 2:
		return domain.CurrencyUSD, nil
	default:
		return "", fmt.Errorf("%w: map currency from db err", domain.ErrInternal)
	}
}

func MapProjectFromDB(p ProjectDB) (domain.Project, error) {
	category, err := MapCategoryFromDB(p.CategoryID)
	if err != nil {
		return domain.Project{}, err
	}

	currency, err := MapCurrencyFromDB(p.CurrencyID)
	if err != nil {
		return domain.Project{}, err
	}

	status, err := MapStatusFromDB(p.StatusID)
	if err != nil {
		return domain.Project{}, err
	}

	return domain.Project{
		ID:            p.ID,
		UserID:        p.UserID,
		Category:      category,
		Name:          p.Name,
		Description:   p.Description,
		Currency:      currency,
		GoalAmount:    p.GoalAmount,
		CurrentAmount: p.CurrentAmount,
		StartedAt:     p.StartedAt,
		DurationDays:  p.DurationDays,
		Status:        status,
		IsBoosted:     p.IsBoosted,
	}, nil
}
