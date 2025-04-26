package transferrepo

import "github.com/google/uuid"

type SearchTransfersWithAccountParams struct {
	AccountID uuid.UUID
	Offset    int32
	Limit     int32
}
