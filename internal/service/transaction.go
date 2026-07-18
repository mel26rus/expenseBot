package service

import (
	"context"
	"expense-bot/internal/model"
	"expense-bot/internal/repository"
)

const showLastCommentsLimit = 5

type TransactionService struct {
	repo *repository.TransactionRepo
}

func NewTransactionService(r *repository.TransactionRepo) *TransactionService {
	return &TransactionService{repo: r}
}

func (s *TransactionService) CreateTx(
	ctx context.Context,
	userID int64,
	accountID int64,
	txType string,
	amount float64,
	comment string,
	txGUID *string,
) error {
	return s.repo.CreateTx(ctx, userID, accountID, txType, amount, comment, txGUID)
}

func (s *TransactionService) GetLastComments(ctx context.Context, userID int64, txType string) ([]model.Comments, error) {
	return s.repo.GetLastComments(ctx, userID, txType, showLastCommentsLimit)
}

func (s *TransactionService) GetCommentByTxID(ctx context.Context, txID int64) (string, error) {
	return s.repo.GetCommentByTxID(ctx, txID)
}
