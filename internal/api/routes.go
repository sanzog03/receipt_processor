package api

import (
	"receiptProcessor/internal/api/handlers"
	"receiptProcessor/internal/services"

	_ "receiptProcessor/docs"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

// setup http routes
func SetupRoutes(receiptService *services.ReceiptService) *mux.Router {
	router := mux.NewRouter()

	receiptHandler := handlers.NewReceiptHandler(receiptService)

	router.HandleFunc("/receipts/process", receiptHandler.ProcessReceipt).Methods("POST")
	router.HandleFunc("/receipts/{id}/points", receiptHandler.ReceiptPoints).Methods("GET")

	// go routes
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	return router
}
