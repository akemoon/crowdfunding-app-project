package project

import (
	"context"
	"log"
	"time"

	"github.com/akemoon/crowdfunding-app-project/metrics"
	"github.com/akemoon/crowdfunding-app-project/repo/project"
)

const (
	defaultFinishInterval = time.Minute * 5
)

type FinishWorker struct {
	repo     project.Repo
	interval time.Duration
	metrics  *metrics.FinishWorkerMetrics
}

func NewFinishWorker(repo project.Repo, interval time.Duration, m *metrics.FinishWorkerMetrics) *FinishWorker {
	if interval <= 0 {
		interval = defaultFinishInterval
	}

	return &FinishWorker{
		repo:     repo,
		interval: interval,
		metrics:  m,
	}
}

func (w *FinishWorker) Run(ctx context.Context) {
	w.runOnce(ctx)

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.runOnce(ctx)
		}
	}
}

func (w *FinishWorker) runOnce(ctx context.Context) {
	w.metrics.LastRunUnix.Set(float64(time.Now().Unix()))

	err := w.repo.FinishProjects(ctx)
	if err != nil {
		if ctx.Err() != nil {
			return
		}
		log.Printf("finish worker error: %v", err)
	}
}
