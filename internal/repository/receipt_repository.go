package repository

import (
	"errors"
	"sync"

	"receiptProcessor/internal/models"
)

type ReceiptStore struct {
	receipts map[string]models.Receipt
	mu       sync.RWMutex
}

// Create a new instance of Receipt Store
func NewReceiptStore() *ReceiptStore {
	return &ReceiptStore{
		receipts: make(map[string]models.Receipt),
	}
}

// Set receipt with respect to its id in the ReceiptStore
func (rs *ReceiptStore) Set(id string, receipt models.Receipt) (string, error) {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	rs.receipts[id] = receipt
	return id, nil
}

// Get receipt with a id from the ReceiptStore
func (rs *ReceiptStore) Get(id string) (models.Receipt, error) {
	rs.mu.RLock()
	defer rs.mu.RUnlock()
	receipt, ok := rs.receipts[id]
	if !ok {
		return models.Receipt{}, errors.New("receipt not found")
	}
	return receipt, nil
}

// Delete the receipt with the id from the ReceiptStore
func (rs *ReceiptStore) delete(id string) (string, error) {
	rs.mu.Lock()
	defer rs.mu.Lock()
	delete(rs.receipts, id)
	return id, nil
}
