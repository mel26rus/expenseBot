package service

import (
	"context"
	"encoding/json"
	"expense-bot/internal/model"
	"expense-bot/internal/repository"
)

type SessionService struct {
	repo *repository.SessionRepo
}

func NewSessionService(r *repository.SessionRepo) *SessionService {
	return &SessionService{repo: r}
}

func (s *SessionService) GetUserSession(ctx context.Context, userID int64) (model.Session, error) {
	return s.repo.GetUserSession(ctx, userID)
}

func (s *SessionService) Set(ctx context.Context, userID int64, state string, payload any) error {
	var data []byte

	if payload != nil {
		data, _ = json.Marshal(payload)
	}

	return s.repo.Save(ctx, model.Session{
		UserID:  userID,
		State:   state,
		Payload: data,
	})
}

func (s *SessionService) SetEditMessageId(ctx context.Context, userID int64, messageID int64) error {
	return s.repo.SetEditMessageId(ctx, userID, messageID)
}

func ParsePayload[T any](data []byte) (T, error) {
	var t T
	err := json.Unmarshal(data, &t)
	return t, err
}
