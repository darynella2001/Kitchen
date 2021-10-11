package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)
//Foods struct which contains an array of foods
type Foods struct {
	Foods []Food `json:"foods"`
}
//Food struct which contains details about the dish
type Food struct {
	Id               int              `json:"id"`
	Name             string           `json:"name"`
	PreparationTime  int              `json:"preparation-time"`
	Complexity       int              `json:"complexity"`
	CookingApparatus string           `json:"cooking-apparatus"`
}
//Order struct which contains details about the generated order
type Order struct {
	Id         int    `json:"id"`
	Items      []int  `json:"items"`
	Priority   int    `json:"priority"`
	MaxWait    int    `json:"max-wait"`
	PickUpTime int    `json:"pick-up-time"`
}

type PreparedOrder struct {
	Id          int    `json:"id"`
	Items       []int  `json:"items"`
	Priority    int    `json:"priority"`
	MaxWait     int    `json:"max-wait"`
	PickUpTime  int    `json:"pick-up-time"`
	CookingTime int    `json:"cooking-time"`
}

type FoodDetails struct{
	FoodId int `json:"food_id"`
	CookId int `json:"cook_id"`
}
func getUnixTimestamp() int64 {
	now := time.Now()
	sec := now.Unix()
	return sec
}

func getJsonRequest(order Order) []byte {
	preparedOrder := &PreparedOrder{
		Id:          order.Id,
		Items:       order.Items,
		Priority:    order.Priority,
		MaxWait:     order.MaxWait,
		PickUpTime:  order.PickUpTime,
		CookingTime: (int(getUnixTimestamp())- order.PickUpTime)}
	ord, err := json.Marshal(preparedOrder)
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
	return ord
}

func makeRequest(ord []byte) {
	req, err := http.NewRequest("POST", "http://localhost:8080/dininghall", bytes.NewBuffer(ord))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Println(string(ord))
	fmt.Println("Request sent")

}

func waiter(order Order){
	request := getJsonRequest(order)
	time.Sleep(time.Second)
	makeRequest(request)
}


func servePage(rw http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var order Order
	err := decoder.Decode(&order)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	log.Println(order)
	fmt.Println("Request Handled")

	var wg sync.WaitGroup
	for i:=1; i <= 100; i++{
		wg.Add(1)

		go func() {
			defer wg.Done()
			waiter(order)
		}()
	}
	wg.Wait()
}


func main() {
	http.HandleFunc("/kitchen", servePage)
	log.Fatal(http.ListenAndServe(":8081", nil))

}
