package pglib

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

// TODO: move to golib

// MapConstraintErr maps a pg constraint violation to a domain error using the provided map.
// Returns fallback if the error is not a pg error or the constraint is not in the map.
func MapConstraintErr(err error, constraints map[string]error, fallback error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if mapped, ok := constraints[pgErr.ConstraintName]; ok {
			return mapped
		}
	}
	return fallback
}