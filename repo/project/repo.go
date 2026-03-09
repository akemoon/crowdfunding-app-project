package project

import (
	"context"

	"github.com/akemoon/crowdfunding-app-project/domain"
	"github.com/google/uuid"
)

type Repo interface {
	CreateProject(ctx context.Context, userID uuid.UUID, req domain.CreateProjectReq) error
	GetProjectByID(ctx context.Context, id uuid.UUID) (domain.Project, error)
	GetProjects(ctx context.Context, req domain.GetProjectsReq) ([]domain.Project, error)
	FinishProjects(ctx context.Context) error
	ListPendingFinishedOutbox(ctx context.Context, limit int) ([]domain.Project, error)
	MarkSentFinishedOutbox(ctx context.Context, projectID uuid.UUID) error
	AddContribution(ctx context.Context, projectID uuid.UUID, amount int64) error
	ApproveProject(ctx context.Context, id uuid.UUID) error
	RejectProject(ctx context.Context, id uuid.UUID, reason string) error
	UpdateProject(ctx context.Context, userID uuid.UUID, id uuid.UUID, req domain.CreateProjectReq) error
	GetApplicationByProjectID(ctx context.Context, projectID uuid.UUID) (domain.Application, error)
	GetPendingApplications(ctx context.Context) ([]domain.Application, error)
}
