package model

import "time"

type ExchangeRate struct {
	Date  time.Time
	Name  string
	Value float64
}
