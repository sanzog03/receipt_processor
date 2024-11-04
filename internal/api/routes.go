package api

import (
	handlers "receiptProcessor/internal/api/handler"
	receipt "receiptProcessor/internal/services"

	_ "receiptProcessor/docs"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

// setup http routes
func SetupRoutes(receiptService *receipt.ReceiptService) *mux.Router {
	router := mux.NewRouter()

	receiptHandler := handlers.NewReceiptHandler(receiptService)

	router.HandleFunc("/receipt/process", receiptHandler.ProcessReceipt).Methods("POST")
	router.HandleFunc("/receipt/{id}/points", receiptHandler.ReceiptPoints).Methods("GET")

	// go routes
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	return router
}
