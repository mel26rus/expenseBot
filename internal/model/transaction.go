package model

type Transaction struct {
	ID        int64
	UserID    int64
	AccountID int64
	Amount    float64
	Comment   string
}
