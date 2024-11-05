package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"receiptProcessor/internal/api/handlers"
	"receiptProcessor/internal/repository"
	"receiptProcessor/internal/services"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
			req, err := http.NewRequest("POST", "/receipts/process", bytes.NewBuffer(testcase.payload))
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
		req, err := http.NewRequest("GET", "/receipts/"+testcase.id+"/points", nil)
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
