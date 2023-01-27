package models

import (
	"github.com/shopspring/decimal"
)

type Order struct {
	Price  decimal.Decimal
	Amount int
}
