package model

type Comments struct {
	TransactionID int64  `json:"transaction_id"`
	Comment       string `json:"comment"`
}
