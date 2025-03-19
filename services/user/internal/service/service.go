package service

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/myacey/jxgercorp-banking/services/shared/converters"
	"github.com/myacey/jxgercorp-banking/services/shared/cstmerr"
	tokenpb "github.com/myacey/jxgercorp-banking/services/shared/proto/token"
	"github.com/myacey/jxgercorp-banking/services/shared/sharedmodels"
	"github.com/myacey/jxgercorp-banking/services/user/internal/confirmation"
	"github.com/myacey/jxgercorp-banking/services/user/internal/models"
	"github.com/myacey/jxgercorp-banking/services/user/internal/repository"
	"go.opentelemetry.io/otel/trace"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var ErrInvalidEmailOrPassword = errors.New("invalid email or password")

const tokenTTL = time.Duration(24 * time.Hour)

type ServiceInterface interface {
	CreateUser(c *gin.Context, newUser *models.UserUnhashed) (*sharedmodels.User, error)
	Login(c *gin.Context, username, password string) (string, error)

	ConfirmUserEmail(c *gin.Context, username, code string) (string, error)

	GetUserByID(c *gin.Context, id int64) (*sharedmodels.User, error)
	GetUserByUsername(c *gin.Context, username string) (*sharedmodels.User, error)
	DeleteUserByUsername(c *gin.Context, username string) error
	UpdateUserInfo(c *gin.Context, username string, newHashedPassword string, newEmail string) (*sharedmodels.User, error)
}

type Service struct {
	userRepo        repository.UserRepository
	tokenServiceRPC tokenpb.TokenServiceClient
	lg              *zap.SugaredLogger

	registerConfirmSrv confirmation.ConfirmationServiceInterface

	tracer trace.Tracer
}

func NewService(ur repository.UserRepository, tk tokenpb.TokenServiceClient, lg *zap.SugaredLogger, registerConfirmSrv confirmation.ConfirmationServiceInterface, tr trace.Tracer) ServiceInterface {
	return &Service{
		userRepo:           ur,
		tokenServiceRPC:    tk,
		lg:                 lg,
		registerConfirmSrv: registerConfirmSrv,
		tracer:             tr,
	}
}

func (s *Service) CreateUser(c *gin.Context, newUser *models.UserUnhashed) (*sharedmodels.User, error) {
	ctx, span := s.tracer.Start(c.Request.Context(), "service: CreateUser")
	defer span.End()
	c.Request = c.Request.WithContext(ctx)

	unhashedPassword := newUser.Password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(unhashedPassword), bcrypt.DefaultCost)
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrHashTooShort):
			return nil, cstmerr.New(http.StatusBadRequest, "password too short", nil)
		case errors.Is(err, bcrypt.ErrPasswordTooLong):
			return nil, cstmerr.New(http.StatusBadRequest, "password too long", nil)
		default:
			return nil, cstmerr.ErrUnknown
		}
	}

	dbUser, err := s.userRepo.CreateUser(c, newUser.Username, newUser.Email, string(hashedPassword))
	if err != nil {
		return nil, fmt.Errorf("cant create new user: %v", err)
	}

	user := converters.ConvertDBUserToModel(dbUser)

	err = s.registerConfirmSrv.GenerateAccountConfirmation(c, newUser.Username, newUser.Email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) Login(c *gin.Context, username, password string) (string, error) {
	ctx, span := s.tracer.Start(c.Request.Context(), "service: Login")
	defer span.End()
	c.Request = c.Request.WithContext(ctx)

	usr, err := s.userRepo.GetUserByUsername(c, username)
	if err != nil {
		return "", err
	}

	if usr.Status == sharedmodels.UserStatusPending {
		return "", cstmerr.New(http.StatusUnauthorized, "check email for account comfirmation", nil)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(usr.HashedPassword), []byte(password)); err != nil { // TODO: change
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return "", ErrInvalidEmailOrPassword
		default:
			return "", cstmerr.ErrUnknown
		}
	}

	genTokenReq := &tokenpb.GenerateTokenRequest{
		Username: username,
		Ttl:      timestamppb.New(time.Now().Add(tokenTTL)),
	}
	resp, err := s.tokenServiceRPC.GenerateToken(c.Request.Context(), genTokenReq)
	if err != nil || resp.Token == "" {
		return "", cstmerr.New(http.StatusUnauthorized, cstmerr.ErrUnknown.Error(), err)
	}

	return resp.Token, nil
}

func (s *Service) ConfirmUserEmail(c *gin.Context, username, code string) (string, error) {
	ctx, span := s.tracer.Start(c.Request.Context(), "service: ConfirmUserEmail")
	defer span.End()
	c.Request = c.Request.WithContext(ctx)

	err := s.registerConfirmSrv.CheckConfirmCode(c, username, code)
	if err != nil {
		return "", err
	}

	_, err = s.userRepo.UpdateUserInfo(c, username, "", "", false)
	if err != nil {
		return "", err
	}

	return "success", nil
}

func (s *Service) GetUserByID(c *gin.Context, id int64) (*sharedmodels.User, error) {
	ctx, span := s.tracer.Start(c.Request.Context(), "service: GetUserByID")
	defer span.End()
	c.Request = c.Request.WithContext(ctx)

	dbUser, err := s.userRepo.GetUserByID(c, id)
	if err != nil {
		return nil, fmt.Errorf("cant find user by id in db: %v", err)
	}

	user := converters.ConvertDBUserToModel(dbUser)
	return user, nil
}

func (s *Service) GetUserByUsername(c *gin.Context, username string) (*sharedmodels.User, error) {
	ctx, span := s.tracer.Start(c.Request.Context(), "service: GetUserByUsername")
	defer span.End()
	c.Request = c.Request.WithContext(ctx)

	dbUser, err := s.userRepo.GetUserByUsername(c, username)
	if err != nil {
		return nil, fmt.Errorf("cant find user by username in db: %v", err)
	}

	user := converters.ConvertDBUserToModel(dbUser)
	return user, nil
}

func (s *Service) DeleteUserByUsername(c *gin.Context, username string) error {
	ctx, span := s.tracer.Start(c.Request.Context(), "service: DeleteUserByUsername")
	defer span.End()
	c.Request = c.Request.WithContext(ctx)

	err := s.userRepo.DeleteUserByUsername(c, username)
	if err != nil {
		return fmt.Errorf("cant delete user by username: %v", err)
	}
	return nil
}

func (s *Service) UpdateUserInfo(c *gin.Context, username string, newEmail string, newPassword string) (*sharedmodels.User, error) {
	ctx, span := s.tracer.Start(c.Request.Context(), "service: UpdateUserInfo")
	defer span.End()
	c.Request = c.Request.WithContext(ctx)

	dbUser, err := s.userRepo.UpdateUserInfo(c, username, newEmail, newPassword, false) // TODO: hash
	if err != nil {
		return nil, fmt.Errorf("cant update user info: %v", err)
	}

	user := converters.ConvertDBUserToModel(dbUser)
	return user, nil
}
