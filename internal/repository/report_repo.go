package repository

import (
	"context"
	"expense-bot/internal/model"
	"fmt"
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
	tgId int64,
	DateStart time.Time,
	DateEnd time.Time,
) ([]*model.TransactionsReport, error) {

	rows, err := r.db.Query(ctx, `
		SELECT
			t.account_id,
			case when t.amount > 0 then 'income' else 'expense' end as expense_type,
			LOWER(t.comment) AS category,
			SUM(t.amount) AS amount
		FROM transactions t
		join users u on u.id = t.user_id 
		WHERE u.telegram_id = $1
		AND t.created_at >= $2
		AND t.created_at <  $3

		GROUP BY
			1,
			2,
			3

		ORDER BY
			t.account_id,
			amount;
	`,
		tgId,
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
		err := rows.Scan(&tr.AccountId, &tr.Expense_type, &tr.Category, &tr.Amount)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, &tr)
	}

	return transactions, nil
}

func (r *ReportRepo) GetUserAccounts(
	ctx context.Context,
	tgId int64,
	StartDate time.Time,
	EndDate time.Time,
) ([]*model.AccountReport, error) {

	rows, err := r.db.Query(ctx, `
		select
			a.id,
			a.name || ' (' || UPPER(c.code) || ')' as title,
			tb.balance,
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
		from users u
		join accounts a on u.id = a.user_id 
		join currencies c on c.id = a.currency_id
		join transactions t on t.account_id = a.id
		join (  select t.account_id, t.user_id, sum(t.amount) as balance 
				from transactions t 
				group by t.account_id, t.user_id) tb on tb.account_id = t.account_id and t.user_id = u.id
		where t.created_at >= $2
			and t.created_at < $3
			and u.telegram_id = $1
		group by
			a.id,
			a.name,
			c.code,
			tb.balance
		order by
			a.id;
	`,
		tgId,
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

func (r *ReportRepo) GetUsersHasTransactionsTgIDs(ctx context.Context, StartDate time.Time, EndDate time.Time) ([]int64, error) {

	duration := EndDate.Sub(StartDate)
	days := int(duration.Hours() / 24)
	typeReportCondition := ``
	if days > 1 {
		typeReportCondition = ` u.ismonhtlyreport = TRUE `
	} else {
		typeReportCondition = ` u.isdailyreport = TRUE `
	}
	sql := fmt.Sprintf(`
		SELECT DISTINCT
			u.telegram_id 
		FROM transactions t
		JOIN users u ON u.id = t.user_id
		WHERE
			%s
			AND t.created_at >= $1
			AND t.created_at < $2
		ORDER BY u.telegram_id
		`,
		typeReportCondition,
	)

	rows, err := r.db.Query(ctx,
		sql,
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
