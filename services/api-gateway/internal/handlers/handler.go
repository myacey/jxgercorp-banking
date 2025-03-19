package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/myacey/jxgercorp-banking/services/shared/cstmerr"
	tokenpb "github.com/myacey/jxgercorp-banking/services/shared/proto/token"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Handler struct {
	tokenSrv tokenpb.TokenServiceClient
	lg       *zap.SugaredLogger

	tracer trace.Tracer
}

func NewHandler(tokenSrv tokenpb.TokenServiceClient, lg *zap.SugaredLogger, tracer trace.Tracer) *Handler {
	return &Handler{
		tokenSrv: tokenSrv,
		lg:       lg,
		tracer:   tracer,
	}
}

func (h *Handler) JSONError(c *gin.Context, err error, opts ...int) {
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
