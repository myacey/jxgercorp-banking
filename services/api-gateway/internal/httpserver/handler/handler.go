package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/myacey/jxgercorp-banking/services/api-gateway/internal/models/dto/response"
	"github.com/myacey/jxgercorp-banking/services/api-gateway/internal/service"
	"github.com/myacey/jxgercorp-banking/services/libs/apperror"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

const (
	HeaderRequestID = "X-Request-Id"
	HeaderUsername  = "X-User-Username"

	CtxKeyRetryAfter = "Retry-After"
	CtxKeyUsername   = "Username"
)

type Handler struct {
	srv service.Service

	tracer trace.Tracer
	// metrics *telemetry.Metrics
}

func NewHandler(srv service.Service) *Handler {
	return &Handler{
		srv:    srv,
		tracer: otel.Tracer("api-gateway"),
	}
}

func wrapCtxWithError(ctx *gin.Context, err error) {
	if httpError, ok := err.(apperror.HTTPError); ok {
		ctx.JSON(httpError.Code, response.Error{
			Code:      httpError.Code,
			Message:   httpError.Message,
			RequestID: ctx.GetHeader(HeaderRequestID),
		})

		if httpError.Code == http.StatusInternalServerError {
			log.Printf("internal error: %v | %v", httpError.Message, httpError.DebugError)
		}
	} else {
		ctx.JSON(http.StatusInternalServerError, response.Error{
			Code:      http.StatusInternalServerError,
			Message:   err.Error(),
			RequestID: ctx.GetHeader(HeaderRequestID),
		})
	}
	ctx.Set(CtxKeyRetryAfter, 10)
}
