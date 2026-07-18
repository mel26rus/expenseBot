package model

type AccountPayload struct {
	Name string `json:"name"`
}

type TxPayload struct {
	Amount      float64 `json:"amount"`
	AmountTo    float64 `json:"amountTo"`
	AccountID   int64   `json:"account_id"`
	AccountToID int64   `json:"account_to_id"`
	ExpenseType int64   `json:"type"`
	Comment     string  `json:"comment"`
}
