package project

import (
	"context"
	"fmt"
	"io"

	"github.com/akemoon/crowdfunding-app-project/client/promocode"
	"github.com/akemoon/crowdfunding-app-project/client/storage"
	"github.com/akemoon/crowdfunding-app-project/domain"
	"github.com/akemoon/crowdfunding-app-project/repo/project"
	"github.com/akemoon/golib/validation"
	"github.com/google/uuid"
)

type Service struct {
	repo          project.Repo
	promo         promocode.Client
	storage       storage.Client
	imagesBaseURL string
}

func NewService(r project.Repo, promo promocode.Client, storage storage.Client, imagesBaseURL string) *Service {
	return &Service{
		repo:          r,
		promo:         promo,
		storage:       storage,
		imagesBaseURL: imagesBaseURL,
	}
}

func (s *Service) CreateProject(ctx context.Context, userID uuid.UUID, req domain.CreateProjectReq) (uuid.UUID, error) {
	err := validateCreateProjectReq(req)
	if err != nil {
		return uuid.UUID{}, err
	}

	id, err := s.repo.CreateProject(ctx, userID, req)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("repo: %w", err)
	}

	return id, nil
}

func (s *Service) GetProjectByID(ctx context.Context, id uuid.UUID, callerID *uuid.UUID) (domain.Project, error) {
	p, err := s.repo.GetProjectByID(ctx, id)
	if err != nil {
		return domain.Project{}, fmt.Errorf("repo: %w", err)
	}
	if callerID == nil || *callerID != p.UserID {
		if p.Status != domain.StatusActive && p.Status != domain.StatusFinished {
			return domain.Project{}, domain.ErrProjectNotFound
		}

		// Hide date
		p.BoostedUntil = nil
	}

	if p.CoverKey != nil {
		p.CoverURL = s.imagesBaseURL + "/" + *p.CoverKey
	}

	for i := range p.Images {
		p.Images[i].URL = s.imagesBaseURL + "/" + p.Images[i].StorageKey
	}

	return p, nil
}

func (s *Service) GetProjects(ctx context.Context, req domain.GetProjectsReq) (domain.GetProjectsResp, error) {
	if req.Status != domain.StatusActive && req.Status != domain.StatusFinished {
		return domain.GetProjectsResp{}, domain.ErrUnknownStatus
	}

	if req.Sort == "" {
		req.Sort = domain.SortDefault
	} else {
		err := domain.ValidateSort(req.Sort)
		if err != nil {
			return domain.GetProjectsResp{}, err
		}
	}

	if req.Category != nil {
		err := domain.ValidateCategory(*req.Category)
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

	for i := range items {
		if items[i].CoverKey != nil {
			items[i].CoverURL = s.imagesBaseURL + "/" + *items[i].CoverKey
		}
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

func (s *Service) SubmitProject(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	p, err := s.repo.GetProjectByID(ctx, id)
	if err != nil {
		return fmt.Errorf("repo: %w", err)
	}

	if p.UserID != userID {
		return domain.ErrForbidden
	}

	if p.CoverKey == nil {
		return domain.ErrProjectCoverRequired
	}

	err = s.repo.SubmitProject(ctx, id, userID)
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

	p, err := s.repo.GetProjectByID(ctx, id)
	if err != nil {
		return fmt.Errorf("repo: %w", err)
	}

	if p.UserID != userID {
		return domain.ErrForbidden
	}

	if p.Status != domain.StatusDraft {
		return domain.ErrProjectNotDraft
	}

	err = s.repo.UpdateProject(ctx, userID, id, req)
	if err != nil {
		return fmt.Errorf("repo: %w", err)
	}

	return nil
}

func (s *Service) GetApplicationByProjectID(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) (domain.Application, error) {
	app, err := s.repo.GetApplicationByProjectID(ctx, projectID, userID)
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

func (s *Service) TakeApplication(ctx context.Context, projectID uuid.UUID, managerID uuid.UUID) error {
	err := s.repo.TakeApplication(ctx, projectID, managerID)
	if err != nil {
		return fmt.Errorf("repo: %w", err)
	}

	return nil
}

func (s *Service) GetMyApplications(ctx context.Context, managerID uuid.UUID) ([]domain.Application, error) {
	apps, err := s.repo.GetMyApplications(ctx, managerID)
	if err != nil {
		return nil, fmt.Errorf("repo: %w", err)
	}

	for i := range apps {
		if apps[i].Project != nil {
			if apps[i].Project.CoverKey != nil {
				apps[i].Project.CoverURL = s.imagesBaseURL + "/" + *apps[i].Project.CoverKey
			}
			for j := range apps[i].Project.Images {
				apps[i].Project.Images[j].URL = s.imagesBaseURL + "/" + apps[i].Project.Images[j].StorageKey
			}
		}
	}

	return apps, nil
}

func (s *Service) GetProjectsByUserID(ctx context.Context, authorID uuid.UUID, isOwner bool) ([]domain.Project, error) {
	projects, err := s.repo.GetProjectsByUserID(ctx, authorID, isOwner)
	if err != nil {
		return nil, fmt.Errorf("repo: %w", err)
	}

	for i := range projects {
		if projects[i].CoverKey != nil {
			projects[i].CoverURL = s.imagesBaseURL + "/" + *projects[i].CoverKey
		}
	}

	return projects, nil
}

func (s *Service) BoostProject(ctx context.Context, userID uuid.UUID, projectID uuid.UUID, promoCode string) error {
	p, err := s.repo.GetProjectByID(ctx, projectID)
	if err != nil {
		return fmt.Errorf("repo: %w", err)
	}

	if p.UserID != userID {
		return domain.ErrForbidden
	}

	if p.Status != domain.StatusActive {
		return domain.ErrProjectNotActive
	}

	// TODO: mayne business rule - consider rejecting boost for projects with little time remaining
	// to reduce the chance of promo code being consumed on a nearly-finished project.

	days, err := s.promo.UsePromoCode(ctx, userID, promoCode, "boost_project")
	if err != nil {
		return err
	}

	// NOTE: non-atomic cross-service operation.
	// If BoostProject fails here, the promo code is already consumed but boost is not applied.
	// Proper fix: saga pattern with compensation (refund promo code via rollback API).
	err = s.repo.BoostProject(ctx, userID, projectID, days)
	if err != nil {
		return fmt.Errorf("repo: %w", err)
	}

	return nil
}

func (s *Service) UploadProjectCover(ctx context.Context, userID, projectID uuid.UUID, r io.Reader, size int64, contentType string) (string, error) {
	if !allowedImageTypes[contentType] {
		return "", domain.ErrUnsupportedFileType
	}

	p, err := s.repo.GetProjectByID(ctx, projectID)
	if err != nil {
		return "", fmt.Errorf("repo: %w", err)
	}

	if p.UserID != userID {
		return "", domain.ErrForbidden
	}

	if p.Status != domain.StatusDraft {
		return "", domain.ErrProjectNotDraft
	}

	key := fmt.Sprintf("projects/%s/cover%s", projectID, imageExtension(contentType))

	_, err = s.storage.Upload(ctx, key, r, size, contentType)
	if err != nil {
		return "", fmt.Errorf("storage: %w", err)
	}

	err = s.repo.SetProjectCover(ctx, projectID, userID, key)
	if err != nil {
		return "", fmt.Errorf("repo: %w", err)
	}

	return s.imagesBaseURL + "/" + key, nil
}

var allowedImageTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/gif":  true,
	"image/webp": true,
}

func imageExtension(contentType string) string {
	switch contentType {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/gif":
		return ".gif"
	case "image/webp":
		return ".webp"
	default:
		return ".bin"
	}
}

func (s *Service) UploadProjectImage(ctx context.Context, userID, projectID uuid.UUID, r io.Reader, size int64, contentType string) (domain.ProjectImage, error) {
	if !allowedImageTypes[contentType] {
		return domain.ProjectImage{}, domain.ErrUnsupportedFileType
	}

	p, err := s.repo.GetProjectByID(ctx, projectID)
	if err != nil {
		return domain.ProjectImage{}, fmt.Errorf("repo: %w", err)
	}

	if p.UserID != userID {
		return domain.ProjectImage{}, domain.ErrForbidden
	}

	if p.Status != domain.StatusDraft {
		return domain.ProjectImage{}, domain.ErrProjectNotDraft
	}

	key := fmt.Sprintf("projects/%s/%s%s", projectID, uuid.New(), imageExtension(contentType))

	_, err = s.storage.Upload(ctx, key, r, size, contentType)
	if err != nil {
		return domain.ProjectImage{}, fmt.Errorf("storage: %w", err)
	}

	// NOTE: non-atomic — file uploaded but DB insert may fail.
	id, err := s.repo.AddProjectImage(ctx, projectID, key)
	if err != nil {
		return domain.ProjectImage{}, fmt.Errorf("repo: %w", err)
	}

	return domain.ProjectImage{ID: id, URL: s.imagesBaseURL + "/" + key}, nil
}

func (s *Service) DeleteProjectImage(ctx context.Context, userID, projectID, imageID uuid.UUID) error {
	p, err := s.repo.GetProjectByID(ctx, projectID)
	if err != nil {
		return fmt.Errorf("repo: %w", err)
	}

	if p.UserID != userID {
		return domain.ErrProjectNotFound
	}

	img, err := s.repo.GetProjectImage(ctx, imageID, projectID)
	if err != nil {
		return fmt.Errorf("repo: %w", err)
	}

	err = s.repo.DeleteProjectImage(ctx, imageID, projectID)
	if err != nil {
		return fmt.Errorf("repo: %w", err)
	}

	// NOTE: non-atomic — DB deleted but MinIO delete may fail.
	err = s.storage.Delete(ctx, img.StorageKey)
	if err != nil {
		return fmt.Errorf("storage: %w", err)
	}

	return nil
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
