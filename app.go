package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/akemoon/crowdfunding-app-project/api"
	"github.com/akemoon/crowdfunding-app-project/client/promocode/resty"
	minioStorage "github.com/akemoon/crowdfunding-app-project/client/storage/minio"
	"github.com/akemoon/crowdfunding-app-project/cluster/contribution"
	"github.com/akemoon/crowdfunding-app-project/metrics"
	"github.com/akemoon/crowdfunding-app-project/publisher/project"
	projectRepo "github.com/akemoon/crowdfunding-app-project/repo/project/postgres"
	projectSvc "github.com/akemoon/crowdfunding-app-project/service/project"
	pgLib "github.com/akemoon/golib/pglib"
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

	KafkaBrokers              []string
	ProjectTopic              string
	ContributionTopic         string
	ContributionConsumerGroup string

	// TODO: add params
	// FinishWorkerInterval  time.Duration
	// PublishWorkerInterval time.Duration
	// PublishWorkerBatch    int

	PromoBaseURL string

	MinioEndpoint  string
	MinioAccessKey string
	MinioSecretKey string
	MinioBucket    string
	MinioPublicURL string

	HTTPAddr string
}

type App struct {
	ctx    context.Context
	config AppConfig

	pg *sql.DB

	publisher *project.Publisher

	finishWorker         *projectSvc.FinishWorker
	publishWorker        *projectSvc.PublishWorker
	contributionConsumer *contribution.Consumer

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

	publishMetrics := metrics.NewPublishWorkerMetrics(prometheus.DefaultRegisterer)

	a.publishWorker = projectSvc.NewPublishWorker(
		repo,
		publisher,
		defaultPublishInterval,
		defaultPublishBatch,
		publishMetrics,
	)

	promoClient := resty.NewPromoClient(a.config.PromoBaseURL)

	storageClient, err := minioStorage.NewStorageClient(
		a.config.MinioEndpoint,
		a.config.MinioAccessKey,
		a.config.MinioSecretKey,
		a.config.MinioBucket,
		a.config.MinioPublicURL,
		false,
	)
	if err != nil {
		return fmt.Errorf("init minio: %w", err)
	}

	err = storageClient.EnsureBucket(a.ctx)
	if err != nil {
		return fmt.Errorf("ensure minio bucket: %w", err)
	}

	a.projectSvc = projectSvc.NewService(repo, promoClient, storageClient, a.config.MinioPublicURL+"/"+a.config.MinioBucket)

	a.contributionConsumer = contribution.New(
		a.config.KafkaBrokers,
		a.config.ContributionTopic,
		a.config.ContributionConsumerGroup,
		a.projectSvc,
	)

	return nil
}

func (a *App) InitServer() {
	server := api.NewServer()
	server.AddProjectHandlers(a.projectSvc)
	server.AddMetrics()
	a.server = server
}

func (a *App) Run() error {
	log.Printf("start finish worker")
	go a.finishWorker.Run(a.ctx)

	log.Printf("start publish worker")
	go a.publishWorker.Run(a.ctx)

	log.Printf("start contribution consumer")
	go a.contributionConsumer.Run(a.ctx)

	log.Printf("http server listening on %s", defaultHTTPAddr)

	a.server.ListenAndServe(defaultHTTPAddr)

	return nil
}
