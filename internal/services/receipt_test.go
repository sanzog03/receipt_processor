package services_test

import (
	"encoding/json"
	"receiptProcessor/internal/api/handlers"
	"receiptProcessor/internal/models"
	"receiptProcessor/internal/repository"
	"receiptProcessor/internal/services"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

type testSetup struct {
	handler     *handlers.ReceiptHandler
	router      *mux.Router
	mockRepo    *repository.ReceiptStore
	mockService *services.ReceiptService
	testID      string
}

func setup() *testSetup {
	mockRepo := repository.NewReceiptStore()
	mockService := services.NewReceiptService(mockRepo)
	handler := handlers.NewReceiptHandler(mockService)

	router := mux.NewRouter()
	router.HandleFunc("/receipts/process", handler.ProcessReceipt).Methods("POST")
	router.HandleFunc("/receipts/{id}/points", handler.ReceiptPoints).Methods("GET")

	testID := uuid.New().String()
	id, _ := mockRepo.SetPoints(testID, 123)

	return &testSetup{
		handler:     handler,
		router:      router,
		mockRepo:    mockRepo,
		mockService: mockService,
		testID:      id,
	}
}

func TestServiceCalculatePoints(t *testing.T) {
	testSetup := setup()

	mockService := testSetup.mockService

	tests := []struct {
		name          string
		receipt       []byte
		expectedPoint int
	}{
		{
			name: "Receipt 1",
			receipt: []byte(`
				{
					"retailer": "Target",
					"purchaseDate": "2022-01-01",
					"purchaseTime": "13:01",
					"items": [
						{
						"shortDescription": "Mountain Dew 12PK",
						"price": "6.49"
						},{
						"shortDescription": "Emils Cheese Pizza",
						"price": "12.25"
						},{
						"shortDescription": "Knorr Creamy Chicken",
						"price": "1.26"
						},{
						"shortDescription": "Doritos Nacho Cheese",
						"price": "3.35"
						},{
						"shortDescription": "   Klarbrunn 12-PK 12 FL OZ  ",
						"price": "12.00"
						}
					],
					"total": "35.35"
					}
				`),
			expectedPoint: 28,
		},
		{
			name: "Receipt 2",
			receipt: []byte(`
				{
					"retailer": "M&M Corner Market",
					"purchaseDate": "2022-03-20",
					"purchaseTime": "14:33",
					"items": [
						{
						"shortDescription": "Gatorade",
						"price": "2.25"
						},{
						"shortDescription": "Gatorade",
						"price": "2.25"
						},{
						"shortDescription": "Gatorade",
						"price": "2.25"
						},{
						"shortDescription": "Gatorade",
						"price": "2.25"
						}
					],
					"total": "9.00"
					}
			`),
			expectedPoint: 109,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var receipt models.Receipt
			json.Unmarshal(test.receipt, &receipt)

			id, _ := mockService.ProcessReceipt(receipt)
			point, _ := mockService.ReceiptPoints(id)

			assert.Equal(t, test.expectedPoint, point, "Calculate Test Failed")
		})
	}
}
