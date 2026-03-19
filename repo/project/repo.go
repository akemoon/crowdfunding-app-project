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
	GetProjectsByUserID(ctx context.Context, authorID uuid.UUID, isOwner bool) ([]domain.Project, error)
	FinishProjects(ctx context.Context) error
	ListPendingFinishedOutbox(ctx context.Context, limit int) ([]domain.Project, error)
	MarkSentFinishedOutbox(ctx context.Context, projectID uuid.UUID) error
	AddContribution(ctx context.Context, projectID uuid.UUID, amount int64) error

	UpdateProject(ctx context.Context, userID uuid.UUID, id uuid.UUID, req domain.CreateProjectReq) error

	ApproveProject(ctx context.Context, id uuid.UUID) error
	RejectProject(ctx context.Context, id uuid.UUID, reason string) error
	GetApplicationByProjectID(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) (domain.Application, error)
	GetPendingApplications(ctx context.Context) ([]domain.Application, error)
	TakeApplication(ctx context.Context, projectID uuid.UUID, managerID uuid.UUID) error
	TakeApplicationForce(ctx context.Context, projectID uuid.UUID, managerID uuid.UUID) error
	GetMyApplications(ctx context.Context, managerID uuid.UUID) ([]domain.Application, error)

	BoostProject(ctx context.Context, userID uuid.UUID, projectID uuid.UUID, days int) error
}
