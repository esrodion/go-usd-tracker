package models

import "time"

type CurrencyRate struct {
	Timestamp time.Time `json:"time"`
	Ask       float64   `json:"ask,string"`
	Bid       float64   `json:"bid,string"`
}
