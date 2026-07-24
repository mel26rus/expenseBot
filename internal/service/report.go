package service

import (
	"context"
	"expense-bot/internal/dateutil"
	"expense-bot/internal/model"
	"expense-bot/internal/repository"
)

type ReportService struct {
	repo *repository.ReportRepo
}

func NewReportService(r *repository.ReportRepo) *ReportService {
	return &ReportService{repo: r}
}

func (r *ReportService) GetUsersForDailyReport(ctx context.Context) ([]int64, error) {
	start, end := dateutil.Today()
	users, err := r.repo.GetUsersHasTransactionsDaily(ctx, start, end)
	return users, err
}

func (s *ReportService) BuildDailyReport(
	ctx context.Context,
	tgID int64,
) ([]*model.AccountReport, error) {

	start, end := dateutil.Today()

	accounts, err := s.repo.GetUserAccounts(ctx, tgID, start, end)
	if err != nil {
		return nil, err
	}

	transactions, err := s.repo.GetAccountTransactions(ctx, tgID, start, end)
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
