package graph

import (
	"github.com/akemoon/crowdfunding-app-project/domain"
	"github.com/akemoon/crowdfunding-app-project/graph/model"
)

func mapApplication(a domain.Application) *model.Application {
	var rejectReason *string
	if a.RejectReason != "" {
		rejectReason = &a.RejectReason
	}
	var project *model.Project
	if a.Project != nil {
		project = mapProject(*a.Project)
	}
	return &model.Application{
		Status:       string(a.Status),
		RejectReason: rejectReason,
		CreatedAt:    a.CreatedAt,
		AssignedAt:   a.AssignedAt,
		ProcessedAt:  a.ProcessedAt,
		Project:      project,
	}
}

func mapProject(p domain.Project) *model.Project {
	return &model.Project{
		ID:            p.ID,
		UserID:        p.UserID,
		Category:      string(p.Category),
		Name:          p.Name,
		Description:   p.Description,
		Currency:      string(p.Currency),
		GoalAmount:    int32(p.GoalAmount),
		CurrentAmount: int32(p.CurrentAmount),
		DurationDays:  int32(p.DurationDays),
		Status:        string(p.Status),
		IsBoosted:     p.IsBoosted,
		StartedAt:     p.StartedAt,
		BoostedUntil:  p.BoostedUntil,
	}
}
