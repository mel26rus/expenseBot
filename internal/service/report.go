package service

import (
	"context"
	"expense-bot/internal/model"
	"expense-bot/internal/repository"
	"time"
)

type ReportService struct {
	repo *repository.ReportRepo
}

func NewReportService(r *repository.ReportRepo) *ReportService {
	return &ReportService{repo: r}
}

func (r *ReportService) GetExistsTxUser(ctx context.Context, start time.Time, end time.Time) ([]model.User, error) {
	return r.repo.GetUsersHasTransactions(ctx, start, end)
}

func (s *ReportService) BuildUserReport(
	ctx context.Context,
	userID int64,
	start time.Time,
	end time.Time,
) ([]*model.AccountReport, error) {

	accounts, err := s.repo.GetUserAccounts(ctx, userID, start, end)
	if err != nil {
		return nil, err
	}

	transactions, err := s.repo.GetAccountTransactions(ctx, userID, start, end)
	if err != nil {
		return nil, err
	}

	txMap := make(map[int64][]*model.TransactionsReport)

	for _, tx := range transactions {
		txMap[tx.AccountId] = append(txMap[tx.AccountId], tx)
	}

	for _, acc := range accounts {
		acc.TransactionsReport = txMap[acc.AccountId]
	}

	return accounts, nil
}
