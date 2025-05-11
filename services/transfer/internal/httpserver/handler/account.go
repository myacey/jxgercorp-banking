package handler

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/myacey/jxgercorp-banking/services/libs/apperror"
	"github.com/myacey/jxgercorp-banking/services/transfer/internal/models/dto/request"
	"github.com/myacey/jxgercorp-banking/services/transfer/internal/models/dto/response"
	"github.com/myacey/jxgercorp-banking/services/transfer/internal/models/entity"
)

type AccountService interface {
	CreateAccount(ctx context.Context, req *request.CreateAccount) (*entity.Account, error)
	SearchAccounts(ctx context.Context, req *request.SearchAccounts) ([]*entity.Account, error)
}

func (h *Handler) CreateAccount(c *gin.Context) {
	var req request.CreateAccount
	if err := c.ShouldBindJSON(&req); err != nil {
		wrapCtxWithError(c, apperror.NewBadReq("invalid req: "+err.Error()))
		return
	}

	req.OwnerUsername = c.GetHeader(HeaderUsername)
	if req.OwnerUsername == "" {
		wrapCtxWithError(c, apperror.NewUnauthorized("invalid token"))
		return
	}

	account, err := h.accountSrv.CreateAccount(c, &req)
	if err != nil {
		wrapCtxWithError(c, err)
		return
	}

	c.JSON(http.StatusCreated, account.ToResponse())
}

func (h *Handler) SearchAccounts(c *gin.Context) {
	var req request.SearchAccounts
	if err := c.ShouldBindJSON(&req); err != nil {
		wrapCtxWithError(c, apperror.NewBadReq("invalid req: "+err.Error()))
		return
	}

	req.Currency = strings.ToUpper(req.Currency)
	accounts, err := h.accountSrv.SearchAccounts(c, &req)
	if err != nil {
		wrapCtxWithError(c, err)
		return
	}
	resp := make([]*response.Account, len(accounts))
	for i, v := range accounts {
		resp[i] = v.ToResponse()
	}

	c.JSON(http.StatusOK, resp)
}
