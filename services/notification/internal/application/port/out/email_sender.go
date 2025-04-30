package out

import "context"

type EmailSender interface {
	Send(ctx context.Context, to, subject, body string) error
}
