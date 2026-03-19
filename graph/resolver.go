package graph

import (
	"context"

	"github.com/google/uuid"

	projectSvc "github.com/akemoon/crowdfunding-app-project/service/project"
)

type contextKey string

const (
	ctxUserID   contextKey = "userID"
	ctxUserRole contextKey = "userRole"
)

type Resolver struct {
	svc *projectSvc.Service
}

func NewResolver(svc *projectSvc.Service) *Resolver {
	return &Resolver{svc: svc}
}

func callerIDFromCtx(ctx context.Context) *uuid.UUID {
	v, ok := ctx.Value(ctxUserID).(uuid.UUID)
	if !ok {
		return nil
	}
	return &v
}

func userIDFromCtx(ctx context.Context) (uuid.UUID, bool) {
	v, ok := ctx.Value(ctxUserID).(uuid.UUID)
	return v, ok
}

func roleFromCtx(ctx context.Context) string {
	v, _ := ctx.Value(ctxUserRole).(string)
	return v
}
