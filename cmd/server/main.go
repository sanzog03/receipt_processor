package main

import (
	"log"
	"net/http"

	api "receiptProcessor/internal/api"
	"receiptProcessor/internal/repository"
	service "receiptProcessor/internal/services"

	"github.com/spf13/viper"

	_ "receiptProcessor/docs"
)

// @title Receipt Processor API
// @version 1.0
// @description This is a receipt processing service API, built as a part of Fech Rewards assessment challenge.
// @host localhost:9080
// @BasePath /
func main() {
	loadconfig()
	receiptRepository := repository.NewReceiptStore()
	receiptService := service.NewReceiptService(receiptRepository)
	router := api.SetupRoutes(receiptService)

	port := viper.GetString("server.port")
	log.Println("Server starting on port " + port + "...")
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func loadconfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./Configs")
	err := viper.ReadInConfig()
	if err != nil {
		// Handle error
		log.Fatalf("Error loading config %s", err)
	}
}
