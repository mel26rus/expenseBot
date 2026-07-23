package repository

import (
	"context"

	"expense-bot/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) GetOrCreate(ctx context.Context, telegramID int64) (model.User, error) {
	var user model.User

	err := r.db.QueryRow(ctx, `
		INSERT INTO users (telegram_id)
		VALUES ($1)
		ON CONFLICT (telegram_id) DO UPDATE SET telegram_id = EXCLUDED.telegram_id
		RETURNING id, telegram_id
	`, telegramID).Scan(&user.ID, &user.TelegramID)

	return user, err
}

func (r *UserRepo) GetUserSettings(ctx context.Context, id int64) (model.UserSettings, error) {
	var userSettings model.UserSettings
	err := r.db.QueryRow(ctx, `
		SELECT id, telegram_id, isdailyreport, ismonhtlyreport
		FROM public.users
		WHERE id=$1
	`, id).Scan(&userSettings.ID, &userSettings.TelegramID, &userSettings.IsDailyReport, &userSettings.IsMonthlyReport)
	return userSettings, err
}

func (r *UserRepo) ChangeDailyReportConfig(ctx context.Context, id int64) (model.UserSettings, error) {
	_, err := r.db.Exec(ctx, `
		update public.users set isdailyreport = not isdailyreport where id = $1
	`, id)
	userSettings, err := r.GetUserSettings(ctx, id)

	return userSettings, err
}

func (r *UserRepo) ChangeMonthlyReportConfig(ctx context.Context, id int64) (model.UserSettings, error) {
	_, err := r.db.Exec(ctx, `
		update public.users set ismonhtlyreport = not ismonhtlyreport where id = $1
	`, id)
	userSettings, err := r.GetUserSettings(ctx, id)

	return userSettings, err
}

const newSqlForGoC = `
WITH ins AS (
    INSERT INTO users (
        telegram_id,
        chat_id
    )
    VALUES (
        $1,
    DO NOTHING
    RETURNING
        id,
        telegram_id,
        chat_id
),
upd AS (
    UPDATE users
    SET chat_id = $2
    WHERE telegram_id = $1
      AND chat_id IS DISTINCT FROM $2
    RETURNING
        id,
        telegram_id,
        chat_id
)
SELECT
    id,
    telegram_id,
    chat_id
FROM ins

UNION ALL

SELECT
    id,
    telegram_id,
    chat_id
FROM upd

UNION ALL

SELECT
    id,
    telegram_id,
    chat_id
FROM users
WHERE telegram_id = $1
LIMIT 1;
`
