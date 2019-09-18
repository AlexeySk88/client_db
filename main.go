package main

import (
	"client_db/config"
	"client_db/service"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

var request service.Request
var districts int

func main() {
	config.Set()

	c := http.Client{}
	request = service.Request{
		Client: c,
	}
	count := 1000000

	dis, err := strconv.Atoi(os.Getenv("DISTRICTS"))
	districts = dis
	if err != nil {
		fmt.Println(err)
	}
	start := time.Now()
	ch := make(chan bool)
	for i := 0; i < count; i++ {
		rand.Seed(int64(i))
		time.Sleep(5 * time.Millisecond)
		go doRequest(request, districts, ch)
	}
	for i := 0; i < count; i++ {
		<-ch
		fmt.Println(i)
	}
	fmt.Println("the number of orders -", count)
	fmt.Println("success", time.Now().Sub(start))
}

func doRequest(request service.Request, districts int, ch chan bool) {
	districtID := rand.Intn(districts) + 1
	price := float64(rand.Intn(20)*10 + 400)
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
}
