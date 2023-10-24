package model

import "time"

type Transaction struct {
	Date     time.Time
	Amount   float64
	IsCredit bool
}
