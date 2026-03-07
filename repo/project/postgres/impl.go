package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	_ "embed"

	"github.com/akemoon/crowdfunding-app-project/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	projectsUserIDNameUnique = "projects_user_id_name_unique"
)

type ProjectRepo struct {
	db *sql.DB
}

func NewProjectRepo(db *sql.DB) *ProjectRepo {
	return &ProjectRepo{
		db: db,
	}
}

func asPostgresError(err error) *pgconn.PgError {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr
	}
	return nil
}

func mapPostgresError(err *pgconn.PgError) error {
	if err.Code == "23505" {
		switch err.ConstraintName {
		case projectsUserIDNameUnique:
			return domain.ErrProjectExists
		default:
			return domain.ErrUnknownConflict
		}
	}

	return domain.ErrInternal
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
		var pgErr *pgconn.PgError

		// Try handle as postgres err
		ok := errors.As(err, &pgErr)
		if ok {
			domainErr := domain.ErrInternal

			if pgErr.Code == "23505" && pgErr.ConstraintName == projectsUserIDNameUnique {
				domainErr = domain.ErrProjectExists
			}

			return fmt.Errorf("%w: %s", domainErr, pgErr.Detail)
		}

		return fmt.Errorf("%w: %s", domain.ErrInternal, err.Error())
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
		&pDB.StartDate,
		&pDB.DurationDays,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Project{}, fmt.Errorf("%w: project with uuid '%s' not found", domain.ErrProjectNotFound, id.String())
		}
		return domain.Project{}, fmt.Errorf("%w: %s", domain.ErrInternal, err.Error())
	}

	p, err := MapProjectFromDB(pDB)
	if err != nil {
		return domain.Project{}, err
	}

	return p, nil
}

//go:embed sql/finish_projects.sql
var finishProjectsSQL string

func (r *ProjectRepo) FinishProjects(ctx context.Context) error {
	_, err := r.db.ExecContext(
		ctx,
		finishProjectsSQL,
	)
	if err != nil {
		return fmt.Errorf("%w: %s", domain.ErrInternal, err.Error())
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
		pgErr := asPostgresError(err)
		if pgErr != nil {
			mappedErr := mapPostgresError(pgErr)
			return nil, fmt.Errorf("%w: %s", mappedErr, pgErr.Detail)
		}
		return nil, fmt.Errorf("%w: %s", domain.ErrInternal, err)
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
			&pDB.StartDate,
			&pDB.DurationDays,
		); err != nil {
			return nil, fmt.Errorf("%w: %s", domain.ErrInternal, err)
		}

		p, err := MapProjectFromDB(pDB)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", domain.ErrInternal, err)
		}

		out = append(out, p)
	}

	if err := rows.Err(); err != nil {
		pgErr := asPostgresError(err)
		if pgErr != nil {
			mappedErr := mapPostgresError(pgErr)
			return nil, fmt.Errorf("%w: %s", mappedErr, pgErr.Detail)
		}
		return nil, fmt.Errorf("%w: %s", domain.ErrInternal, err)
	}

	return out, nil
}

//go:embed sql/mark_sent_finished_outbox.sql
var markSentFinishedOutboxSQL string

func (r *ProjectRepo) MarkSentFinishedOutbox(ctx context.Context, projectID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, markSentFinishedOutboxSQL, projectID)
	if err != nil {
		pgErr := asPostgresError(err)
		if pgErr != nil {
			mappedErr := mapPostgresError(pgErr)
			return fmt.Errorf("%w: %s", mappedErr, pgErr.Detail)
		}
		return fmt.Errorf("%w: %s", domain.ErrInternal, err)
	}

	return nil
}
