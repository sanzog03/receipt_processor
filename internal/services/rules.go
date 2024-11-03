package service

import (
	"math"
	"receiptProcessor/internal/models"
	"strings"
	"time"
	"unicode"
)

type Rule func(models.Receipt) int

func alphanumericRule(receipt models.Receipt) int {
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

func roundDollarRule(receipt models.Receipt) int {
	// 50 points if the total is a round dollar amount with no cents.
	receiptTotal := receipt.Total
	if math.Mod(receiptTotal, 1.0) == 0 {
		return 50
	}
	return 0
}

func multipleOfQuaterRule(receipt models.Receipt) int {
	// 25 points if the total is a multiple of 0.25
	receiptTotal := receipt.Total
	if math.Mod(receiptTotal, 0.25) == 0 {
		return 25
	}
	return 0
}

func pairRule(receipt models.Receipt) int {
	// 5 points for every two items on the receipt
	receiptItems := receipt.Items
	pairs := len(receiptItems) / 2
	points := pairs * 5
	return points
}

func multipleOfThreeRule(receipt models.Receipt) int {
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

func oddPurchaseDateRule(receipt models.Receipt) int {
	// 6 points if the day in the purchase date is odd.
	purchaseDate := receipt.PurchaseDate
	day := purchaseDate.Day()
	if day%2 == 1 {
		return 6
	}
	return 0
}

func offTimePurchaseRule(receipt models.Receipt) int {
	// 10 points if the time of purchase is after 2:00pm and before 4:00pm.
	// ct.Time.Format("15:04")
	// For time, set date to default 0000-01-01
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
