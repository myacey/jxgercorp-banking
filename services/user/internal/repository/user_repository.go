package repository

import (
	"github.com/gin-gonic/gin"
	db "github.com/myacey/jxgercorp-banking/db/sqlc"
	"github.com/myacey/jxgercorp-banking/user/internal/models"
)

type UserRepository interface {
	// postgres
	CreateUser(c *gin.Context, newUser *models.UserUnhashed) (*db.User, error)
	GetUserByID(c *gin.Context, id int64) (*db.User, error)
	GetUserByUsername(c *gin.Context, username string) (*db.User, error)
	DeleteUserByUsername(c *gin.Context, username string) error
	UpdateUserInfo(c *gin.Context, username string, newEmail string, newHashedPassword string) (*db.User, error)
}
