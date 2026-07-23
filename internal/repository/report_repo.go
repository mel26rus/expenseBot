package repository

import (
	"context"
	"expense-bot/internal/model"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ReportRepo struct {
	db *pgxpool.Pool
}

func NewReportRepo(db *pgxpool.Pool) *ReportRepo {
	return &ReportRepo{db: db}
}

func (r *ReportRepo) GetAccountTransactions(
	ctx context.Context,
	UserID int64,
	DateStart time.Time,
	DateEnd time.Time,
) ([]*model.TransactionsReport, error) {

	rows, err := r.db.Query(ctx, `
		SELECT
			t.account_id,
			LOWER(t.comment) AS category,
			SUM(t.amount) AS amount

		FROM transactions t

		WHERE t.user_id = $1
		AND t.created_at >= $2
		AND t.created_at <  $3

		GROUP BY
			t.account_id,
			LOWER(t.comment)

		ORDER BY
			t.account_id,
			amount;
	`,
		UserID,
		DateStart,
		DateEnd,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*model.TransactionsReport
	for rows.Next() {
		var tr model.TransactionsReport
		err := rows.Scan(&tr.AccountId, &tr.Category, &tr.Amount)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, &tr)
	}

	return transactions, nil
}

func (r *ReportRepo) GetUserAccounts(
	ctx context.Context,
	UserId int64,
	StartDate time.Time,
	EndDate time.Time,
) ([]*model.AccountReport, error) {

	rows, err := r.db.Query(ctx, `
			select
				a.id,
				a.name || ' (' || UPPER(c.code) || ')' as title,
				SUM(t.amount) as balance,
				coalesce(SUM(t.amount)
					filter (
						where t.created_at >= $2
						and t.created_at < $3
						and t.amount > 0
					), 0) as income,
				coalesce(ABS(SUM(t.amount)
					filter (
						where t.created_at >= $2
						and t.created_at < $3
						and t.amount < 0
					)), 0) as expense
			from
				accounts a
			join currencies c on c.id = a.currency_id
			left join transactions t on t.account_id = a.id
			where
				a.user_id = $1
			group by
				a.id,
				a.name,
				c.code
			order by
				a.id;
	`,
		UserId,
		StartDate,
		EndDate,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []*model.AccountReport
	for rows.Next() {
		var ac model.AccountReport
		err := rows.Scan(&ac.AccountId, &ac.Title, &ac.Balance, &ac.Income, &ac.Expense)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, &ac)
	}

	return accounts, nil
}

func (r *ReportRepo) GetUsersHasTransactionsDaily(ctx context.Context, StartDate time.Time, EndDate time.Time) ([]int64, error) {
	rows, err := r.db.Query(ctx, `
		SELECT DISTINCT
			t.user_id
		FROM transactions t
		JOIN users u ON u.id = t.user_id
		WHERE
			u.isdailyreport = TRUE
			AND t.created_at >= $1
			AND t.created_at < $2
		ORDER BY t.user_id
	`,
		StartDate,
		EndDate,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []int64

	for rows.Next() {

		var id int64

		if err := rows.Scan(&id); err != nil {
			return nil, err
		}

		users = append(users, id)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *ReportRepo) GetUsersHasTransactionsMonthly(ctx context.Context, StartDate time.Time, EndDate time.Time) ([]int64, error) {
	rows, err := r.db.Query(ctx, `
		SELECT DISTINCT
			t.user_id
		FROM transactions t
		JOIN users u ON u.id = t.user_id
		WHERE
			u.ismonhtlyreport = TRUE
			AND t.created_at >= $1
			AND t.created_at < $2
		ORDER BY t.user_id
	`,
		StartDate,
		EndDate,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []int64

	for rows.Next() {

		var id int64

		if err := rows.Scan(&id); err != nil {
			return nil, err
		}

		users = append(users, id)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
