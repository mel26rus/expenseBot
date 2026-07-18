package service

import (
	"context"

	"expense-bot/internal/model"
	"expense-bot/internal/repository"
)

type UserService struct {
	repo *repository.UserRepo
}

func NewUserService(r *repository.UserRepo) *UserService {
	return &UserService{repo: r}
}

func (s *UserService) GetOrCreate(ctx context.Context, telegramID int64) (model.User, error) {
	return s.repo.GetOrCreate(ctx, telegramID)
}

func (s *UserService) GetUserSettings(ctx context.Context, id int64) (model.UserSettings, error) {
	return s.repo.GetUserSettings(ctx, id)
}

func (s *UserService) ChangeDailyReportConfig(ctx context.Context, id int64) (model.UserSettings, error) {
	return s.repo.ChangeDailyReportConfig(ctx, id)
}

func (s *UserService) ChangeMonthlyReportConfig(ctx context.Context, id int64) (model.UserSettings, error) {
	return s.repo.ChangeMonthlyReportConfig(ctx, id)
}
