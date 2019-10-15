package service

import (
	"bytes"
	"client_db/orders"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Request struct {
	Client http.Client
	Host   string
	Port   int
	Query  string
}

type OrderValue struct {
	DistrictId int
	Price      []float64
}

type OrderResult struct {
	Res      string
	Order_id int
	Entry_id []int
}

type ClickValue struct {
	EntryID    int
	DistrictID int
}

type ClickResult struct {
	EntryID int
	Status  string
}

type ReceiptValue struct {
	OrderID    int
	DistrictID int
	Price      []struct {
		Payment string
		Value   float64
	}
}

type ReceiptResult struct {
	Res        string
	Receipt_id []int64
}

type DeliveryValue struct {
	OrderID    int
	DistrictID int
}

type DeliveryResult struct {
	Res      string
	Order_id int
}

func (r Request) AddOrder(order orders.Order) *OrderResult {
	ord := OrderValue{order.DistrictID, []float64{order.Price}}

	jsonBody, err := json.Marshal(ord)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	req, err := http.NewRequest(
		r.Query, fmt.Sprintf("http://%s:%d/order", r.Host, r.Port), bytes.NewBuffer(jsonBody),
	)
	req.Header.Set("Content-Type", "application/json")
	resp, err := r.Client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	reqJSON := []byte(body)
	o := OrderResult{}
	errType := json.Unmarshal(reqJSON, &o)
	if errType != nil {
		fmt.Println(err)
		return nil
	}
	return &o
}

func (r Request) Pay(order orders.Order) bool {
	message := map[string]interface{}{
		"orderId":    order.OrderID,
		"districtID": order.DistrictID,
		"price":      []map[string]interface{}{{"payment": "cash", "value": order.Price}},
	}
	jsonBody, err := json.Marshal(message)
	if err != nil {
		fmt.Println(err)
		return false
	}

	req, err := http.NewRequest(
		r.Query, fmt.Sprintf("http://%s:%d/pay", r.Host, r.Port), bytes.NewBuffer(jsonBody),
	)
	req.Header.Set("Content-Type", "application/json")
	resp, err := r.Client.Do(req)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return false
	}

	reqJSON := []byte(body)
	rr := ReceiptResult{}
	errType := json.Unmarshal(reqJSON, &rr)
	if errType != nil {
		fmt.Println(err)
		return false
	}
	if rr.Res != "success" {
		return false
	}
	return true
}

func (r Request) Click(entryId int, districtID int) string {
	cv := ClickValue{entryId, districtID}
	jsonBody, err := json.Marshal(cv)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	req, err := http.NewRequest(
		r.Query, fmt.Sprintf("http://%s:%d/click", r.Host, r.Port), bytes.NewBuffer(jsonBody),
	)
	req.Header.Set("Content-Type", "application/json")
	resp, err := r.Client.Do(req)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	reqJSON := []byte(body)
	c := ClickResult{}
	errType := json.Unmarshal(reqJSON, &c)
	if errType != nil {
		fmt.Println(err)
		return ""
	}
	return c.Status
}

func (r Request) Delivered(order orders.Order) bool {
	dv := DeliveryValue{order.OrderID, order.DistrictID}

	jsonBody, err := json.Marshal(dv)
	if err != nil {
		fmt.Println(err)
		return false
	}

	req, err := http.NewRequest(
		r.Query, fmt.Sprintf("http://%s:%d/delivered", r.Host, r.Port), bytes.NewBuffer(jsonBody),
	)
	req.Header.Set("Content-Type", "application/json")
	resp, err := r.Client.Do(req)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return false
	}
	reqJSON := []byte(body)
	d := DeliveryResult{}
	errType := json.Unmarshal(reqJSON, &d)
	if errType != nil {
		fmt.Println(err)
		return false
	}
	if d.Res != "success" {
		return false
	}
	return true
}