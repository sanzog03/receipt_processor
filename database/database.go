package database

import (
	"errors"
	"fmt"
	"strings"
	"sync"
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

type ReceiptStore struct {
	receipts map[string]Receipt
	mu       sync.RWMutex
}

func NewReceiptStore() *ReceiptStore {
	return &ReceiptStore{
		receipts: make(map[string]Receipt),
	}
}

func (rs *ReceiptStore) Set(id string, receipt Receipt) {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	rs.receipts[id] = receipt
}

func (rs *ReceiptStore) Get(id string) (Receipt, error) {
	rs.mu.RLock()
	defer rs.mu.RUnlock()
	receipt, ok := rs.receipts[id]
	if !ok {
		return Receipt{}, errors.New("Receipt not found")
	}
	return receipt, nil
}

func (rs *ReceiptStore) delete(id string) {
	rs.mu.Lock()
	defer rs.mu.Lock()
	delete(rs.receipts, id)
}
