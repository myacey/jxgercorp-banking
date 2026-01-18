package userrepo

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/google/uuid"
	"github.com/myacey/jxgercorp-banking/services/user/internal/models/dto/request"
	"github.com/myacey/jxgercorp-banking/services/user/internal/models/entity"
	"github.com/myacey/jxgercorp-banking/services/user/internal/repository"
	db "github.com/myacey/jxgercorp-banking/services/user/internal/repository/sqlc"
)

var (
	ErrCredentialsDuplicate = errors.New("username or email already exists")
	ErrUserNotFound         = errors.New("user not found")
)

type PostgresUser struct {
	store *db.Queries

	tracer trace.Tracer
}

func NewUserRepo(queries *db.Queries) *PostgresUser {
	return &PostgresUser{
		store:  queries,
		tracer: otel.Tracer("repository-user"),
	}
}

func (r *PostgresUser) CreateUser(ctx context.Context, req *request.Register) (*entity.User, error) {
	ctx, span := r.tracer.Start(ctx, "repository: CreateUser")
	span.SetAttributes(
		attribute.String("db.operation", "CreateUser"),
	)
	defer span.End()

	arg := db.CreateUserParams{
		ID:             uuid.New(),
		Username:       req.Username,
		Email:          req.Email,
		HashedPassword: req.Password,
	}
	dbUser, err := r.store.CreateUser(ctx, arg)
	if err != nil {
		switch {
		case repository.IsUniqueViolation(err):
			return nil, ErrCredentialsDuplicate
		default:
			return nil, err
		}
	}
	return &entity.User{
		ID:             dbUser.ID,
		Username:       dbUser.Username,
		Email:          dbUser.Email,
		HashedPassword: dbUser.HashedPassword,
		CreatedAt:      dbUser.CreatedAt,
		Status:         dbUser.Status,
	}, nil
}

func (r *PostgresUser) GetUserByID(ctx context.Context, req *request.GetUserByID) (*entity.User, error) {
	ctx, span := r.tracer.Start(ctx, "repository: CreateUser")
	span.SetAttributes(
		attribute.String("db.operation", "GetUserByID"),
	)
	defer span.End()

	start := time.Now()
	dbUser, err := r.store.GetUserByID(ctx, req.ID)
	elapsed := time.Since(start)
	span.SetAttributes(attribute.Float64("sql.time.ms", float64(elapsed.Milliseconds())))

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrUserNotFound
		default:
			return nil, err
		}
	}
	return &entity.User{
		ID:             dbUser.ID,
		Username:       dbUser.Username,
		Email:          dbUser.Email,
		HashedPassword: dbUser.HashedPassword,
		CreatedAt:      dbUser.CreatedAt,
		Status:         dbUser.Status,
	}, nil
}

func (r *PostgresUser) GetUserByUsername(ctx context.Context, username string) (*entity.User, error) {
	ctx, span := r.tracer.Start(ctx, "repository: GetUserByUsername")
	span.SetAttributes(
		attribute.String("db.operation", "GetUserByUsername"),
	)
	defer span.End()

	start := time.Now()
	dbUser, err := r.store.GetUserByUsername(ctx, username)
	elapsed := time.Since(start)
	span.SetAttributes(attribute.Float64("sql.time.ms", float64(elapsed.Milliseconds())))

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrUserNotFound
		default:
			return nil, err
		}
	}
	return &entity.User{
		ID:             dbUser.ID,
		Username:       dbUser.Username,
		Email:          dbUser.Email,
		HashedPassword: dbUser.HashedPassword,
		CreatedAt:      dbUser.CreatedAt,
		Status:         dbUser.Status,
	}, nil
}

func (r *PostgresUser) DeleteUserByUsername(ctx context.Context, username string) error {
	ctx, span := r.tracer.Start(ctx, "repository: DeleteUserByUsername")
	span.SetAttributes(
		attribute.String("db.operation", "DeleteUserByUsername"),
	)
	defer span.End()

	start := time.Now()
	err := r.store.DeleteUserByUsername(ctx, username)
	elapsed := time.Since(start)
	span.SetAttributes(attribute.Float64("sql.time.ms", float64(elapsed.Milliseconds())))

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrUserNotFound
		default:
			return err
		}
	}

	return nil
}

func (r *PostgresUser) UpdateUserInfo(ctx context.Context, req *request.UpdateUserInfo, username string) (*entity.User, error) {
	ctx, span := r.tracer.Start(ctx, "repository: UpdateUserInfo")
	span.SetAttributes(
		attribute.String("db.operation", "UpdateUserInfo"),
	)
	defer span.End()

	arg := db.UpdateUserInfoParams{
		Username:       username,
		HashedPassword: sql.NullString{String: req.NewPassword, Valid: req.NewPassword != ""},
		Email:          sql.NullString{String: req.NewEmail, Valid: req.NewEmail != ""},
	}
	start := time.Now()
	dbUser, err := r.store.UpdateUserInfo(ctx, arg)
	elapsed := time.Since(start)
	span.SetAttributes(attribute.Float64("sql.time.ms", float64(elapsed.Milliseconds())))
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrUserNotFound
		default:
			return nil, err
		}
	}

	return &entity.User{
		ID:             dbUser.ID,
		Username:       dbUser.Username,
		Email:          dbUser.Email,
		HashedPassword: dbUser.HashedPassword,
		CreatedAt:      dbUser.CreatedAt,
		Status:         dbUser.Status,
	}, nil
}

func (r *PostgresUser) UpdateUserStatus(ctx context.Context, newStatus entity.UserStatus, username string) (*entity.User, error) {
	ctx, span := r.tracer.Start(ctx, "repository: UpdateUserStatus")
	span.SetAttributes(
		attribute.String("db.operation", "UpdateUserStatus"),
	)
	defer span.End()

	arg := db.UpdateUserStatusParams{
		Username: username,
		Status:   newStatus,
	}
	start := time.Now()
	dbUser, err := r.store.UpdateUserStatus(ctx, arg)
	elapsed := time.Since(start)
	span.SetAttributes(attribute.Float64("sql.time.ms", float64(elapsed.Milliseconds())))
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrUserNotFound
		default:
			return nil, err
		}
	}

	return &entity.User{
		ID:             dbUser.ID,
		Username:       dbUser.Username,
		Email:          dbUser.Email,
		HashedPassword: dbUser.HashedPassword,
		CreatedAt:      dbUser.CreatedAt,
		Status:         dbUser.Status,
	}, nil
}
