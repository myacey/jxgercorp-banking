package handler

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/myacey/jxgercorp-banking/services/libs/apperror"
	"github.com/myacey/jxgercorp-banking/services/libs/telemetry"
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

type Config struct {
	AppDomain string `mapstructure:"app_domain"`
}

type UserService interface {
	ConfirmUserEmail(ctx context.Context, req *request.ConfirmUserEmail) (string, error)
	CreateUser(ctx context.Context, req *request.Register) (*entity.User, error)
	DeleteUserByUsername(ctx context.Context, username string) error
	GetUserByID(ctx context.Context, req *request.GetUserByID) (*entity.User, error)
	GetUserByUsername(ctx context.Context, username string) (*entity.User, error)
	Login(ctx context.Context, req *request.Login) (string, error)
	UpdateUserInfo(ctx context.Context, req *request.UpdateUserInfo, username string) (*entity.User, error)
	SearchUsersByUsernameSimilarity(ctx context.Context, searchUsername string) ([]string, error)
}

type Handler struct {
	cfg Config

	userSrv UserService

	tracer  trace.Tracer
	metrics *telemetry.HTTPMetrics
}

func NewHandler(userSrv UserService, cfg Config, metrics telemetry.HTTPMetrics) *Handler {
	return &Handler{
		userSrv: userSrv,
		tracer:  otel.Tracer("handler"),
		// metrics: metrics,
		cfg: cfg,
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
