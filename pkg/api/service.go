package api

import (
	"context"

	"github.com/eddie023/wex-tag/ent"
	"github.com/eddie023/wex-tag/pkg/api/service"
	"github.com/eddie023/wex-tag/pkg/types"
)

//go:generate go run go.uber.org/mock/mockgen -destination=mocks/mock_transaction.go -package=mocks . TransactionService
type TransactionService interface {
	CreateNewPurchaseTransaction(ctx context.Context, payload types.CreateNewPurchaseTransaction) (types.Transaction, error)
}

//go:generate go run go.uber.org/mock/mockgen -destination=mocks/mock_exchange_rate.go -package=mocks . ExchangeRateService
type ExchangeRateService interface {
	GetExchangeRate(ctx context.Context, payload service.ExchangeRatePayload) (service.ExchangeRateResponse, error)
	ConvertCurrency(requestConversionPayload service.ExchangeRatePayload, transactionInfo *ent.Transaction, exchangeRateInfo service.ExchangeRateResponse) (types.GetPurchaseTransaction, error)
}
