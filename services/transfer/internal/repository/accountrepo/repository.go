package accountrepo

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/google/uuid"
	"github.com/myacey/jxgercorp-banking/services/transfer/internal/models/entity"
	"github.com/myacey/jxgercorp-banking/services/transfer/internal/repository"
	db "github.com/myacey/jxgercorp-banking/services/transfer/internal/repository/sqlc"
)

var (
	ErrAccountAlreadyExists = errors.New("account already exists")
	ErrAccountNotFound      = errors.New("account not found")
)

type PostgresAccount struct {
	store *db.Queries

	tracer trace.Tracer
}

func NewPostgresAccount(store *db.Queries) *PostgresAccount {
	return &PostgresAccount{
		store:  store,
		tracer: otel.Tracer("repository-account"),
	}
}

func (r *PostgresAccount) CreateAccount(ctx context.Context, acc *entity.Account) (*entity.Account, error) {
	ctx, span := r.tracer.Start(ctx, "repository: CreateAccount")
	defer span.End()

	arg := db.CreateAccountParams{
		ID:            acc.ID,
		OwnerUsername: acc.OwnerUsername,
		Balance:       acc.Balance,
		Currency:      acc.Currency,
		CreatedAt:     acc.CreatedAt,
	}
	start := time.Now()
	res, err := r.store.CreateAccount(ctx, arg)
	if err != nil {
		switch {
		case repository.IsUniqueViolation(err):
			return nil, ErrAccountAlreadyExists
		default:
			return nil, err
		}
	}
	span.SetAttributes(attribute.Float64("sql.time.ms", float64(time.Since(start).Milliseconds())))

	return &entity.Account{
		ID:            res.ID,
		OwnerUsername: res.OwnerUsername,
		Balance:       res.Balance,
		Currency:      res.Currency,
		CreatedAt:     res.CreatedAt,
	}, nil
}

func (r *PostgresAccount) SearchAccounts(ctx context.Context, p *SearchAccountsParams) ([]*entity.Account, error) {
	ctx, span := r.tracer.Start(ctx, "repository: SearchAccounts")
	defer span.End()

	arg := db.SearchAccountsParams{
		Username: p.OwnerUsername,
		Currency: db.NullCurrencyEnum{CurrencyEnum: db.CurrencyEnum(p.Currency), Valid: p.Currency != ""},
	}
	start := time.Now()
	res, err := r.store.SearchAccounts(ctx, arg)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return []*entity.Account{}, nil
		default:
			return nil, err
		}
	}
	span.SetAttributes(attribute.Float64("sql.time.ms", float64(time.Since(start).Milliseconds())))

	acc := make([]*entity.Account, len(res))
	for i, v := range res {
		acc[i] = &entity.Account{
			ID:            v.ID,
			OwnerUsername: v.OwnerUsername,
			Balance:       v.Balance,
			Currency:      v.Currency,
			CreatedAt:     v.CreatedAt,
		}
	}

	return acc, nil
}

func (r *PostgresAccount) GetAccountByID(ctx context.Context, id uuid.UUID) (*entity.Account, error) {
	ctx, span := r.tracer.Start(ctx, "repository: GetAccountByID")
	defer span.End()

	start := time.Now()
	res, err := r.store.GetAccountByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, nil
		default:
			return nil, err
		}
	}
	span.SetAttributes(attribute.Float64("sql.time.ms", float64(time.Since(start).Milliseconds())))

	return &entity.Account{
		ID:            res.ID,
		OwnerUsername: res.OwnerUsername,
		Balance:       res.Balance,
		Currency:      res.Currency,
		CreatedAt:     res.CreatedAt,
	}, nil
}

func (r *PostgresAccount) AddAccountBalance(ctx context.Context, p *AddAccountBalance) (*entity.Account, error) {
	ctx, span := r.tracer.Start(ctx, "repository: AddAccountBalance")
	defer span.End()

	arg := db.AddAccountBalanceParams{
		ID:         p.AccountID,
		AddBalance: p.AddBalance,
	}
	start := time.Now()
	res, err := r.store.AddAccountBalance(ctx, arg)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrAccountNotFound
		default:
			return nil, err
		}
	}
	span.SetAttributes(attribute.Float64("sql.time.ms", float64(time.Since(start).Milliseconds())))

	return &entity.Account{
		ID:            res.ID,
		OwnerUsername: res.OwnerUsername,
		Balance:       res.Balance,
		Currency:      res.Currency,
		CreatedAt:     res.CreatedAt,
	}, nil
}

func (r *PostgresAccount) AddTwoAccountsBalance(ctx context.Context, p *AddTwoAccountsBalance) ([]*entity.Account, error) {
	ctx, span := r.tracer.Start(ctx, "repository: AddTwoAccountsBalance")
	defer span.End()

	arg := db.AddTwoAccountsBalanceParams{
		FromAccountID: p.FromAccountID,
		ToAccountID:   p.ToAccountUD,
		Balance:       p.Amount,
	}
	start := time.Now()
	res, err := r.store.AddTwoAccountsBalance(ctx, arg)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrAccountNotFound
		default:
			return nil, err
		}
	}
	span.SetAttributes(attribute.Float64("sql.time.ms", float64(time.Since(start).Milliseconds())))

	acc := make([]*entity.Account, len(res))
	for i, v := range res {
		acc[i] = &entity.Account{
			ID:            v.ID,
			OwnerUsername: v.OwnerUsername,
			Balance:       v.Balance,
			Currency:      v.Currency,
			CreatedAt:     v.CreatedAt,
		}
	}

	return acc, nil
}
