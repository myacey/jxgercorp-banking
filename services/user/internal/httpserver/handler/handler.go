package handler

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/myacey/jxgercorp-banking/services/libs/apperror"
	"github.com/myacey/jxgercorp-banking/services/user/internal/models/dto/request"
	"github.com/myacey/jxgercorp-banking/services/user/internal/models/dto/response"
	"github.com/myacey/jxgercorp-banking/services/user/internal/models/entity"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

const (
	HeaderUsername  = "X-User-Username"
	HeaderRequestID = "X-Request-Id"

	CtxKeyRetryAfter = "Retry-After"
)

type UserService interface {
	ConfirmUserEmail(ctx context.Context, req *request.ConfirmUserEmail) (string, error)
	CreateUser(ctx context.Context, req *request.Register) (*entity.User, error)
	DeleteUserByUsername(ctx context.Context, username string) error
	GetUserByID(ctx context.Context, req *request.GetUserByID) (*entity.User, error)
	GetUserByUsername(ctx context.Context, username string) (*entity.User, error)
	Login(ctx context.Context, req *request.Login) (string, error)
	UpdateUserInfo(ctx context.Context, req *request.UpdateUserInfo, username string) (*entity.User, error)
}

type Handler struct {
	userSrv UserService

	tracer trace.Tracer
	// metrics telemetry.UserMetrics
}

func NewHandler(userSrv UserService) *Handler {
	return &Handler{
		userSrv: userSrv,
		tracer:  otel.Tracer("handler"),
		// metrics: metrics,
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
