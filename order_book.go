package main

import (
	"fmt"
	"sort"
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

type Orders []*Order

func (o Orders) Len() int           { return len(o) }
func (o Orders) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }
func (o Orders) Less(i, j int) bool { return o[i].Timestamp < o[j].Timestamp }

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
	Orders      Orders
	TotalVolume float64
}

type Limits []*Limit
type ByBestAsk struct{ Limits }

func (b ByBestAsk) Len() int           { return len(b.Limits) }
func (b ByBestAsk) Swap(i, j int)      { b.Limits[i], b.Limits[i] = b.Limits[j], b.Limits[i] }
func (b ByBestAsk) Less(i, j int) bool { return b.Limits[i].Price < b.Limits[j].Price }

type ByBestBid struct{ Limits }

func (b ByBestBid) Len() int           { return len(b.Limits) }
func (b ByBestBid) Swap(i, j int)      { b.Limits[i], b.Limits[i] = b.Limits[j], b.Limits[i] }
func (b ByBestBid) Less(i, j int) bool { return b.Limits[i].Price > b.Limits[j].Price }

func NewLimit(price float64) *Limit {
	return &Limit{
		Price:       price,
		Orders:      []*Order{},
		TotalVolume: 0,
	}
}

func (l *Limit) String() string {
	return fmt.Sprintf("[price:%.2f |volume:%.2f]", l.Price, l.TotalVolume)
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

	sort.Sort(l.Orders)
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
	var limit *Limit

	if order.Bid {
		limit = ob.BidLimits[price]
	} else {
		limit = ob.AskLimits[price]
	}
	if limit == nil {
		limit = NewLimit(price)
		limit.AddOrder(order)
		if order.Bid {
			ob.Bids = append(ob.Bids, limit)
			ob.BidLimits[price] = limit
		} else {
			ob.Asks = append(ob.Asks, limit)
			ob.AskLimits[price] = limit
		}
	} else {
		limit.AddOrder(order)
	}
}
