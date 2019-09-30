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
	count     int    = 1000000
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

	var wg sync.WaitGroup
	wg.Add(chunk)

	start := time.Now()
	ch := make(chan bool)
	for i := 0; i < count; i++ {
		rand.Seed(time.Now().UnixNano())
		order := orders.Order{
			DistrictID: rand.Intn(districts) + 1,
			Price:      float64(rand.Intn(20)*10 + 400),
		}
		o := doOrder(request, wg, order)
	}
	for i := 0; i < count; i++ {
		<-ch
		fmt.Println(i)
	}
	fmt.Println("the number of orders -", count)
	fmt.Println("success", time.Now().Sub(start))
}

func doOrder(request service.Request, wg sync.WaitGroup, order orders.Order) request.OrderResult {
	defer wg.Done()
	res := request.AddOrder(order.Price, order.DistrictID)
	return res
}

/*func doRequest(request service.Request, districts int, ch chan bool) {
	order := request.AddOrder(price, districtID)
	request.Pay(order.Order_id, districtID, price)
	for _, v := range order.Entry_id {
		for {
			click := request.Click(v, districtID)
			if click == "done" {
				break
			}
		}
	}
	request.Delivered(order.Order_id, districtID)
	ch <- true
}*/
