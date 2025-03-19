package postgresrepo

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/myacey/jxgercorp-banking/services/db/sqlc"
	"github.com/myacey/jxgercorp-banking/services/shared/cstmerr"
	"github.com/myacey/jxgercorp-banking/services/user/internal/repository"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type PostgresUserRepository struct {
	store *db.Queries
	lg    *zap.SugaredLogger

	tracer trace.Tracer
}

func NewUserRepo(queries *db.Queries, logger *zap.SugaredLogger, tr trace.Tracer) repository.UserRepository {
	return &PostgresUserRepository{
		store:  queries,
		lg:     logger,
		tracer: tr,
	}
}

func (r *PostgresUserRepository) CreateUser(c *gin.Context, username, email, hashedPassword string) (*db.User, error) {
	ctx, span := r.tracer.Start(c.Request.Context(), "repository: CreateUser")
	span.SetAttributes(
		attribute.String("db.operation", "CreateUser"),
		attribute.String("username", username),
		attribute.String("email", email),
	)
	defer span.End()
	c.Request = c.Request.WithContext(ctx)

	arg := db.CreateUserParams{
		Username:       username,
		Email:          email,
		HashedPassword: hashedPassword,
	}
	dbUser, err := r.store.CreateUser(c, arg)
	if err != nil {
		switch {
		case isUniqueViolation(err):
			return nil, cstmerr.New(http.StatusBadRequest, "email or username not unique", nil)
		default:
			return nil, cstmerr.New(http.StatusInternalServerError, cstmerr.ErrUnknown.Error(), err)
		}
	}
	return &dbUser, nil
}

func (r *PostgresUserRepository) GetUserByID(c *gin.Context, id int64) (*db.User, error) {
	ctx, span := r.tracer.Start(c.Request.Context(), "repository: CreateUser")
	span.SetAttributes(
		attribute.String("db.operation", "GetUserByID"),
		attribute.Int64("id", id),
	)
	defer span.End()
	c.Request = c.Request.WithContext(ctx)

	start := time.Now()
	dbUser, err := r.store.GetUserByID(c, id)
	elapsed := time.Since(start)
	span.SetAttributes(attribute.Float64("sql.time", float64(elapsed.Milliseconds())))

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, cstmerr.New(http.StatusNotFound, "no user with this ID found", nil)
		default:
			return nil, cstmerr.New(http.StatusInternalServerError, cstmerr.ErrUnknown.Error(), err)
		}
	}
	return &dbUser, nil
}

func (r *PostgresUserRepository) GetUserByUsername(c *gin.Context, username string) (*db.User, error) {
	ctx, span := r.tracer.Start(c.Request.Context(), "repository: GetUserByUsername")
	span.SetAttributes(
		attribute.String("db.operation", "GetUserByUsername"),
		attribute.String("username", username),
	)
	defer span.End()
	c.Request = c.Request.WithContext(ctx)

	start := time.Now()
	dbUser, err := r.store.GetUserByUsername(c, username)
	elapsed := time.Since(start)
	span.SetAttributes(attribute.Float64("sql.time", float64(elapsed.Milliseconds())))

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, cstmerr.New(http.StatusBadRequest, "no user found", nil)
		default:
			return nil, cstmerr.New(http.StatusInternalServerError, cstmerr.ErrUnknown.Error(), err)
		}
	}
	return &dbUser, nil
}

func (r *PostgresUserRepository) DeleteUserByUsername(c *gin.Context, username string) error {
	ctx, span := r.tracer.Start(c.Request.Context(), "repository: DeleteUserByUsername")
	span.SetAttributes(
		attribute.String("db.operation", "DeleteUserByUsername"),
		attribute.String("username", username),
	)
	defer span.End()
	c.Request = c.Request.WithContext(ctx)

	start := time.Now()
	err := r.store.DeleteUserByUsername(c, username)
	elapsed := time.Since(start)
	span.SetAttributes(attribute.Float64("sql.time", float64(elapsed.Milliseconds())))

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return cstmerr.New(http.StatusBadRequest, "no user with such username", nil)
		default:
			return cstmerr.New(http.StatusInternalServerError, cstmerr.ErrUnknown.Error(), err)
		}
	}

	return nil
}

func (r *PostgresUserRepository) UpdateUserInfo(c *gin.Context, username string, newEmail string, newHashedPassword string, pendingStatus bool) (*db.User, error) {
	ctx, span := r.tracer.Start(c.Request.Context(), "repository: UpdateUserInfo")
	span.SetAttributes(
		attribute.String("db.operation", "UpdateUserInfo"),
		attribute.String("username", username),
		attribute.String("newEmail", newEmail),
		attribute.Bool("pendingStatus", pendingStatus),
	)
	defer span.End()
	c.Request = c.Request.WithContext(ctx)

	arg := db.UpdateUserInfoParams{
		Username:       username,
		HashedPassword: sql.NullString{String: newHashedPassword, Valid: newHashedPassword != ""},
		Email:          sql.NullString{String: newEmail, Valid: newEmail != ""},
		Pending:        sql.NullBool{Bool: pendingStatus, Valid: !pendingStatus}, // pending can only move from true->false
	}
	start := time.Now()
	dbUser, err := r.store.UpdateUserInfo(c, arg)
	elapsed := time.Since(start)
	span.SetAttributes(attribute.Float64("sql.time", float64(elapsed.Milliseconds())))
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, cstmerr.New(http.StatusBadRequest, "no user with such username", nil)
		default:
			return nil, cstmerr.New(http.StatusInternalServerError, cstmerr.ErrUnknown.Error(), err)
		}
	}

	return &dbUser, nil
}
