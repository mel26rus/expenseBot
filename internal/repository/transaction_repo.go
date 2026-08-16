package repository

import (
	"context"
	"expense-bot/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionRepo struct {
	db *pgxpool.Pool
}

func NewTransactionRepo(db *pgxpool.Pool) *TransactionRepo {
	return &TransactionRepo{db: db}
}

func (r *TransactionRepo) CreateTx(
	ctx context.Context,
	userID int64,
	accountID int64,
	txType string,
	amount float64,
	comment string,
	txGUID *string,
) error {

	_, err := r.db.Exec(ctx, `
		INSERT INTO transactions (user_id, account_id, tx_type, amount, comment, tx_guid)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, userID, accountID, txType, amount, comment, txGUID)

	return err
}

func (r *TransactionRepo) GetLastComments(
	ctx context.Context,
	userID int64,
	txType string,
	limit int,
) ([]model.Comments, error) {

	rows, err := r.db.Query(ctx, `
		select
			max(id) as id,
			comment
		from
			transactions
		where
			user_id = $1
			and comment is not null
			and comment != ''
			and $2 = tx_type
		group by
			comment
		order by
			id desc
		limit $3
	`, userID, txType, limit)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.Comments

	for rows.Next() {
		var c model.Comments
		rows.Scan(&c.TransactionID, &c.Comment)
		result = append(result, c)
	}

	return result, nil
}

func (r *TransactionRepo) GetCommentByTxID(ctx context.Context, txID int64) (string, error) {
	var result string
	err := r.db.QueryRow(ctx, `
		SELECT comment
		FROM transactions
		WHERE id = $1
	`, txID).Scan(&result)
	if err != nil {
		return "", err
	}
	return result, nil
}

func (r *TransactionRepo) GetLastUserTx(ctx context.Context, userID int64) ([]*model.Transaction, error) {
	rows, err := r.db.Query(ctx, `
		select
			t.id,
			'💳 ' || a."name" || CASE WHEN t.comment IS NOT NULL AND t.comment != '' THEN ' 💬 ' || t.comment ELSE '' END as title,
			t.amount 
		from
			transactions t
			join accounts a on a.id = t.account_id 
		where
			t.user_id = $1
			and t.tx_guid is null
		order by
			id desc
		limit 5
	`, userID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*model.Transaction

	for rows.Next() {
		var t model.Transaction
		rows.Scan(&t.ID, &t.Comment, &t.Amount)
		result = append(result, &t)
	}

	return result, nil
}

func (r *TransactionRepo) DeleteTx(ctx context.Context, txID int64) error {
	_, err := r.db.Exec(ctx, `
		delete from transactions where id = $1
	`, txID)
	return err
}
