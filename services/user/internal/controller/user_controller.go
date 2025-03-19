package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/myacey/jxgercorp-banking/services/shared/cstmerr"
	"github.com/myacey/jxgercorp-banking/services/shared/ctxkeys"
	"github.com/myacey/jxgercorp-banking/services/user/internal/models"
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
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login checks user data with existing in db and sends auth token
func (h *Controller) Login(c *gin.Context) {
	ctx, span := h.tracer.Start(c.Request.Context(), "controller: Login")
	defer span.End()
	c.Request = c.Request.WithContext(ctx)

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

	c.SetCookie("authToken", token, secondsPerDay, "/", "localhost", false, false)
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

// GetUserBalance checks user token and return user's balance
func (h *Controller) GetUserBalance(c *gin.Context) {
	ctx, span := h.tracer.Start(c.Request.Context(), "controller: GetUserBalance")
	defer span.End()
	c.Request = c.Request.WithContext(ctx)

	username := c.GetHeader("X-User-Username")
	if username == "" {
		h.JSONError(c, cstmerr.ErrInvalidToken)
		return
	}

	usr, err := h.srv.GetUserByUsername(c, username)
	if err != nil {
		h.JSONError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"balance": usr.Balance})
}

type ConfirmUserEmailReq struct {
	Username string `json:"username"`
	Code     string `json:"code"`
}

func (h *Controller) ConfirmUserEmail(c *gin.Context) {
	ctx, span := h.tracer.Start(c.Request.Context(), "controller: ConfirmUserEmail")
	defer span.End()
	c.Request = c.Request.WithContext(ctx)

	var req ConfirmUserEmailReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.lg.Debug(err)
		h.JSONError(c, err)
		return
	}

	h.lg.Debugw("got confirm code req",
		"username", req.Username,
		"code", req.Code,
	)

	msg, err := h.srv.ConfirmUserEmail(c, req.Username, req.Code)
	if err != nil {
		h.lg.Debug(err)
		h.JSONError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": msg})
}

type GetUserByIDReq struct {
	ID int64 `json:"id" binding:"required,number"`
}

func (h *Controller) GetUserByID(c *gin.Context) {
	ctx, span := h.tracer.Start(c.Request.Context(), "controller: GetUserByID")
	defer span.End()
	c.Request = c.Request.WithContext(ctx)

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
	ctx, span := h.tracer.Start(c.Request.Context(), "controller: GetUserByUsername")
	defer span.End()
	c.Request = c.Request.WithContext(ctx)

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
	ctx, span := h.tracer.Start(c.Request.Context(), "controller: DeleteUserByUsername")
	defer span.End()
	c.Request = c.Request.WithContext(ctx)

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

// TODO: change
func (h *Controller) UpdateUserInfo(c *gin.Context) {
	ctx, span := h.tracer.Start(c.Request.Context(), "controller: UpdateUserInfo")
	defer span.End()
	c.Request = c.Request.WithContext(ctx)

	var req UpdateUserInfoReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.JSONError(c, err)
		return
	}

	usrnameCtx, exists := c.Get(string(ctxkeys.UsernameKey))
	if !exists {
		h.JSONError(c, cstmerr.ErrInvalidToken)
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
