package api

import (
	"context"

	"github.com/eddie023/wex-tag/pkg/types"
)

//go:generate go run go.uber.org/mock/mockgen -destination=mocks/mock_transaction.go -package=mocks . TransactionService
type TransactionService interface {
	CreateNewPurchase(ctx context.Context, payload types.CreateNewPurchaseTransaction) (types.Transaction, error)
}
