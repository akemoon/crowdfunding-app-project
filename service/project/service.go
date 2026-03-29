package project

import (
	"context"
	"fmt"

	"github.com/akemoon/crowdfunding-app-project/domain"
	projectRepo "github.com/akemoon/crowdfunding-app-project/repo/project"
	"github.com/akemoon/golib/validation"
	"github.com/google/uuid"
)

type Service struct {
	repo projectRepo.Repo
}

func NewService(r projectRepo.Repo) *Service {
	return &Service{repo: r}
}

func (s *Service) CreateProject(ctx context.Context, req domain.CreateProjectReq) (domain.Project, error) {
	err := validateCreateProjectReq(req)
	if err != nil {
		return domain.Project{}, err
	}

	p, err := s.repo.CreateProject(ctx, req)
	if err != nil {
		return domain.Project{}, fmt.Errorf("repo: %w", err)
	}

	return p, nil
}

func (s *Service) GetProjects(ctx context.Context) ([]domain.Project, error) {
	projects, err := s.repo.GetProjects(ctx)
	if err != nil {
		return nil, fmt.Errorf("repo: %w", err)
	}
	return projects, nil
}

func (s *Service) GetProjectByID(ctx context.Context, id uuid.UUID) (domain.Project, error) {
	p, err := s.repo.GetProjectByID(ctx, id)
	if err != nil {
		return domain.Project{}, fmt.Errorf("repo: %w", err)
	}
	return p, nil
}

func validateCreateProjectReq(req domain.CreateProjectReq) error {
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
		ve.Add("goal_amount", err.Error())
	}
	if err := domain.ValidateDurationDays(req.DurationDays); err != nil {
		ve.Add("duration_days", err.Error())
	}

	if ve.HasErrors() {
		return ve
	}
	return nil
}
