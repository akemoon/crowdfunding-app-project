package postgres

import (
	"fmt"

	"github.com/akemoon/crowdfunding-app-project/domain"
)

func MapStatusToDB(s domain.Status) (int, error) {
	switch s {
	case domain.StatusReview:
		return 1, nil
	case domain.StatusActive:
		return 2, nil
	case domain.StatusFinished:
		return 3, nil
	case domain.StatusDraft:
		return 4, nil
	default:
		return 0, fmt.Errorf("%w: map status to db err", domain.ErrInternal)
	}
}

func MapStatusFromDB(id int) (domain.Status, error) {
	switch id {
	case 1:
		return domain.StatusReview, nil
	case 2:
		return domain.StatusActive, nil
	case 3:
		return domain.StatusFinished, nil
	case 4:
		return domain.StatusDraft, nil
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
	case domain.CategoryArt:
		return 6, nil
	case domain.CategoryFilm:
		return 7, nil
	case domain.CategoryGames:
		return 8, nil
	case domain.CategoryEducation:
		return 9, nil
	case domain.CategoryFood:
		return 10, nil
	case domain.CategoryFashion:
		return 11, nil
	case domain.CategoryHealth:
		return 12, nil
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
	case 6:
		return domain.CategoryArt, nil
	case 7:
		return domain.CategoryFilm, nil
	case 8:
		return domain.CategoryGames, nil
	case 9:
		return domain.CategoryEducation, nil
	case 10:
		return domain.CategoryFood, nil
	case 11:
		return domain.CategoryFashion, nil
	case 12:
		return domain.CategoryHealth, nil
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

func MapApplicationStatusFromDB(id int) (domain.ApplicationStatus, error) {
	switch id {
	case 1:
		return domain.ApplicationStatusPending, nil
	case 2:
		return domain.ApplicationStatusInReview, nil
	case 3:
		return domain.ApplicationStatusRejected, nil
	case 4:
		return domain.ApplicationStatusApproved, nil
	default:
		return "", fmt.Errorf("%w: map application status from db err", domain.ErrInternal)
	}
}

func MapApplicationStatusToDB(s domain.ApplicationStatus) (int, error) {
	switch s {
	case domain.ApplicationStatusPending:
		return 1, nil
	case domain.ApplicationStatusInReview:
		return 2, nil
	case domain.ApplicationStatusRejected:
		return 3, nil
	case domain.ApplicationStatusApproved:
		return 4, nil
	default:
		return 0, fmt.Errorf("%w: map application status to db err", domain.ErrInternal)
	}
}

func MapApplicationFromDB(a ApplicationDB, p ProjectDB) (domain.Application, error) {
	status, err := MapApplicationStatusFromDB(a.StatusID)
	if err != nil {
		return domain.Application{}, err
	}

	project, err := MapProjectFromDB(p)
	if err != nil {
		return domain.Application{}, err
	}
	project.StartedAt = nil

	return domain.Application{
		Status:       status,
		RejectReason: a.RejectReason,
		CreatedAt:    a.CreatedAt,
		AssignedAt:   a.AssignedAt,
		ProcessedAt:  a.ProcessedAt,
		Project:      &project,
	}, nil
}

func MapProjectImageFromDB(img ProjectImageDB) domain.ProjectImage {
	return domain.ProjectImage{
		ID:         img.ID,
		StorageKey: img.StorageKey,
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
		BoostedUntil:  p.BoostedUntil,
		CoverKey:      p.CoverKey,
	}, nil
}
