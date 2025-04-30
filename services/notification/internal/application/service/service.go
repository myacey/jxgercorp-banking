package service

import (
	"context"

	"github.com/myacey/jxgercorp-banking/services/notification/internal/application/port/in"
	"github.com/myacey/jxgercorp-banking/services/notification/internal/application/port/out"
	"github.com/myacey/jxgercorp-banking/services/notification/internal/domain"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type Notification struct {
	EmailSender out.EmailSender

	tracer trace.Tracer
}

func NewNotificationService(sender out.EmailSender) in.NotificationUseCaese {
	return &Notification{EmailSender: sender, tracer: otel.Tracer("service")}
}

func (s *Notification) Handle(ctx context.Context, n domain.Notification) error {
	ctx, span := s.tracer.Start(ctx, "service: HandleSMTPNotification")
	defer span.End()

	subject := n.Subject
	body := n.Text

	return s.EmailSender.Send(ctx, n.Email, subject, body)
}
