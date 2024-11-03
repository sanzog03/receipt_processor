package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/receipt/process", ProcessReceipt).Methods("POST")
	r.HandleFunc("/receipt/{id}/points", ReceiptPoints).Methods("GET")

	log.Println("Server starting on port 9080...")
	log.Fatal(http.ListenAndServe(":9080", r))
}

func ReceiptPoints(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	receipt, ok := receipts[id]
	if !ok {
		http.Error(w, "Receipt not found", http.StatusNotFound)
		return
	}
	// find the points for the receipt and return the point
	fmt.Println("recepit", receipt)
	points := 30
	response := map[string]int{"points": points}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func ProcessReceipt(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Println(vars)
	var receipt Receipt
	err := json.NewDecoder(r.Body).Decode(&receipt)
	// receipts = append(receipts, receipt)
	if err != nil {
		log.Println("error is: ", err)
		http.Error(w, "Invalid receipt data", http.StatusBadRequest)
		return
	}
	// create a id
	id := generateId()
	receipts[id] = receipt
	fmt.Println("new receipt with id:", id, " and value ", receipts)
}

func generateId() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

//  type definations

type Item struct {
	ShortDescription string  `json:"shortDescription"`
	Price            float64 `json:"price,string"`
}

type Receipt struct {
	Retailer     string     `json:"retailer"`
	PurchaseDate CustomTime `json:"purchaseDate"`
	PurchaseTime CustomTime `json:"purchaseTime"`
	Items        []Item     `json:"items"`
	Total        float64    `json:"total,string"`
}

// for special type date and time
type CustomTime struct {
	time.Time
}

func (ct *CustomTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		t, err = time.Parse("15:04", s)
		if err != nil {
			return err
		}
	}
	ct.Time = t
	return nil
}

func (ct CustomTime) MarshalJSON() ([]byte, error) {
	if ct.Time.Hour() == 0 && ct.Time.Minute() == 0 {
		return []byte(fmt.Sprintf("\"%s\"", ct.Time.Format("2006-01-02"))), nil
	}
	return []byte(fmt.Sprintf("\"%s\"", ct.Time.Format("15:04"))), nil
}

// mock database. A hashmap for O(1) retrival

var receipts = make(map[string]Receipt)
