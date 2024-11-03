package main

import (
	"log"
	"net/http"

	api "receiptProcessor/internal/api"
	"receiptProcessor/internal/repository"
	service "receiptProcessor/internal/services"
)

func main() {
	receiptRepository := repository.NewReceiptStore()
	receiptService := service.NewReceiptService(receiptRepository)
	router := api.SetupRoutes(receiptService)

	log.Println("Server starting on port 9080...")
	log.Fatal(http.ListenAndServe(":9080", router))
}
