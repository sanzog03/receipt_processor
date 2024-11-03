package models

import (
	"fmt"
	"strings"
	"time"
)

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

// Structure of the receipt reward points
type ReceiptPoints struct {
	Points int `json:"points"`
}

// Response structure of the processReceipt
type ReceiptResult struct {
	Id string `json:"id"`
}

// Response structure of the get points for a receipt
type PointsResult struct {
	Points int `json:"points"`
}

// utils
// Time struct to handle the input date time format in receipt
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
