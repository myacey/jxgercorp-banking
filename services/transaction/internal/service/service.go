package service

import (
	"github.com/gin-gonic/gin"
	"github.com/myacey/jxgercorp-banking/services/shared/converters"
	"github.com/myacey/jxgercorp-banking/services/shared/sharedmodels"
	"github.com/myacey/jxgercorp-banking/services/transaction/internal/repository"
	"go.uber.org/zap"
)

type ServiceInterface interface {
	CreateNewTransaction(c *gin.Context, fromUser, toUser string, amount int64) (*sharedmodels.Trx, error)
	// SearchTransactions(c *gin.Context, username string, offset int, limit int) ([]*sharedmodels.Trx, error)
	SearchEntriesForUser(c *gin.Context, username string, offset int, limit int) ([]*sharedmodels.Entry, error)
}

type Service struct {
	trxRepo repository.TransactionRepository
	lg      *zap.SugaredLogger
}

func NewService(tr repository.TransactionRepository, lg *zap.SugaredLogger) ServiceInterface {
	return &Service{
		trxRepo: tr,
		lg:      lg,
	}
}

func (s *Service) CreateNewTransaction(c *gin.Context, fromUser, toUser string, amount int64) (*sharedmodels.Trx, error) {
	dbTrx, err := s.trxRepo.CreateTransactionTX(c, fromUser, toUser, amount)
	if err != nil {
		return nil, err
	}

	trx := converters.ConvertDBTrxToModel(dbTrx)
	return trx, nil
}

func (s *Service) SearchEntriesForUser(c *gin.Context, username string, offset int, limit int) ([]*sharedmodels.Entry, error) {
	dbTrxs, err := s.trxRepo.SearchTransactionsWithUser(c, username, offset, limit)
	if err != nil {
		return nil, err
	}

	// trxs := make([]*sharedmodels.Trx, len(dbTrxs))
	// for i := range dbTrxs {
	// 	trxs[i] = converters.ConvertDBTrxToModel(dbTrxs[i])
	// }

	entries := make([]*sharedmodels.Entry, len(dbTrxs))
	for i := range dbTrxs {
		withUser := ""
		amount := int64(0)
		if dbTrxs[i].FromUser == username {
			withUser = dbTrxs[i].ToUser
			amount = -dbTrxs[i].Amount
		} else {
			withUser = dbTrxs[i].FromUser
			amount = +dbTrxs[i].Amount
		}
		entries[i] = &sharedmodels.Entry{
			WithUser:  withUser,
			Amount:    amount,
			CreatedAt: dbTrxs[i].CreatedAt,
		}
	}

	return entries, nil
}
