package orderbook

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"math/rand"
	"sort"
	"sync"
	"time"
)

type Trade struct {
	Price     float64
	Size      float64
	Bid       bool
	Timestamp int64
}

type Match struct {
	Ask        *Order
	Bid        *Order
	SizeFilled float64
	Price      float64
}

type Order struct {
	ID        int64
	UserID    int64
	Size      float64
	Bid       bool
	Limit     *Limit
	Timestamp int64
}

type Orders []*Order

func (o Orders) Len() int           { return len(o) }
func (o Orders) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }
func (o Orders) Less(i, j int) bool { return o[i].Timestamp < o[j].Timestamp }

func NewOrder(bid bool, size float64, userID int64) *Order {
	return &Order{
		UserID:    userID,
		ID:        int64(rand.Intn(10000000)),
		Size:      size,
		Bid:       bid,
		Timestamp: time.Now().UnixNano(),
	}
}

func (o *Order) String() string { return fmt.Sprintf("[size:%.2f] | [id: %d]", o.Size, o.ID) }

func (o *Order) Type() string {
	if o.Bid {
		return "BID"
	}
	return "ASK"
}

func (o *Order) IsFilled() bool { return o.Size == 0.0 }

// Limit 限价单
type Limit struct {
	Price       float64
	Orders      Orders
	TotalVolume float64
}

type Limits []*Limit
type ByBestAsk struct{ Limits }

func (b ByBestAsk) Len() int           { return len(b.Limits) }
func (b ByBestAsk) Swap(i, j int)      { b.Limits[i], b.Limits[j] = b.Limits[j], b.Limits[i] }
func (b ByBestAsk) Less(i, j int) bool { return b.Limits[i].Price < b.Limits[j].Price }

type ByBestBid struct{ Limits }

func (b ByBestBid) Len() int           { return len(b.Limits) }
func (b ByBestBid) Swap(i, j int)      { b.Limits[i], b.Limits[j] = b.Limits[j], b.Limits[i] }
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

// Fill 成交订单
func (l *Limit) Fill(order *Order) []Match {
	var (
		matches        []Match
		ordersToDelete []*Order
	)

	for _, o := range l.Orders {
		if order.IsFilled() {
			break
		}
		match := l.fillOrder(o, order)
		matches = append(matches, match)
		l.TotalVolume -= match.SizeFilled
		if o.IsFilled() {
			ordersToDelete = append(ordersToDelete, o)
		}
	}

	for _, o := range ordersToDelete {
		l.DeleteOrder(o)
	}

	return matches
}

// 成交订单
func (l *Limit) fillOrder(a, b *Order) Match {
	var (
		bid        *Order
		ask        *Order
		sizeFilled float64
	)
	if a.Bid {
		bid = a
		ask = b
	} else {
		ask = a
		bid = b
	}

	if a.Size >= b.Size {
		a.Size -= b.Size
		sizeFilled = b.Size
		b.Size = 0.0
	} else {
		b.Size -= a.Size
		sizeFilled = a.Size
		a.Size = 0.0
	}

	return Match{
		Ask:        ask,
		Bid:        bid,
		SizeFilled: sizeFilled,
		Price:      l.Price,
	}
}

// OrderBook 订单簿
type OrderBook struct {
	asks      []*Limit // 卖
	bids      []*Limit // 买
	Trades    []*Trade
	mu        sync.RWMutex
	AskLimits map[float64]*Limit
	BidLimits map[float64]*Limit
	Orders    map[int64]*Order
}

func NewOrderBook() *OrderBook {
	return &OrderBook{
		asks:      []*Limit{},
		bids:      []*Limit{},
		Trades:    []*Trade{},
		AskLimits: make(map[float64]*Limit),
		BidLimits: make(map[float64]*Limit),
		Orders:    make(map[int64]*Order),
	}
}

// PlaceMarketOrder 下市价单
func (ob *OrderBook) PlaceMarketOrder(order *Order) []Match {
	ob.mu.Lock()
	defer ob.mu.Unlock()

	var matches []Match
	if order.Bid { // bid
		if order.Size > ob.AskTotalVolume() {
			panic(fmt.Sprintf("not enough volume [size: %.2f] for market orders [size: %.2f]", ob.AskTotalVolume(), order.Size))
		}
		for _, limit := range ob.Asks() {
			limitMatches := limit.Fill(order)
			matches = append(matches, limitMatches...)

			if len(limit.Orders) == 0 {
				ob.clearLimit(false, limit)
			}
		}
	} else { // ask
		if order.Size > ob.BidTotalVolume() {
			panic("not enough volume sitting in the books")
		}
		for _, limit := range ob.Bids() {
			limitMatches := limit.Fill(order)
			matches = append(matches, limitMatches...)

			if len(limit.Orders) == 0 {
				ob.clearLimit(true, limit)
			}
		}
	}

	for _, match := range matches {
		trade := &Trade{
			Price:     match.Price,
			Size:      match.SizeFilled,
			Timestamp: time.Now().UnixNano(),
			Bid:       order.Bid,
		}
		ob.Trades = append(ob.Trades, trade)
	}

	logrus.WithFields(logrus.Fields{
		"currentPrice": ob.Trades[len(ob.Trades)-1].Price,
	}).Info()

	return matches
}

// PlaceLimitOrder 下限价单
func (ob *OrderBook) PlaceLimitOrder(price float64, order *Order) {
	ob.mu.Lock()
	defer ob.mu.Unlock()

	var limit *Limit

	if order.Bid {
		limit = ob.BidLimits[price]
	} else {
		limit = ob.AskLimits[price]
	}
	if limit == nil {
		limit = NewLimit(price)
		if order.Bid {
			ob.bids = append(ob.bids, limit)
			ob.BidLimits[price] = limit
		} else {
			ob.asks = append(ob.asks, limit)
			ob.AskLimits[price] = limit
		}
	}

	logrus.WithFields(logrus.Fields{
		"price":  limit.Price,
		"type":   order.Type(),
		"size":   order.Size,
		"userID": order.UserID,
	}).Info("new limit order")

	ob.Orders[order.ID] = order
	limit.AddOrder(order)
}

func (ob *OrderBook) clearLimit(bid bool, l *Limit) {
	if bid {
		delete(ob.BidLimits, l.Price)
		for i := 0; i < len(ob.bids); i++ {
			if ob.bids[i] == l {
				ob.bids[i] = ob.bids[len(ob.bids)-1]
				ob.bids = ob.bids[:len(ob.bids)-1]
			}
		}
	} else {
		delete(ob.AskLimits, l.Price)
		for i := 0; i < len(ob.asks); i++ {
			if ob.asks[i] == l {
				ob.asks[i] = ob.asks[len(ob.asks)-1]
				ob.asks = ob.asks[:len(ob.asks)-1]
			}
		}
	}

	fmt.Printf("clearing limit price level [%.2f]\n", l.Price)
}

func (ob *OrderBook) CancelOrder(order *Order) {
	limit := order.Limit
	limit.DeleteOrder(order)
	delete(ob.Orders, order.ID)
	if len(limit.Orders) == 0 {
		ob.clearLimit(order.Bid, limit)
	}
}

// BidTotalVolume 买单(深度)挂单量
func (ob *OrderBook) BidTotalVolume() float64 {
	totalVolume := 0.0
	for _, bid := range ob.bids {
		totalVolume += bid.TotalVolume
	}
	return totalVolume
}

// AskTotalVolume 卖单(深度)挂单量
func (ob *OrderBook) AskTotalVolume() float64 {
	totalVolume := 0.0
	for _, ask := range ob.asks {
		totalVolume += ask.TotalVolume
	}
	return totalVolume
}

func (ob *OrderBook) Asks() []*Limit {
	sort.Sort(ByBestAsk{ob.asks})
	return ob.asks
}

func (ob *OrderBook) Bids() []*Limit {
	sort.Sort(ByBestBid{ob.bids})
	return ob.bids
}
