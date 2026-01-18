package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/myacey/jxgercorp-banking/services/libs/apperror"
	"github.com/myacey/jxgercorp-banking/services/user/internal/models/dto/request"
)

const (
	secondsPerDay = 86400
	CookieAuth    = "authToken"

	QueryUsername    = "username"
	QueryConfirmCode = "code"
)

// CreateUser checks request's json and send req so service
func (h *Handler) CreateUser(c *gin.Context) {
	var req request.Register
	if err := c.ShouldBindJSON(&req); err != nil {
		wrapCtxWithError(c, apperror.NewBadReq("invalid req: "+err.Error()))
		return
	}

	usr, err := h.userSrv.CreateUser(c, &req)
	if err != nil {
		wrapCtxWithError(c, err)
		return
	}

	// h.metrics.RecordRegister(c.Request.Context(), attrs)

	c.JSON(http.StatusCreated, usr.ToResponse())
}

// Login checks user data with existing in db and sends auth token
func (h *Handler) Login(c *gin.Context) {
	ctx, span := h.tracer.Start(c.Request.Context(), "controller: Login")
	defer span.End()
	c.Request = c.Request.WithContext(ctx)

	var req request.Login
	if err := c.ShouldBindJSON(&req); err != nil {
		wrapCtxWithError(c, apperror.NewBadReq("invalid req: "+err.Error()))
		return
	}

	token, err := h.userSrv.Login(c, &req)
	if err != nil {
		wrapCtxWithError(c, err)
		return
	}

	c.SetCookie(
		CookieAuth,
		token,
		secondsPerDay,
		"/",
		h.cfg.AppDomain,
		h.cfg.AppDomain != "localhost",
		false,
	)
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *Handler) ConfirmUserEmail(c *gin.Context) {
	ctx, span := h.tracer.Start(c.Request.Context(), "controller: ConfirmUserEmail")
	defer span.End()
	c.Request = c.Request.WithContext(ctx)

	username, ok := c.GetQuery(QueryUsername)
	if !ok {
		wrapCtxWithError(c, apperror.NewBadReq("invalid req: no username query key"))
		return
	}

	code, ok := c.GetQuery(QueryConfirmCode)
	if !ok {
		wrapCtxWithError(c, apperror.NewBadReq("invalid req: no code query key"))
		return
	}

	msg, err := h.userSrv.ConfirmUserEmail(c, &request.ConfirmUserEmail{username, code})
	if err != nil {
		wrapCtxWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": msg})
}

func (h *Handler) GetUserByID(c *gin.Context) {
	ctx, span := h.tracer.Start(c.Request.Context(), "controller: GetUserByID")
	defer span.End()
	c.Request = c.Request.WithContext(ctx)

	var req request.GetUserByID
	if err := c.ShouldBindJSON(&req); err != nil {
		wrapCtxWithError(c, apperror.NewBadReq("invalid req: "+err.Error()))
		return
	}

	usr, err := h.userSrv.GetUserByID(c, &req)
	if err != nil {
		wrapCtxWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, usr.ToResponse())
}

func (h *Handler) GetUserByUsername(c *gin.Context) {
	ctx, span := h.tracer.Start(c.Request.Context(), "controller: GetUserByUsername")
	defer span.End()
	c.Request = c.Request.WithContext(ctx)

	var req request.GetUserByUsername
	if err := c.ShouldBindJSON(&req); err != nil {
		wrapCtxWithError(c, apperror.NewBadReq("invalid req: "+err.Error()))
		return
	}

	usr, err := h.userSrv.GetUserByUsername(c, req.Username)
	if err != nil {
		wrapCtxWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, usr.ToResponse())
}

func (h *Handler) DeleteUserByUsername(c *gin.Context) {
	ctx, span := h.tracer.Start(c.Request.Context(), "controller: DeleteUserByUsername")
	defer span.End()
	c.Request = c.Request.WithContext(ctx)

	username := c.GetHeader(HeaderUsername)
	if username == "" {
		wrapCtxWithError(c, apperror.NewInternal(
			"failed to delete user",
			errors.New("should have Username in ctx, but dont: "+c.GetHeader(HeaderRequestID)),
		))
		return
	}

	err := h.userSrv.DeleteUserByUsername(c, username)
	if err != nil {
		wrapCtxWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (h *Handler) UpdateUserInfo(c *gin.Context) {
	ctx, span := h.tracer.Start(c.Request.Context(), "controller: UpdateUserInfo")
	defer span.End()
	c.Request = c.Request.WithContext(ctx)

	var req request.UpdateUserInfo
	if err := c.ShouldBindJSON(&req); err != nil {
		wrapCtxWithError(c, apperror.NewBadReq("invalid req: "+err.Error()))
		return
	}

	username := c.GetHeader(HeaderUsername)

	usr, err := h.userSrv.UpdateUserInfo(c, &req, username)
	if err != nil {
		wrapCtxWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, usr)
}
