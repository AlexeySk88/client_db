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
	chunk = 100
	count = 10000
	districts = 5
	host      string = "localhost"
	port      int    = 9000
	query     string = "POST"
)

func main() {
	c := http.Client{}
	request = service.Request{
		Client: c,
		Host:   host,
		Port:   port,
		Query:  query,
	}
	wg := sync.WaitGroup{}
	pipe := make(chan orders.Order, count)

	start := time.Now()

	for i := 0; i < chunk; i++ {

		go doRequest(request, pipe, &wg)
	}

	for i := 0; i < count; i++ {
		wg.Add(1)
		rand.Seed(time.Now().UnixNano())
		order := orders.Order{
			DistrictID: rand.Intn(districts) + 1,
			Price:      float64(rand.Intn(20)*10 + 400),
		}
		pipe <- order
	}
	wg.Wait()


	fmt.Println("the number of orders -", count)
	fmt.Println("success", time.Now().Sub(start))
}

func doRequest(request service.Request, ch chan orders.Order, wg *sync.WaitGroup) {
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
