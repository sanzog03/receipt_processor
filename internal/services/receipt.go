package service

import (
	"receiptProcessor/internal/models"
	"receiptProcessor/internal/repository"

	"github.com/google/uuid"
)

// Provide services related to processing the receipts
type ReceiptService struct {
	store *repository.ReceiptStore
}

// create a new instance of ReceiptService
func NewReceiptService(store *repository.ReceiptStore) *ReceiptService {
	return &ReceiptService{store: store}
}

// Processes the receipt and return the receipt id
func (s *ReceiptService) ProcessReceipt(receipt models.Receipt) (string, error) {
	// create a id
	id := s.generateId()
	reward := s.calculateReward(receipt)
	_, err := s.store.Set(id, receipt)
	if err != nil {
		return id, err
	}
	_, err = s.store.SetPoints(id, reward)
	return id, err
}

func (s *ReceiptService) ReceiptPoints(id string) (int, error) {
	// find the points for the receipt and return the point
	points, err := s.store.GetPoints(id)
	if err != nil {
		return -1, err
	}
	return points, nil
}

func (s *ReceiptService) generateId() string {
	return uuid.New().String()
}

func (s *ReceiptService) calculateReward(receipt models.Receipt) int {
	rules := []Rule{
		Rule1{},
		Rule2{},
		Rule3{},
		Rule4{},
		Rule5{},
		Rule6{},
		Rule7{},
	}

	total_reward := 0
	for _, rule := range rules {
		reward := rule.process(receipt)
		total_reward += reward
	}

	return total_reward
}
