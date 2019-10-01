package main

import (
	"client_db/orders"
	"client_db/service"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
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
	//tr := &http.Transport{
	//	MaxIdleConns:       20,
	//	MaxIdleConnsPerHost:  20,
	//}
	//c := http.Client{Transport: tr}
	c := http.Client{}
	request = service.Request{
		Client: c,
		Host:   host,
		Port:   port,
		Query:  query,
	}

	var wg sync.WaitGroup
	wg.Add(chunk)

	start := time.Now()
	orderChan := make(chan orders.Order, chunk)
	payChan := make(chan orders.Order, chunk)
	clickChan := make(chan orders.Order, chunk)
	deliveryChan := make(chan bool, chunk)

	for i := 0; i < count; i++ {
		rand.Seed(time.Now().UnixNano())
		//time.Sleep(5 * time.Millisecond)
		order := orders.Order{
			DistrictID: rand.Intn(districts) + 1,
			Price:      float64(rand.Intn(20)*10 + 400),
		}
		go doOrder(request, orderChan, order)
		go doPay(request, payChan, <-orderChan)
		go doClick(request, clickChan, <-payChan)
		go doDelivery(request, deliveryChan, <-clickChan)
		fmt.Println(i)
		<-deliveryChan
	}

	fmt.Println("the number of orders -", count)
	fmt.Println("success", time.Now().Sub(start))
}

func doOrder(request service.Request, ch chan orders.Order, order orders.Order) {
	res := request.AddOrder(order)
	ch <- orders.Order{
		OrderID: res.Order_id,
		DistrictID: order.DistrictID,
		Price: order.Price,
		EntryIDs: res.Entry_id,
	}
}

func doPay(request service.Request, ch chan orders.Order, order orders.Order) {
	request.Pay(order)
	ch <- order
}

func doClick(request service.Request, ch chan orders.Order, order orders.Order) {
	for _, v := range order.EntryIDs {
		for {
			click := request.Click(v, order.DistrictID)
			if click == "done" {
				break
			}
		}
	}
	ch <- order
}

func doDelivery(request service.Request, ch chan bool, order orders.Order) {
	request.Delivered(order)
	ch <- true
}