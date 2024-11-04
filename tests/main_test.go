package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	handlers "receiptProcessor/internal/api/handler"
	"receiptProcessor/internal/models"
	"receiptProcessor/internal/repository"
	service "receiptProcessor/internal/services"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testSetup struct {
	handler     *handlers.ReceiptHandler
	router      *mux.Router
	mockRepo    *repository.ReceiptStore
	mockService *service.ReceiptService
	testID      string
}

func setup() *testSetup {
	mockRepo := repository.NewReceiptStore()
	mockService := service.NewReceiptService(mockRepo)
	handler := handlers.NewReceiptHandler(mockService)

	router := mux.NewRouter()
	router.HandleFunc("/receipt/process", handler.ProcessReceipt).Methods("POST")
	router.HandleFunc("/receipt/{id}/points", handler.ReceiptPoints).Methods("GET")

	testID := fmt.Sprintf("%d", time.Now().UnixNano())
	id, _ := mockRepo.SetPoints(testID, 123)

	return &testSetup{
		handler:     handler,
		router:      router,
		mockRepo:    mockRepo,
		mockService: mockService,
		testID:      id,
	}
}

func TestReceiptHandlerProcessReceipt(t *testing.T) {
	testSetup := setup()

	tests := []struct {
		name           string
		payload        []byte
		expectedStatus int
		expectedBody   string
		checkJSONField string
	}{
		{
			name: "Valid Receipt",
			payload: []byte(`
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
			expectedStatus: http.StatusOK,
			checkJSONField: "id",
		},
		{
			name:           "Invalid Receipt",
			payload:        []byte(`{"invalid_json":`),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid receipt data\n",
		},
	}

	for _, testcase := range tests {
		t.Run(testcase.name, func(t *testing.T) {
			req, err := http.NewRequest("POST", "/receipt/process", bytes.NewBuffer(testcase.payload))
			require.NoError(t, err)

			responseRecorder := httptest.NewRecorder()
			testSetup.router.ServeHTTP(responseRecorder, req)

			assert.Equal(t, testcase.expectedStatus, responseRecorder.Code)

			if testcase.expectedBody != "" {
				assert.Equal(t, testcase.expectedBody, responseRecorder.Body.String())
			}

			if testcase.checkJSONField != "" {
				var responseBody map[string]interface{}
				err = json.Unmarshal(responseRecorder.Body.Bytes(), &responseBody)
				require.NoError(t, err, "Invalid receipt data")
				_, exists := responseBody[testcase.checkJSONField]
				assert.True(t, exists, "Response should contain a '%s' field", testcase.checkJSONField)
			}
		})
	}
}

func TestReceiptHandlerReceiptPointsByID(t *testing.T) {
	testSetup := setup()

	tests := []struct {
		name           string
		id             string
		expectedStatus int
		expectedBody   string
		checkJSONField string
	}{
		{
			name:           "Valid ID",
			id:             testSetup.testID,
			expectedStatus: http.StatusOK,
			checkJSONField: "points",
		},
		{
			name:           "Invalid ID",
			id:             "invalidID123Random",
			expectedStatus: http.StatusNotFound,
			expectedBody:   "Receipt not found\n"},
	}

	for _, testcase := range tests {
		req, err := http.NewRequest("GET", "/receipt/"+testcase.id+"/points", nil)
		require.NoError(t, err)

		responseRecorder := httptest.NewRecorder()
		testSetup.router.ServeHTTP(responseRecorder, req)

		assert.Equal(t, testcase.expectedStatus, responseRecorder.Code)

		if testcase.expectedBody != "" {
			assert.Equal(t, testcase.expectedBody, responseRecorder.Body.String())
		}

		if testcase.checkJSONField != "" {
			var responseBody map[string]interface{}
			err = json.Unmarshal(responseRecorder.Body.Bytes(), &responseBody)
			require.NoError(t, err, "Response body should be valid JSON")
			_, exists := responseBody[testcase.checkJSONField]
			assert.True(t, exists, "Response should contain a '%s' field", testcase.checkJSONField)
		}
	}
}

func TestRepositoryReceiptStoreGetPoints(t *testing.T) {
	testSetup := setup()

	tests := []struct {
		name     string
		id       string
		expected int
	}{
		{
			name:     "Repo Valid Id test",
			id:       testSetup.testID,
			expected: 123,
		},
		{
			name:     "Repo Invalid Id test",
			id:       "2asdf",
			expected: -1,
		},
	}

	for _, testcase := range tests {
		points, _ := testSetup.mockRepo.GetPoints(testcase.id)
		assert.Equal(t, testcase.expected, points, "Repo get points test failed")
	}
}

func TestRepositoryReceiptStoreSetPoints(t *testing.T) {
	testSetup := setup()

	mockRepo := testSetup.mockRepo
	testID := fmt.Sprintf("%d", time.Now().UnixNano())
	id, err := mockRepo.SetPoints(testID, 20)
	t.Run("Receipt Store Set Points test", func(t *testing.T) {
		assert.NoError(t, err, "Expected no error for a valid input")
		assert.NotEmpty(t, id, "Expected a valid id which is not empty")
	})
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
