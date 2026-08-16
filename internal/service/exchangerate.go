package service

import (
	"context"
	"expense-bot/internal/model"
	"expense-bot/internal/repository"
	"time"
)

type ExchangeRateService struct {
	repo *repository.ExchangeRateRepo
}

func NewExchangeRateService(r *repository.ExchangeRateRepo) *ExchangeRateService {
	return &ExchangeRateService{repo: r}
}

func (s *ExchangeRateService) GetExchangeRates(ctx context.Context) ([]*model.ExchangeRate, error) {
	return s.repo.GetExchangeRates(ctx)
}

func (s *ExchangeRateService) SaveRate(ctx context.Context, date time.Time, name string, value float64) (int64, error) {
	return s.repo.SaveRate(ctx, date, name, value)
}
