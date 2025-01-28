package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/myacey/jxgercorp-banking/shared/cstmerr"
	"github.com/myacey/jxgercorp-banking/user/internal/service"
	"go.uber.org/zap"
)

// ControllerInterface represents all functions
// that cant be called threw API
type ControllerInterface interface {
	// users
	CreateUser(c *gin.Context)
	GetUserByID(c *gin.Context)
	GetUserByUsername(c *gin.Context)
	DeleteUserByUsername(c *gin.Context)
	UpdateUserInfo(c *gin.Context)

	AuthMiddleware() gin.HandlerFunc
}

type Controller struct {
	srv service.ServiceInterface
	lg  *zap.SugaredLogger
}

func NewController(srv service.ServiceInterface, lg *zap.SugaredLogger) *Controller {
	return &Controller{
		srv: srv,
		lg:  lg,
	}
}

func (h *Controller) JSONError(c *gin.Context, err error, opts ...int) {
	h.lg.Debug(err)
	if httpError, ok := err.(*cstmerr.HTTPErr); ok {
		if httpError.Err != nil {
			h.lg.Error(httpError.Err)
		}
		c.JSON(httpError.Code, gin.H{"error": httpError})
	} else {
		statusCode := http.StatusBadRequest
		if len(opts) > 0 {
			statusCode = opts[0]
		}
		c.JSON(statusCode, gin.H{"error": gin.H{
			"code":    statusCode,
			"message": err.Error(),
			"details": "",
		}})
	}
}
