package service

import (
	"fmt"
	"receiptProcessor/internal/models"
	"receiptProcessor/internal/repository"
	"time"
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
	fmt.Println("new receipt with id:", id, " and value ", receipt)
	fmt.Println("You have got: ", reward, " reward points.")
	return id, err
}

func (s *ReceiptService) ReceiptPoints(id string) (int, error) {
	receipt, error := s.store.Get(id)
	if error != nil {
		return -1, error
	}
	// find the points for the receipt and return the point
	fmt.Println("recepit", receipt)
	points, err := s.store.GetPoints(id)
	if err != nil {
		return -1, err
	}
	return points, nil
}

func (s *ReceiptService) generateId() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func (s *ReceiptService) calculateReward(receipt models.Receipt) int {
	rules := []Rule{
		alphanumericRule,
		roundDollarRule,
		multipleOfQuaterRule,
		pairRule,
		multipleOfThreeRule,
		oddPurchaseDateRule,
		offTimePurchaseRule,
	}

	total_reward := 0
	for _, rule := range rules {
		reward := rule(receipt)
		total_reward += reward
	}

	return total_reward
}
