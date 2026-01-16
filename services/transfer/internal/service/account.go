package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"github.com/myacey/jxgercorp-banking/services/libs/apperror"
	"github.com/myacey/jxgercorp-banking/services/transfer/internal/models/dto/request"
	"github.com/myacey/jxgercorp-banking/services/transfer/internal/models/entity"
	"github.com/myacey/jxgercorp-banking/services/transfer/internal/repository/accountrepo"
)

// start balance of user account.
const accountStartBalance = 100000

type AccountRepo interface {
	CreateAccount(ctx context.Context, acc *entity.Account) (*entity.Account, error)
	GetAccountByID(ctx context.Context, id uuid.UUID) (*entity.Account, error)
	SearchAccounts(ctx context.Context, p *accountrepo.SearchAccountsParams) ([]*entity.Account, error)
	DeleteAccount(ctx context.Context, id uuid.UUID) error

	AddTwoAccountsBalance(ctx context.Context, p *accountrepo.AddTwoAccountsBalance) ([]*entity.Account, error)

	GetCurrencies(ctx context.Context) ([]*entity.Currency, error)
}

type Account struct {
	repo AccountRepo

	tracer trace.Tracer
}

func NewAccount(repo AccountRepo) *Account {
	return &Account{
		repo:   repo,
		tracer: otel.Tracer("service-account"),
	}
}

func (s *Account) CreateAccount(ctx context.Context, req *request.CreateAccount) (*entity.Account, error) {
	ctx, span := s.tracer.Start(ctx, "service: CreateAccount")
	defer span.End()

	arg := &entity.Account{
		ID:            uuid.New(),
		OwnerUsername: req.OwnerUsername,
		Balance:       accountStartBalance,
		CurrencyCode:  req.Currency,
		CreatedAt:     time.Now(),
	}
	account, err := s.repo.CreateAccount(ctx, arg)
	if err != nil {
		switch {
		case errors.Is(err, accountrepo.ErrAccountAlreadyExists):
			return nil, apperror.NewBadReq(accountrepo.ErrAccountAlreadyExists.Error())
		default:
			return nil, apperror.NewInternal("failed to create account", err)
		}
	}

	return account, nil
}

func (s *Account) SearchAccounts(ctx context.Context, req *request.SearchAccounts) ([]*entity.Account, error) {
	ctx, span := s.tracer.Start(ctx, "service: SearchAccounts")
	defer span.End()

	arg := accountrepo.SearchAccountsParams{
		OwnerUsername: req.Username,
		Currency:      req.Currency,
	}
	accounts, err := s.repo.SearchAccounts(ctx, &arg)
	if err != nil {
		switch {
		case errors.Is(err, accountrepo.ErrAccountAlreadyExists):
			return nil, apperror.NewBadReq(accountrepo.ErrAccountAlreadyExists.Error())
		default:
			return nil, apperror.NewInternal("failed to search accounts", err)
		}
	}

	return accounts, nil
}

func (s *Account) GetAccountByID(ctx context.Context, id uuid.UUID) (*entity.Account, error) {
	ctx, span := s.tracer.Start(ctx, "service: GetAccountByID")
	defer span.End()

	account, err := s.repo.GetAccountByID(ctx, id)
	if account == nil {
		return nil, apperror.NewNotFound(fmt.Sprintf("account not found: %s", id))
	}
	if err != nil {
		return nil, apperror.NewInternal(fmt.Sprintf("failed to get account: %s", id), err)
	}

	return account, nil
}

func (s *Account) DeleteAccount(ctx context.Context, req *request.DeleteAccount) error {
	ctx, span := s.tracer.Start(ctx, "service: DeleteAccount")
	defer span.End()

	id, err := uuid.Parse(req.AccountID)
	if err != nil {
		return apperror.NewBadReq(fmt.Sprintf("invalid account id: %s", id))
	}

	err = s.repo.DeleteAccount(ctx, id)
	if err != nil {
		switch err {
		case accountrepo.ErrAccountNotFound:
			return apperror.NewBadReq(fmt.Sprintf("account not found: %s", req.AccountID))
		default:
			return apperror.NewInternal("failed to delete account", err)
		}
	}

	return nil
}

func (s *Account) AddTwoAccountsBalance(ctx context.Context, req *request.CreateTransfer) ([]*entity.Account, error) {
	ctx, span := s.tracer.Start(ctx, "service: AddTwoAccountsBalance")
	defer span.End()

	arg := accountrepo.AddTwoAccountsBalance{
		FromAccountID: req.FromAccountID,
		ToAccountUD:   req.ToAccountID,
		Amount:        int64(req.Amount * 100),
	}
	accounts, err := s.repo.AddTwoAccountsBalance(ctx, &arg)
	if err != nil {
		return nil, apperror.NewInternal("failed to get account", err)
	}

	return accounts, nil
}

func (s *Account) GetCurrencies(ctx context.Context) ([]*entity.Currency, error) {
	ctx, span := s.tracer.Start(ctx, "service: SearchAccounts")
	defer span.End()

	currencies, err := s.repo.GetCurrencies(ctx)
	if err != nil {
		return nil, apperror.NewInternal("failed to get accounts", err)
	}

	return currencies, nil
}
