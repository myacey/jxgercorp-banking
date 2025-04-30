package service

import (
	"context"

	"github.com/myacey/jxgercorp-banking/services/notification/internal/application/port/in"
	"github.com/myacey/jxgercorp-banking/services/notification/internal/application/port/out"
	"github.com/myacey/jxgercorp-banking/services/notification/internal/domain"
)

type Notification struct {
	EmailSender out.EmailSender
}

func NewNotificationService(sender out.EmailSender) in.NotificationUseCaese {
	return &Notification{EmailSender: sender}
}

func (s *Notification) Handle(ctx context.Context, n domain.Notification) error {
	subject := n.Subject
	body := n.Text

	return s.EmailSender.Send(ctx, n.Email, subject, body)
}
