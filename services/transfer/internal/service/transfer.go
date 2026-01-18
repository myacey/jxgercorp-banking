package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/myacey/jxgercorp-banking/services/libs/apperror"
	"github.com/myacey/jxgercorp-banking/services/transfer/internal/models/dto/request"
	"github.com/myacey/jxgercorp-banking/services/transfer/internal/models/entity"
	"github.com/myacey/jxgercorp-banking/services/transfer/internal/repository/transferrepo"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type TransferRepo interface {
	CraeteTransfer(ctx context.Context, req *entity.Transfer) (*entity.Transfer, error)
	SearchTransfersWithAccount(ctx context.Context, req *transferrepo.SearchTransfersWithAccountParams) ([]*entity.Transfer, error)
}

type Transfer struct {
	accountSrv *Account
	repo       TransferRepo
	conn       *pgxpool.Pool

	tracer trace.Tracer
}

func NewTransfer(conn *pgxpool.Pool, accountSrv *Account, repo TransferRepo) *Transfer {
	return &Transfer{
		accountSrv: accountSrv,
		repo:       repo,
		conn:       conn,
		tracer:     otel.Tracer("service-transfer"),
	}
}

func (s *Transfer) CreateTransfer(ctx context.Context, req *request.CreateTransfer) (*entity.Transfer, error) {
	ctx, span := s.tracer.Start(ctx, "service: CreateTransfer")
	defer span.End()

	tx, err := s.conn.Begin(ctx)
	if err != nil {
		return nil, apperror.NewInternal("failed to create transfer", err)
	}
	defer tx.Rollback(ctx)

	fromAccount, err := s.accountSrv.GetAccountByID(ctx, req.FromAccountID)
	if err != nil {
		return nil, err
	}

	if fromAccount.Balance-int64(req.Amount*100) < 0 {
		return nil, apperror.NewBadReq("not enough money")
	}

	toAccount, err := s.accountSrv.GetAccountByID(ctx, req.ToAccountID)
	if err != nil {
		return nil, err
	}

	if fromAccount.CurrencyCode != toAccount.CurrencyCode {
		return nil, apperror.NewBadReq(fmt.Sprintf("not same currency: from %s to %s", fromAccount.CurrencyCode, toAccount.CurrencyCode))
	}

	transfer := &entity.Transfer{
		ID:                  uuid.New(),
		FromAccountID:       req.FromAccountID,
		FromAccountUsername: req.FromAccountUsername,
		ToAccountID:         req.ToAccountID,
		ToAccountUsername:   req.ToAccountUsername,
		Amount:              int64(req.Amount * 100),
		CurrencyCode:        req.Currency,
	}

	_, err = s.accountSrv.AddTwoAccountsBalance(ctx, req)
	if err != nil {
		return nil, err
	}

	res, err := s.repo.CraeteTransfer(ctx, transfer)
	if err != nil {
		return nil, apperror.NewInternal("failed to create transfer", err)
	}

	tx.Commit(ctx)
	return res, nil
}

func (s *Transfer) SearchTransfersWithAccount(ctx context.Context, req *request.SearchTransfersWithAccount) ([]*entity.Transfer, error) {
	ctx, span := s.tracer.Start(ctx, "service: SearchTransfersWithAccount")
	defer span.End()

	currentAccountID, err := uuid.Parse(req.CurrentAccountID)
	if err != nil {
		return nil, apperror.NewBadReq(fmt.Sprintf("invalid account id: %v", req.CurrentAccountID))
	}

	withAccountID := uuid.Nil
	if req.WithAccountID != "" {
		withAccountID, err = uuid.Parse(req.WithAccountID)
		if err != nil {
			return nil, apperror.NewBadReq(fmt.Sprintf("invalid account id: %v", req.CurrentAccountID))
		}
	}

	arg := transferrepo.SearchTransfersWithAccountParams{
		CurrentAccountID: currentAccountID,
		WithUsername:     req.WithUsername,
		WithAccountID:    withAccountID,
		Currency:         req.CurrencyCode,
		Offset:           req.Offset,
		Limit:            req.Limit,
	}
	transfers, err := s.repo.SearchTransfersWithAccount(ctx, &arg)
	if err != nil {
		return nil, err
	}

	return transfers, nil
}
