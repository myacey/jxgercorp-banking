package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/myacey/jxgercorp-banking/shared/cstmerr"
	"github.com/myacey/jxgercorp-banking/shared/ctxkeys"
	"github.com/myacey/jxgercorp-banking/user/internal/models"
)

const secondsPerDay = 86400

type CreateUserReq struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"email,required"`
	Password string `json:"password" binding:"required"`
} // TODO: sign from new device

// CreateUser checks request's json and send req so service
func (h *Controller) CreateUser(c *gin.Context) {
	h.lg.Info("AOOOOOOOOOo")
	var req CreateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.JSONError(c, err)
		return
	}

	usr := &models.UserUnhashed{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}
	_, err := h.srv.CreateUser(c, usr)
	if err != nil {
		h.JSONError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "success"})
}

type LoginReq struct {
	Username string `json:"username" bindding:"required"`
	Password string `json:"password" bingding:"required"`
}

func (h *Controller) Login(c *gin.Context) {
	var req LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.JSONError(c, err)
		return
	}

	token, err := h.srv.Login(c, req.Username, req.Password)
	if err != nil {
		h.JSONError(c, err)
		return
	}

	c.SetCookie("authToken", token, secondsPerDay, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

type GetUserByIDReq struct {
	ID int64 `json:"id", binding="required,number"`
}

func (h *Controller) GetUserByID(c *gin.Context) {
	var req GetUserByIDReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.JSONError(c, err)
		return
	}

	usr, err := h.srv.GetUserByID(c, req.ID)
	if err != nil {
		h.JSONError(c, err)
		return
	}

	c.JSON(http.StatusOK, &usr)
}

type GetUserByUsernameReq struct {
	Username string `json:"username" binding:"required"`
}

func (h *Controller) GetUserByUsername(c *gin.Context) {
	var req GetUserByUsernameReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.JSONError(c, err)
		return
	}

	usr, err := h.srv.GetUserByUsername(c, req.Username)
	if err != nil {
		h.JSONError(c, err)
		return
	}

	c.JSON(http.StatusOK, &usr)
}

type DeleteUserByUsernameReq struct {
	Username string `json:"id" binding:"required"`
}

func (h *Controller) DeleteUserByUsername(c *gin.Context) {
	var req DeleteUserByUsernameReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.JSONError(c, err)
		return
	}

	err := h.srv.DeleteUserByUsername(c, req.Username)
	if err != nil {
		h.JSONError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

type UpdateUserInfoReq struct {
	NewEmail    string `json:"email" binding:"email"`
	NewPassword string `json:"password"`
}

func (h *Controller) UpdateUserInfo(c *gin.Context) {
	var req UpdateUserInfoReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.JSONError(c, err)
		return
	}

	usrnameCtx, exists := c.Get(string(ctxkeys.UsernameKey))
	if !exists {
		h.JSONError(c, fmt.Errorf("unauthorized"), http.StatusUnauthorized)
		return
	}
	username, ok := usrnameCtx.(string)
	if !ok {
		h.JSONError(c, cstmerr.ErrUnknown, http.StatusInternalServerError)
		return
	}

	_, err := h.srv.UpdateUserInfo(c, username, req.NewEmail, req.NewPassword)
	if err != nil {
		h.JSONError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}
