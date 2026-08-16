package scheduler

import (
	"context"
	"expense-bot/internal/exchangerate"
	"expense-bot/internal/service"
	"fmt"
	"log/slog"
	"time"
)

type GetRatesJob struct {
	CurrencyService     *service.CurrencyService
	ExchangeRateService *service.ExchangeRateService
}

func NewGetRatesJob(currencyService *service.CurrencyService, exchangeRateService *service.ExchangeRateService) *GetRatesJob {
	return &GetRatesJob{
		CurrencyService:     currencyService,
		ExchangeRateService: exchangeRateService,
	}
}

func (j *GetRatesJob) Name() string {
	return "Get Rates Job"
}

func (j *GetRatesJob) NextRun(now time.Time) time.Time {

	next := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		1, 0, 0, 0,
		now.Location(),
	)

	if !next.After(now) {
		next = next.AddDate(0, 0, 1)
	}
	// next := now.Add(time.Minute * 1)
	slog.Debug("j.NextRun", "name", j.Name(), "Next", next)
	return next

}

func (j *GetRatesJob) Run(ctx context.Context) error {
	fn_name := "GetRatesJob.Run"
	slog.Debug(fmt.Sprintf("+%s", fn_name))
	result := exchangerate.FetchRates(ctx)
	cur, err := j.CurrencyService.GetAll(ctx)
	if err != nil {
		slog.Error("GetRatesJob.Run_1", "Error", err)
		return err
	}

	for _, c := range cur {
		if val, ok := result.Values[c.Code]; ok {
			// Сохраняем: передаем текущую дату, код валюты и её курс
			_, err := j.ExchangeRateService.SaveRate(ctx, result.Date, c.Code, val)
			if err != nil {
				slog.Error("GetRatesJob.SaveRate", "Currency", c.Code, "Error", err)
			}
		}
	}

	slog.Debug(fmt.Sprintf("-%s", fn_name))
	return nil
}
