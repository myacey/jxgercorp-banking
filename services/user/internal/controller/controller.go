package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/myacey/jxgercorp-banking/services/shared/cstmerr"
	tokenpb "github.com/myacey/jxgercorp-banking/services/shared/proto/token"
	"github.com/myacey/jxgercorp-banking/services/user/internal/service"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// ControllerInterface represents all functions
// that can be called threw API
type ControllerInterface interface {
	// users
	CreateUser(c *gin.Context)
	Login(c *gin.Context)
	GetUserBalance(c *gin.Context)

	GetUserByID(c *gin.Context)
	GetUserByUsername(c *gin.Context)
	DeleteUserByUsername(c *gin.Context)
	UpdateUserInfo(c *gin.Context)

	AuthMiddleware() gin.HandlerFunc
}

type Controller struct {
	srv      service.ServiceInterface
	lg       *zap.SugaredLogger
	tokenSrv tokenpb.TokenServiceClient // for creating tokens in /register; TODO: move login to api-gateway??

	tracer trace.Tracer
}

func NewController(srv service.ServiceInterface, tokenSrv tokenpb.TokenServiceClient, lg *zap.SugaredLogger, tracer trace.Tracer) *Controller {
	return &Controller{
		srv:      srv,
		lg:       lg,
		tokenSrv: tokenSrv,
		tracer:   tracer,
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
		if errors.Is(err, cstmerr.ErrInvalidToken) {
			statusCode = http.StatusUnauthorized
		}

		c.JSON(statusCode, gin.H{"error": gin.H{
			"code":    statusCode,
			"message": err.Error(),
			"details": "",
		}})
	}
}
