package project

import (
	"context"

	"github.com/akemoon/crowdfunding-app-project/domain"
	"github.com/google/uuid"
)

type Repo interface {
	CreateProject(ctx context.Context, req domain.CreateProjectReq) (domain.Project, error)
	GetProjectByID(ctx context.Context, id uuid.UUID) (domain.Project, error)
	GetProjects(ctx context.Context) ([]domain.Project, error)
}
