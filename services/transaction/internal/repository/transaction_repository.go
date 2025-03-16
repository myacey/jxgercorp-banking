package repository

import (
	"github.com/gin-gonic/gin"
	db "github.com/myacey/jxgercorp-banking/services/db/sqlc"
)

type TransactionRepository interface {
	CreateTransactionTX(c *gin.Context, fromUser, toUser string, amount int64) (*db.Transaction, error)
	// SearchOutcomeTransactions(c *gin.Context, fromUser string, offset int, limit int) ([]*db.Transaction, error)
	// SearchIncomeTransactions(c *gin.Context, toUser string, offset int, limit int) ([]*db.Transaction, error)
	SearchTransactionsWithUser(c *gin.Context, username string, offset, limit int) ([]*db.Transaction, error)
}
