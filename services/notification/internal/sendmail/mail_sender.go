package sendmail

import (
	"net/http"

	"github.com/myacey/jxgercorp-banking/services/shared/cstmerr"
	"go.uber.org/zap"
	"gopkg.in/gomail.v2"
)

type MainSenderInterface interface {
	SendMail(to []string, subject, body string, parametrs *SendMailParams) error
}

type MailSender struct {
	fromAdress string
	dialerConn *gomail.Dialer
	lg         *zap.SugaredLogger
}

func NewMailSender(fromAdress string, conn *gomail.Dialer, lg *zap.SugaredLogger) MainSenderInterface {
	return &MailSender{
		fromAdress: fromAdress,
		dialerConn: conn,
		lg:         lg,
	}
}

type SendMailParams struct {
	BodyType string   // can be 'text/html', 'text/plain' etc...
	Attach   []string // paths for files
}

func (ms *MailSender) SendMail(to []string, subject, body string, parametrs *SendMailParams) error {
	msg := gomail.NewMessage()

	bodyType := "text/html"
	var attach []string
	if parametrs != nil {
		if parametrs.BodyType != "" {
			bodyType = parametrs.BodyType
		}

		if len(parametrs.Attach) != 0 {
			attach = parametrs.Attach
		}
	}

	msg.SetHeader("From", ms.fromAdress)
	msg.SetHeader("To", to...)
	msg.SetHeader("Subject", subject)
	msg.SetBody(bodyType, body)
	for _, v := range attach {
		msg.Attach(v)
	}

	if err := ms.dialerConn.DialAndSend(msg); err != nil {
		return cstmerr.New(http.StatusInternalServerError, "cant send mail", err)
	}

	return nil
}
