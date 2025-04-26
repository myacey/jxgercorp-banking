package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/myacey/jxgercorp-banking/services/libs/apperror"
	"github.com/myacey/jxgercorp-banking/services/transfer/internal/models/dto/request"
	"github.com/myacey/jxgercorp-banking/services/transfer/internal/models/dto/response"
	"github.com/myacey/jxgercorp-banking/services/transfer/internal/models/entity"
)

type TransferService interface {
	CreateTransfer(ctx context.Context, req *request.CreateTransfer) (*entity.Transfer, error)
	SearchTransfersWithAccount(ctx context.Context, req *request.SearchTransfersWithAccount) ([]*entity.Transfer, error)
}

func (h *Handler) CreateTransfer(c *gin.Context) {
	var req request.CreateTransfer
	if err := c.ShouldBindJSON(&req); err != nil {
		wrapCtxWithError(c, apperror.NewBadReq("invalid req: "+err.Error()))
		return
	}

	if req.FromAccountID == req.ToAccountID {
		wrapCtxWithError(c, apperror.NewBadReq("cant send money on same account"))
		return
	}

	if req.Amount <= 0 {
		wrapCtxWithError(c, apperror.NewBadReq("invalid money amount"))
		return
	}

	username := c.GetHeader(HeaderUsername)
	if username == "" {
		wrapCtxWithError(c, apperror.NewUnauthorized("invalid token"))
		return
	}

	transfer, err := h.transferSrv.CreateTransfer(c, &req)
	if err != nil {
		wrapCtxWithError(c, err)
		return
	}

	c.JSON(http.StatusCreated, transfer.ToResponse())
}

func (h *Handler) SearchTransfersWithAccount(c *gin.Context) {
	var req request.SearchTransfersWithAccount
	if err := c.ShouldBindJSON(&req); err != nil {
		wrapCtxWithError(c, apperror.NewBadReq("invalid req: "+err.Error()))
		return
	}

	transfers, err := h.transferSrv.SearchTransfersWithAccount(c, &req)
	if err != nil {
		wrapCtxWithError(c, err)
		return
	}

	resp := make([]*response.Transfer, len(transfers))
	for i, v := range transfers {
		resp[i] = v.ToResponse()
	}

	c.JSON(http.StatusOK, resp)
}
