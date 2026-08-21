package model

import "time"

type UserReport struct {
	UserID        int64
	TgUserID      int64
	AccountReport []*AccountReport
}

type AccountReport struct {
	UserID             int64
	AccountId          int64
	Title              string
	Balance            float64
	Income             float64
	Expense            float64
	TransactionsReport []*TransactionsReport
	CurrencyName       string
	ExRate             float64
	ExDate             time.Time
	USDBalance         float64
	RUBBalance         float64
}

type TransactionsReport struct {
	AccountId    int64
	Expense_type string
	Category     string
	Amount       float64
}
