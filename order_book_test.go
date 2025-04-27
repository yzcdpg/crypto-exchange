package main

import (
	"fmt"
	"testing"
)

func TestLimit(t *testing.T) {
	l := NewLimit(10_000)
	buyOrderA := NewOrder(true, 5)
	buyOrderB := NewOrder(true, 8)
	buyOrderC := NewOrder(true, 10)
	l.AddOrder(buyOrderA)
	l.AddOrder(buyOrderB)
	l.AddOrder(buyOrderC)
	fmt.Println(l)
	l.DeleteOrder(buyOrderB)
	fmt.Println(l)
}

func TestNewOrderBook(t *testing.T) {
	ob := NewOrderBook()
	buyOrderA := NewOrder(true, 10)
	buyOrderB := NewOrder(true, 2000)
	buyOrderC := NewOrder(true, 5)
	ob.PlaceOrder(18_000, buyOrderA)
	ob.PlaceOrder(19_000, buyOrderB)
	ob.PlaceOrder(18_000, buyOrderC)
	fmt.Println(ob.Bids)
}
