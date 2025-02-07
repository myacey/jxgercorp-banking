package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/myacey/jxgercorp-banking/services/notification/internal/sendmail"
	"github.com/myacey/jxgercorp-banking/services/shared/sharedmodels"
	"github.com/segmentio/kafka-go"
)

type RegisterEmailKafka struct {
	KafkaTopic   string
	KafkaPartion int

	confirmationLink string
}

// generateConfirmationLink formats a link using appDomain
func (re *RegisterEmailKafka) generateConfirmationLink(appDomain string) {
	re.confirmationLink = appDomain + "/#/user/confirm?username=%s&code=%s"
}

type Controller struct {
	emailSender sendmail.MainSenderInterface

	registerEmailKafka *RegisterEmailKafka
}

func NewController(es sendmail.MainSenderInterface, re *RegisterEmailKafka, appDomain string) *Controller {
	re.generateConfirmationLink(appDomain)
	return &Controller{
		emailSender:        es,
		registerEmailKafka: re,
	}
}

// ListemEmailRegistrConfirm starts listen to Kafka notif.email.register.confirm
func (h *Controller) ListenEmailRegisterConfirm() error {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{"localhost:9092"},
		Topic:     "notif.email.register.confirm",
		Partition: 0,
		MaxBytes:  10e6, // 10MB
	})
	r.SetOffsetAt(context.Background(), time.Now()) // start listen from now

	msgFormat := `
	<html>
		<body style=\"font-family: Arial, sans-serif;\">
			<p>Visit the following link to confirm your account:</p>
			<a href=\"%s\" style=\"color: #7b3f98; font-weight: bold;\">%s</a>
			<p>If you did not request this, please ignore this email.</p>
		</body>
	</html>
	`

	log.Print("START")
	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			break
		}
		log.Printf("message at offset %d: %s=%s\n", m.Offset, string(m.Key), string(m.Value))

		var msg sharedmodels.RegisterConfirmMsgEmail
		if err := json.Unmarshal(m.Value, &msg); err != nil {
			log.Fatal("cant unmarshal err:", err)
		}

		lnk := fmt.Sprintf(h.registerEmailKafka.confirmationLink, msg.Username, msg.ConfirmCode)

		if err := h.emailSender.SendMail(
			[]string{msg.Email},
			"Complete Registration Process",
			fmt.Sprintf(msgFormat, lnk, lnk),
			&sendmail.SendMailParams{
				BodyType: "text/html",
			},
		); err != nil {
			log.Fatal("cant send mail:", err)
		}
	}

	if err := r.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
	return nil
}
