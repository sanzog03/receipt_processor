package main

import (
	"log"
	"net/http"

	api "receiptProcessor/internal/api"
	"receiptProcessor/internal/repository"
	service "receiptProcessor/internal/services"

	_ "receiptProcessor/docs"
)

// @title Receipt Processor API
// @version 1.0
// @description This is a receipt processing service API, built as a part of Fech Rewards assessment challenge.
// @host localhost:9080
// @BasePath /
func main() {
	receiptRepository := repository.NewReceiptStore()
	receiptService := service.NewReceiptService(receiptRepository)
	router := api.SetupRoutes(receiptService)

	log.Println("Server starting on port 9080...")
	log.Fatal(http.ListenAndServe(":9080", router))
}
