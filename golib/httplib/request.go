package httplib

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

// TODO: move to golib

func ParseUUIDHeader(r *http.Request, header string) (uuid.UUID, error) {
	val := r.Header.Get(header)
	id, err := uuid.Parse(val)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("parse uuid header %q: %w", header, err)
	}
	return id, nil
}