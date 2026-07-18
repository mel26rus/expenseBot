package model

type User struct {
	ID         int64
	TelegramID int64
}

type UserSettings struct {
	ID              int64
	TelegramID      int64
	IsDailyReport   bool
	IsMonthlyReport bool
}
