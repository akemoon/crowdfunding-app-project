package project

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
)

type Publisher struct {
	writer *kafka.Writer
}

func NewPublisher(brokers []string, topic string) (*Publisher, error) {
	if len(brokers) == 0 {
		return nil, fmt.Errorf("empty kafka brokers")
	}
	if topic == "" {
		return nil, fmt.Errorf("empty kafka topic")
	}

	writer := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireAll,
	}

	return &Publisher{writer: writer}, nil
}

func (p *Publisher) Publish(ctx context.Context, event Event) error {
	if event.Type == EventTypeUnknown {
		return fmt.Errorf("unknown event type")
	}

	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}

	key := []byte(event.Project.ID.String())
	msg := kafka.Message{
		Key:   key,
		Value: payload,
	}

	return p.writer.WriteMessages(ctx, msg)
}

func (p *Publisher) Close() error {
	err := p.writer.Close()
	p.writer = nil
	return err
}
