package api

import (
	handlers "receiptProcessor/internal/api/handler"
	receipt "receiptProcessor/internal/services"

	"github.com/gorilla/mux"
)

func SetupRoutes(receiptService receipt.ReceiptService) *mux.Router {
	router := mux.NewRouter()

	receiptHandler := handlers.NewReceiptHandler(&receiptService)

	router.HandleFunc("/receipt/process", receiptHandler.ProcessReceipt).Methods("POST")
	router.HandleFunc("/receipt/{id}/points", receiptHandler.ReceiptPoints).Methods("GET")

	return router
}
