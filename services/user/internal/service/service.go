package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/myacey/jxgercorp-banking/shared/converters"
	tokenpb "github.com/myacey/jxgercorp-banking/shared/proto/token"
	"github.com/myacey/jxgercorp-banking/shared/sharedmodels"
	"github.com/myacey/jxgercorp-banking/user/internal/models"
	"github.com/myacey/jxgercorp-banking/user/internal/repository"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var ErrInvalidEmailOrPassword = errors.New("invalid email or password")

const tokenTTL = time.Duration(24 * time.Hour)

type ServiceInterface interface {
	CreateUser(c *gin.Context, newUser *models.UserUnhashed) (*sharedmodels.User, error)
	Login(c *gin.Context, username, password string) (string, error)

	GetUserByID(c *gin.Context, id int64) (*sharedmodels.User, error)
	GetUserByUsername(c *gin.Context, username string) (*sharedmodels.User, error)
	DeleteUserByUsername(c *gin.Context, username string) error
	UpdateUserInfo(c *gin.Context, username string, newHashedPassword string, newEmail string) (*sharedmodels.User, error)
}

type Service struct {
	userRepo        repository.UserRepository
	tokenServiceRPC tokenpb.TokenServiceClient
	lg              *zap.SugaredLogger
}

func NewService(ur repository.UserRepository, tk tokenpb.TokenServiceClient, lg *zap.SugaredLogger) ServiceInterface {
	return &Service{
		userRepo:        ur,
		tokenServiceRPC: tk,
		lg:              lg,
	}
}

func (s *Service) CreateUser(c *gin.Context, newUser *models.UserUnhashed) (*sharedmodels.User, error) {
	dbUser, err := s.userRepo.CreateUser(c, newUser)
	if err != nil {
		return nil, fmt.Errorf("cant create new user: %v", err)
	}

	user := converters.ConvertDBUserToModel(dbUser)
	return user, nil
}

func (s *Service) Login(c *gin.Context, username, password string) (string, error) {
	usr, err := s.userRepo.GetUserByUsername(c, username)
	if err != nil {
		return "", err
	}

	if password != usr.HashedPassword { // TODO: change
		return "", ErrInvalidEmailOrPassword
	}

	genTokenReq := &tokenpb.GenerateTokenRequest{
		Username: username,
		Ttl:      timestamppb.New(time.Now().Add(tokenTTL)),
	}
	resp, err := s.tokenServiceRPC.GenerateToken(c.Request.Context(), genTokenReq)

	return resp.Token, err
}

func (s *Service) GetUserByID(c *gin.Context, id int64) (*sharedmodels.User, error) {
	dbUser, err := s.userRepo.GetUserByID(c, id)
	if err != nil {
		return nil, fmt.Errorf("cant find user by id in db: %v", err)
	}

	user := converters.ConvertDBUserToModel(dbUser)
	return user, nil
}

func (s *Service) GetUserByUsername(c *gin.Context, username string) (*sharedmodels.User, error) {
	dbUser, err := s.userRepo.GetUserByUsername(c, username)
	if err != nil {
		return nil, fmt.Errorf("cant find user by username in db: %v", err)
	}

	user := converters.ConvertDBUserToModel(dbUser)
	return user, nil
}

func (s *Service) DeleteUserByUsername(c *gin.Context, username string) error {
	err := s.userRepo.DeleteUserByUsername(c, username)
	if err != nil {
		return fmt.Errorf("cant delete user by username: %v", err)
	}
	return nil
}

func (s *Service) UpdateUserInfo(c *gin.Context, username string, newEmail string, newPassword string) (*sharedmodels.User, error) {
	dbUser, err := s.userRepo.UpdateUserInfo(c, username, newEmail, newPassword) // TODO: hash
	if err != nil {
		return nil, fmt.Errorf("cant update user info: %v", err)
	}

	user := converters.ConvertDBUserToModel(dbUser)
	return user, nil
}
