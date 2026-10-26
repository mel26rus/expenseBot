package service

import (
	"context"
	"errors"
	"expense-bot/internal/exchangerate"
	"expense-bot/internal/model"
	"expense-bot/internal/repository"
	"strings"
	"time"
)

var ErrUnsupportedCurrency = errors.New("unsupported currency")
var ErrRatesUnavailable = errors.New("exchange rates unavailable")

type currencyRepository interface {
	GetAll(context.Context) ([]*model.Currency, error)
	FindWithLatestRate(context.Context, string) (int64, bool, error)
	SaveWithRate(context.Context, string, time.Time, float64) (int64, error)
}

type CurrencyService struct {
	repo       currencyRepository
	fetchRates func(context.Context) *exchangerate.Rates
}

func NewCurrencyService(r *repository.CurrencyRepo) *CurrencyService {
	return &CurrencyService{repo: r, fetchRates: exchangerate.FetchRates}
}

func (s *CurrencyService) GetAll(ctx context.Context) ([]*model.Currency, error) {
	return s.repo.GetAll(ctx)
}

func (s *CurrencyService) Create(ctx context.Context, name string) (int64, error) {
	name = strings.ToUpper(strings.TrimSpace(name))
	id, hasRate, err := s.repo.FindWithLatestRate(ctx, name)
	if err != nil {
		return 0, err
	}
	if id != 0 && hasRate {
		return id, nil
	}
	rates := s.fetchRates(ctx)
	if rates == nil {
		return 0, ErrRatesUnavailable
	}
	value, ok := rates.Values[name]
	if !ok {
		return 0, ErrUnsupportedCurrency
	}
	return s.repo.SaveWithRate(ctx, name, rates.Date, value)
}
