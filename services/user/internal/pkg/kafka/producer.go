package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/segmentio/kafka-go"
)

type Config struct {
	Brokers []string `mapstructure:"brokers"`
	Topic   string   `mapstructure:"topic"`
}

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(cfg Config) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(cfg.Brokers...),
			Topic:        cfg.Topic,
			Balancer:     &kafka.LeastBytes{},
			RequiredAcks: kafka.RequireOne,
		},
	}
}

func (p *Producer) Send(ctx context.Context, val interface{}) error {
	b, err := json.Marshal(val)
	if err != nil {
		return err
	}

	msg := kafka.Message{
		Value: b,
		Time:  time.Now(),
	}

	err = p.writer.WriteMessages(ctx, msg)
	if err != nil {
		switch {
		case errors.Is(err, kafka.ErrGenerationEnded):
			return nil
		default:
			return err
		}
	}
	return p.writer.WriteMessages(ctx, msg)
}
