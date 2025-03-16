package confirmation

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/myacey/jxgercorp-banking/services/shared/cstmerr"
	"github.com/myacey/jxgercorp-banking/services/shared/sharedmodels"
	"github.com/myacey/jxgercorp-banking/services/shared/util"
	"github.com/myacey/jxgercorp-banking/services/user/internal/repository"
	"github.com/segmentio/kafka-go"
)

type ConfirmationServiceInterface interface {
	GenerateAccountConfirmation(c context.Context, username, email string) error
	CheckConfirmCode(c context.Context, username, confirmCode string) error
}

type ConfirmationService struct {
	confirmCodesRepo repository.ConfirmCodesRepository

	registerCodesTopic   string
	registerCodesPartion int
}

func NewConfirmationService(repo repository.ConfirmCodesRepository, topic string, partion int) ConfirmationServiceInterface {
	return &ConfirmationService{
		confirmCodesRepo:     repo,
		registerCodesTopic:   topic,
		registerCodesPartion: partion,
	}
}

func (cs *ConfirmationService) GenerateAccountConfirmation(c context.Context, username, email string) error {
	confirmCode := util.RandomString(32)

	err := cs.confirmCodesRepo.CreateCode(c, username, confirmCode)
	if err != nil {
		return cstmerr.New(http.StatusInternalServerError, cstmerr.ErrUnknown.Error(), err)
	}

	msg := &sharedmodels.RegisterConfirmMsgEmail{
		Username:    username,
		Email:       email,
		ConfirmCode: confirmCode,
	}

	msgMarshalled, err := json.Marshal(msg)
	if err != nil {
		return cstmerr.New(http.StatusInternalServerError, cstmerr.ErrUnknown.Error(), err)
	}

	// kafka
	conn, err := kafka.DialLeader(c, "tcp", "localhost:9092", cs.registerCodesTopic, cs.registerCodesPartion)
	if err != nil {
		return cstmerr.New(http.StatusInternalServerError, cstmerr.ErrUnknown.Error(), err)
	}

	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	_, err = conn.WriteMessages(
		kafka.Message{Value: msgMarshalled},
	)
	if err != nil {
		return cstmerr.New(http.StatusInternalServerError, cstmerr.ErrUnknown.Error(), err)
	}

	if err = conn.Close(); err != nil {
		return cstmerr.New(http.StatusInternalServerError, cstmerr.ErrUnknown.Error(), err)
	}
	return nil
}

func (cs *ConfirmationService) CheckConfirmCode(c context.Context, username, confirmCode string) error {
	dbConfirmCode, err := cs.confirmCodesRepo.GetCode(c, username)
	log.Print("DB CONFIRMATION CODE:'", dbConfirmCode, "', HAS: '", confirmCode, "'\n")
	if err != nil || dbConfirmCode != confirmCode {
		return cstmerr.New(http.StatusBadRequest, "invalid confirm code", nil)
	}

	return nil
}
