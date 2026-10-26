package service

import (
	"context"
	"errors"
	"expense-bot/internal/exchangerate"
	"expense-bot/internal/model"
	"testing"
	"time"
)

type currencyRepoStub struct {
	id      int64
	hasRate bool
	saved   bool
	code    string
	date    time.Time
	value   float64
	err     error
}

func (r *currencyRepoStub) GetAll(context.Context) ([]*model.Currency, error) { return nil, nil }
func (r *currencyRepoStub) FindWithLatestRate(context.Context, string) (int64, bool, error) {
	return r.id, r.hasRate, r.err
}
func (r *currencyRepoStub) SaveWithRate(_ context.Context, code string, date time.Time, value float64) (int64, error) {
	r.saved, r.code, r.date, r.value = true, code, date, value
	return 7, r.err
}
func TestCreateCurrency(t *testing.T) {
	date := time.Date(2026, 10, 2, 0, 0, 0, 0, time.UTC)
	for _, tc := range []struct {
		name      string
		code      string
		id        int64
		hasRate   bool
		rates     *exchangerate.Rates
		wantErr   error
		wantFetch bool
		wantSave  bool
	}{
		{name: "existing current rate", code: "USD", id: 4, hasRate: true},
		{name: "new currency", code: " eur ", rates: &exchangerate.Rates{Date: date, Values: map[string]float64{"EUR": 0.9, "USD": 1}}, wantFetch: true, wantSave: true},
		{name: "missing rate", code: "EUR", id: 4, rates: &exchangerate.Rates{Date: date, Values: map[string]float64{"EUR": 0.9}}, wantFetch: true, wantSave: true},
		{name: "unsupported", code: "XYZ", rates: &exchangerate.Rates{Date: date, Values: map[string]float64{"USD": 1}}, wantErr: ErrUnsupportedCurrency, wantFetch: true},
		{name: "API unavailable", code: "EUR", wantErr: ErrRatesUnavailable, wantFetch: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			repo := &currencyRepoStub{id: tc.id, hasRate: tc.hasRate}
			fetched := false
			s := &CurrencyService{repo: repo, fetchRates: func(context.Context) *exchangerate.Rates { fetched = true; return tc.rates }}
			id, err := s.Create(context.Background(), tc.code)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("error = %v, want %v", err, tc.wantErr)
			}
			if fetched != tc.wantFetch || repo.saved != tc.wantSave {
				t.Fatalf("fetch=%v save=%v", fetched, repo.saved)
			}
			if tc.wantSave && (repo.code != "EUR" || !repo.date.Equal(date) || repo.value != 0.9 || id != 7) {
				t.Fatalf("incorrect saved currency: %+v, id=%d", repo, id)
			}
			if !tc.wantFetch && id != tc.id {
				t.Fatalf("id=%d", id)
			}
		})
	}
}
