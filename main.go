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

const(
	chunk = 500
	count = 10000
	districts = 5
	//host      string = "dbs1.dc.pizzasoft.ru"
	host      string = "localhost"
	port      int    = 9000
	query     string = "POST"
)

func main() {
	wg := sync.WaitGroup{}
	pipe := make(chan orders.Order, count)

	start := time.Now()
	for i := 0; i < count; i++ {
		rand.Seed(time.Now().UnixNano())
		order := orders.Order{
			DistrictID: rand.Intn(districts) + 1,
			Price:      float64(rand.Intn(20)*10 + 400),
		}
		pipe <- order
	}

	for i := 0; i < count/chunk; i++ {
		wg.Add(chunk)
		go doRequest(pipe, &wg)
		wg.Wait()
	}

	fmt.Println("the number of orders -", count)
	fmt.Println("success", time.Now().Sub(start))
}

func doRequest(ch chan orders.Order, wg *sync.WaitGroup) {
	c := http.Client{}
	request = service.Request{
		Client: c,
		Host:   host,
		Port:   port,
		Query:  query,
	}

	for curOrder := range ch {
		order := request.AddOrder(curOrder)
		curOrder.OrderID = order.Order_id
		curOrder.EntryIDs = order.Entry_id
		fmt.Println(curOrder.OrderID)

		request.Pay(curOrder)
		for _, v := range order.Entry_id {
			for {
				click := request.Click(v, curOrder.DistrictID)
				if click == "done" {
					break
				}
			}
		}
		request.Delivered(curOrder)
		wg.Done()
	}
}
