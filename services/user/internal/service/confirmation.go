package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/myacey/jxgercorp-banking/services/libs/apperror"
	"github.com/myacey/jxgercorp-banking/services/libs/util"
	"github.com/myacey/jxgercorp-banking/services/user/internal/models/entity"
)

type ConfirmRepo interface {
	CreateCode(ctx context.Context, username string, code string) error
	GetCode(ctx context.Context, username string) (string, error)
}

type ConfirmationSender interface {
	Send(ctx context.Context, val interface{}) error
}

type Confirmation struct {
	repo   ConfirmRepo
	sender ConfirmationSender

	tracer trace.Tracer
}

func NewConfirmationService(repo ConfirmRepo, sender ConfirmationSender) *Confirmation {
	return &Confirmation{
		repo:   repo,
		sender: sender,

		tracer: otel.Tracer("service-confirmation"),
	}
}

// GenerateAccountConfirmation generates confirm code and sends
// it to kafka. Returns apperror
func (cs *Confirmation) generateAccountConfirmation(ctx context.Context, username, email string) error {
	ctx, span := cs.tracer.Start(ctx, "email-confirmation: GenerateAccountConfirmation")
	defer span.End()

	confirmCode := util.RandomString(32)

	err := cs.repo.CreateCode(ctx, username, confirmCode)
	if err != nil {
		return apperror.NewInternal("failed to create confirmation code", err)
	}

	link := fmt.Sprintf("localhost:80/api/v1/user/confirm?username=%s&code=%s", username, confirmCode)

	msg := fmt.Sprintf("Dear %s.\nTo confirm your account, proceed the link: %v", username, link)
	n := entity.Notification{
		ID:        uuid.New(),
		Username:  username,
		Email:     email,
		Type:      "mail",
		Subject:   "Account Confirmation",
		Text:      msg,
		CreatedAt: time.Now(),
	}

	start := time.Now()
	err = cs.sender.Send(ctx, n)
	if err != nil {
		return apperror.NewInternal("failed to send confirmation", err)
	}
	span.SetAttributes(attribute.Float64("kafka.creation.time.ms", float64(time.Since(start).Milliseconds())))

	return nil
}

// CheckConfirmCode checks confirm code
func (cs *Confirmation) checkConfirmCode(ctx context.Context, username, confirmCode string) error {
	c, span := cs.tracer.Start(ctx, "email-confirmation: CheckConfirmCode")
	defer span.End()

	dbConfirmCode, err := cs.repo.GetCode(c, username)
	// log.Print("DB CONFIRMATION CODE:'", dbConfirmCode, "', HAS: '", confirmCode, "'\n")
	if err != nil || dbConfirmCode != confirmCode {
		return apperror.NewBadReq("invalid confirm code")
	}

	return nil
}
