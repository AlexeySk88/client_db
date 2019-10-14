package main

import (
	"client_db/orders"
	"client_db/service"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

var request service.Request

const (
	//host      string = "dbs1.dc.pizzasoft.ru"
	host      string = "localhost"
	port      int    = 9000
	query     string = "POST"
	districts int    = 5
	count     int    = 1000
	chunk int = 1024
)

func main() {
	c := http.Client{}
	request = service.Request{
		Client: c,
		Host:   host,
		Port:   port,
		Query:  query,
	}

	newOrderChan := make(chan orders.Order, chunk)
	orderChan := make(chan orders.Order, chunk)
	payChan := make(chan orders.Order, chunk)
	clickChan := make(chan orders.Order, chunk)
	deliveryChan := make(chan bool, chunk)

	for i := 0; i < count; i++ {
		go doOrder(request, orderChan, newOrderChan)
		go doPay(request, payChan, orderChan)
		go doClick(request, clickChan, payChan)
		go doDelivery(request, deliveryChan, clickChan)
	}

	start := time.Now()
	for i := 0; i < count; i++ {
		rand.Seed(time.Now().UnixNano())
		order := orders.Order{
			DistrictID: rand.Intn(districts) + 1,
			Price:      float64(rand.Intn(20)*10 + 400),
		}
		newOrderChan <- order
	}

	for i := 0; i < count; i++ {
		fmt.Println(i)
		<-orderChan
	}

	fmt.Println("the number of orders -", count)
	fmt.Println("success", time.Now().Sub(start))
}

func doOrder(request service.Request, ch chan orders.Order, order chan orders.Order) {
	for {
		curOrder := <-order
		res := request.AddOrder(curOrder)
		ch <- orders.Order{
			OrderID: res.Order_id,
			DistrictID: curOrder.DistrictID,
			Price: curOrder.Price,
			EntryIDs: res.Entry_id,
		}
	}
}

func doPay(request service.Request, ch chan orders.Order, order chan orders.Order) {
	o := <- order
	request.Pay(o)
	ch <- o
}

func doClick(request service.Request, ch chan orders.Order, order chan orders.Order) {
	o := <- order
	for _, v := range o.EntryIDs {
		for {
			click := request.Click(v, o.DistrictID)
			if click == "done" {
				break
			}
		}
	}
	ch <- o
}

func doDelivery(request service.Request, ch chan bool, order chan orders.Order) {
	o := <- order
	request.Delivered(o)
	ch <- true
}