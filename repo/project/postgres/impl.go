package postgres

import (
	"context"
	"database/sql"
	"errors"

	_ "embed"

	"github.com/akemoon/crowdfunding-app-project/domain"
	"github.com/akemoon/golib/pglib"
	"github.com/google/uuid"
)

var projectConstraints = map[string]error{
	"projects_user_id_name_unique": domain.ErrProjectExists,
}

type ProjectRepo struct {
	db *sql.DB
}

func NewProjectRepo(db *sql.DB) *ProjectRepo {
	return &ProjectRepo{db: db}
}

//go:embed sql/create_project.sql
var createProjectSQL string

func (r *ProjectRepo) CreateProject(ctx context.Context, req domain.CreateProjectReq) (domain.Project, error) {
	categoryID, err := MapCategoryToDB(req.Category)
	if err != nil {
		return domain.Project{}, err
	}

	currencyID, err := MapCurrencyToDB(req.Currency)
	if err != nil {
		return domain.Project{}, err
	}

	var p ProjectDB
	err = r.db.QueryRowContext(ctx, createProjectSQL,
		req.UserID,
		categoryID,
		req.Name,
		req.Description,
		currencyID,
		req.GoalAmount,
		req.DurationDays,
	).Scan(
		&p.ID,
		&p.UserID,
		&p.CategoryID,
		&p.Name,
		&p.Description,
		&p.CurrencyID,
		&p.GoalAmount,
		&p.StartedAt,
		&p.DurationDays,
		&p.StatusID,
	)
	if err != nil {
		return domain.Project{}, pglib.MapConstraintErr(err, projectConstraints, err)
	}

	return MapProjectFromDB(p)
}

//go:embed sql/get_project_by_id.sql
var getProjectByIDSQL string

func (r *ProjectRepo) GetProjectByID(ctx context.Context, id uuid.UUID) (domain.Project, error) {
	var p ProjectDB

	err := r.db.QueryRowContext(ctx, getProjectByIDSQL, id).Scan(
		&p.ID,
		&p.UserID,
		&p.CategoryID,
		&p.Name,
		&p.Description,
		&p.CurrencyID,
		&p.GoalAmount,
		&p.StartedAt,
		&p.DurationDays,
		&p.StatusID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Project{}, domain.ErrProjectNotFound
		}
		return domain.Project{}, err
	}

	return MapProjectFromDB(p)
}

//go:embed sql/get_projects.sql
var getProjectsSQL string

func (r *ProjectRepo) GetProjects(ctx context.Context) ([]domain.Project, error) {
	rows, err := r.db.QueryContext(ctx, getProjectsSQL)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.Project
	for rows.Next() {
		var p ProjectDB
		if err := rows.Scan(
			&p.ID,
			&p.UserID,
			&p.CategoryID,
			&p.Name,
			&p.Description,
			&p.CurrencyID,
			&p.GoalAmount,
			&p.StartedAt,
			&p.DurationDays,
			&p.StatusID,
		); err != nil {
			return nil, err
		}
		project, err := MapProjectFromDB(p)
		if err != nil {
			return nil, err
		}
		out = append(out, project)
	}

	return out, rows.Err()
}
