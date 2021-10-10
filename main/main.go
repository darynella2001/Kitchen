package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"crypto/rand"
	"net/http"
	"os"
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
	MaxWait    int    `json:"maxWait"`
}

type FoodDetails struct{
	FoodId int `json:"food_id"`
	CookId int `json:"cook_id"`
}
func genRandNum(min, max int64) int64 {
	// calculate the max we will be using
	bg := big.NewInt(max - min)

	n, err := rand.Int(rand.Reader, bg)
	if err != nil {
		panic(err)
	}

	// add n to min to support the passed in range
	return n.Int64() + min
}
func processOrder(order Order){
	// some code

	// sleep for 3-10 seconds
	preparation_time := int(genRandNum(1,7)) + 3
	time.Sleep(time.Duration(preparation_time) * time.Second)

	// finished
	cookOrder(order, preparation_time)
}

func cookOrder(order Order, prepatationTime int){
	//some code
}

func servePage(rw http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var order Order
	err := decoder.Decode(&order)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Print(time.Now().Clock())
	fmt.Printf(": Cooking order number: %d\n", order.Id)
	go processOrder(order)
}


func main() {
	jsonFile, err := os.Open("foods.json")
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println("Successfully opened food.json")
	defer jsonFile.Close()

	//read our opened jsonFile as a byte array
	byteValue, _ := ioutil.ReadAll(jsonFile)

	//we initialize our Foods array
	var foods Foods

	//we unmarshal our byteArray which contains our
	//jsonFile's content info 'foods' which we defined above

	json.Unmarshal(byteValue, &foods)


	http.HandleFunc("/kitchen", servePage)
	log.Fatal(http.ListenAndServe(":8081", nil))

}
