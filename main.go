package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	envPgDSN    = "POSTGRES_DSN"
	envHTTPPort = "HTTP_PORT"
)

// @title           Project Service API
// @version         1.0
// @description     REST API for the crowdfunding projects service.
// @host            localhost:10001
// @BasePath        /
func main() {
	mainCtx, stop := context.WithCancel(context.Background())
	defer stop()

	cfg, err := loadConfigFromEnv()
	if err != nil {
		log.Fatal(err)
	}

	app := NewApp(mainCtx, cfg)

	err = app.Init()
	if err != nil {
		log.Fatalf("app init: %v", err)
	}

	err = app.Run()
	if err != nil {
		log.Fatalf("run app: %v", err)
	}
}

func loadConfigFromEnv() (AppConfig, error) {
	dsn, err := getRequiredEnv(envPgDSN)
	if err != nil {
		return AppConfig{}, err
	}

	port := strings.TrimSpace(os.Getenv(envHTTPPort))
	if port == "" {
		port = "10001"
	}

	return AppConfig{
		PostgresDSN: dsn,
		HTTPAddr:    ":" + port,
	}, nil
}

func getRequiredEnv(key string) (string, error) {
	val := strings.TrimSpace(os.Getenv(key))
	if val == "" {
		return "", fmt.Errorf("env %s is empty", key)
	}
	return val, nil
}
