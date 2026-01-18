package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"github.com/myacey/jxgercorp-banking/services/libs/apperror"
	"github.com/myacey/jxgercorp-banking/services/transfer/internal/models/dto/response"
)

const (
	HeaderUsername  = "X-User-Username"
	HeaderRequestID = "X-Request-Id"

	CtxKeyRetryAfter = "Retry-After"
)

type Handler struct {
	transferSrv TransferService
	accountSrv  AccountService

	tracer trace.Tracer
}

func NewHandler(transferSrv TransferService, accountSrv AccountService) *Handler {
	return &Handler{
		transferSrv: transferSrv,
		accountSrv:  accountSrv,

		tracer: otel.Tracer("handler"),
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
			log.Printf("internal error. responseID: %v; user message: %v; debug message: %v", ctx.GetHeader(HeaderRequestID), httpError.Message, httpError.DebugError)
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
