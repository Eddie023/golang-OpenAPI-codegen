package api

import (
	"context"

	"github.com/eddie023/wex-tag/ent"
)

//go:generate go run go.uber.org/mock/mockgen -destination=mocks/mock_transaction.go -package=mocks . TransactionService
type TransactionService interface {
	CreatePurchase(ctx context.Context) (*ent.Transaction, error)
}
