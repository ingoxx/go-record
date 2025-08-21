package main

import "fmt"

type OrderState interface {
	Handle(order *Order)
	String() string
}

type PendingState struct{}

func (s *PendingState) Handle(order *Order) {
	order.Status = "Paid"
	order.State = &PaidState{}
}

func (s *PendingState) String() string {
	return "Pending"
}

type PaidState struct{}

func (s *PaidState) Handle(order *Order) {
	order.Status = "Shipped"
	order.State = &ShippedState{}
}

func (s *PaidState) String() string {
	return "Paid"
}

type ShippedState struct{}

func (s *ShippedState) Handle(order *Order) {
	order.Status = "Delivered"
	order.State = &DeliveredState{}
}

func (s *ShippedState) String() string {
	return "Shipped"
}

type DeliveredState struct{}

func (s *DeliveredState) Handle(order *Order) {
	// Do nothing
	fmt.Println("订单已完成")
}

func (s *DeliveredState) String() string {
	return "Delivered"
}

type Order struct {
	Status string
	State  OrderState
}

func NewOrder() *Order {
	return &Order{
		Status: "Pending",
		State:  &PendingState{},
	}
}

func main() {
	order := NewOrder()
	order.State.Handle(order)
	fmt.Println(order.State.String())
}
