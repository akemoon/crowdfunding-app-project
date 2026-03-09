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

	err := validateCreateProjectReq(req)
	if err != nil {
		return err
	}

	err = s.repo.CreateProject(ctx, userID, req)
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
	if p.Status == domain.StatusReview {
		return domain.Project{}, domain.ErrProjectNotFound
	}

	return p, nil
}

func (s *Service) GetProjects(ctx context.Context, req domain.GetProjectsReq) (domain.GetProjectsResp, error) {
	err := domain.ValidateStatus(req.Status)
	if err != nil {
		return domain.GetProjectsResp{}, err
	}

	if req.Sort == "" {
		req.Sort = domain.SortDefault
	} else {
		err = domain.ValidateSort(req.Sort)
		if err != nil {
			return domain.GetProjectsResp{}, err
		}
	}

	if req.Category != nil {
		err = domain.ValidateCategory(*req.Category)
		if err != nil {
			return domain.GetProjectsResp{}, err
		}
	}

	if req.Limit <= 0 || req.Limit > 100 {
		req.Limit = 20
	}
	if req.Offset < 0 {
		req.Offset = 0
	}

	items, err := s.repo.GetProjects(ctx, req)
	if err != nil {
		return domain.GetProjectsResp{}, fmt.Errorf("repo: %w", err)
	}

	return domain.GetProjectsResp{Items: items}, nil
}

func (s *Service) AddContribution(ctx context.Context, projectID uuid.UUID, amount int64) error {
	err := s.repo.AddContribution(ctx, projectID, amount)
	if err != nil {
		return fmt.Errorf("repo: %w", err)
	}

	return nil
}

func (s *Service) ApproveProject(ctx context.Context, id uuid.UUID) error {
	err := s.repo.ApproveProject(ctx, id)
	if err != nil {
		return fmt.Errorf("repo: %w", err)
	}

	return nil
}

func (s *Service) RejectProject(ctx context.Context, id uuid.UUID, reason string) error {
	err := s.repo.RejectProject(ctx, id, reason)
	if err != nil {
		return fmt.Errorf("repo: %w", err)
	}

	return nil
}

func (s *Service) UpdateProject(ctx context.Context, userID uuid.UUID, id uuid.UUID, req domain.CreateProjectReq) error {
	err := validateCreateProjectReq(req)
	if err != nil {
		return err
	}

	err = s.repo.UpdateProject(ctx, userID, id, req)
	if err != nil {
		return fmt.Errorf("repo: %w", err)
	}

	return nil
}

func (s *Service) GetApplicationByProjectID(ctx context.Context, projectID uuid.UUID) (domain.Application, error) {
	app, err := s.repo.GetApplicationByProjectID(ctx, projectID)
	if err != nil {
		return domain.Application{}, fmt.Errorf("repo: %w", err)
	}

	return app, nil
}

func (s *Service) GetPendingApplications(ctx context.Context) ([]domain.Application, error) {
	apps, err := s.repo.GetPendingApplications(ctx)
	if err != nil {
		return nil, fmt.Errorf("repo: %w", err)
	}

	return apps, nil
}

// TODO: cycle
func validateCreateProjectReq(req domain.CreateProjectReq) error {
	ve := &validation.Error{}

	err := domain.ValidateCategory(req.Category)
	if err != nil {
		ve.Add("category", err.Error())
	}
	err = domain.ValidateName(req.Name)
	if err != nil {
		ve.Add("name", err.Error())
	}
	err = domain.ValidateDescription(req.Description)
	if err != nil {
		ve.Add("description", err.Error())
	}
	err = domain.ValidateCurrency(req.Currency)
	if err != nil {
		ve.Add("currency", err.Error())
	}
	err = domain.ValidateGoalAmount(req.GoalAmount)
	if err != nil {
		ve.Add("goalAmount", err.Error())
	}
	err = domain.ValidateDurationDays(req.DurationDays)
	if err != nil {
		ve.Add("durationDays", err.Error())
	}
	if ve.HasErrors() {
		return ve
	}

	return nil
}
