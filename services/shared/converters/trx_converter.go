package converters

import (
	db "github.com/myacey/jxgercorp-banking/services/db/sqlc"
	"github.com/myacey/jxgercorp-banking/services/shared/sharedmodels"
)

func ConvertDBTrxToModel(dbTrx *db.Transaction) *sharedmodels.Trx {
	return &sharedmodels.Trx{
		ID:        dbTrx.ID,
		FromUser:  dbTrx.FromUser,
		ToUser:    dbTrx.ToUser,
		Amount:    dbTrx.Amount,
		CreatedAt: dbTrx.CreatedAt,
	}
}
