package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/eddie023/wex-tag/ent"
	"github.com/eddie023/wex-tag/pkg/apiout"
	"github.com/eddie023/wex-tag/pkg/types"
	"github.com/google/uuid"
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
		return types.Transaction{}, apiout.BadRequest(err.Error())
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
		AmountInUSD: fmt.Sprintf("%.2f", roundedAmount.InexactFloat64()),
		Date:        transaction.Date,
		Description: transaction.Description,
		Id:          transaction.ID.String(),
	}, nil
}

func ParseStringToUUID(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}

// RoundToNearestCent will round the given decimal number to nearest cent.
func RoundToNearestCent(n decimal.Decimal) decimal.Decimal {
	return n.Round(2)
}
