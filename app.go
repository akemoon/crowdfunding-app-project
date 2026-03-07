package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/akemoon/crowdfunding-app-project/api"
	"github.com/akemoon/crowdfunding-app-project/metrics"
	"github.com/akemoon/crowdfunding-app-project/publisher/project"
	projectRepo "github.com/akemoon/crowdfunding-app-project/repo/project/postgres"
	projectSvc "github.com/akemoon/crowdfunding-app-project/service/project"
	pgLib "github.com/akemoon/golib/postgres"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	defaultHTTPAddr = ":80"

	defaultFinishInterval  = time.Minute * 5
	defaultPublishInterval = time.Minute * 5
	defaultPublishBatch    = 20
)

type AppConfig struct {
	PostgresDSN           string
	PostgresMigrationsDir string

	KafkaBrokers []string
	ProjectTopic string

	// TODO: add params
	// FinishWorkerInterval  time.Duration
	// PublishWorkerInterval time.Duration
	// PublishWorkerBatch    int

	HTTPAddr string
}

type App struct {
	ctx    context.Context
	config AppConfig

	pg *sql.DB

	publisher *project.Publisher

	finishWorker  *projectSvc.FinishWorker
	publishWorker *projectSvc.PublishWorker

	projectSvc *projectSvc.Service

	server *api.Server
}

func NewApp(ctx context.Context, config AppConfig) *App {
	return &App{
		ctx:    ctx,
		config: config,
	}
}

func (a *App) Init() error {
	err := a.InitDB()
	if err != nil {
		return fmt.Errorf("init db: %w", err)
	}

	err = a.InitServices()
	if err != nil {
		return fmt.Errorf("init services: %w", err)
	}

	a.InitServer()

	return nil
}

func (a *App) InitDB() error {
	var err error

	a.pg, err = pgLib.Connect(a.ctx, a.config.PostgresDSN)
	if err != nil {
		return fmt.Errorf("postgres connection: %w", err)
	}

	err = pgLib.Migrate(a.ctx, a.pg, a.config.PostgresMigrationsDir)
	if err != nil {
		return fmt.Errorf("postgres migration: %w", err)
	}

	return nil
}

func (a *App) InitServices() error {
	repo := projectRepo.NewProjectRepo(a.pg)

	finishMetrics := metrics.NewFinishWorkerMetrics(prometheus.DefaultRegisterer)

	a.finishWorker = projectSvc.NewFinishWorker(
		repo,
		defaultFinishInterval,
		finishMetrics,
	)

	// TODO: change name
	publisher, err := project.NewPublisher(a.config.KafkaBrokers, a.config.ProjectTopic)
	if err != nil {
		return fmt.Errorf("init kafka publisher: %w", err)
	}
	a.publisher = publisher

	a.publishWorker = projectSvc.NewPublishWorker(
		repo,
		publisher,
		defaultPublishInterval,
		defaultPublishBatch,
	)

	a.projectSvc = projectSvc.NewService(repo)

	return nil
}

func (a *App) InitServer() {
	server := api.NewServer()
	server.AddProjectHandlers(a.projectSvc)
	server.AddSwaggerUI()
	server.AddMetrics()
	a.server = server
}

func (a *App) Run() error {
	log.Printf("start finish worker")
	go a.finishWorker.Run(a.ctx)

	log.Printf("start publish worker")
	go a.publishWorker.Run(a.ctx)

	log.Printf("http server listening on %s", defaultHTTPAddr)

	a.server.ListenAndServe(defaultHTTPAddr)

	return nil
}
