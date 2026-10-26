package repository

import (
	"context"
	"expense-bot/internal/model"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CurrencyRepo struct {
	db *pgxpool.Pool
}

// FindWithLatestRate checks whether this currency has a rate in the latest snapshot.
func (r *CurrencyRepo) FindWithLatestRate(ctx context.Context, code string) (int64, bool, error) {
	var id int64
	var hasRate bool
	err := r.db.QueryRow(ctx, `SELECT id, EXISTS (
 SELECT 1 FROM exchange_rates WHERE name = $1 AND date = (SELECT max(date) FROM exchange_rates)
 ) FROM currencies WHERE code = $1`, code).Scan(&id, &hasRate)
	if err == pgx.ErrNoRows {
		return 0, false, nil
	}
	return id, hasRate, err
}

// SaveWithRate stores only the requested currency and its API-dated rate atomically.
func (r *CurrencyRepo) SaveWithRate(ctx context.Context, code string, date time.Time, value float64) (int64, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)
	var id int64
	err = tx.QueryRow(ctx, `INSERT INTO currencies (code) VALUES ($1)
 ON CONFLICT (code) DO UPDATE SET code = EXCLUDED.code RETURNING id`, code).Scan(&id)
	if err != nil {
		return 0, err
	}
	if _, err = saveRate(ctx, tx, date, code, value); err != nil {
		return 0, err
	}
	if err = tx.Commit(ctx); err != nil {
		return 0, err
	}
	return id, nil
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

	return currencies, rows.Err()
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
