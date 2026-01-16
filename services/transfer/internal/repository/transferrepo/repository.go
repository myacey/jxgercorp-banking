package transferrepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
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
		ID:                  req.ID,
		FromAccountID:       req.FromAccountID,
		FromAccountUsername: req.FromAccountUsername,
		ToAccountID:         req.ToAccountID,
		ToAccountUsername:   req.ToAccountUsername,
		Amount:              req.Amount,
		CurrencyCode:        req.CurrencyCode,
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
		ID:                  res.ID,
		FromAccountID:       res.FromAccountID,
		FromAccountUsername: res.FromAccountUsername,
		ToAccountID:         res.ToAccountID,
		ToAccountUsername:   res.ToAccountUsername,
		Amount:              res.Amount,
		CreatedAt:           res.CreatedAt,
		CurrencyCode:        res.CurrencyCode,
	}, nil
}

func (r *PostgresTranser) SearchTransfersWithAccount(ctx context.Context, req *SearchTransfersWithAccountParams) ([]*entity.Transfer, error) {
	ctx, span := r.tracer.Start(ctx, "repository: SearchTransfersWithAccount")
	defer span.End()

	arg := db.SearchTransfersParams{
		CurrentAccountID: req.CurrentAccountID,
		WithUsername: pgtype.Text{
			String: req.WithUsername,
			Valid:  req.WithUsername != "",
		},
		WithAccountID: pgtype.UUID{
			Bytes: [16]byte(req.WithAccountID[:]),
			Valid: req.WithAccountID != uuid.Nil,
		},
		Currency: pgtype.Text{
			String: req.Currency,
			Valid:  req.Currency != "",
		},
		Offset: req.Offset,
		Limit:  req.Limit,
	}
	fmt.Printf("%+v", arg)

	res, err := r.store.SearchTransfers(ctx, arg)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return []*entity.Transfer{}, nil
		default:
			return nil, err
		}
	}

	transfers := make([]*entity.Transfer, len(res))
	for i, v := range res {
		transfers[i] = &entity.Transfer{
			ID:                  v.ID,
			FromAccountID:       v.FromAccountID,
			FromAccountUsername: v.FromAccountUsername,
			ToAccountID:         v.ToAccountID,
			ToAccountUsername:   v.ToAccountUsername,
			Amount:              v.Amount,
			CurrencyCode:        v.CurrencyCode,
			CreatedAt:           v.CreatedAt,
		}
	}
	return transfers, nil
}
