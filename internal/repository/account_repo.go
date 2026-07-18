package repository

import (
	"context"
	"expense-bot/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AccountRepo struct {
	db *pgxpool.Pool
}

func NewAccountRepo(db *pgxpool.Pool) *AccountRepo {
	return &AccountRepo{db: db}
}

func (r *AccountRepo) GetAccountsByUserID(ctx context.Context, userID int64) ([]*model.Account, error) {

	rows, err := r.db.Query(ctx, `
		SELECT accounts.id, name||' ('||currencies.code||')' as name, currency_id
		FROM accounts
		join currencies on accounts.currency_id = currencies.id
		WHERE user_id = $1
	`, userID)
	if err != nil {
		return []*model.Account{}, err
	}
	defer rows.Close()

	var accounts []*model.Account
	for rows.Next() {
		var a model.Account
		// slog.Debug("GetAccountsByUserID.Scanning user account row")
		err := rows.Scan(&a.ID, &a.Name, &a.CurrencyID)
		// slog.Debug("GetAccountsByUserID.Scanned user account row", "id", a.ID, "name", a.Name, "currencyID", a.CurrencyID, "err", err)
		if err != nil {
			return []*model.Account{}, err
		}
		// slog.Debug("GetAccountsByUserID.Scanned user account row", "id", a.ID, "name", a.Name, "currencyID", a.CurrencyID)
		accounts = append(accounts, &a)
	}

	return accounts, nil
}

func (r *AccountRepo) CreateAccount(ctx context.Context, userID int64, name string, currencyID int64) (int64, error) {
	var accountID int64

	err := r.db.QueryRow(ctx, `
		INSERT INTO accounts (user_id, name, currency_id)
		VALUES ($1, $2, $3)
		RETURNING id
	`, userID, name, currencyID).Scan(&accountID)

	if err != nil {
		return 0, err
	}

	return accountID, nil
}

func (r *AccountRepo) GetAccountTitleByID(ctx context.Context, accountID int64) (model.Account, error) {
	var a model.Account
	err := r.db.QueryRow(ctx, `
		SELECT accounts.id, accounts.user_id, accounts.currency_id, name||' ('||currencies.code||')' as name
		FROM accounts
		join currencies on accounts.currency_id = currencies.id
		WHERE accounts.id = $1
	`, accountID).Scan(&a.ID, &a.UserID, &a.CurrencyID, &a.Name)
	return a, err
}

func (r *AccountRepo) GetAccountBalanceByID(ctx context.Context, accountID int64) (float64, error) {
	var balance float64
	err := r.db.QueryRow(ctx, `
		SELECT SUM(COALESCE(amount, 0))
		FROM transactions
		WHERE account_id = $1
	`, accountID).Scan(&balance)
	return balance, err
}

func (r *AccountRepo) GetAccountCurrencyByID(ctx context.Context, accountID int64) (int64, error) {
	var currencyId int64
	err := r.db.QueryRow(ctx, `
	 select currency_id
	 from accounts
	 where id = $1
	`, accountID).Scan(&currencyId)
	return currencyId, err
}
