package service

import (
	"context"

	"github.com/eddie023/wex-tag/ent"
)

type Service struct {
	Ent *ent.Client
}

func (s *Service) CreatePurchase(ctx context.Context) (*ent.Transaction, error) {
	return nil, nil
}
