package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/myacey/jxgercorp-banking/services/notification/internal/adapter/inbound/kafka"
	"github.com/myacey/jxgercorp-banking/services/notification/internal/adapter/outbound/smtp"
	"github.com/myacey/jxgercorp-banking/services/notification/internal/application/service"
	"github.com/myacey/jxgercorp-banking/services/notification/internal/config"
)

var cfgPath = flag.String("f", "./configs/config.yaml", "path to the transfer config")

func main() {
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// config
	cfg, err := config.LoadConfig(*cfgPath)
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	sender := smtp.NewSender(cfg.SMTPConfig)
	usecase := service.NewNotificationService(sender)
	consumer := kafka.NewConsumer(cfg.KafkaConfig, usecase)

	consumer.Start(ctx)
	log.Println("notification service started...")

	<-ctx.Done()
	consumer.Stop()
	log.Println("notification service stopped...")
}
