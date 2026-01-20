package kafka

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/myacey/jxgercorp-banking/services/notification/internal/application/port/in"
	"github.com/myacey/jxgercorp-banking/services/notification/internal/domain"
	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type Config struct {
	Brokers []string `mapstruct:"brokers"`
	GroupID string   `mapstructure:"group_id"`
	Topic   string   `mapstructure:"topic"`
}

type Consumer struct {
	reader  *kafka.Reader
	handler in.NotificationUseCaese

	tracer trace.Tracer
}

func NewConsumer(cfg Config, handler in.NotificationUseCaese) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: cfg.Brokers,
		GroupID: cfg.GroupID,
		Topic:   cfg.Topic,
	})

	return &Consumer{reader: reader, handler: handler, tracer: otel.Tracer("adapter-inbound-consumer")}
}

func (c *Consumer) Start(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Println("consumer context canceled, stopping")
				return
			default:
			}

			ctx, span := c.tracer.Start(ctx, "adapter: ProceedKafkaMsg")

			m, err := c.reader.ReadMessage(ctx)
			if err != nil {
			}

			if err != nil {
				if err == context.Canceled {
					log.Println("consumer stopped by context")
					return
				}

				log.Printf("kafka read error: %v", err)
				time.Sleep(time.Second) // backoff
				continue
			}

			var n domain.Notification
			if err := json.Unmarshal(m.Value, &n); err != nil {
				log.Printf("invalid message: %v", err)
				span.End()
				continue
			}
			log.Printf("got message: %v", n.Subject)
			log.Printf("message data: %v", n)
			if err := c.handler.Handle(ctx, n); err != nil {
				log.Printf("handler error: %v", err)
			}

			span.End()
		}
	}()
}

func (c *Consumer) Stop() {
	c.reader.Close()
}
