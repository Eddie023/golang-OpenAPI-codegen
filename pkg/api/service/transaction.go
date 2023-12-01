package service

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/eddie023/wex-tag/ent"
	"github.com/eddie023/wex-tag/ent/transaction"
	"github.com/eddie023/wex-tag/pkg/apiout"
	"github.com/eddie023/wex-tag/pkg/types"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type Service struct {
	Ent *ent.Client
}

// CreatePurchase will store the request payload into database and return a new purchase transaction.
func (s *Service) CreateNewPurchaseTransaction(ctx context.Context, payload types.CreateNewPurchaseTransaction) (types.Transaction, error) {
	slog.Info("creating new purchase transaction", "amount", payload.Amount)

	amount, err := decimal.NewFromString(payload.Amount)
	if err != nil {
		return types.Transaction{}, apiout.BadRequest(errors.Wrap(err, fmt.Sprintf("unable to parse '%s'", payload.Amount)).Error())
	}

	// we are passing amount type as string for precision. Thus, we need to check for case
	// such as when user passes negative integer
	if amount.IsNegative() && !amount.IsZero() {
		return types.Transaction{}, apiout.BadRequest("amount cannot be negative number")
	}

	roundedAmount := RoundToNearestCent(amount)
	transaction, err := s.Ent.Transaction.Create().SetAmountInUsd(roundedAmount).SetDate(time.Now().UTC()).SetDescription(payload.Description).Save(ctx)
	if err != nil {
		return types.Transaction{}, err
	}

	slog.Info("successfully processed new purchase transaction", "transaction_id", transaction.ID)

	return types.Transaction{
		AmountInUSD: roundedAmount.String(),
		Date:        transaction.Date.UTC(),
		Description: transaction.Description,
		Id:          transaction.ID.String(),
	}, nil
}

// GetPurchaseDetailsByTransactionId will query the database to see if the purchase order with provided transaction id exist.
func (s *Service) GetPurchaseDetailsByTransactionId(ctx context.Context, id uuid.UUID) (*ent.Transaction, error) {
	slog.Info("fetching transaction details", "transaction_id", id)

	transaction, err := s.Ent.Transaction.Query().Where(transaction.ID(id)).First(ctx)
	if err != nil {
		return nil, apiout.NewRequestError(errors.New("given transaction id not found"), http.StatusNotFound)
	}

	return transaction, nil
}

// ParseStringToUUID will try to parse the provided string to UUID
func ParseStringToUUID(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}

// RoundToNearestCent will round the given decimal number to nearest cent.
func RoundToNearestCent(n decimal.Decimal) decimal.Decimal {
	return n.Round(2)
}
