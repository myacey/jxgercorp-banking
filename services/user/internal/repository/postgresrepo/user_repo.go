package postgresrepo

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/myacey/jxgercorp-banking/services/db/sqlc"
	"github.com/myacey/jxgercorp-banking/services/shared/cstmerr"
	"github.com/myacey/jxgercorp-banking/services/user/internal/repository"
	"go.uber.org/zap"
)

type PostgresUserRepository struct {
	store *db.Queries
	lg    *zap.SugaredLogger
}

func NewUserRepo(queries *db.Queries, logger *zap.SugaredLogger) repository.UserRepository {
	return &PostgresUserRepository{
		store: queries,
		lg:    logger,
	}
}

func (r *PostgresUserRepository) CreateUser(c *gin.Context, username, email, hashedPassword string) (*db.User, error) {
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
	dbUser, err := r.store.GetUserByID(c, id)
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
	dbUser, err := r.store.GetUserByUsername(c, username)
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
	err := r.store.DeleteUserByUsername(c, username)
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
	arg := db.UpdateUserInfoParams{
		Username:       username,
		HashedPassword: sql.NullString{String: newHashedPassword, Valid: newHashedPassword != ""},
		Email:          sql.NullString{String: newEmail, Valid: newEmail != ""},
		Pending:        sql.NullBool{Bool: pendingStatus, Valid: pendingStatus != true}, // pending can only move from true->false
	}
	dbUser, err := r.store.UpdateUserInfo(c, arg)
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
