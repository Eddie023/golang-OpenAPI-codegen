package service

import (
	"context"

	"github.com/eddie023/wex-tag/ent"
	"github.com/eddie023/wex-tag/pkg/types"
)

type Service struct {
	Ent *ent.Client
}

func (s *Service) CreatePurchase(ctx context.Context, payload types.CreateNewPurchaseTransaction) (types.Transaction, error) {
	return types.Transaction{}, nil
}
