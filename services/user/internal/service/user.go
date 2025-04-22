package service

import (
	"context"
	"errors"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/myacey/jxgercorp-banking/services/libs/apperror"
	"github.com/myacey/jxgercorp-banking/services/user/internal/models/dto/request"
	"github.com/myacey/jxgercorp-banking/services/user/internal/models/entity"
	"github.com/myacey/jxgercorp-banking/services/user/internal/pkg/hasher"
	"github.com/myacey/jxgercorp-banking/services/user/internal/repository/userrepo"
)

var ErrInvalidEmailOrPassword = errors.New("invalid email or password")

const tokenTTL = time.Duration(24 * time.Hour)

type UserRepo interface {
	CreateUser(ctx context.Context, req *request.Register) (*entity.User, error)
	DeleteUserByUsername(ctx context.Context, username string) error
	GetUserByID(ctx context.Context, req *request.GetUserByID) (*entity.User, error)
	GetUserByUsername(ctx context.Context, username string) (*entity.User, error)
	UpdateUserInfo(ctx context.Context, req *request.UpdateUserInfo, username string) (*entity.User, error)
	UpdateUserStatus(ctx context.Context, newStatus entity.UserStatus, username string) (*entity.User, error)
}

type TokenService interface {
	GenerateToken(ctx context.Context, username string, ttl time.Duration) (string, error)
	ValidateToken(ctx context.Context, token string) (bool, error)
}

type Hasher interface {
	Compare(hash string, msg string) error
	Generate(pswd string) (string, error)
}

type User struct {
	repo UserRepo

	// confirmCodeSrv used as service to create and send confirm codes.
	confirmCodeSrv Confirmation

	// tokenSrv is used to create and user tokens for /login.
	tokenSrv TokenService

	// to hash password of user.
	hasher Hasher

	// tracer realize trace logics into service.
	tracer trace.Tracer
}

func NewUserSrv(repo UserRepo, confirmCodeSrv Confirmation, tokenSrv TokenService, hasher Hasher) *User {
	return &User{
		repo:           repo,
		confirmCodeSrv: confirmCodeSrv,
		tokenSrv:       tokenSrv,
		hasher:         hasher,

		tracer: otel.Tracer("service-user"),
	}
}

func (s *User) CreateUser(ctx context.Context, req *request.Register) (*entity.User, error) {
	ctx, span := s.tracer.Start(ctx, "service: CreateUser")
	defer span.End()

	start := time.Now()
	hashedPassword, err := s.hasher.Generate(req.Password)
	if err != nil {
		switch {
		case errors.Is(err, hasher.ErrTooShort):
			return nil, apperror.NewBadReq("password too short")
		case errors.Is(err, hasher.ErrTooLong):
			return nil, apperror.NewBadReq("password too long")
		default:
			return nil, apperror.NewInternal("failed to register", err)
		}
	}
	span.SetAttributes(attribute.Float64("hash.create.ms", float64(time.Since(start).Milliseconds())))

	req.Password = hashedPassword
	user, err := s.repo.CreateUser(ctx, req)
	if err != nil {
		switch {
		case errors.Is(err, userrepo.ErrCredentialsDuplicate):
			return nil, apperror.NewBadReq("duplicate email or username")
		default:
			return nil, apperror.NewInternal("failed to register", err)
		}
	}

	err = s.confirmCodeSrv.generateAccountConfirmation(ctx, req.Username, req.Email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *User) Login(ctx context.Context, req *request.Login) (string, error) {
	ctx, span := s.tracer.Start(ctx, "service: Login")
	defer span.End()

	usr, err := s.repo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		switch {
		case errors.Is(err, userrepo.ErrUserNotFound):
			return "", apperror.NewNotFound("user not found")
		default:
			return "", apperror.NewInternal("failed to login", err)
		}
	}

	if usr.Status == entity.UserStatusPending {
		return "", apperror.NewUnauthorized("check email for account confirmation")
	}

	start := time.Now()
	if err := s.hasher.Compare(usr.HashedPassword, req.Password); err != nil {
		switch {
		case errors.Is(err, hasher.ErrDontCompare):
			return "", apperror.NewUnauthorized("user not found")
		default:
			return "", apperror.NewInternal("failed to login", err)
		}
	}
	span.SetAttributes(attribute.Float64("hash.create.ms", float64(time.Since(start).Milliseconds())))

	token, err := s.tokenSrv.GenerateToken(ctx, req.Username, tokenTTL)
	if err != nil || token == "" {
		return "", apperror.NewUnauthorized("failed to login")
	}

	return token, nil
}

func (s *User) ConfirmUserEmail(ctx context.Context, req *request.ConfirmUserEmail) (string, error) {
	ctx, span := s.tracer.Start(ctx, "service: ConfirmUserEmail")
	defer span.End()

	err := s.confirmCodeSrv.checkConfirmCode(ctx, req.Username, req.Code)
	if err != nil {
		return "", err
	}

	_, err = s.repo.UpdateUserStatus(ctx, entity.UserStatusActive, req.Username)
	if err != nil {
		return "", err
	}

	return "success", nil
}

func (s *User) GetUserByID(ctx context.Context, req *request.GetUserByID) (*entity.User, error) {
	ctx, span := s.tracer.Start(ctx, "service: GetUserByID")
	defer span.End()

	user, err := s.repo.GetUserByID(ctx, req)
	if err != nil {
		switch {
		case errors.Is(err, userrepo.ErrUserNotFound):
			return nil, apperror.NewNotFound("user not found")
		default:
			return nil, apperror.NewInternal("failed to find user", err)
		}
	}

	return user, nil
}

func (s *User) GetUserByUsername(ctx context.Context, username string) (*entity.User, error) {
	ctx, span := s.tracer.Start(ctx, "service: GetUserByUsername")
	defer span.End()

	user, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		switch {
		case errors.Is(err, userrepo.ErrUserNotFound):
			return nil, apperror.NewNotFound("user not found")
		default:
			return nil, apperror.NewInternal("failed to find user", err)
		}
	}

	return user, nil
}

func (s *User) DeleteUserByUsername(ctx context.Context, username string) error {
	ctx, span := s.tracer.Start(ctx, "service: DeleteUserByUsername")
	defer span.End()

	err := s.repo.DeleteUserByUsername(ctx, username)
	if err != nil {
		switch {
		case errors.Is(err, userrepo.ErrUserNotFound):
			return apperror.NewNotFound("user not found")
		default:
			return apperror.NewInternal("failed to find user", err)
		}
	}
	return nil
}

func (s *User) UpdateUserInfo(ctx context.Context, req *request.UpdateUserInfo, username string) (*entity.User, error) {
	ctx, span := s.tracer.Start(ctx, "service: UpdateUserInfo")
	defer span.End()

	user, err := s.repo.UpdateUserInfo(ctx, req, username) // TODO: hash
	if err != nil {
		switch {
		case errors.Is(err, userrepo.ErrUserNotFound):
			return nil, apperror.NewNotFound("user not found")
		default:
			return nil, apperror.NewInternal("failed to find user", err)
		}
	}

	return user, nil
}
