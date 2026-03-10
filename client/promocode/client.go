package promocode

import (
	"context"

	"github.com/google/uuid"
)

type Client interface {
	UsePromoCode(ctx context.Context, userID uuid.UUID, code string, expectedType string) (int, error)
}
