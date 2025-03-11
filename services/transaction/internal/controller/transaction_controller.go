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
	// var req GetUserTrxReq
	// if err := c.ShouldBindJSON(&req); err != nil {
	// 	h.JSONError(c, err)
	// }

	h.lg.Info("1")
	offset, err := strconv.Atoi(c.Query("offset"))
	if err != nil {
		h.lg.Info("2")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid offset"})
		return
	}
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		h.lg.Info("3")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit"})
		return
	}
	h.lg.Info("4")

	username := c.GetHeader("X-User-Username")
	if username == "" {
		h.JSONError(c, cstmerr.ErrInvalidToken)
		return
	}
	h.lg.Info("5")

	trxs, err := h.srv.SearchEntriesForUser(c, username, offset, limit)
	if err != nil {
		h.lg.Info("6")
		h.JSONError(c, err)
		return
	}

	c.JSON(http.StatusOK, trxs)
}
