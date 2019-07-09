package service

import (
	"fmt"
	"net/http"
	"strconv"
	"math/rand"
	"encoding/json"
	"bytes"
	"io/ioutil"
	"os"
)

type Request struct{
	Client http.Client
}

type OrderValue struct {
	DistrictId int
	Price      []float64
}

type OrderResult struct{
	Res 	string
	Order_id int
	Entry_id []int
}

type ClickValue struct {
	EntryID int
}

type ClickResult struct {
	EntryID int
	Status  string
}

type ReceiptValue struct {
	OrderID int
	Price   []struct {
		Payment string
		Value   float64
	}
}

type ReceiptResult struct {
	Res        string
	Receipt_id []int64
}

type DeliveryValue struct {
	OrderID int
}

type DeliveryResult struct {
	Res      string
	Order_id int
}

func (r Request) AddOrder(price float64) *OrderResult{
	districts, err := strconv.Atoi(os.Getenv("DISTRICTS"))
	if err != nil {
		fmt.Println(err)
		return nil
	}

	ord := OrderValue{rand.Intn(districts)+1, []float64{price}}

	jsonBody, err := json.Marshal(ord)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	req, err := http.NewRequest( 
		os.Getenv("METHOD"), fmt.Sprintf("http://%s:%s/order", os.Getenv("SERVER"), os.Getenv("PORT")), bytes.NewBuffer(jsonBody),
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

func (r Request) Pay(orderId int, price float64) bool{
	message := map[string]interface{} {
		"orderId":orderId,
		"price": []map[string]interface{}{{"payment":"cash", "value": price}},
	}
	jsonBody, err := json.Marshal(message)
	if err != nil {
		fmt.Println(err)
		return false
	}

	req, err := http.NewRequest( 
		os.Getenv("METHOD"), fmt.Sprintf("http://%s:%s/pay", os.Getenv("SERVER"), os.Getenv("PORT")), bytes.NewBuffer(jsonBody),
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
	if(rr.Res != "success") {
		return false
	}
	return true
}

func (r Request) Click(entryId int) string {
	cv := ClickValue{entryId}
	jsonBody, err := json.Marshal(cv)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	req, err := http.NewRequest( 
		os.Getenv("METHOD"), fmt.Sprintf("http://%s:%s/click", os.Getenv("SERVER"), os.Getenv("PORT")), bytes.NewBuffer(jsonBody),
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

func (r Request) Delivered(orderId int) bool {
	dv := DeliveryValue{orderId}

	jsonBody, err := json.Marshal(dv)
	if err != nil {
		fmt.Println(err)
		return false
	}
	
	req, err := http.NewRequest( 
		os.Getenv("METHOD"), fmt.Sprintf("http://%s:%s/delivered", os.Getenv("SERVER"), os.Getenv("PORT")), bytes.NewBuffer(jsonBody),
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
	if(d.Res != "success") {
		return false
	}
	return true
}