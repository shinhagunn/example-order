package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type Trade struct {
	CreatedAt time.Time
	Price     decimal.Decimal
	Amount    int
}
