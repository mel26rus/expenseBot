package service

import (
	"context"
	"expense-bot/internal/model"
	"expense-bot/internal/repository"
)

type CurrencyService struct {
	repo *repository.CurrencyRepo
}

func NewCurrencyService(r *repository.CurrencyRepo) *CurrencyService {
	return &CurrencyService{repo: r}
}

func (s *CurrencyService) GetAll(ctx context.Context) ([]*model.Currency, error) {
	return s.repo.GetAll(ctx)
}

func (s *CurrencyService) Create(ctx context.Context, name string) (int64, error) {
	return s.repo.CreateCurrency(ctx, name)
}
