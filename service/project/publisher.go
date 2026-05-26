package project

import (
	"context"
	"log"
	"time"

	"github.com/akemoon/crowdfunding-app-project/metrics"
	projectPublisher "github.com/akemoon/crowdfunding-app-project/publisher/project"
	projectRepo "github.com/akemoon/crowdfunding-app-project/repo/project"
)

const (
	defaultPublishInterval = time.Minute * 5
	defaultPublishBatch    = 20
)

type PublishWorker struct {
	repo      projectRepo.Repo
	publisher *projectPublisher.Publisher
	interval  time.Duration
	batch     int
	metrics   *metrics.PublishWorkerMetrics
}

func NewPublishWorker(repo projectRepo.Repo, publisher *projectPublisher.Publisher, interval time.Duration, batch int, m *metrics.PublishWorkerMetrics) *PublishWorker {
	if interval <= 0 {
		interval = defaultPublishInterval
	}

	if batch <= 0 {
		batch = defaultPublishBatch
	}

	return &PublishWorker{
		repo:      repo,
		publisher: publisher,
		interval:  interval,
		batch:     batch,
		metrics:   m,
	}
}

func (w *PublishWorker) Run(ctx context.Context) {
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

func (w *PublishWorker) runOnce(ctx context.Context) {
	w.metrics.LastRunUnix.Set(float64(time.Now().Unix()))

	projects, err := w.repo.ListPendingFinishedOutbox(ctx, defaultPublishBatch)
	if err != nil {
		if ctx.Err() != nil {
			return
		}
		log.Printf("publisher worker error: %v", err)
		return
	}

	for _, p := range projects {
		event := projectPublisher.Event{
			Type:    projectPublisher.EventTypeFinished,
			Project: p,
		}
		err := w.publisher.Publish(ctx, event)
		if err != nil {
			log.Printf("publisher worker publish error: %s", err)
			continue
		}

		err = w.repo.MarkSentFinishedOutbox(ctx, p.ID)
		if err != nil {
			log.Printf("publisher worker mark sent error: %s", err)
			continue
		}
	}
}
