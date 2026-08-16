package exchangerate

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"
)

const constReatesUrl = "https://www.cbr-xml-daily.ru/latest.js"

type Rates struct {
	Date   time.Time
	Values map[string]float64
}

type cbrLatestResponse struct {
	Timestamp int64              `json:"timestamp"`
	Base      string             `json:"base"`
	Rates     map[string]float64 `json:"rates"`
}

// FetchRates принимает контекст, делает запрос, пересчитывает базу на USD и возвращает мапу
func FetchRates(ctx context.Context) *Rates {

	client := &http.Client{Timeout: 5 * time.Second}

	req, err := http.NewRequestWithContext(ctx, "GET", constReatesUrl, nil)
	if err != nil {
		slog.Error("fetch_rates.create_request", "Error", err)
		return nil
	}

	resp, err := client.Do(req)
	if err != nil {
		slog.Error("fetch_rates.do_request", "Error", err)
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.Error("fetch_rates.status_code", "Status", resp.Status)
		return nil
	}

	var rawData cbrLatestResponse
	if err := json.NewDecoder(resp.Body).Decode(&rawData); err != nil {
		slog.Error("fetch_rates.decode_json", "Error", err)
		return nil
	}

	apiDateTime := time.Unix(rawData.Timestamp, 0).UTC()

	// Находим стоимость доллара в рублях (разворачиваем обратный курс)
	usdInRubRate, hasUsd := rawData.Rates["USD"]
	if !hasUsd || usdInRubRate == 0 {
		slog.Error("fetch_rates.convert", "Error", "USD rate not found or zero")
		return nil
	}
	usdPriceInRub := 1.0 / usdInRubRate

	// Создаем результирующую мапу с базой USD
	usdBaseRates := make(map[string]float64)

	// Пересчитываем все валюты в базу USD
	for currency, rateToRub := range rawData.Rates {
		usdBaseRates[currency] = rateToRub * usdPriceInRub
	}

	// Явно задаем базовые валюты для точности
	usdBaseRates["USD"] = 1.0
	usdBaseRates["RUB"] = usdPriceInRub
	usdBaseRates["USDT"] = 1.0

	return &Rates{
		Date:   apiDateTime,
		Values: usdBaseRates,
	}
}
