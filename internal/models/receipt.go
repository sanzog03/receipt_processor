package models

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

type Item struct {
	ShortDescription string  `json:"shortDescription" validate:"required,ValidateReceiptItemShortDesc"`
	Price            float64 `json:"price,string" validate:"required,ValidateReceiptItemPrice"`
}

type Receipt struct {
	Retailer     string     `json:"retailer" validate:"required,ValidateRetailerName"`
	PurchaseDate CustomTime `json:"purchaseDate" swaggertype:"string" validate:"required,datetime=2006-01-02"`
	PurchaseTime CustomTime `json:"purchaseTime" swaggertype:"string" validate:"required,datetime=15:04"`
	Items        []Item     `json:"items" validate:"required,min=1,dive,required"`
	Total        float64    `json:"total,string" validate:"required,ValidateReceiptTotal"`
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

// VALIDATORS

func validateRegex(pattern string, fieldToValidate string) bool {
	matched, _ := regexp.MatchString(pattern, fieldToValidate)
	return matched
}

// custom field validator for retailer name
func ValidateRetailerName(f validator.FieldLevel) bool {
	retailer := f.Field().String()
	pattern := `^[\w\s\-&]+$`
	return validateRegex(pattern, retailer)
}

// custom field validator for the total of a receipts
func ValidateReceiptTotal(f validator.FieldLevel) bool {
	total := f.Field().String()
	pattern := `^\d+\.\d{2}$`
	return validateRegex(pattern, total)
}

// custom field validator for short description field
func ValidateReceiptItemShortDesc(f validator.FieldLevel) bool {
	desc := f.Field().String()
	pattern := `^[\w\s\-]+$`
	return validateRegex(pattern, desc)
}

// custom field validator for the receipt item price
func ValidateReceiptItemPrice(f validator.FieldLevel) bool {
	price := f.Field().String()
	pattern := `\d+\.\d{2}$`
	return validateRegex(pattern, price)
}
