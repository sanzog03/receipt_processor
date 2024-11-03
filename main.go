package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"receiptProcessor/database"

	"github.com/gorilla/mux"
)

var receiptStore *database.ReceiptStore

func main() {
	r := mux.NewRouter()
	receiptStore = database.NewReceiptStore()

	r.HandleFunc("/receipt/process", ProcessReceipt).Methods("POST")
	r.HandleFunc("/receipt/{id}/points", ReceiptPoints).Methods("GET")

	log.Println("Server starting on port 9080...")
	log.Fatal(http.ListenAndServe(":9080", r))
}

func ReceiptPoints(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	receipt, error := receiptStore.Get(id)
	if error != nil {
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
	var receipt database.Receipt
	err := json.NewDecoder(r.Body).Decode(&receipt)
	// receipts = append(receipts, receipt)
	if err != nil {
		log.Println("error is: ", err)
		http.Error(w, "Invalid receipt data", http.StatusBadRequest)
		return
	}
	// create a id
	id := generateId()
	receiptStore.Set(id, receipt)
	fmt.Println("new receipt with id:", id, " and value ", receipt)
}

func generateId() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
