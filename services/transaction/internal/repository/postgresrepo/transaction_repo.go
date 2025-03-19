package postgresrepo

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/myacey/jxgercorp-banking/services/db/sqlc"
	"github.com/myacey/jxgercorp-banking/services/shared/cstmerr"
	"github.com/myacey/jxgercorp-banking/services/transaction/internal/repository"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type PostgresTransactionRepo struct {
	store  *db.Queries
	dbConn *sql.DB
	lg     *zap.SugaredLogger

	tracer trace.Tracer
}

func NewPostgresTransactionRepo(store *db.Queries, dbConn *sql.DB, lg *zap.SugaredLogger, tr trace.Tracer) repository.TransactionRepository {
	return &PostgresTransactionRepo{
		store:  store,
		dbConn: dbConn,
		lg:     lg,

		tracer: tr,
	}
}

func (r *PostgresTransactionRepo) CreateTransactionTX(c *gin.Context, fromUser, toUser string, amount int64) (*db.Transaction, error) {
	ctx, span := r.tracer.Start(c.Request.Context(), "repository: CreateTransactionTX")
	defer span.End()
	c.Request = c.Request.WithContext(ctx)
	span.SetAttributes(
		attribute.String("fromUser", fromUser),
		attribute.String("toUser", toUser),
		attribute.Int64("amount", amount),
	)

	tx, err := r.dbConn.Begin()
	if err != nil {
		return nil, cstmerr.New(http.StatusInternalServerError, cstmerr.ErrUnknown.Error(), err)
	}
	defer tx.Rollback()

	q := r.store.WithTx(tx)

	res, err := q.UpdateTwoUserBalance(c, db.UpdateTwoUserBalanceParams{
		FromUsername: fromUser,
		ToUsername:   toUser,
		Balance:      amount,
	})
	if err != nil {
		return nil, cstmerr.New(http.StatusInternalServerError, cstmerr.ErrUnknown.Error(), err)
	}

	for _, usr := range res {
		if usr.Username == fromUser && usr.Balance < 0 {
			return nil, cstmerr.New(http.StatusBadRequest, "not enough money", nil)
		}
	}

	arg := db.CreateTransactionParams{
		FromUser: fromUser,
		ToUser:   toUser,
		Amount:   amount,
	}
	transaction, err := q.CreateTransaction(c, arg)
	if err != nil {
		return nil, cstmerr.New(http.StatusInternalServerError, cstmerr.ErrUnknown.Error(), err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, cstmerr.New(http.StatusInternalServerError, cstmerr.ErrUnknown.Error(), err)
	}

	return &transaction, nil
}

func (r *PostgresTransactionRepo) SearchTransactionsWithUser(c *gin.Context, username string, offset, limit int) ([]*db.Transaction, error) {
	ctx, span := r.tracer.Start(c.Request.Context(), "repository: SearchTransactionsWithUser")
	defer span.End()
	c.Request = c.Request.WithContext(ctx)
	span.SetAttributes(
		attribute.String("username", username),
		attribute.Int("offset", offset),
		attribute.Int("limit", limit),
	)

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
