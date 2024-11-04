package repository

import (
	"errors"
	"sync"

	"receiptProcessor/internal/models"
)

type ReceiptStore struct {
	receipts      map[string]models.Receipt
	receiptsPoint map[string]int
	mu            sync.RWMutex
}

// Create a new instance of Receipt Store
func NewReceiptStore() *ReceiptStore {
	return &ReceiptStore{
		receipts:      make(map[string]models.Receipt),
		receiptsPoint: make(map[string]int),
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
func (rs *ReceiptStore) Delete(id string) (string, error) {
	rs.mu.Lock()
	defer rs.mu.Lock()
	delete(rs.receipts, id)
	return id, nil
}

func (rs *ReceiptStore) SetPoints(id string, points int) (string, error) {
	rs.receiptsPoint[id] = points
	return id, nil
}

func (rs *ReceiptStore) GetPoints(id string) (int, error) {
	points, exists := rs.receiptsPoint[id]
	if exists {
		return points, nil
	} else {
		return -1, errors.New("key does not exist")
	}
}
