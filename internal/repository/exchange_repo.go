package repository

import (
	"context"
	"expense-bot/internal/model"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ExchangeRateRepo struct {
	db *pgxpool.Pool
}

func NewExchangeRateRepo(db *pgxpool.Pool) *ExchangeRateRepo {
	return &ExchangeRateRepo{db: db}
}

func (r *ExchangeRateRepo) getLastActualRateDate(ctx context.Context) (time.Time, error) {
	sql := `
	 select max(date) from exchange_rates
	`
	var res time.Time
	err := r.db.QueryRow(ctx, sql).Scan(res)
	return res, err
}

func (r *ExchangeRateRepo) GetExchangeRates(ctx context.Context) ([]*model.ExchangeRate, error) {
	date, _ := r.getLastActualRateDate(ctx)
	rows, err := r.db.Query(ctx, `
		SELECT date, name, value
		FROM exchange_rates
		where date = $1
	`, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var exRate []*model.ExchangeRate
	for rows.Next() {
		var e model.ExchangeRate
		err := rows.Scan(&e.Date, &e.Name, &e.Value)
		if err != nil {
			return nil, err
		}
		exRate = append(exRate, &e)
	}

	return exRate, nil
}

func (r *ExchangeRateRepo) SaveRate(ctx context.Context, date time.Time, name string, value float64) (int64, error) {

	var id int64
	err := r.db.QueryRow(ctx, `
		INSERT INTO exchange_rates (date, name, value)
		VALUES ($1, $2, $3)
		RETURNING id
	`, date, name, value).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, err
}
