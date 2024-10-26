package models

import "time"

type CurrenceyRate struct {
	Timestamp time.Time
	Ask       float64
	Bid       float64
}
