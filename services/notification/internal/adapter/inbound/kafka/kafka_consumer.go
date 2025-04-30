package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/myacey/jxgercorp-banking/services/notification/internal/application/port/in"
	"github.com/myacey/jxgercorp-banking/services/notification/internal/domain"
	"github.com/segmentio/kafka-go"
)

type Config struct {
	Brokers []string `mapstruct:"brokers"`
	GroupID string   `mapstructure:"group_id"`
	Topic   string   `mapstructure:"topic"`
}

type Consumer struct {
	reader  *kafka.Reader
	handler in.NotificationUseCaese
}

func NewConsumer(cfg Config, handler in.NotificationUseCaese) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: cfg.Brokers,
		GroupID: cfg.GroupID,
		Topic:   cfg.Topic,
	})

	return &Consumer{reader: reader, handler: handler}
}

func (c *Consumer) Start(ctx context.Context) {
	go func() {
		for {
			m, err := c.reader.ReadMessage(ctx)
			if err != nil {
				log.Printf("ERROR: kafka read: %v", err)
				continue
			}

			var n domain.Notification
			if err := json.Unmarshal(m.Value, &n); err != nil {
				log.Printf("ERROR: Invalid message: %v", err)
				continue
			}
			log.Printf("got message: %v", n)
			if err := c.handler.Handle(ctx, n); err != nil {
				log.Printf("ERROR: Handler error: %v", err)
			}
		}
	}()
}

func (c *Consumer) Stop() {
	c.reader.Close()
}
