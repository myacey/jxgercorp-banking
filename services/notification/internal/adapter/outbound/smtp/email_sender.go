package smtp

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"gopkg.in/gomail.v2"
)

type Config struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type Sender struct {
	from       string
	dialerConn *gomail.Dialer

	tracer trace.Tracer
}

func NewSender(cfg Config) *Sender {
	d := gomail.NewDialer(cfg.Host, cfg.Port, cfg.Username, cfg.Password)

	return &Sender{dialerConn: d, from: cfg.Username, tracer: otel.Tracer("adapter-outbound-smtp")}
}

func (s *Sender) Send(ctx context.Context, to, subject, body string) error {
	_, span := s.tracer.Start(ctx, "adapter: SendMail")
	defer span.End()

	m := gomail.NewMessage()
	m.SetHeader("From", s.from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	return s.dialerConn.DialAndSend(m)
}
