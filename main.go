package main
import (
	"fmt"
	"net/http"
	"client_db/config"
	"client_db/service"
	"math/rand"
	"time"
)

func main() {
	config.Set()

	c := http.Client{}
	request := service.Request{
		Client: c,
	}

	start := time.Now()
	for i:=0; i < 10000000; i++ {
		rand.Seed(int64(i))
		price := float64(rand.Intn(20)*10+400)
		order := request.AddOrder(price)
		request.Pay(order.Order_id, price)
		for _, v := range order.Entry_id {
			for {
				click := request.Click(v)
				if click == "done" {
					break
				}
			}
		}
		request.Delivered(order.Order_id)
	}
	fmt.Println("success", time.Now().Sub(start))
 }
