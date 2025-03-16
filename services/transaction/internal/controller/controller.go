package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/myacey/jxgercorp-banking/services/shared/cstmerr"
	"github.com/myacey/jxgercorp-banking/services/transaction/internal/service"
	"go.uber.org/zap"
)

// ControllerInterfca represents all functons
// that can be called threw API
type ControllerInterface interface {
	CreateNewTransaction(c *gin.Context)
	SearchEntriesForUser(c *gin.Context)
}

type Controller struct {
	srv service.ServiceInterface
	lg  *zap.SugaredLogger
}

func NewController(srv service.ServiceInterface, lg *zap.SugaredLogger) ControllerInterface {
	return &Controller{
		srv: srv,
		lg:  lg,
	}
}

func (h *Controller) JSONError(c *gin.Context, err error, opts ...int) {
	h.lg.Debug(err)

	var httpError *cstmerr.HTTPErr
	if errors.As(err, &httpError) {
		h.lg.Debug("error can be cast to cstmerr:", httpError)
		if httpError.Err != nil {
			h.lg.Error(httpError.Err)
		}
		c.JSON(httpError.Code, gin.H{"error": httpError})
		return
	}

	statusCode := http.StatusBadRequest
	if len(opts) > 0 {
		statusCode = opts[0]
	}

	if errors.Is(err, cstmerr.ErrInvalidToken) {
		statusCode = http.StatusUnauthorized
	}

	c.JSON(statusCode, gin.H{"error": gin.H{
		"code":    statusCode,
		"message": err.Error(),
		"details": "",
	}})
}
