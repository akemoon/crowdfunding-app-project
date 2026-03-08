package project

import (
	"context"
	"fmt"

	"github.com/akemoon/crowdfunding-app-project/domain"
	"github.com/akemoon/crowdfunding-app-project/golib/validation"
	"github.com/akemoon/crowdfunding-app-project/repo/project"
	"github.com/google/uuid"
)

type Service struct {
	repo project.Repo
}

func NewService(r project.Repo) *Service {
	return &Service{
		repo: r,
	}
}

func (s *Service) CreateProject(ctx context.Context, userID uuid.UUID, req domain.CreateProjectReq) error {
	// TODO: check author account

	ve := &validation.Error{}

	if err := domain.ValidateCategory(req.Category); err != nil {
		ve.Add("category", err.Error())
	}
	if err := domain.ValidateName(req.Name); err != nil {
		ve.Add("name", err.Error())
	}
	if err := domain.ValidateDescription(req.Description); err != nil {
		ve.Add("description", err.Error())
	}
	if err := domain.ValidateCurrency(req.Currency); err != nil {
		ve.Add("currency", err.Error())
	}
	if err := domain.ValidateGoalAmount(req.GoalAmount); err != nil {
		ve.Add("goalAmount", err.Error())
	}
	if err := domain.ValidateDurationDays(req.DurationDays); err != nil {
		ve.Add("durationDays", err.Error())
	}
	if ve.HasErrors() {
		return ve
	}

	err := s.repo.CreateProject(ctx, userID, req)
	if err != nil {
		return fmt.Errorf("repo: %w", err)
	}

	return nil
}

func (s *Service) GetProjectByID(ctx context.Context, id uuid.UUID) (domain.Project, error) {
	p, err := s.repo.GetProjectByID(ctx, id)
	if err != nil {
		return domain.Project{}, fmt.Errorf("repo: %w", err)
	}

	return p, nil
}

// TODO: use search keys (categorie, status, name, user_id)
//
//func (s *Service) GetProjects() ([]project.Project, error) {
//}
