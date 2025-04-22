package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/myacey/jxgercorp-banking/services/libs/apperror"
	"github.com/myacey/jxgercorp-banking/services/libs/sharedmodels"
	"github.com/myacey/jxgercorp-banking/services/libs/util"
	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type ConfirmationKafkaConfig struct {
	Network      string `mapstructure:"network"`
	ListenAdress string `mapstructure:"adress"`
	Topic        string `mapstructure:"topic"`
	Partion      int    `mapstructure:"partion"`
}

type ConfirmRepo interface {
	CreateCode(ctx context.Context, username string, code string) error
	GetCode(ctx context.Context, username string) (string, error)
}

type Confirmation struct {
	repo ConfirmRepo
	cfg  ConfirmationKafkaConfig

	tracer trace.Tracer
}

func NewConfirmationService(repo ConfirmRepo, cfg ConfirmationKafkaConfig) *Confirmation {
	return &Confirmation{
		repo:   repo,
		cfg:    cfg,
		tracer: otel.Tracer("service-confirmation"),
	}
}

// GenerateAccountConfirmation generates confirm code and sends
// it to kafka. Returns apperror
func (cs *Confirmation) generateAccountConfirmation(c context.Context, username, email string) error {
	c, span := cs.tracer.Start(c, "email-confirmation: GenerateAccountConfirmation")
	defer span.End()

	confirmCode := util.RandomString(32)

	err := cs.repo.CreateCode(c, username, confirmCode)
	if err != nil {
		return apperror.NewInternal("failed to create confirm code", err)
	}

	msg := &sharedmodels.RegisterConfirmMsgEmail{
		Username:    username,
		Email:       email,
		ConfirmCode: confirmCode,
	}

	msgMarshalled, err := json.Marshal(msg)
	if err != nil {
		return apperror.NewInternal("failed to create confirm code", err)
	}

	// kafka
	conn, err := kafka.DialLeader(c, cs.cfg.Network, cs.cfg.ListenAdress, cs.cfg.Topic, cs.cfg.Partion)
	if err != nil {
		return apperror.NewInternal("failed to create confirm code", err)
	}

	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	_, err = conn.WriteMessages(
		kafka.Message{Value: msgMarshalled},
	)
	if err != nil {
		return apperror.NewInternal("failed to craete confirm code", err)
	}

	if err = conn.Close(); err != nil {
		return apperror.NewInternal("failed to craete confirm code", err)
	}
	return nil
}

// CheckConfirmCode checks confirm code
func (cs *Confirmation) checkConfirmCode(c context.Context, username, confirmCode string) error {
	c, span := cs.tracer.Start(c, "email-confirmation: CheckConfirmCode")
	defer span.End()

	dbConfirmCode, err := cs.repo.GetCode(c, username)
	// log.Print("DB CONFIRMATION CODE:'", dbConfirmCode, "', HAS: '", confirmCode, "'\n")
	if err != nil || dbConfirmCode != confirmCode {
		return apperror.NewBadReq("invalid confirm code")
	}

	return nil
}
