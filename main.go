package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strings"
	"time"
	"unicode"

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
	reward := calculateReward(receipt)
	receiptStore.Set(id, receipt)
	fmt.Println("new receipt with id:", id, " and value ", receipt)
	fmt.Println("You have got: ", reward, " reward points.")
}

func generateId() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func calculateReward(receipt database.Receipt) int {

	rules := []Rule{
		alphanumericRule,
		roundDollarRule,
		multipleOfQuaterRule,
		pairRule,
		multipleOfThreeRule,
		oddPurchaseDateRule,
		offTimePurchaseRule,
	}
	// pass the receipt through various rules
	// for each passes, collect the rewarad
	// find the accumulated reward.
	// return that reward
	total_reward := 0
	for idx, rule := range rules {
		reward := rule(receipt)
		if reward > 0 {
			fmt.Println("reward from: rule_", idx, " is :", reward)
		}
		total_reward += reward
	}

	return total_reward

	// To Note:
	// we have multiple rules right now
	// new rules can be added in the future
	// so we need a mechanism:
	// we pass the receipt to all the rules, not manually but automatically
	// collect the amount automatically

	// For the above given purpose, the rules is a list. we can add the rules to the list.
	// we will go through each elements in the list and pass our receipt to them
	// we will collect the reward value and add to the global total reward.
	// we return the global total reward.

}

type Rule func(database.Receipt) int

func alphanumericRule(receipt database.Receipt) int {
	// One point for every alphanumeric character in the retailer name.
	retailerName := receipt.Retailer

	count := 0
	for _, char := range retailerName {
		if unicode.IsLetter(char) || unicode.IsDigit(char) {
			count++
		}
	}

	return count
}

func roundDollarRule(receipt database.Receipt) int {
	// 50 points if the total is a round dollar amount with no cents.
	receiptTotal := receipt.Total
	if math.Mod(receiptTotal, 1.0) == 0 {
		return 50
	}
	return 0
}

func multipleOfQuaterRule(receipt database.Receipt) int {
	// 25 points if the total is a multiple of 0.25
	receiptTotal := receipt.Total
	if math.Mod(receiptTotal, 0.25) == 0 {
		return 25
	}
	return 0
}

func pairRule(receipt database.Receipt) int {
	// 5 points for every two items on the receipt
	receiptItems := receipt.Items
	pairs := len(receiptItems) / 2
	points := pairs * 5
	return points
}

func multipleOfThreeRule(receipt database.Receipt) int {
	// If the trimmed length of the item description is a multiple of 3, multiply the price by 0.2 and round up to the nearest integer. The result is the number of points earned.
	receiptItems := receipt.Items
	reward := 0
	for _, item := range receiptItems {
		trimmedDescription := strings.TrimSpace(item.ShortDescription)
		lenTrimmedDescription := len(trimmedDescription)
		if math.Mod(float64(lenTrimmedDescription), 3) == 0 {
			reward += int(math.Ceil(item.Price * 0.2))
		}
	}
	return reward
}

func oddPurchaseDateRule(receipt database.Receipt) int {
	// 6 points if the day in the purchase date is odd.
	purchaseDate := receipt.PurchaseDate
	day := purchaseDate.Day()
	if day%2 == 1 {
		return 6
	}
	return 0
}

func offTimePurchaseRule(receipt database.Receipt) int {
	// 10 points if the time of purchase is after 2:00pm and before 4:00pm.
	// ct.Time.Format("15:04")
	// 0000-01-01
	purchaseTime := receipt.PurchaseTime

	year := 0
	month := time.January
	day := 1
	twoPm := time.Date(year, month, day, 14, 0, 0, 0, time.UTC)
	fourPm := time.Date(year, month, day, 16, 0, 0, 0, time.UTC)

	if purchaseTime.After(twoPm) && purchaseTime.Before(fourPm) {
		return 10
	}
	return 0
}
