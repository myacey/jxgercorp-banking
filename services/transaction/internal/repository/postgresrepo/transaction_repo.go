package postgresrepo

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/myacey/jxgercorp-banking/services/db/sqlc"
	"github.com/myacey/jxgercorp-banking/services/shared/cstmerr"
	"github.com/myacey/jxgercorp-banking/services/transaction/internal/repository"
	"go.uber.org/zap"
)

type PostgresTransactionRepo struct {
	store  *db.Queries
	dbConn *sql.DB
	lg     *zap.SugaredLogger
}

func NewPostgresTransactionRepo(store *db.Queries, dbConn *sql.DB, lg *zap.SugaredLogger) repository.TransactionRepository {
	return &PostgresTransactionRepo{
		store:  store,
		dbConn: dbConn,
		lg:     lg,
	}
}

func (r *PostgresTransactionRepo) CreateTransactionTX(c *gin.Context, fromUser, toUser string, amount int64) (*db.Transaction, error) {
	tx, err := r.dbConn.Begin()
	if err != nil {
		log.Print("1")
		return nil, cstmerr.New(http.StatusInternalServerError, cstmerr.ErrUnknown.Error(), err)
	}
	defer tx.Rollback()

	q := r.store.WithTx(tx)

	usr, err := q.GetUserByUsername(c, fromUser)
	if err != nil {
		log.Print("2")
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, cstmerr.New(http.StatusBadRequest, "no user found", nil)
		default:
			return nil, cstmerr.New(http.StatusInternalServerError, cstmerr.ErrUnknown.Error(), err)
		}
	}

	if usr.Balance-amount < 0 {
		log.Print("3")
		return nil, cstmerr.New(http.StatusBadRequest, "not enough money", nil)
	}

	arg := db.CreateTransactionParams{
		FromUser: fromUser,
		ToUser:   toUser,
		Amount:   amount,
	}
	transaction, err := q.CreateTransaction(c, arg)
	if err != nil {
		log.Print("4")
		return nil, cstmerr.New(http.StatusInternalServerError, cstmerr.ErrUnknown.Error(), err)
	}

	fromUserBalanceArg := db.ChangeUserBalanceParams{
		Username:   fromUser,
		AddBalance: -amount,
	}
	_, err = q.ChangeUserBalance(c, fromUserBalanceArg)
	if err != nil {
		log.Print("5")
		return nil, cstmerr.New(http.StatusInternalServerError, cstmerr.ErrUnknown.Error(), err)
	}

	toUserBalanceArg := db.ChangeUserBalanceParams{
		Username:   toUser,
		AddBalance: +amount,
	}
	_, err = q.ChangeUserBalance(c, toUserBalanceArg)
	if err != nil {
		log.Print("6")
		return nil, cstmerr.New(http.StatusInternalServerError, cstmerr.ErrUnknown.Error(), err)
	}

	err = tx.Commit()
	if err != nil {
		log.Print("7")
		return nil, cstmerr.New(http.StatusInternalServerError, cstmerr.ErrUnknown.Error(), err)
	}

	return &transaction, nil
}

func (r *PostgresTransactionRepo) SearchTransactionsWithUser(c *gin.Context, username string, offset, limit int) ([]*db.Transaction, error) {
	arg := db.SearchTransactionsWithUserParams{
		SearchUser: username,
		Offset:     int32(offset),
		Limit:      int32(limit),
	}

	dbTrxs, err := r.store.SearchTransactionsWithUser(c, arg)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return []*db.Transaction{}, nil
		default:
			return nil, cstmerr.New(http.StatusInternalServerError, cstmerr.ErrUnknown.Error(), err)
		}
	}

	dbTrxsC := make([]*db.Transaction, len(dbTrxs))
	for i := range dbTrxs {
		dbTrxsC[i] = &dbTrxs[i]
	}

	return dbTrxsC, nil
}
