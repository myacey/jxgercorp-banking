package in

import (
	"context"

	"github.com/myacey/jxgercorp-banking/services/notification/internal/domain"
)

type NotificationUseCaese interface {
	Handle(ctx context.Context, n domain.Notification) error
}
