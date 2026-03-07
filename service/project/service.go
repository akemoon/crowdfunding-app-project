package project

import (
	"context"
	"fmt"

	"github.com/akemoon/crowdfunding-app-project/domain"
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

	err := domain.ValidateCreateProjectReq(req)
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

	return p, nil
}

// TODO: use search keys (categorie, status, name, user_id)
//
//func (s *Service) GetProjects() ([]project.Project, error) {
//}
