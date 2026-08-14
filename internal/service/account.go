package service

import (
	"context"
	"expense-bot/internal/model"
	"expense-bot/internal/repository"
	"log/slog"
)

type AccountService struct {
	repo *repository.AccountRepo
}

func NewAccountService(r *repository.AccountRepo) *AccountService {
	return &AccountService{repo: r}
}

func (s *AccountService) GetAccountsByUserID(ctx context.Context, userID int64) ([]*model.Account, error) {
	return s.repo.GetAccountsByUserID(ctx, userID)
}

func (s *AccountService) CreateAccount(ctx context.Context, userID int64, name string, currencyID int64) (int64, error) {
	return s.repo.CreateAccount(ctx, userID, name, currencyID)
}

func (s *AccountService) GetAccountTitle(ctx context.Context, accountID int64) string {
	account, err := s.repo.GetAccountTitleByID(ctx, accountID)
	if err != nil {
		slog.Error("Error getting account by ID", "accountID", accountID, "error", err)
		return "Unknown Account"
	}
	return account.Name
}

func (s *AccountService) GetAccountBalance(ctx context.Context, accountID int64) float64 {
	balance, err := s.repo.GetAccountBalanceByID(ctx, accountID)
	if err != nil {
		slog.Error("Error getting account balance", "accountID", accountID, "error", err)
		return 0
	}
	return balance
}

func (s *AccountService) GetAccountCurrencyByID(ctx context.Context, accountID int64) int64 {
	currencyID, err := s.repo.GetAccountCurrencyByID(ctx, accountID)
	if err != nil {
		slog.Error("Error getting account currency id", "accountID", accountID, "error", err)
		return 0
	}
	return currencyID
}

func (s *AccountService) GetMenuAccountsByUserID(ctx context.Context, userID int64) ([]*model.Account, error) {
	return s.repo.GetMenuAccountsByUserID(ctx, userID)
}

func (s *AccountService) ChangeAccountIshidden(ctx context.Context, accountId int64) error {
	return s.repo.ChangeAccountIshidden(ctx, accountId)
}
