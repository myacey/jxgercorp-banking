package converters

import (
	db "github.com/myacey/jxgercorp-banking/services/db/sqlc"
	"github.com/myacey/jxgercorp-banking/services/shared/sharedmodels"
)

func ConvertDBUserToModel(dbUser *db.User) *sharedmodels.User {
	return &sharedmodels.User{
		ID:        dbUser.ID,
		Username:  dbUser.Username,
		Email:     dbUser.Email,
		Balance:   dbUser.Balance,
		CreatedAt: dbUser.CreatedAt,
	}
}
