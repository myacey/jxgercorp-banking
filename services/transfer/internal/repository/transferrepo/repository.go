package transferrepo

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/myacey/jxgercorp-banking/services/transfer/internal/models/entity"
	"github.com/myacey/jxgercorp-banking/services/transfer/internal/repository"
	db "github.com/myacey/jxgercorp-banking/services/transfer/internal/repository/sqlc"
)

var ErrTransferAlreadyExists = errors.New("transfer already exists")

type PostgresTranser struct {
	store *db.Queries

	tracer trace.Tracer
}

func NewPostgresTransfer(store *db.Queries) *PostgresTranser {
	return &PostgresTranser{
		store:  store,
		tracer: otel.Tracer("repository-transfer"),
	}
}

func (r *PostgresTranser) CraeteTransfer(ctx context.Context, req *entity.Transfer) (*entity.Transfer, error) {
	ctx, span := r.tracer.Start(ctx, "repository: CraeteTransfer")
	defer span.End()

	arg := db.CreateTransferParams{
		ID:            req.ID,
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}
	start := time.Now()
	res, err := r.store.CreateTransfer(ctx, arg)
	if err != nil {
		switch {
		case repository.IsUniqueViolation(err):
			return nil, ErrTransferAlreadyExists
		default:
			return nil, err
		}
	}
	span.SetAttributes(attribute.Float64("sql.time.ms", float64(time.Since(start).Microseconds())))

	return &entity.Transfer{
		ID:            res.ID,
		FromAccountID: res.FromAccountID,
		ToAccountID:   res.ToAccountID,
		Amount:        res.Amount,
		CreatedAt:     res.CreatedAt,
	}, nil
}

func (r *PostgresTranser) SearchTransfersWithAccount(ctx context.Context, req *SearchTransfersWithAccountParams) ([]*entity.Transfer, error) {
	ctx, span := r.tracer.Start(ctx, "repository: SearchTransfersWithAccount")
	defer span.End()

	arg := db.SearchTransfersWithAccountParams{
		SearchAccount: req.AccountID,
		Offset:        req.Offset,
		Limit:         req.Limit,
	}
	start := time.Now()
	res, err := r.store.SearchTransfersWithAccount(ctx, arg)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return []*entity.Transfer{}, nil
		default:
			return nil, err
		}
	}
	span.SetAttributes(attribute.Float64("sql.time.ms", float64(time.Since(start).Microseconds())))

	acc := make([]*entity.Transfer, len(res))
	for i, v := range res {
		acc[i] = &entity.Transfer{
			ID:            v.ID,
			FromAccountID: v.FromAccountID,
			ToAccountID:   v.ToAccountID,
			Amount:        v.Amount,
			CreatedAt:     v.CreatedAt,
		}
	}

	return acc, nil
}
