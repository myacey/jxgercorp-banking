package main

import (
	"sync"

	"github.com/myacey/jxgercorp-banking/services/notification/internal/controller"
	"github.com/myacey/jxgercorp-banking/services/notification/internal/sendmail"
	"github.com/myacey/jxgercorp-banking/services/shared/backconfig"
	"github.com/myacey/jxgercorp-banking/services/shared/logging"
	"gopkg.in/gomail.v2"
)

func main() {
	cfg, err := backconfig.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	logger, err := logging.ConfigureLogger()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	dialerConn := gomail.NewDialer("smtp.gmail.com", 587, cfg.GoogleMailAdress, cfg.GoogleAppPassword)
	mailSender := sendmail.NewMailSender(cfg.GoogleMailAdress, dialerConn, logger)

	ctrller := controller.NewController(mailSender, &controller.RegisterEmailKafka{KafkaTopic: "notif.email.register.confirm", KafkaPartion: 0}, cfg.AppDomain)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		err := ctrller.ListenEmailRegisterConfirm()
		wg.Done()
		if err != nil {
			panic(err)
		}
	}()

	wg.Wait()
}
