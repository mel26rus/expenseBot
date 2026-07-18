package repository

import (
	"context"
	"expense-bot/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CurrencyRepo struct {
	db *pgxpool.Pool
}

func NewCurrencyRepo(db *pgxpool.Pool) *CurrencyRepo {
	return &CurrencyRepo{db: db}
}

func (r *CurrencyRepo) GetAll(ctx context.Context) ([]*model.Currency, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, code
		FROM currencies
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var currencies []*model.Currency
	for rows.Next() {
		var c model.Currency
		err := rows.Scan(&c.ID, &c.Code)
		if err != nil {
			return nil, err
		}
		currencies = append(currencies, &c)
	}

	return currencies, nil
}

func (r *CurrencyRepo) CreateCurrency(ctx context.Context, name string) (int64, error) {
	var id int64

	err := r.db.QueryRow(ctx, `
		INSERT INTO currencies (code)
		VALUES ($1)
		ON CONFLICT (code) DO UPDATE SET code = EXCLUDED.code
		RETURNING id
	`, name).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}
