package models

import "time"

type CurrenceyRate struct {
	Timestamp time.Time `json:"time"`
	Ask       float64   `json:"ask,string"`
	Bid       float64   `json:"bid,string"`
}
