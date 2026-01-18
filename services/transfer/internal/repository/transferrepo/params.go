package transferrepo

import "github.com/google/uuid"

type SearchTransfersWithAccountParams struct {
	CurrentAccountID uuid.UUID
	WithUsername     string
	WithAccountID    uuid.UUID
	Currency         string
	Offset           int32
	Limit            int32
}
