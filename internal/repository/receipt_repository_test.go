package repository_test

import (
	handlers "receiptProcessor/internal/api/handler"
	"receiptProcessor/internal/repository"
	service "receiptProcessor/internal/services"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
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
	testID := uuid.New().String()
	id, err := mockRepo.SetPoints(testID, 20)
	t.Run("Receipt Store Set Points test", func(t *testing.T) {
		assert.NoError(t, err, "Expected no error for a valid input")
		assert.NotEmpty(t, id, "Expected a valid id which is not empty")
	})
}
