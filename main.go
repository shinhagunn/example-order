package main

import (
	"bai-2/models"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/shopspring/decimal"
)

func GenerateOrder(askChan chan *models.Order, bidChan chan *models.Order) {
	price := decimal.NewFromFloat(223.5)
	period := time.Second
	for {
		isPositive := rand.Intn(2)
		randSide := rand.Intn(2)
		order := &models.Order{}

		if isPositive == 1 {
			order.Price = price.Add(decimal.NewFromFloat(rand.Float64() * 1).Round(1))
		} else {
			order.Price = price.Sub(decimal.NewFromFloat(rand.Float64() * 1).Round(1))
		}
		order.Amount = rand.Intn(10) + 1

		if randSide == 1 {
			askChan <- order
		} else {
			bidChan <- order
		}

		time.Sleep(period)
	}
}

func HandleTrade(book *models.Book, orders []models.Order, bookType models.BookType, order *models.Order, trade chan *models.Trade) {
	isFinded := false
	opsType := models.BookTypeAsk
	if bookType == models.BookTypeAsk {
		opsType = models.BookTypeBid
	} else {
		opsType = models.BookTypeAsk
	}

	exist, index := book.FindIndex(orders, order.Price)
	if exist {
		isFinded = true
		if orders[index].Amount < order.Amount {
			trade <- &models.Trade{
				CreatedAt: time.Now().UTC(),
				Price:     orders[index].Price,
				Amount:    orders[index].Amount,
			}
			book.DeleteOrder(opsType, index)
			book.AddOrder(bookType, order.Price, order.Amount-orders[index].Amount)
		} else if orders[index].Amount == order.Amount {
			book.DeleteOrder(opsType, index)
			trade <- &models.Trade{
				CreatedAt: time.Now().UTC(),
				Price:     orders[index].Price,
				Amount:    orders[index].Amount,
			}
		} else {
			orders[index].Amount -= order.Amount
			trade <- &models.Trade{
				CreatedAt: time.Now().UTC(),
				Price:     orders[index].Price,
				Amount:    order.Amount,
			}
		}
	}

	if !isFinded {
		book.AddOrder(bookType, order.Price, order.Amount)
	}
}

func ReadBook(orderBook *models.Book) {
	for {
		log.Println(orderBook.Ask)
		log.Println(orderBook.Bid)

		time.Sleep(time.Second)
	}
}

func ReadHistoryTrade(trade chan *models.Trade) {
	log.Println("HISTORY TRADE")
	for {
		select {
		case t := <-trade:
			fmt.Printf("Created At: %v Price: %v Amount %v\n", t.CreatedAt, t.Price, t.Amount)
		}
	}
}

func main() {
	orderBook := new(models.Book)
	trade := make(chan *models.Trade)
	askChan := make(chan *models.Order)
	bidChan := make(chan *models.Order)

	go GenerateOrder(askChan, bidChan)
	go ReadHistoryTrade(trade)
	// go ReadBook(orderBook)

	for {
		select {
		case orderAsk := <-askChan:
			HandleTrade(orderBook, orderBook.Bid, models.BookTypeAsk, orderAsk, trade)
		case orderBid := <-bidChan:
			HandleTrade(orderBook, orderBook.Ask, models.BookTypeBid, orderBid, trade)
		}
	}
}
