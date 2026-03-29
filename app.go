package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/akemoon/crowdfunding-app-project/api"
	projectRepo "github.com/akemoon/crowdfunding-app-project/repo/project/postgres"
	projectSvc "github.com/akemoon/crowdfunding-app-project/service/project"
	pgLib "github.com/akemoon/golib/pglib"
)

type AppConfig struct {
	PostgresDSN string
	HTTPAddr    string
}

type App struct {
	ctx    context.Context
	config AppConfig

	pg *sql.DB

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

	a.InitServices()
	a.InitServer()

	return nil
}

func (a *App) InitDB() error {
	var err error

	a.pg, err = pgLib.Connect(a.ctx, a.config.PostgresDSN)
	if err != nil {
		return fmt.Errorf("postgres connection: %w", err)
	}

	return nil
}

func (a *App) InitServices() {
	repo := projectRepo.NewProjectRepo(a.pg)
	a.projectSvc = projectSvc.NewService(repo)
}

func (a *App) InitServer() {
	server := api.NewServer()
	server.AddProjectHandlers(a.projectSvc)
	a.server = server
}

func (a *App) Run() error {
	log.Printf("http server listening on %s", a.config.HTTPAddr)
	return a.server.ListenAndServe(a.config.HTTPAddr)
}
