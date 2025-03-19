package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/myacey/jxgercorp-banking/services/shared/cstmerr"
)

type CreateTrxReq struct {
	ToUser string `json:"to_user" binding:"required"`
	Amount int64  `json:"amount" binding:"required,number"`
}

func (h *Controller) CreateNewTransaction(c *gin.Context) {
	ctx, span := h.tracer.Start(c.Request.Context(), "controller: CreateNewTransaction")
	defer span.End()
	c.Request = c.Request.WithContext(ctx)

	var req CreateTrxReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.JSONError(c, err)
		return
	}

	fromUser := c.GetHeader("X-User-Username")
	if fromUser == "" {
		h.JSONError(c, cstmerr.ErrInvalidToken)
		return
	}

	if fromUser == req.ToUser {
		h.JSONError(c, cstmerr.New(http.StatusBadRequest, "cant send money to yourself", nil))
		return
	}

	trx, err := h.srv.CreateNewTransaction(c, fromUser, req.ToUser, req.Amount)
	if err != nil {
		h.JSONError(c, err)
		return
	}

	c.JSON(http.StatusCreated, trx)
}

type GetUserTrxReq struct {
	Offset int `json:"offset" biniding:"number"`
	Limit  int `json:"limit" biniding:"number"`
}

func (h *Controller) SearchEntriesForUser(c *gin.Context) {
	ctx, span := h.tracer.Start(c.Request.Context(), "controller: SearchEntriesForUser")
	defer span.End()
	c.Request = c.Request.WithContext(ctx)

	// var req GetUserTrxReq
	// if err := c.ShouldBindJSON(&req); err != nil {
	// 	h.JSONError(c, err)
	// }

	offset, err := strconv.Atoi(c.Query("offset"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid offset"})
		return
	}
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit"})
		return
	}

	username := c.GetHeader("X-User-Username")
	if username == "" {
		h.JSONError(c, cstmerr.ErrInvalidToken)
		return
	}

	trxs, err := h.srv.SearchEntriesForUser(c, username, offset, limit)
	if err != nil {
		h.JSONError(c, err)
		return
	}

	c.JSON(http.StatusOK, trxs)
}
