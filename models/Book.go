package models

import "github.com/shopspring/decimal"

type BookType string

var (
	BookTypeAsk = BookType("ask")
	BookTypeBid = BookType("bid")
)

type Book struct {
	Ask []Order
	Bid []Order
}

func (b *Book) FindIndex(orders []Order, price decimal.Decimal) (bool, int) {
	for i, v := range orders {
		if v.Price.Equal(price) {
			return true, i
		}
	}

	return false, 0
}

func (b *Book) AddOrder(t BookType, price decimal.Decimal, amount int) {
	if t == BookTypeAsk {
		exist, index := b.FindIndex(b.Ask, price)

		if exist {
			b.Ask[index].Amount += amount
		} else {
			b.Ask = append(b.Ask, Order{price, amount})
		}
	} else {
		exist, index := b.FindIndex(b.Bid, price)

		if exist {
			b.Bid[index].Amount += amount
		} else {
			b.Bid = append(b.Bid, Order{price, amount})
		}
	}
}

func (b *Book) DeleteOrder(t BookType, index int) {
	if t == BookTypeAsk {
		b.Ask = append(b.Ask[:index], b.Ask[index+1:]...)
	} else {
		b.Bid = append(b.Bid[:index], b.Bid[index+1:]...)
	}
}
