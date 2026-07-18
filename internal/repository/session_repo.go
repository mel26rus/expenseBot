package repository

import (
	"context"
	"expense-bot/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SessionRepo struct {
	db *pgxpool.Pool
}

func NewSessionRepo(db *pgxpool.Pool) *SessionRepo {
	return &SessionRepo{db: db}
}

func (r *SessionRepo) GetUserSession(ctx context.Context, userID int64) (model.Session, error) {
	var s model.Session

	err := r.db.QueryRow(ctx, `
		SELECT user_id, state, payload, edit_message_id
		FROM sessions
		WHERE user_id = $1
	`, userID).Scan(&s.UserID, &s.State, &s.Payload, &s.EditMessageId)

	if err != nil {
		// если нет записи → возвращаем пустую сессию
		return model.Session{UserID: userID, State: "", EditMessageId: 0}, nil
	}

	return s, nil
}

func (r *SessionRepo) Save(ctx context.Context, s model.Session) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO sessions (user_id, state, payload, edit_message_id)
		VALUES ($1,$2,$3,$4)
		ON CONFLICT (user_id) DO UPDATE
		SET state=$2, payload=$3, edit_message_id=$4
	`, s.UserID, s.State, s.Payload, s.EditMessageId)

	return err
}

func (r *SessionRepo) SetEditMessageId(ctx context.Context, userID int64, messageID int64) error {
	_, err := r.db.Exec(ctx, `
		UPDATE sessions
		SET edit_message_id = $1
		WHERE user_id = $2
	`, messageID, userID)

	return err
}
