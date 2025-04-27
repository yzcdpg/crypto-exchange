package main

import (
	"fmt"
	"time"
)

type Match struct {
	Ask        *Order
	Bid        *Order
	SizeFilled float64
	Price      float64
}

type Order struct {
	Size      float64
	Bid       bool
	Limit     *Limit
	Timestamp int64
}

func (o *Order) String() string {
	return fmt.Sprintf("[size:%.2f]", o.Size)
}

func NewOrder(bid bool, size float64) *Order {
	return &Order{
		Size:      size,
		Bid:       bid,
		Timestamp: time.Now().UnixNano(),
	}
}

// Limit 限价单
type Limit struct {
	Price       float64
	Orders      []*Order
	TotalVolume float64
}

func NewLimit(price float64) *Limit {
	return &Limit{
		Price:       price,
		Orders:      []*Order{},
		TotalVolume: 0,
	}
}

func (l *Limit) AddOrder(order *Order) {
	order.Limit = l
	l.Orders = append(l.Orders, order)
	l.TotalVolume += order.Size
}

func (l *Limit) DeleteOrder(order *Order) {
	for i := 0; i < len(l.Orders); i++ {
		if l.Orders[i] == order {
			// 用最后一个覆盖
			l.Orders[i] = l.Orders[len(l.Orders)-1]
			l.Orders = l.Orders[:len(l.Orders)-1]
		}
	}

	order.Limit = nil
	l.TotalVolume -= order.Size

	// TODO: resort the whole resting orders
}

// OrderBook 订单簿
type OrderBook struct {
	Asks      []*Limit // 卖
	Bids      []*Limit // 买
	AskLimits map[float64]*Limit
	BidLimits map[float64]*Limit
}

func NewOrderBook() *OrderBook {
	return &OrderBook{
		Asks:      []*Limit{},
		Bids:      []*Limit{},
		AskLimits: make(map[float64]*Limit),
		BidLimits: make(map[float64]*Limit),
	}
}

// PlaceOrder 下单
func (ob *OrderBook) PlaceOrder(price float64, order *Order) []Match {
	// 1. try to match the order
	// 2. add the rest of the order to the books
	if order.Size > 0.0 {
		ob.add(price, order)
	}

	return []Match{}
}

func (ob *OrderBook) add(price float64, order *Order) {



}
