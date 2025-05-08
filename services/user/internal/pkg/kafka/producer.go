package kafka

import (
	"context"
	"encoding/json"
	"time"

	"github.com/segmentio/kafka-go"
)

type Config struct {
	Network string `mapstructure:"network"`
	Broker  string `mapstructure:"broker"`
	Topic   string `mapstructure:"topic"`
	Partion int    `mapstructure:"partion"`

	WriteDeadline int `mapstructure:"write_deadline"`
}

type Producer struct {
	cfg Config
}

func NewProducer(cfg Config) *Producer {
	return &Producer{cfg}
}

func (p *Producer) Send(ctx context.Context, val interface{}) error {
	data, err := json.Marshal(val)
	if err != nil {
		return err
	}

	conn, err := kafka.DialLeader(ctx, p.cfg.Network, p.cfg.Broker, p.cfg.Topic, p.cfg.Partion)
	if err != nil {
		return err
	}
	defer conn.Close()

	conn.SetWriteDeadline(time.Now().Add(time.Millisecond * time.Duration(p.cfg.WriteDeadline)))
	msg := kafka.Message{
		Value: data,
	}
	_, err = conn.WriteMessages(msg)
	if err != nil {
		return err
	}

	return nil
}
