package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/akemoon/crowdfunding-app-project/docs"
)

const (
	envPgDSN           = "POSTGRES_DSN"
	envPgMigrationsDir = "POSTGRES_MIGRATIONS_DIR"

	envKafkaBrokers = "KAFKA_BROKERS"
	envProjectTopic = "PROJECT_TOPIC"

	// TODO: add envs
	// envFinishWorkerInterval  = "FINISH_WORKER_INTERVAL"
	// envPublishWorkerInterval = "PUBLISH_WORKER_INTERVAL"
	// envPublishWorkerBatch    = "PUBLISH_WORKER_BATCH"
	// envHTTPAddr = "HTTP_ADDR"
)

// @title Project service API
// @version	1.0
// @description	Projects service API for a crowdfunding app.
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

	// TODO: graceful shutdown...
}

func loadConfigFromEnv() (AppConfig, error) {
	dsn, err := getRequiredEnv(envPgDSN)
	if err != nil {
		return AppConfig{}, err
	}

	migrationsDir, err := getRequiredEnv(envPgMigrationsDir)
	if err != nil {
		return AppConfig{}, err
	}

	brokersValue, err := getRequiredEnv(envKafkaBrokers)
	if err != nil {
		return AppConfig{}, err
	}

	brokers := parseBrokers(brokersValue)
	if len(brokers) == 0 {
		return AppConfig{}, fmt.Errorf("env %s is empty", envKafkaBrokers)
	}

	topic, err := getRequiredEnv(envProjectTopic)
	if err != nil {
		return AppConfig{}, err
	}

	return AppConfig{
		PostgresDSN:           dsn,
		PostgresMigrationsDir: migrationsDir,
		KafkaBrokers:          brokers,
		ProjectTopic:          topic,
	}, nil
}

func getRequiredEnv(key string) (string, error) {
	val := strings.TrimSpace(os.Getenv(key))
	if val == "" {
		return "", fmt.Errorf("env %s is empty", key)
	}
	return val, nil
}

func parseBrokers(value string) []string {
	parts := strings.Split(value, ",")

	brokers := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			brokers = append(brokers, trimmed)
		}
	}

	return brokers
}
