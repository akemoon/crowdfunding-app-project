package postgres

import (
	"context"
	"database/sql"
	"errors"

	_ "embed"

	"github.com/akemoon/crowdfunding-app-project/domain"
	"github.com/akemoon/crowdfunding-app-project/golib/pglib"
	"github.com/google/uuid"
)

var projectConstraints = map[string]error{
	"projects_user_id_name_unique": domain.ErrProjectExists,
}

type ProjectRepo struct {
	db *sql.DB
}

func NewProjectRepo(db *sql.DB) *ProjectRepo {
	return &ProjectRepo{
		db: db,
	}
}

//go:embed sql/create_project.sql
var createProjectSQL string

func (r *ProjectRepo) CreateProject(ctx context.Context, userID uuid.UUID, req domain.CreateProjectReq) error {
	dbCategory, err := MapCategoryToDB(req.Category)
	if err != nil {
		return err
	}

	dbCurrency, err := MapCurrencyToDB(req.Currency)
	if err != nil {
		return err
	}

	// Prevent toctou: https://cwe.mitre.org/data/definitions/367.html
	_, err = r.db.ExecContext(ctx, createProjectSQL,
		userID,
		dbCategory,
		req.Name,
		req.Description,
		dbCurrency,
		req.GoalAmount,
		req.DurationDays,
	)
	if err != nil {
		return pglib.MapConstraintErr(err, projectConstraints, err)
	}

	return nil
}

//go:embed sql/get_project_by_id.sql
var getProjectByIDSQL string

func (r *ProjectRepo) GetProjectByID(ctx context.Context, id uuid.UUID) (domain.Project, error) {
	var pDB ProjectDB

	err := r.db.QueryRowContext(ctx, getProjectByIDSQL, id).Scan(
		&pDB.ID,
		&pDB.UserID,
		&pDB.CategoryID,
		&pDB.Name,
		&pDB.Description,
		&pDB.CurrencyID,
		&pDB.GoalAmount,
		&pDB.CurrentAmount,
		&pDB.StartedAt,
		&pDB.DurationDays,
		&pDB.StatusID,
		&pDB.IsBoosted,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Project{}, domain.ErrProjectNotFound
		}
		return domain.Project{}, err
	}

	p, err := MapProjectFromDB(pDB)
	if err != nil {
		return domain.Project{}, err
	}

	return p, nil
}

//go:embed sql/get_projects.sql
var getProjectsSQL string

func (r *ProjectRepo) GetProjects(ctx context.Context, req domain.GetProjectsReq) ([]domain.Project, error) {
	statusID, err := MapStatusToDB(req.Status)
	if err != nil {
		return nil, err
	}

	var categoryID interface{}
	if req.Category != nil {
		id, err := MapCategoryToDB(*req.Category)
		if err != nil {
			return nil, err
		}
		categoryID = id
	}

	rows, err := r.db.QueryContext(ctx, getProjectsSQL,
		statusID,
		categoryID,
		req.Search,
		req.Limit,
		req.Offset,
		string(req.Sort),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.Project

	for rows.Next() {
		var pDB ProjectDB
		if err := rows.Scan(
			&pDB.ID,
			&pDB.UserID,
			&pDB.CategoryID,
			&pDB.Name,
			&pDB.Description,
			&pDB.CurrencyID,
			&pDB.GoalAmount,
			&pDB.CurrentAmount,
			&pDB.StartedAt,
			&pDB.DurationDays,
			&pDB.StatusID,
			&pDB.IsBoosted,
		); err != nil {
			return nil, err
		}

		p, err := MapProjectFromDB(pDB)
		if err != nil {
			return nil, err
		}

		out = append(out, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

//go:embed sql/finish_projects.sql
var finishProjectsSQL string

func (r *ProjectRepo) FinishProjects(ctx context.Context) error {
	_, err := r.db.ExecContext(
		ctx,
		finishProjectsSQL,
	)
	if err != nil {
		return err
	}

	return nil
}

//go:embed sql/list_pending_finished_outbox.sql
var listPendingFinishedOutboxSQL string

func (r *ProjectRepo) ListPendingFinishedOutbox(ctx context.Context, limit int) ([]domain.Project, error) {
	if limit <= 0 {
		limit = 10
	}

	rows, err := r.db.QueryContext(ctx, listPendingFinishedOutboxSQL, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.Project

	for rows.Next() {
		var pDB ProjectDB
		if err := rows.Scan(
			&pDB.ID,
			&pDB.UserID,
			&pDB.CategoryID,
			&pDB.Name,
			&pDB.Description,
			&pDB.CurrencyID,
			&pDB.GoalAmount,
			&pDB.CurrentAmount,
			&pDB.StartedAt,
			&pDB.DurationDays,
			&pDB.StatusID,
			&pDB.IsBoosted,
		); err != nil {
			return nil, err
		}

		p, err := MapProjectFromDB(pDB)
		if err != nil {
			return nil, err
		}

		out = append(out, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

//go:embed sql/mark_sent_finished_outbox.sql
var markSentFinishedOutboxSQL string

func (r *ProjectRepo) MarkSentFinishedOutbox(ctx context.Context, projectID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, markSentFinishedOutboxSQL, projectID)
	if err != nil {
		return err
	}

	return nil
}

//go:embed sql/add_contribution.sql
var addContributionSQL string

// TODO: idempotency - add processed_events table to skip duplicate contribution.created events
func (r *ProjectRepo) AddContribution(ctx context.Context, projectID uuid.UUID, amount int64) error {
	_, err := r.db.ExecContext(ctx, addContributionSQL, projectID, amount)
	if err != nil {
		return err
	}

	return nil
}

//go:embed sql/approve_project.sql
var approveProjectSQL string

func (r *ProjectRepo) ApproveProject(ctx context.Context, id uuid.UUID) error {
	res, err := r.db.ExecContext(ctx, approveProjectSQL, id)
	if err != nil {
		return err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return domain.ErrProjectNotOnReview
	}

	return nil
}

//go:embed sql/reject_project.sql
var rejectProjectSQL string

func (r *ProjectRepo) RejectProject(ctx context.Context, id uuid.UUID, reason string) error {
	res, err := r.db.ExecContext(ctx, rejectProjectSQL, id, reason)
	if err != nil {
		return err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return domain.ErrProjectNotOnReview
	}

	return nil
}

//go:embed sql/update_project.sql
var updateProjectSQL string

func (r *ProjectRepo) UpdateProject(ctx context.Context, userID uuid.UUID, id uuid.UUID, req domain.CreateProjectReq) error {
	dbCategory, err := MapCategoryToDB(req.Category)
	if err != nil {
		return err
	}

	dbCurrency, err := MapCurrencyToDB(req.Currency)
	if err != nil {
		return err
	}

	res, err := r.db.ExecContext(ctx, updateProjectSQL,
		id, userID, dbCategory, req.Name, req.Description, dbCurrency, req.GoalAmount, req.DurationDays,
	)
	if err != nil {
		return pglib.MapConstraintErr(err, projectConstraints, err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return domain.ErrApplicationNotFound
	}

	return nil
}

//go:embed sql/get_application_by_project_id.sql
var getApplicationByProjectIDSQL string

func (r *ProjectRepo) GetApplicationByProjectID(ctx context.Context, projectID uuid.UUID) (domain.Application, error) {
	var aDB ApplicationDB
	var pDB ProjectDB

	err := r.db.QueryRowContext(ctx, getApplicationByProjectIDSQL, projectID).Scan(
		&aDB.StatusID,
		&aDB.RejectReason,
		&aDB.CreatedAt,
		&pDB.ID,
		&pDB.UserID,
		&pDB.CategoryID,
		&pDB.Name,
		&pDB.Description,
		&pDB.CurrencyID,
		&pDB.GoalAmount,
		&pDB.CurrentAmount,
		&pDB.DurationDays,
		&pDB.StatusID,
		&pDB.IsBoosted,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Application{}, domain.ErrApplicationNotFound
		}
		return domain.Application{}, err
	}

	app, err := MapApplicationFromDB(aDB, pDB)
	if err != nil {
		return domain.Application{}, err
	}

	return app, nil
}

//go:embed sql/get_pending_applications.sql
var getPendingApplicationsSQL string

func (r *ProjectRepo) GetPendingApplications(ctx context.Context) ([]domain.Application, error) {
	rows, err := r.db.QueryContext(ctx, getPendingApplicationsSQL)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.Application

	for rows.Next() {
		var aDB ApplicationDB
		var pDB ProjectDB

		if err := rows.Scan(
			&aDB.StatusID,
			&aDB.RejectReason,
			&aDB.CreatedAt,
			&pDB.ID,
			&pDB.UserID,
			&pDB.CategoryID,
			&pDB.Name,
			&pDB.Description,
			&pDB.CurrencyID,
			&pDB.GoalAmount,
			&pDB.CurrentAmount,
			&pDB.DurationDays,
			&pDB.StatusID,
			&pDB.IsBoosted,
		); err != nil {
			return nil, err
		}

		app, err := MapApplicationFromDB(aDB, pDB)
		if err != nil {
			return nil, err
		}

		out = append(out, app)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

//go:embed sql/boost_project.sql
var boostProjectSQL string

func (r *ProjectRepo) BoostProject(ctx context.Context, userID uuid.UUID, projectID uuid.UUID, days int) error {
	res, err := r.db.ExecContext(ctx, boostProjectSQL, projectID, userID, days)
	if err != nil {
		return err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return domain.ErrProjectNotFound
	}

	return nil
}
