package contribution

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/akemoon/crowdfunding-app-project/service/project"
	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
	svc    *project.Service
}

func New(brokers []string, topic, groupID string, svc *project.Service) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: groupID,
	})
	return &Consumer{reader: reader, svc: svc}
}

func (c *Consumer) Run(ctx context.Context) error {
	for {
		msg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return nil
			}
			log.Printf("contribution consumer fetch: %v", err)
			continue
		}

		err = c.handle(ctx, msg)
		if err != nil {
			log.Printf("contribution consumer %v", err)
		}

		err = c.reader.CommitMessages(ctx, msg)
		if err != nil {
			log.Printf("contribution consumer commit: %v", err)
		}
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}

func (c *Consumer) handle(ctx context.Context, msg kafka.Message) error {
	var event Event
	err := json.Unmarshal(msg.Value, &event)
	if err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}

	if event.Type != EventTypeCreated {
		return nil
	}

	err = c.svc.AddContribution(ctx, event.ProjectID, event.Amount)
	if err != nil {
		return fmt.Errorf("add contribution: %w", err)
	}

	log.Printf("contribution consumer: added %d to project %s", event.Amount, event.ProjectID)

	return nil
}
